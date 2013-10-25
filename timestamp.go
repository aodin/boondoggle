package boondoggle

import (
	"fmt"
	"time"
)

type Timestamp struct {
	time.Time
}

func (t Timestamp) String() string {
	// TODO Ugh, how to check for a nil time?
	n := Timestamp{}
	if t == n {
		return ""
	}
	return fmt.Sprintf("%s, %s %s", t.Weekday(), t.Month(), t.Format("2, 2006"))
}

func CreateTimestamp(input string) (Timestamp, error) {
	t, err := time.Parse("2006-01-02", input)
	if err != nil {
		return Timestamp{}, err
	}
	return Timestamp{t}, nil
}

func MustCreate(input string) Timestamp {
	t, err := time.Parse("2006-01-02", input)
	if err != nil {
		panic(err)
	}
	return Timestamp{t}
}
