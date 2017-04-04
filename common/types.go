package common

import (
    "time"
    "fmt"
)

type MetricContext interface {
	GetEnv() string
	GetNameSpace() string
	GetRole() string
	SendId() bool
}

type MetricsGatherer interface {
	GetMetrics() ([]Metric, error)
}

type MetricsRepeater interface {
    ProcessMetrics(context MetricContext, metrics []Metric) error
    RepeatForContext() bool
    RepeatForNoIdContext() bool
    Verify() bool
}

type MetricsAlarmer interface {
    ProcessMetrics(context MetricContext, metrics []Metric) error
    Verify() bool
}

type Metric struct {
	Type       MetricType
    Source     MetricValueSource
	IntValue   int64
	FloatValue float64
	StrValue   string
	Name       string
	Provider   string
    Alarm      bool
    When       time.Time
}

func (m Metric) MetricFormattedString() string {
    w := m.When.Format(STD_TIME_FORMAT)
    a := "false"; if m.Alarm { a = "true" }
    return fmt.Sprintf("Metric {\r\n    Type:%d,\r\n    Source:%d,\r\n    IntValue:%d,\r\n    FloatValue:%2.2f,\r\n    StrValue:\"%s\",\r\n    Name:\"%s\",\r\n    Provider:\"%s\",\r\n    Alarm: %s,\r\n    When:%s\r\n}",
        m.Type, m.Source, m.IntValue, m.FloatValue, m.StrValue, m.Name, m.Provider, a, w)
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
