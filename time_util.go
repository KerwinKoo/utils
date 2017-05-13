package utils

import (
	"fmt"
	"time"
)

// taken from time/format.go
var conversion = map[string]string{
	/*stdLongMonth      */ "B": "January",
	/*stdMonth          */ "b": "Jan",
	// stdNumMonth       */ "m": "1",
	/*stdZeroMonth      */ "m": "01",
	/*stdLongWeekDay    */ "A": "Monday",
	/*stdWeekDay        */ "a": "Mon",
	// stdDay            */ "d": "2",
	// stdUnderDay       */ "d": "_2",
	/*stdZeroDay        */ "d": "02",
	/*stdHour           */ "H": "15",
	// stdHour12         */ "I": "3",
	/*stdZeroHour12     */ "I": "03",
	// stdMinute         */ "M": "4",
	/*stdZeroMinute     */ "M": "04",
	// stdSecond         */ "S": "5",
	/*stdZeroSecond     */ "S": "05",
	/*stdLongYear       */ "Y": "2006",
	/*stdYear           */ "y": "06",
	/*stdPM             */ "p": "PM",
	// stdpm             */ "p": "pm",
	/*stdTZ             */ "Z": "MST",
	// stdISO8601TZ      */ "z": "Z0700",  // prints Z for UTC
	// stdISO8601ColonTZ */ "z": "Z07:00", // prints Z for UTC
	/*stdNumTZ          */ "z": "-0700", // always numeric
	// stdNumShortTZ     */ "b": "-07",    // always numeric
	// stdNumColonTZ     */ "b": "-07:00", // always numeric
	"%": "%",
}

// StrfTime This is an alternative to time.Format because no one knows
// what date 040305 is supposed to create when used as a 'layout' string
// this takes standard strftime format options. For a complete list
// of format options see http://strftime.org/
func StrfTime(format string, t time.Time) string {
	layout := ""
	length := len(format)
	for i := 0; i < length; i++ {
		if format[i] == '%' && i <= length-2 {
			if layoutCmd, ok := conversion[format[i+1:i+2]]; ok {
				layout = layout + layoutCmd
				i++
				continue
			}
		}
		layout = layout + format[i:i+1]
	}
	return t.Format(layout)
}

// GetTimeStempAdd change the timestamp to match server timeer
func GetTimeStempAdd(addm, adds int64) string {
	t := time.Now()
	// changedt := t.Unix() + addm*60 + 30
	changedt := t.Unix() + addm*60 + adds
	tm := time.Unix(changedt, 0)
	timeStamp := StrfTime("%Y%m%d%H%M%S", tm)
	return timeStamp
}

// GetTimeStempYmdHMS get timeStamp format: "%Y%m%d%H%M%S"
// return type as yyyyMMddHHmmss
func GetTimeStempYmdHMS() string {
	t := time.Now()
	timeStamp := StrfTime("%Y%m%d%H%M%S", t)
	return timeStamp
}

// GetTimeStempUnix get 13 bits (to 0.001s) length of unix timestamp
func GetTimeStempUnix() string {
	t := time.Now()
	ts := fmt.Sprintf("%d", t.UnixNano())[:13]

	return ts
}

// GetTimeStempUnixInt64 get 13 bits (to 0.001s) length of unix timestamp
// return int64
func GetTimeStempUnixInt64() int64 {
	t := time.Now()
	ts := t.UnixNano() / 1000000

	return ts
}

// GetTimeStempUnixInt get unix timestamp
// return int64
func GetTimeStempUnixInt() int {
	t := time.Now()
	ts := t.Unix()

	return int(ts)
}
