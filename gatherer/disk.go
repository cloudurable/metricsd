package gatherer

import (
	l "github.com/cloudurable/simplelog/logging"
	c "github.com/cloudurable/metricsd/common"
	"strings"
)
var DiskFields_all = "all"

const (
	DiskField_total        = "total"
	DiskField_used         = "used"
	DiskField_available    = "available"
	DiskField_usedpct      = "usedpct"
	DiskField_availablepct = "availablepct"
	DiskField_capacitypct  = "capacitypct"
	DiskField_mount        = "mount"
)

var allFields = []string {
    DiskField_total,
    DiskField_used,
    DiskField_available,
    DiskField_usedpct,
    DiskField_availablepct,
    DiskField_capacitypct,
    DiskField_mount,
}

type DiskMetricsGatherer struct {
	logger       l.Logger
	command      string
	diskIncludes []diskInclude
	fields       []string
    diskAlarmThreshold float64
}

type diskInclude struct {
	starts bool
	value string
}

func NewDiskMetricsGatherer(logger l.Logger, config *c.Config) *DiskMetricsGatherer {

	logger = c.EnsureLogger(logger, config.Debug, c.GATHERER_DISK)
	command := 	c.ReadConfigString("df command", config.DiskCommand, "/usr/bin/df", logger)
	fields := c.ReadConfigStringArray("disk fields", config.DiskFields, []string{DiskField_availablepct}, logger, true)
    if fields[0] == DiskFields_all {
        fields = allFields
    }
	dfses := c.ReadConfigStringArray("disk file systems", config.DiskFileSystems, []string{"/dev/*"}, logger, false)

	var diskIncludes = []diskInclude{}
	for _, dfs := range dfses {
		if strings.HasSuffix(dfs, "*") {
			diskIncludes = append(diskIncludes, diskInclude{true, dfs[:len(dfs)-1]})
		} else {
			diskIncludes = append(diskIncludes, diskInclude{false, dfs})
		}
	}

    diskAlarmThreshold := float64(config.DiskAlarmThreshold)
    if diskAlarmThreshold <= 0 { diskAlarmThreshold = 101 }
	return &DiskMetricsGatherer{
		logger:       logger,
		command:      command,
		diskIncludes: diskIncludes,
		fields:       fields,
        diskAlarmThreshold: diskAlarmThreshold,
	}
}

func (disk *DiskMetricsGatherer) GetMetrics() ([]c.Metric, error) {

	output, err := c.ExecCommand(disk.command, "-P", "-k", "-l") // P for posix compatibility output, k for 1K blocks, l for local only
	if err != nil {
		return nil, err
	}

	var metrics = []c.Metric{}
	first := true // skip first line
	for _, line := range strings.Split(output, c.NEWLINE) {
		if first {
			first = false
		} else if disk.shouldReportDisk(line) {
			metrics = disk.appendDf(metrics, line)
		}
	}

	return metrics, nil

}

func (disk *DiskMetricsGatherer) shouldReportDisk(line string) bool {
	fsname := c.SplitGetFieldByIndex(line, 0)
	for _,include := range disk.diskIncludes {
		if include.starts {
			if strings.HasPrefix(fsname, include.value) {
				return true
			}
		} else {
			if fsname == include.value {
				return true
			}
		}
	}
	return false
}

func (disk *DiskMetricsGatherer) appendDf(metrics []c.Metric, line string) []c.Metric {

	// Filesystem     1024-blocks    Used Available Capacity Mounted on
	// udev               4019524       0   4019524       0% /dev
	// tmpfs               808140    9648    798492       2% /run
	// /dev/sda5         88339720 9322112  74507144      12% /
	// tmpfs              4040700  119244   3921456       3% /dev/shm

	valuesOnly := strings.Fields(line)
	name := valuesOnly[0]
	total := c.ToInt64(valuesOnly[1], 0)
	used := c.ToInt64(valuesOnly[2], 0)
	available := c.ToInt64(valuesOnly[3], 0)
	capacity := c.ToInt64( valuesOnly[4][0:len(valuesOnly[4])-1], 0)
	mount := valuesOnly[5]

	var totalF = float64(total)

	var upct = c.PercentRoundDp(float64(used), totalF, 2)
	var apct = c.PercentRoundDp(float64(available), totalF, 2)

	disk.logger.Debugf("name %s, total %d, used %d, available %d, usedpct %2.2f, availablepct %2.2f, capacity %d, mount %s\n",
		                name,    total,    used,    available,    upct,          apct,               capacity,    mount)

    alarm := upct >= disk.diskAlarmThreshold
    if alarm {
        m := *c.NewMetricFloat(c.MT_PERCENT, upct, "diskUsedPct:" + name, c.PROVIDER_DISK)
        m.Alarm = true
        metrics = append(metrics, m)
    }

	for _,field := range disk.fields {
		switch field {
		case DiskField_total:
			metrics = append(metrics, *c.NewMetricInt(c.MT_SIZE_KB, total, "diskTotal:" + name, c.PROVIDER_DISK))
		case DiskField_used:
			metrics = append(metrics, *c.NewMetricInt(c.MT_SIZE_KB, used, "diskUsed:" + name, c.PROVIDER_DISK))
		case DiskField_available:
            metrics = append(metrics, *c.NewMetricInt(c.MT_SIZE_KB, available, "diskAvailable:"+name, c.PROVIDER_DISK))
		case DiskField_usedpct:
            if !alarm { // otherwise it was already taken care of
                metrics = append(metrics, *c.NewMetricFloat(c.MT_PERCENT, upct, "diskUsedPct:" + name, c.PROVIDER_DISK))
            }
		case DiskField_availablepct:
			metrics = append(metrics, *c.NewMetricFloat(c.MT_PERCENT, apct, "diskAvailPct:" + name, c.PROVIDER_DISK))
		case DiskField_capacitypct:
			metrics = append(metrics, *c.NewMetricInt(c.MT_PERCENT, capacity, "diskCapacityPct:" + name, c.PROVIDER_DISK))
		case DiskField_mount:
			metrics = append(metrics, *c.NewMetricString(mount, "diskAvailMount:" + name, c.PROVIDER_DISK))
		}
	}

	return metrics
}
