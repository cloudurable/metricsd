package common

import (
    "time"
    "fmt"
    "strings"
)

type MetricsGatherer interface {
	Gather() ([]Metric, error)
    Name() string
}

type MetricsRepeater interface {
    Repeat(metrics []Metric) error
    Verify() bool
    Name() string
}

type MetricsAlarmer interface {
    Alarm(metrics []Metric) error
    Verify() bool
    Name() string
}

type MetricAlarm struct {
    Description       string
    Threshold         float64
    PeriodSeconds     int64
    EvaluationPeriods int64
    Comparison        MetricAlarmComparisonType
    Statistic         MetricAlarmStatType
}

func (this MetricAlarm) Empty() bool {
    return this.Description == EMPTY && this.Threshold == 0 && this.PeriodSeconds == 0 && this.EvaluationPeriods == 0 &&
        this.Comparison == MACT_NONE && this.Statistic == MAST_NONE
}

type Metric struct {
	Type       MetricType
    Source     MetricValueSource
	IntValue   int64
	FloatValue float64
	StrValue   string
	Name       string
	Provider   string
    When       time.Time
    Alarm      MetricAlarm
}

func (this Metric) MetricString() string {
    s := this.MetricStringFormatted()
    s = strings.Replace(s, ",\r\n    ", ", ", -1)
    s = strings.Replace(s, "\r\n    ", "", -1)
    return strings.Replace(s, "\r\n", "", -1)
}

func (this Metric) MetricStringFormatted() string {
    w := this.When.Format(STD_TIME_FORMAT)
    return fmt.Sprintf("Metric {\r\n    Name: \"%s\",\r\n    Type: %d,\r\n    Source: %d,\r\n    IntValue: %d,\r\n    FloatValue: %2.2f,\r\n    StrValue: \"%s\",\r\n    Provider: \"%s\",\r\n   When: %s,\r\n    Alarm: %s\r\n}",
        this.Type, this.Source, this.IntValue, this.FloatValue, this.StrValue, this.Name, this.Provider, w, this.Alarm.Description)
}

func newMetric(mt MetricType, mvs MetricValueSource, name string, provider string) *Metric {
	return &Metric{
		Type:     mt,
		Source:   mvs,
		Name:     name,
		Provider: provider,
        When:     time.Now(),
	}
}

func NewMetricInt(mt MetricType, value int64, name string, provider string) *Metric {
	m := newMetric(mt, MVS_INT, name, provider)
	m.IntValue = value
	m.FloatValue = float64(value)
	m.StrValue = Int64ToString(value)
	return m
}

func NewMetricFloat(mt MetricType, value float64, name string, provider string) *Metric {
	m := newMetric(mt, MVS_FLOAT, name, provider)
	m.IntValue = int64(RoundDp(value, 0))
	m.FloatValue = value
	m.StrValue = Float64ToString(value)
	return m
}

func NewMetricIntString(mt MetricType, value string, name string, provider string) *Metric {
	return newMetricString(mt, MVS_INT, value, name, provider, VALUE_ERROR)
}

func NewMetricFloatString(mt MetricType, value string, name string, provider string) *Metric {
	return newMetricString(mt, MVS_FLOAT, value, name, provider, VALUE_ERROR)
}

func NewMetricString(value string, name string, provider string) *Metric {
	return newMetricString(MT_NONE, MVS_STRING, value, name, provider, VALUE_N_A)
}

func newMetricString(mt MetricType, mvs MetricValueSource, value string, name string, provider string, errorValue int64) *Metric {
	m := newMetric(mt, mvs, name, provider)
	m.StrValue = value
	if value == IN_VALUE_N_A {
		m.IntValue = VALUE_N_A
		m.FloatValue = float64(VALUE_N_A)
	} else if value == IN_VALUE_NAN {
		m.IntValue = VALUE_NAN
		m.FloatValue = float64(VALUE_NAN)
	} else if mvs == MVS_INT {
		m.IntValue = ToInt64(value, errorValue)
		m.FloatValue = float64(m.IntValue)
	} else if mvs == MVS_FLOAT {
		m.FloatValue = ToFloat64(value, float64(errorValue))
		m.IntValue = int64(m.FloatValue)
	} else {
		m.IntValue = ToInt64(value, errorValue)
		m.FloatValue = ToFloat64(value, float64(errorValue))
	}
	return m
}

func NewMetricStringCode(mt MetricType, value string, code int64, name string, provider string) *Metric {
	m := newMetric(mt, MVS_STRING, name, provider)
	m.StrValue = value
	m.IntValue = code
	m.FloatValue = float64(code)
	return m
}
