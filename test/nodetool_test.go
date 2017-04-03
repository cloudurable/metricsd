package test

import (
	g "github.com/cloudurable/metricsd/gatherer"
	c "github.com/cloudurable/metricsd/common"
	nt "github.com/cloudurable/metricsd/gatherer/nodetool"
	"testing"
)

func TestNodetool(test *testing.T) {

	config := c.Config{
		Debug: true,
		NodetoolFunctions: []string{"notfound"},
	}

	gatherers := g.NewNodetoolMetricGatherers(nil, &config)
	if len(gatherers) > 0 {
		test.Errorf("Expecting no gatherers")
	}

	for _,gatherer := range gatherers {
		StandardTest(test, gatherer)
	}

	config = c.Config{
		Debug: false,
		NodetoolFunctions: nt.NodetoolAllSupportedFunctions,
	}

	gatherers = g.NewNodetoolMetricGatherers(nil, &config)
	for _,gatherer := range gatherers {
		StandardTest(test, gatherer)
	}
}
