package errors

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/stackrox/rox/pkg/errorhelpers"
	"github.com/stackrox/rox/pkg/sac"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorToGrpcCodeInterceptor translates common errors defined in errorhelpers to GRPC codes.
func ErrorToGrpcCodeInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	return resp, ErrToGrpcStatus(err).Err()
}

// ErrorToGrpcCodeStreamInterceptor translates common errors defined in errorhelpers to GRPC codes.
func ErrorToGrpcCodeStreamInterceptor(srv interface{}, ss grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	err := handler(srv, ss)
	return ErrToGrpcStatus(err).Err()
}

// ErrToHTTPStatus maps known internal and gRPC errors to the appropriate
// HTTP status code.
func ErrToHTTPStatus(err error) int {
	return runtime.HTTPStatusFromCode(ErrToGrpcStatus(err).Code())
}

// ErrToGrpcStatus wraps an error into a gRPC status with code.
func ErrToGrpcStatus(err error) *status.Status {
	if s, ok := status.FromError(err); ok {
		// `err` is either nil or status.Status.
		return s
	}
	code := errorTypeToGrpcCode(err)
	return status.New(code, err.Error())
}

func errorTypeToGrpcCode(err error) codes.Code {
	switch {
	case errors.Is(err, errorhelpers.ErrNotFound):
		return codes.NotFound
	case errors.Is(err, errorhelpers.ErrInvalidArgs):
		return codes.InvalidArgument
	case errors.Is(err, errorhelpers.ErrAlreadyExists):
		return codes.AlreadyExists
	case errors.Is(err, errorhelpers.ErrReferencedByAnotherObject):
		return codes.FailedPrecondition
	case errors.Is(err, errorhelpers.ErrInvariantViolation):
		return codes.Internal
	case errors.Is(err, errorhelpers.ErrNotAuthenticated):
		return codes.Unauthenticated
	case errors.Is(err, errorhelpers.ErrNotAuthorized):
		return codes.PermissionDenied
	case errors.Is(err, sac.ErrResourceAccessDenied):
		return codes.PermissionDenied
	default:
		return codes.Internal
	}
}
