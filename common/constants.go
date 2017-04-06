package common

const STD_TIME_FORMAT = "2006-01-02T15:04:05.999-07:00"

const (
	EMPTY                   = ""
	SPACE                   = " "
    NEWLINE                 = "\n"
	UNDER                   = "_"
	DOT                     = "."
	QUOTE                   = "\""
	COMMA                   = ","
	COMMA_SPACE             = ", "
	QUOTE_COLON_SPACE       = "\": "
	QUOTE_COLON_SPACE_QUOTE = "\": \""
	QUOTE_COMMA_SPACE       = "\", "
	OPEN_BRACE              = "["
	CLOSE_BRACE             = "]"
	COLON                   = ":"
)

const (
	DEFAULT_LABEL = "Default"
	CONFIG_LABEL  = "Config"
)

const (
    PROVIDER_DISK      = "disk"
    PROVIDER_CPU       = "cpu"
    PROVIDER_FREE      = "free"
    PROVIDER_CASSANDRA = "cass"
)

const (
	GATHERER_DISK      = PROVIDER_DISK
    GATHERER_CPU       = PROVIDER_CPU
	GATHERER_FREE      = PROVIDER_FREE
	GATHERER_CASSANDRA = PROVIDER_CASSANDRA
)

const (
    REPEATER        = "repeater-"
    REPEATER_AWS    = "aws"
    REPEATER_LOG    = "log"
)

const (
    ALARMER          = "alarmer-"
    ALARMER_AWS      = REPEATER_AWS
    ALARMER_EMAIL    = "email"
)

type MetricValueSource int8
const (
	MVS_INT MetricValueSource = iota
	MVS_FLOAT
	MVS_STRING
)

func (mvs *MetricValueSource) Name() string {
	switch *mvs {
	case MVS_INT: return "int"
	case MVS_FLOAT: return "float"
	case MVS_STRING: return "str"
	}
	return EMPTY
}

type MetricType int8
const (
    MT_NONE MetricType = iota
    MT_COUNT
	MT_PERCENT
	MT_MICROS
	MT_MILLIS
	MT_SIZE_BYTE
	MT_SIZE_KB
	MT_SIZE_MB
	MT_SIZE_GB
	MT_SIZE_TB
)

func (mt MetricType) Name() string { // EXACT MATCHES FOR CLOUDWATCH CONSTANTS
	switch mt {
	case MT_COUNT:     return "Count"
	case MT_PERCENT:   return "Percent"
	case MT_MICROS:    return "Microseconds"
	case MT_MILLIS:    return "Milliseconds"
	case MT_SIZE_BYTE: return "Byte"
	case MT_SIZE_MB:   return "Megabytes"
	case MT_SIZE_KB:   return "Kilobytes"
	}
	return "None"
}

type MetricAlarmComparisonType int8
const (
    MACT_NONE MetricAlarmComparisonType = iota
    MACT_GTE
    MACT_GT
    MACT_LT
    MACT_LTE
)

func (mact MetricAlarmComparisonType) Name() string { // EXACT MATCHES FOR CLOUDWATCH CONSTANTS
    switch mact {
    case MACT_GTE: return "GreaterThanThreshold"
    case MACT_GT:  return "GreaterThanThreshold"
    case MACT_LT:  return "LessThanThreshold"
    case MACT_LTE: return "LessThanOrEqualToThreshold"
    }
    return "None"
}

type MetricAlarmStatType int8
const (
    MAST_NONE MetricAlarmStatType = iota
    MAST_SAMPLE
    MAST_AVG
    MAST_SUM
    MAST_MIN
    MAST_MAX
)

func (mast MetricAlarmStatType) Name() string {  // EXACT MATCHES FOR CLOUDWATCH CONSTANTS
    switch mast {
    case MAST_SAMPLE: return "SampleCount"
    case MAST_AVG:    return "Average"
    case MAST_SUM:    return "Sum"
    case MAST_MIN:    return "Minimum"
    case MAST_MAX:    return "Maximum"
    }
    return "None"
}

const (
	VALUE_N_A   int64 = -125
	VALUE_NAN   int64 = -126
	VALUE_ERROR int64 = -127
)

const (
	IN_VALUE_N_A = "n/a"
	IN_VALUE_NAN = "NaN"
)
