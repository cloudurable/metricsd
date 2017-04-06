package test

import (
	g "github.com/cloudurable/metricsd/gatherer"
	c "github.com/cloudurable/metricsd/common"
	"testing"
)

func TestDisk(test *testing.T) {

	config := c.Config{
		Debug: true,
        DiskConfig: c.DiskGathererConfig{
            Command: "df",
            FileSystems: []string{"/dev/*", "udev", "C:*"},
            Fields: []string{"total", "used", "available", "usedpct", "availablepct", "capacitypct", "mount"},
        },
	}

	StandardTest(test, g.NewDiskMetricsGatherer(nil, &config))
}
