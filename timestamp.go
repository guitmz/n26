package n26

import (
	"strconv"
	"strings"
	"time"
)

type TimeStamp struct {
	time.Time
}

const millisAsNanos = int64(time.Millisecond)

var loc, _ = time.LoadLocation("Europe/Berlin")

func (ts *TimeStamp) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		return
	}
	value, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return
	}
	//Adjust the timezone information as timestamp is given in local TZ
	unadjusted := time.Unix(value/1000, (value%1000)*millisAsNanos).In(loc)
	_, offset := unadjusted.Zone()
	ts.Time = unadjusted.Add(time.Duration(-offset) * time.Second)
	return
}

// convert the timestamp to an integer - millis since epoch
func (ts *TimeStamp) AsMillis() int64 {
	return ts.UnixNano() / int64(time.Millisecond)
}
