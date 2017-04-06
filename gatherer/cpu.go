package gatherer

import (
	l "github.com/cloudurable/simplelog/logging"
	c "github.com/cloudurable/metricsd/common"
	"bufio"
	"os"
	"strings"
)

type CPUMetricsGatherer struct {
	procStatPath string
	lastStats    *CpuStats
	logger       l.Logger
	reportZeros  bool
}

type CpuTimes struct {
	User      int64
	Nice      int64
	System    int64
	Idle      int64
	IoWait    int64
	Irq       int64
	SoftIrq   int64
	Steal     int64
	Guest     int64
	GuestNice int64
}

type CpuStats struct {
	CpuMap              map[string]CpuTimes
	ContextSwitchCount  int64
	BootTime            int64
	ProcessCount        int64
	ProcessRunningCount int64
	ProcessBlockCount   int64
	InterruptCount      int64
	SoftInterruptCount  int64
}

func NewCPUMetricsGatherer(logger l.Logger, config *c.Config) *CPUMetricsGatherer {

	logger = c.EnsureLogger(logger, config.Debug, c.GATHERER_CPU)
	procStatPath := c.ReadConfigString("proc/stat path", config.CpuConfig.ProcStat, "/proc/stat", logger)

	return &CPUMetricsGatherer{
		procStatPath: procStatPath,
		logger:       logger,
		reportZeros:  config.CpuConfig.ReportZeros,
	}
}

func (this *CPUMetricsGatherer) TestingChangeProcStatPath(inProcStatPath string) {
	this.procStatPath = inProcStatPath
}

func (this *CPUMetricsGatherer) Name() string {
    return c.GATHERER_CPU
}

func (this *CPUMetricsGatherer) Gather() ([]c.Metric, error) {

	var cpuStats *CpuStats
	var err error

	if cpuStats, err = this.readProcStat(); err != nil {
		return nil, err
	}

	metrics := this.convertToMetrics(this.lastStats, cpuStats)
	this.lastStats = cpuStats
	return metrics, nil
}

func (this *CPUMetricsGatherer) appendCount(metrics []c.Metric, name string, count int64) []c.Metric {
	if this.reportZeros || count > 0 {
		metrics = append(metrics, *c.NewMetricInt(c.MT_COUNT, count, name, c.PROVIDER_CPU))
	}
	return metrics
}

func (this *CPUMetricsGatherer) convertToMetrics(lastTimeStats *CpuStats, nowStats *CpuStats) []c.Metric {
	var metrics = []c.Metric{}

	if lastTimeStats != nil {

		metrics = this.appendCount(metrics, "cpuSoftIrqCnt", int64(nowStats.SoftInterruptCount - lastTimeStats.SoftInterruptCount))
		metrics = this.appendCount(metrics, "cpuIntrCnt", int64(nowStats.InterruptCount - lastTimeStats.InterruptCount))
		metrics = this.appendCount(metrics, "cpuCtxtCnt", int64(nowStats.ContextSwitchCount - lastTimeStats.ContextSwitchCount))
		metrics = this.appendCount(metrics, "cpuProcessesStrtCnt", int64(nowStats.ProcessCount - lastTimeStats.ProcessCount))

		for cpuName, nowCpuTimes := range nowStats.CpuMap {
			lastCpuTimes, found := lastTimeStats.CpuMap[cpuName]
			if !found {
				lastCpuTimes = CpuTimes{}
			}
			metrics = this.appendCount(metrics, "cpuGuestJif", int64(nowCpuTimes.Guest - lastCpuTimes.Guest))
			metrics = this.appendCount(metrics, "cpuUsrJif", int64(nowCpuTimes.User - lastCpuTimes.User))
			metrics = this.appendCount(metrics, "cpuIdleJif", int64(nowCpuTimes.Idle - lastCpuTimes.Idle))
			metrics = this.appendCount(metrics, "cpuIowaitJif", int64(nowCpuTimes.IoWait - lastCpuTimes.IoWait))
			metrics = this.appendCount(metrics, "cpuIrqJif", int64(nowCpuTimes.Irq - lastCpuTimes.Irq))
			metrics = this.appendCount(metrics, "cpuGuestniceJif", int64(nowCpuTimes.GuestNice - lastCpuTimes.GuestNice))
			metrics = this.appendCount(metrics, "cpuStealJif", int64(nowCpuTimes.Steal - lastCpuTimes.Steal))
			metrics = this.appendCount(metrics, "cpuNiceJif", int64(nowCpuTimes.Nice - lastCpuTimes.Nice))
			metrics = this.appendCount(metrics, "cpuSysJif", int64(nowCpuTimes.System - lastCpuTimes.System))
			metrics = this.appendCount(metrics, "cpuSoftIrqJif", int64(nowCpuTimes.SoftIrq - lastCpuTimes.SoftIrq))
		}
	}

	metrics = append(metrics, *c.NewMetricInt(c.MT_COUNT, int64(nowStats.ProcessRunningCount), "cpuProcsRunning", c.PROVIDER_CPU))
	metrics = append(metrics, *c.NewMetricInt(c.MT_COUNT, int64(nowStats.ProcessBlockCount), "cpuProcsBlocked", c.PROVIDER_CPU))

	return metrics
}

func (this *CPUMetricsGatherer) readProcStat() (*CpuStats, error) {
	org, err := os.Open(this.procStatPath)
	fd := bufio.NewReader(org)
	if err != nil {
		this.logger.Emergencyf("Error reading file %v", err)
	}

	stats := CpuStats{}
	stats.CpuMap = make(map[string]CpuTimes)

	for {
		theLine := c.EMPTY
		bytes, _, err := fd.ReadLine()
		if err != nil {
			if err.Error() != "EOF" { // EOF error is ok, other errors are not ok
				this.logger.PrintError("Error reading line from /proc/stat", err)
				return nil, err
			}
			break // EOF, leave the for loop
		}
		if len(bytes) == 0 {
			break
		}

		theLine = string(bytes)

		valuesOnly := strings.Fields(theLine)
		lineName := valuesOnly[0]
		value := c.ToInt64(valuesOnly[1], 0)

		switch lineName {
		case "ctxt":          stats.ContextSwitchCount = value
		case "btime":         stats.BootTime = value
		case "processes":     stats.ProcessCount = value
		case "procs_running": stats.ProcessRunningCount = value
		case "procs_blocked": stats.ProcessBlockCount = value
		case "intr":          stats.InterruptCount = value
		case "softirq":       stats.SoftInterruptCount = value
		default:
			if strings.HasPrefix(lineName, "cpu") {
				cpuTimes := CpuTimes{}
				for i := 1; i < len(valuesOnly); i++ {
					value = c.ToInt64(valuesOnly[i], 0)
					switch i {
					case  1: cpuTimes.User = value
					case  2: cpuTimes.Nice = value
					case  3: cpuTimes.System = value
					case  4: cpuTimes.Idle = value
					case  5: cpuTimes.IoWait = value
					case  6: cpuTimes.Irq = value
					case  7: cpuTimes.SoftIrq = value
					case  8: cpuTimes.Steal = value
					case  9: cpuTimes.Guest = value
					case 10: cpuTimes.GuestNice = value
					default:
                        this.logger.Debug("Unknown cpu time column, index:", i, "found in", theLine)
					}
				}
				stats.CpuMap[lineName] = cpuTimes
			} else {
                this.logger.Debug("Unknown Data", theLine)
			}
		}
	}

	return &stats, nil
}

	/*
	cpu  5017 1 3356 7561462 1674 2 53 3 4 5
	cpu0 1105 0 1113 1890502 345 0 35 0 0 0
	cpu1 1291 0 792 1889318 496 0 6 0 0 0
	cpu2 1251 0 713 1890968 482 0 5 0 0 0
	cpu3 1370 0 738 1890674 351 0 7 0 0 0
	intr 1488221 27 0 0 0 348 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 159 0 12047 0 22568 1 0 3143 0 76 0 91494 83893 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
	ctxt 2710734
	btime 1480893191
	processes 9277
	procs_running 3
	procs_blocked 0
	softirq 655105 0 221348 21 39766 0 0 1 215075 0 178894
	*/
