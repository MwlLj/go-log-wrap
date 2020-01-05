package logwrap

import (
    "testing"
    "time"
)

func TestLog(t *testing.T) {
    l := NewLog("test", &Config{
        MaxSizeKb: 8,
        RemainDays: 3,
    })
    for {
        l.Logger().Infof("hello, %s\n", "world")
        <-time.After(time.Millisecond * 100)
    }
}
