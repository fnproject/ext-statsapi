package stats

import (
	"github.com/fnproject/fn/api/server"
	"github.com/fnproject/fn/fnext"
)

// These functions are provided to support loading this extension via a custom ext.yaml file

// StatisticsExtensionName returns the name of the Stats API extension
func StatisticsExtensionName() string {
	return "github.com/fnproject/ext-statsapi/stats"
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
