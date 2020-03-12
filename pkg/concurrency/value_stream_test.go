package concurrency

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stackrox/rox/pkg/sync"
	"github.com/stretchr/testify/assert"
)

func TestValueStream_SequentialSync(t *testing.T) {
	t.Parallel()

	vs := NewValueStream(1)
	vs.Push(2)
	vs.Push(3)
	vs.Push(4)

	iter := vs.Iterator(true) // start observing
	vs.Push(5)
	vs.Push(6)
	vs.Push(7)
	vs.Push(8)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	expect := 4
	var err error
	for err == nil {
		fmt.Println("Value is", iter.Value())
		assert.Equal(t, expect, iter.Value())
		expect++
		if expect == 8 {
			break
		}
		iter, err = iter.Next(ctx)
	}
	assert.NoError(t, err)
}

func TestValueStream_SequentialAsync(t *testing.T) {
	t.Parallel()

	vs := NewValueStream(1)
	iter := vs.Iterator(true) // start observing

	go func() {
		time.Sleep(1 * time.Second)
		for val := 2; val < 8; val++ {
			vs.Push(val)
			time.Sleep(50 * time.Millisecond)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	expect := 1
	var err error
	for err == nil {
		fmt.Println("Value is", iter.Value())
		assert.Equal(t, expect, iter.Value())
		expect++
		if expect == 8 {
			break
		}
		iter, err = iter.Next(ctx)
	}
	assert.NoError(t, err)
}

func receive(ctx context.Context, t *testing.T, start ValueStreamIter, num int, store map[int]struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	it := start
	for i := 0; i < num; i++ {
		var err error
		it, err = it.Next(ctx)
		if !assert.NoError(t, err) {
			return
		}
		store[it.Value().(int)] = struct{}{}
	}
}

func TestValueStream_ParallelAsync(t *testing.T) {
	t.Parallel()

	vs := NewValueStream(0)
	start := vs.Iterator(true) // start observing

	var wg sync.WaitGroup
	wg.Add(2)

	recvd1 := make(map[int]struct{})
	recvd2 := make(map[int]struct{})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	go receive(ctx, t, start, 200, recvd1, &wg)
	go receive(ctx, t, start, 200, recvd2, &wg)

	for _, even := range []bool{false, true} {
		go func(even bool) {
			ofs := 1
			if even {
				ofs = 2
			}
			for i := 0; i < 100; i++ {
				go vs.Push(2*i + ofs)
			}
		}(even)
	}

	wg.Wait()

	assert.Len(t, recvd1, 200)
	assert.Equal(t, recvd1, recvd2)

	for i := 1; i <= 200; i++ {
		delete(recvd1, i)
	}
	assert.Empty(t, recvd1)
}

func TestValueStream_NonStrict(t *testing.T) {
	t.Parallel()

	vs := NewValueStream(0)
	it := vs.Iterator(false)

	go func() {
		for i := 1; i <= 10; i++ {
			time.Sleep(50 * time.Millisecond)
			vs.Push(i)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var err error
	lastVal := -2
	for err == nil {
		// Make sure we skip at least one value (but we can never exceed 10)
		minNextVal := lastVal + 2
		if minNextVal > 10 {
			minNextVal = 10
		}
		assert.GreaterOrEqual(t, it.Value(), minNextVal)
		lastVal = it.Value().(int)

		if lastVal >= 10 {
			break
		}

		time.Sleep(150 * time.Millisecond)
		it, err = it.Next(ctx)
	}

	assert.NoError(t, err)
}

func TestValueStream_SubscribeChan(t *testing.T) {
	t.Parallel()

	vs := NewValueStream(0)

	ch := make(chan interface{})
	errSig := NewErrorSignal()
	defer errSig.Signal()

	go errSig.SignalWithErrorWhen(context.DeadlineExceeded, TimeoutOr(1*time.Second, &errSig), &errSig)

	subscribeErrC := make(chan error)
	go func() {
		subscribeErrC <- vs.SubscribeChan(&errSig, ch, true)
	}()

	go func() {
		for i := 1; i <= 10; i++ {
			time.Sleep(10 * time.Millisecond)
			vs.Push(i)
		}
	}()

	for i := 0; i <= 10; i++ {
		select {
		case val := <-ch:
			assert.Equal(t, i, val)
		case <-errSig.Done():
			assert.Fail(t, "error signal should not expire")
		}
	}
	errSig.Signal()

	subscribeErr := <-subscribeErrC
	assert.NoError(t, subscribeErr)
}

func TestValueStream_SubscribeChanTyped(t *testing.T) {
	t.Parallel()

	vs := NewValueStream(0)

	ch := make(chan int)
	errSig := NewErrorSignal()
	defer errSig.Signal()

	go errSig.SignalWithErrorWhen(context.DeadlineExceeded, TimeoutOr(1*time.Second, &errSig), &errSig)

	subscribeErrC := make(chan error)
	go func() {
		subscribeErrC <- vs.SubscribeChanTyped(&errSig, ch, true)
	}()

	go func() {
		for i := 1; i <= 10; i++ {
			time.Sleep(10 * time.Millisecond)
			vs.Push(i)
		}
	}()

	for i := 0; i <= 10; i++ {
		select {
		case val := <-ch:
			assert.Equal(t, i, val)
		case <-errSig.Done():
			assert.Fail(t, "error signal should not expire")
		}
	}
	errSig.Signal()

	subscribeErr := <-subscribeErrC
	assert.NoError(t, subscribeErr)
}
