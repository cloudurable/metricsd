package metric

import (
	"os"
	"strings"
	"bufio"
	l "github.com/advantageous/metrics/logger"
	"fmt"
	"errors"
	"runtime"
)

type CpuTimeType byte
type CpuTime uint64
type CpuCount uint64

type CPUMetricsGatherer struct {
	path     string
	lastTime *CpuStats
	logger   l.Logger
}

type CpuStats struct {
	CpuTimeList         []CpuTimes
	ContextSwitchCount  CpuCount
	BootTime            CpuTime
	ProcessCount        CpuCount
	ProcessRunningCount CpuCount
	ProcessBlockCount   CpuCount
	InterruptCount      CpuCount
	SoftInterruptCount  CpuCount
}

type CpuTimes struct {
	Name      string
	User      CpuTime
	Nice      CpuTime
	System    CpuTime
	Idle      CpuTime
	IoWait    CpuTime
	Irq       CpuTime
	SoftIrq   CpuTime
	Steal     CpuTime
	Guest     CpuTime
	GuestNice CpuTime
}

func NewCPUMetricsGathererWithPath(path string, logger l.Logger) *CPUMetricsGatherer {

	if logger == nil {
		logger = l.GetSimpleLogger("MT_CPU_DEBUG", "cpu")
	}
	return &CPUMetricsGatherer{
		path: path,
		logger: logger,
	}
}

func NewCPUMetricsGatherer(logger l.Logger) *CPUMetricsGatherer {
	statFile := "/proc/stat"
	if logger == nil {
		logger = l.GetSimpleLogger("MT_CPU_DEBUG", "cpu")
	}
	if runtime.GOOS == "darwin" {
		dir, _ := os.Getwd()
		logger.Println("DIR", dir)
		statFile = dir + "/metric/test-data/proc/stat"
	}
	return NewCPUMetricsGathererWithPath(statFile, logger)
}


func (cpu *CPUMetricsGatherer) SetPath(path string) {
	cpu.path = path
}

func (cpu *CPUMetricsGatherer) GetMetrics() ([]Metric, error) {

	var cpuStats *CpuStats
	var err error

	if cpuStats, err = cpu.readCpuStats(); err != nil {
		return nil, err
	}

	metrics := convertToMetrics(cpu.lastTime, cpuStats)
	cpu.lastTime = cpuStats
	cpu.logger.Debugf("%+v \n", cpuStats)
	return metrics, nil

}

func convertToMetrics(lastTimeStats *CpuStats, nowStats *CpuStats) []Metric {
	var metrics = []Metric{}

	if lastTimeStats != nil {

		softInterruptCount := nowStats.SoftInterruptCount - lastTimeStats.SoftInterruptCount
		if softInterruptCount > 0 {
			metrics = append(metrics, metric{COUNT, MetricValue(softInterruptCount), "softirq"})
		}

		interruptCount := nowStats.InterruptCount - lastTimeStats.InterruptCount
		if interruptCount > 0 {
			metrics = append(metrics, metric{COUNT, MetricValue(interruptCount), "intr"})
		}

		contextSwitchCount := nowStats.ContextSwitchCount - lastTimeStats.ContextSwitchCount
		if contextSwitchCount > 0 {
			metrics = append(metrics, metric{COUNT, MetricValue(contextSwitchCount), "ctxt"})
		}

		processCount := nowStats.ProcessCount - lastTimeStats.ProcessCount
		if processCount > 0 {
			metrics = append(metrics, metric{COUNT, MetricValue(processCount), "processes"})
		}

		for index, c := range nowStats.CpuTimeList {

			guest := c.Guest - lastTimeStats.CpuTimeList[index].Guest
			if guest > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(guest), c.Name + "-guest"})
			}

			user := c.User - lastTimeStats.CpuTimeList[index].User
			if user > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(user), c.Name + "-usr"})
			}

			idle := c.Idle - lastTimeStats.CpuTimeList[index].Idle
			if idle > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(idle), c.Name + "-idle"})
			}

			IoWait := c.IoWait - lastTimeStats.CpuTimeList[index].IoWait
			if IoWait > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(IoWait), c.Name + "-iowait"})
			}

			Irq := c.Irq - lastTimeStats.CpuTimeList[index].Irq
			if Irq > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(Irq), c.Name + "-irq"})
			}

			GuestNice := c.GuestNice - lastTimeStats.CpuTimeList[index].GuestNice
			if GuestNice > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(GuestNice), c.Name + "-guestnice"})
			}

			Steal := c.Steal - lastTimeStats.CpuTimeList[index].Steal
			if Steal > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(Steal), c.Name + "-steal"})
			}

			Nice := c.Nice - lastTimeStats.CpuTimeList[index].Nice
			if Nice > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(Nice), c.Name + "-nice"})
			}

			System := c.System - lastTimeStats.CpuTimeList[index].System
			if System > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(System), c.Name + "-sys"})
			}

			SoftIrq := c.SoftIrq - lastTimeStats.CpuTimeList[index].SoftIrq
			if SoftIrq > 0 {
				metrics = append(metrics, metric{COUNT, MetricValue(SoftIrq), c.Name + "-softirq"})
			}

		}

	} else {
		metrics = append(metrics, metric{COUNT, MetricValue(nowStats.SoftInterruptCount), "softirq"})
		metrics = append(metrics, metric{COUNT, MetricValue(nowStats.InterruptCount), "intr"})
		metrics = append(metrics, metric{COUNT, MetricValue(nowStats.ContextSwitchCount), "ctxt"})
		metrics = append(metrics, metric{COUNT, MetricValue(nowStats.ContextSwitchCount), "processes"})

		for _, c := range nowStats.CpuTimeList {
			metrics = append(metrics, metric{COUNT, MetricValue(c.Guest), c.Name + "-guest"})
			metrics = append(metrics, metric{COUNT, MetricValue(c.User), c.Name + "-usr"})
			metrics = append(metrics, metric{COUNT, MetricValue(c.Idle), c.Name + "-idle"})
			metrics = append(metrics, metric{COUNT, MetricValue(c.IoWait), c.Name + "-iowait"})
			metrics = append(metrics, metric{COUNT, MetricValue(c.Irq), c.Name + "-irq"})
			metrics = append(metrics, metric{COUNT, MetricValue(c.GuestNice), c.Name + "-guestnice"})
			metrics = append(metrics, metric{COUNT, MetricValue(c.Steal), c.Name + "-steal"})
			metrics = append(metrics, metric{COUNT, MetricValue(c.Nice), c.Name + "-nice"})
			metrics = append(metrics, metric{COUNT, MetricValue(c.System), c.Name + "-sys"})
			metrics = append(metrics, metric{COUNT, MetricValue(c.SoftIrq), c.Name + "-softirq"})
		}
	}

	metrics = append(metrics, metric{COUNT, MetricValue(nowStats.ContextSwitchCount), "procs_running"})
	metrics = append(metrics, metric{COUNT, MetricValue(nowStats.ContextSwitchCount), "procs_blocked"})

	return metrics
}

func (cpu *CPUMetricsGatherer) readCpuStats() (*CpuStats, error) {
	org, err := os.Open(cpu.path)

	fd := bufio.NewReader(org)

	if err != nil {
		cpu.logger.Emergencyf("Error reading file %v", err)
	}

	stats := CpuStats{}
	stats.CpuTimeList = make([]CpuTimes, 0)

	for {
		var name string
		var value uint64
		var line string

		if bytes, _, err := fd.ReadLine(); err == nil {
			line = string(bytes)
		} else if err.Error() == "EOF" { //Error EOF is ok
			cpu.logger.Debug("EOF while reading /proc/stat file")
			break
		} else { //Others errors are not ok
			cpu.logger.PrintError("Error reading line from /proc/stat", err)
			return nil, err
		}

		if count, err := fmt.Sscanf(line, "%s %d ", &name, &value); err != nil {
			cpu.logger.PrintError("Error scanning name / value from a line from /proc/stat", err)
			return nil, err
		} else if count == 0 {
			cpu.logger.Debug("Count was 0 when scanning /proc/stat line")
			break

		}

		if err = cpu.parseLine(name, value, line, &stats); err != nil {
			return nil, err
		}

		cpu.logger.Debugf("%s %d", name, value)
	}
	cpu.logger.Debugf("%+v \n", stats)
	return &stats, nil
}

func (cpu *CPUMetricsGatherer) parseLine(name string, value uint64, line string, stats *CpuStats) error {

	switch (name) {
	case "ctxt":
		stats.ContextSwitchCount = CpuCount(value)
	case "btime":
		stats.BootTime = CpuTime(value)
	case "processes":
		stats.ProcessCount = CpuCount(value)
	case "procs_running":
		stats.ProcessBlockCount = CpuCount(value)
	case "procs_blocked":
		stats.ProcessBlockCount = CpuCount(value)
	case "intr":
		stats.InterruptCount = CpuCount(value)
	case "softirq":
		stats.SoftInterruptCount = CpuCount(value)
	default:
		if strings.HasPrefix(name, "cpu") {
			t := CpuTimes{}
			t.Name = name
			t.User = CpuTime(value)
			count, err := fmt.Sscanf(line, "%s %d %d %d %d %d %d %d %d %d %d", &t.Name,
				&t.User, &t.Nice, &t.System,
				&t.Idle, &t.IoWait, &t.Irq,
				&t.SoftIrq, &t.Steal, &t.Guest,
				&t.GuestNice)

			if err != nil {
				cpu.logger.PrintError("Failure parsing cpu stats", err)
				return err
			}

			if count != 11 {
				cpu.logger.Errorf("cpu scan amount is off expected 11, but got %d", count)
				return errors.New("Unable to scan cpu times")
			}
			stats.CpuTimeList = append(stats.CpuTimeList, t)
		} else {
			return fmt.Errorf("Not sure what this is %s", name)
		}
	}
	return nil
}