package test

import (
	g "github.com/cloudurable/metricsd/gatherer"
	c "github.com/cloudurable/metricsd/common"
	"testing"
)

func TestFree(test *testing.T) {
	config := c.Config{ Debug: true, FreeCommand: "free"}
	StandardTest(test, g.NewFreeMetricGatherer(nil, &config))
}
