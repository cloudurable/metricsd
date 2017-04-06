package common

import (
	l "github.com/cloudurable/simplelog/logging"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"math"
    "encoding/json"
    "reflect"
)

// ========================================================================================================================
// LOGGING HELPERS
// ========================================================================================================================
func GetLogger(debug bool, name string) l.Logger {
    return EnsureLogger(nil, debug, name)
}

func EnsureLogger(logger l.Logger, debug bool, name string) l.Logger {
	if logger == nil {
		if debug {
			return l.NewSimpleDebugLogger(name)
		}
        return l.NewSimpleLogger(name)
	}
	return logger
}

// ========================================================================================================================
// EXEC HELPERS
// ========================================================================================================================
func ExecCommand(name string, arg ...string) (string, error) {
	if out, err := exec.Command(name, arg...).Output(); err != nil {
		return EMPTY, err
	} else {
		return string(out), nil
	}
}

// ========================================================================================================================
// DEBUG HELPERS
// ========================================================================================================================
func Dump(logger l.Logger, arr []string, label string) {
	for _,s := range arr {
		logger.Debug(label + " -->" + s + "<--")
	}

}

// ========================================================================================================================
// STRING TO NUMBER CONVERSIONS
// ========================================================================================================================
func ToInt64(i string, dflt int64) int64 {
	i64, err := strconv.ParseInt(i, 10, 0)
	if err != nil {
		return dflt
	}
	return i64
}

func ToFloat64(f string, dflt float64) float64 {
	f64, err := strconv.ParseFloat(f, 64)
	if err != nil {
		return dflt
	}
	return f64
}

// ========================================================================================================================
// OBJECT TO STRING CONVERSIONS
// ========================================================================================================================
func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

func DurationToString(dur time.Duration) string {
	return strconv.FormatInt(int64(dur), 10)
}

func ByteToString(b byte) string {
	return strconv.FormatInt(int64(b), 10)
}

func IntToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func Float64ToStringPrecise(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}

// ========================================================================================================================
// MATH
// ========================================================================================================================
func Round(f float64) int64 {
	t := math.Trunc(f)
	x := math.Trunc( (f - t) * 100 )
	if x < 50 {
		return int64(t)
	}
	return  int64(t) + 1
}

func TruncDp(value float64, decimalPoints int) float64 {
	var factor float64 = math.Pow10(decimalPoints)
	return float64(int(value * factor)) / factor
}

func RoundDp(value float64, decimalPoints int) float64 {
	var round float64 = math.Pow10(-decimalPoints) * .5
	return TruncDp(value + round, decimalPoints)
}

func Percent(top float64, bot float64) float64 {
    return top * 100 / bot
}

func PercentTruncDp(top float64, bot float64, decimalPoints int) float64 {
    return TruncDp(Percent(top, bot), decimalPoints)
}

func PercentRoundDp(top float64, bot float64, decimalPoints int) float64 {
    return RoundDp(Percent(top, bot), decimalPoints)
}

// ========================================================================================================================
// STRING READING
// ========================================================================================================================
func GetLastIndex(a []string) int {
	if a != nil && len(a) > 0 {
		return len(a) - 1
	}
	return -1
}

func GetLastField(a []string) string {
	if a != nil && len(a) > 0 {
		return a[len(a) - 1]
	}
	return EMPTY
}

func GetFieldByIndex(a []string, columnIndex int) string {
	if a != nil && len(a) > columnIndex {
		return a[columnIndex]
	}
	return EMPTY
}

func SplitGetFieldByIndex(text string, columnIndex int) string {
	return GetFieldByIndex(strings.Fields(text), columnIndex)
}

func SplitGetLastField(text string) string {
	return GetLastField(strings.Fields(text))
}

func StringArraysEqual(sa1 []string, sa2 []string) bool {
	saLen := len(sa1)
	if saLen != len(sa2) {
		return false
	}

	for i := 0; i < saLen; i++ {
		if sa1[i] != sa2[i] {
			return false
		}
	}

	return true
}

// ========================================================================================================================
// STRING MANIPULATION
// ========================================================================================================================

// ========================================================================================================================
// JSON STUFF
// ========================================================================================================================
func ToJson(v interface{}) string{
    j, err := json.Marshal(v)
    if err == nil {
        return string(j)
    }

    panic(err)
}

func ToJsonFormatted(v interface{}) string{
    j, err := json.MarshalIndent(v, "", "    ")
    if err == nil {
        return string(j)
    }

    panic(err)
}

func ToJsonLabeledString(v interface{}) string {
    return "{\"" + ObjectPrefix(v) + "\": " + ToJson(v) + "}"
}

func ToJsonLabeledFormattedString(v interface{}) string {
    return "{\"" + ObjectPrefix(v) + "\": " + ToJsonFormatted(v) + "}"
}

func UpFirst(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func ArrayToString(a []string) string {
	result := OPEN_BRACE
	for _, s := range a {
		if result == OPEN_BRACE {
			result = result + QUOTE + s + QUOTE
		} else {
			result = result + COMMA + SPACE + QUOTE + s + QUOTE
		}
	}
	return result + CLOSE_BRACE
}

func ObjectPrefix(v interface{}) string {
    typeOf := reflect.TypeOf(v)
    if typeOf == nil {
        return "<i>"
    }
    str := typeOf.String()
    at := strings.LastIndex(str, DOT)
    return str[at+1:]
}

func ToSizeMetricType(size string) MetricType {
	switch strings.ToUpper(size) {
	case "BYTE":  return MT_SIZE_BYTE
	case "BYTES": return MT_SIZE_BYTE
	case "KB":    return MT_SIZE_KB
	case "MB":    return MT_SIZE_MB
	case "GB":    return MT_SIZE_GB
	case "TB":    return MT_SIZE_TB
	}
	return MT_NONE
}