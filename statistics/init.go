package statistics

import (
	"github.com/fnproject/fn/api/server"
	"github.com/fnproject/fn/fnext"
)

// These functions are provided to support loading this extension via a custom ext.yaml file

func StatisticsExtensionName() string {
	return "github.com/fnproject/ext-metrics/statistics"
}

func init() {
	server.RegisterExtension(&statisticsExt{})
}

type statisticsExt struct {
}

func (e *statisticsExt) Name() string {
	return StatisticsExtensionName()
}

func (e *statisticsExt) Setup(s fnext.ExtServer) error {
	AddEndpoints(s)
	return nil
}
