package scrapeupdate

import (
	"github.com/stackrox/rox/central/scrape/sensor/accept"
	"github.com/stackrox/rox/central/sensor/service/pipeline"
	"github.com/stackrox/rox/generated/internalapi/central"
	"github.com/stackrox/rox/pkg/logging"
)

var (
	log = logging.LoggerForModule()
)

// Template design pattern. We define control flow here and defer logic to subclasses.
//////////////////////////////////////////////////////////////////////////////////////

// NewPipeline returns a new instance of Pipeline.
func NewPipeline() pipeline.Fragment {
	return &pipelineImpl{
		accepter: accept.SingletonAccepter(),
	}
}

type pipelineImpl struct {
	accepter accept.Accepter
}

func (s *pipelineImpl) Match(msg *central.MsgFromSensor) bool {
	return msg.GetScrapeUpdate() != nil
}

// Run runs the pipeline template on the input and returns the output.
func (s *pipelineImpl) Run(msg *central.MsgFromSensor, _ pipeline.MsgInjector) (err error) {
	s.accepter.AcceptUpdate(msg.GetScrapeUpdate())
	return nil
}

func (s *pipelineImpl) OnFinish() {
	s.accepter.OnFinish()
}
