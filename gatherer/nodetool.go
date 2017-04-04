package gatherer

import (
	l "github.com/cloudurable/simplelog/logging"
	c "github.com/cloudurable/metricsd/common"
	nt "github.com/cloudurable/metricsd/gatherer/nodetool"
	"strings"
)

type NodetoolMetricGatherer struct {
	logger          l.Logger
    cqlshCommand    string
	nodetoolCommand string
	nodeFunction    string
}

func nodetoolFunctionSupported(nodeFunction string) bool {
	lower := strings.ToLower(nodeFunction)
	for _,supported := range nt.NodetoolAllSupportedFunctions {
		if supported == lower {
			return true
		}
	}
	return false
}

func NewNodetoolMetricGatherers(logger l.Logger, config *c.Config) []*NodetoolMetricGatherer {

	if config.NodetoolFunctions == nil || len(config.NodetoolFunctions) == 0 {
		return nil
	}

	logger = c.EnsureLogger(logger, config.Debug, c.GATHERER_NODETOOL)

	gatherers := []*NodetoolMetricGatherer{}
	for _, nodeFunction := range config.NodetoolFunctions {
		if nodetoolFunctionSupported(nodeFunction) {
			gatherers = append(gatherers, newNodetoolMetricGatherer(logger, config, nodeFunction))
		} else {
			logger.Warn("Unsupported or unknown Nodetool function", &nodeFunction)
		}
	}

	return gatherers
}

func newNodetoolMetricGatherer(logger l.Logger, config *c.Config, nodeFunction string) *NodetoolMetricGatherer {
    cqlshCommand := c.ReadConfigString("cqlsh command", config.CqlshCommand, "/usr/bin/cqlsh", logger)
    nodetoolCommand := c.ReadConfigString("nodetool command", config.CqlshCommand, "/usr/bin/nodetool", logger)

	return &NodetoolMetricGatherer{
		logger:          logger,
        cqlshCommand:    cqlshCommand,
        nodetoolCommand: nodetoolCommand,
		nodeFunction:    strings.ToLower(nodeFunction),
	}
}

func (gatherer *NodetoolMetricGatherer) GetMetrics() ([]c.Metric, error) {

	switch gatherer.nodeFunction {
	case nt.NtFunc_netstats:		    return nt.Netstats(gatherer.nodetoolCommand)
	case nt.NtFunc_gcstats:			    return nt.Gcstats(gatherer.nodetoolCommand)
	case nt.NtFunc_tpstats:			    return nt.Tpstats(gatherer.nodetoolCommand)
	case nt.NtFunc_getlogginglevels:    return nt.Getlogginglevels(gatherer.nodetoolCommand)
	case nt.NtFunc_gettimeout:	        return nt.Gettimeout(gatherer.nodetoolCommand)
	case nt.NtFunc_cfstats:	            return nt.Cfstats(gatherer.nodetoolCommand)
	case nt.NtFunc_proxyhistograms:     return nt.ProxyHistograms(gatherer.nodetoolCommand)
	case nt.NtFunc_listsnapshots:       return nt.ListSnapshots(gatherer.nodetoolCommand)
	case nt.NtFunc_statuses:            return nt.Statuses(gatherer.nodetoolCommand)
    case nt.NtFunc_getstreamthroughput: return nt.GetStreamThroughput(gatherer.nodetoolCommand)
	}

	return nil, nil
}
