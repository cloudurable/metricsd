package nodetool

import (
	"strings"
	c "github.com/cloudurable/metricsd/common"
)

func Gcstats(nodetoolCommand string) ([]c.Metric, error) {
	output, err := c.ExecCommand(nodetoolCommand, NtFunc_gcstats)
	if err != nil {
		return nil, err
	}

	// -- sample gcstats output --
	//Interval (ms) Max GC Elapsed (ms)Total GC Elapsed (ms)Stdev GC Elapsed (ms)   GC Reclaimed (MB)         Collections      Direct Memory Bytes
	//3491665                   0                   0                 NaN                   0                   0                       -1

	lines := strings.Split(output, c.NEWLINE)
	values := strings.Fields(lines[1])

	var metrics = []c.Metric{}
	metrics = append(metrics, *c.NewMetricIntString(c.MT_MILLIS, values[0], "ntGcInterval", c.EMPTY, c.PROVIDER_CASSANDRA))
	metrics = append(metrics, *c.NewMetricIntString(c.MT_MILLIS, values[1], "ntGcMaxElapsed", c.EMPTY, c.PROVIDER_CASSANDRA))
	metrics = append(metrics, *c.NewMetricIntString(c.MT_MILLIS, values[2], "ntGcTotalElapsed", c.EMPTY, c.PROVIDER_CASSANDRA))
	metrics = append(metrics, *c.NewMetricIntString(c.MT_MILLIS, values[3], "ntGcStdevElapsed", c.EMPTY, c.PROVIDER_CASSANDRA))
	metrics = append(metrics, *c.NewMetricIntString(c.MT_SIZE_MB, values[4], "ntGcReclaimed", c.EMPTY, c.PROVIDER_CASSANDRA))
	metrics = append(metrics, *c.NewMetricIntString(c.MT_COUNT, values[5], "ntGcCollections", c.EMPTY, c.PROVIDER_CASSANDRA))
	metrics = append(metrics, *c.NewMetricIntString(c.MT_SIZE_BYTE, values[6], "ntGcDirectMemoryBytes", c.EMPTY, c.PROVIDER_CASSANDRA))

	return metrics, nil
}

