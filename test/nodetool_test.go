package test

import (
	g "github.com/cloudurable/metricsd/gatherer"
	c "github.com/cloudurable/metricsd/common"
	"testing"
)

func TestNodetool(test *testing.T) {

	config := c.Config{
		Debug: false,
		//NodetoolFunctions: []string{"cfstats", "tpstats", "gcstats", "getlogginglevels", "netstats", "gettimeout", "proxyhistograms", "listsnapshots", "statuses"},
		NodetoolFunctions: []string{"statuses"},
	}

	gatherers := g.NewNodetoolMetricGatherers(nil, &config)
	for _,gatherer := range gatherers {
		StandardTest(test, gatherer)
	}
}
