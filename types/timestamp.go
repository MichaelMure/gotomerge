package types

import "time"

// Timestamp represents a time in millisecond since Epoch.
type Timestamp int64

func Now() Timestamp {
	return Timestamp(time.Now().UnixMilli())
}

func FromTime(t time.Time) Timestamp {
	return Timestamp(t.UnixMilli())
}

func (t Timestamp) Time() time.Time {
	return time.Unix(int64(t/1000), int64(t%1000)*1_000_000)
}

func (t Timestamp) String() string {
	return t.Time().String()
}
