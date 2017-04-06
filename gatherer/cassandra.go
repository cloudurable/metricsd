package gatherer

import (
	l "github.com/cloudurable/simplelog/logging"
	c "github.com/cloudurable/metricsd/common"
	nt "github.com/cloudurable/metricsd/gatherer/nodetool"
	"strings"
)

type CassandraMetricGatherer struct {
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

func NewCassandraMetricGatherers(logger l.Logger, config *c.Config) []*CassandraMetricGatherer {

	if config.CassandraConfig.NodetoolFunctions == nil || len(config.CassandraConfig.NodetoolFunctions) == 0 {
		return nil
	}

	logger = c.EnsureLogger(logger, config.Debug, c.GATHERER_CASSANDRA)

	gatherers := []*CassandraMetricGatherer{}
	for _, nodeFunction := range config.CassandraConfig.NodetoolFunctions {
		if nodetoolFunctionSupported(nodeFunction) {
			gatherers = append(gatherers, newNodetoolMetricGatherer(logger, config, nodeFunction))
		} else {
			logger.Warn("Unsupported or unknown Nodetool function", &nodeFunction)
		}
	}

	return gatherers
}

func newNodetoolMetricGatherer(logger l.Logger, config *c.Config, nodeFunction string) *CassandraMetricGatherer {
    cqlshCommand := c.ReadConfigString("cqlsh command", config.CassandraConfig.CqlshCommand, "/usr/bin/cqlsh", logger)
    nodetoolCommand := c.ReadConfigString("nodetool command", config.CassandraConfig.CqlshCommand, "/usr/bin/nodetool", logger)

	return &CassandraMetricGatherer{
		logger:          logger,
        cqlshCommand:    cqlshCommand,
        nodetoolCommand: nodetoolCommand,
		nodeFunction:    strings.ToLower(nodeFunction),
	}
}

func (this *CassandraMetricGatherer) Gather() ([]c.Metric, error) {

	switch this.nodeFunction {
	case nt.NtFunc_netstats:		    return nt.Netstats(this.nodetoolCommand)
	case nt.NtFunc_gcstats:			    return nt.Gcstats(this.nodetoolCommand)
	case nt.NtFunc_tpstats:			    return nt.Tpstats(this.nodetoolCommand)
	case nt.NtFunc_getlogginglevels:    return nt.Getlogginglevels(this.nodetoolCommand)
	case nt.NtFunc_gettimeout:	        return nt.Gettimeout(this.nodetoolCommand)
	case nt.NtFunc_cfstats:	            return nt.Cfstats(this.nodetoolCommand)
	case nt.NtFunc_proxyhistograms:     return nt.ProxyHistograms(this.nodetoolCommand)
	case nt.NtFunc_listsnapshots:       return nt.ListSnapshots(this.nodetoolCommand)
	case nt.NtFunc_statuses:            return nt.Statuses(this.nodetoolCommand)
    case nt.NtFunc_getstreamthroughput: return nt.GetStreamThroughput(this.nodetoolCommand)
	}

	return nil, nil
}

func (this *CassandraMetricGatherer) Name() string {
    return c.GATHERER_CASSANDRA + c.DOT + this.nodeFunction
}
