package test

import (
	"testing"
	c "github.com/cloudurable/metricsd/common"
	"fmt"
)

func TestScratch(test *testing.T) {

	logger := GetTestLogger(test, "scratch")
	config := GetTestConfig(logger)

	fmt.Println(c.ToJsonLabeledString(config))

	config = &c.Config{
		Debug: false,
        DiskConfig: c.DiskGathererConfig{
            Command: "df",
            FileSystems: []string{"/dev/*", "udev"},
            Fields: []string{"total", "used", "available", "usedpct", "availablepct", "capacitypct", "mount"},
        },
	}

	fmt.Println(c.ToJsonLabeledString(config))
}
