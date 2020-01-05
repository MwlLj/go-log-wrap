package logwrap

import (
    "testing"
    "fmt"
    "time"
)

func TestLog(t *testing.T) {
    l, err := NewLog("test", &Config{
        MaxSizeKb: 8,
        RemainDays: 3,
    })
    if err != nil {
        fmt.Println(err)
        return
    }
    for {
        l.Logger().Infof("hello, %s\n", "world")
        <-time.After(time.Millisecond * 100)
    }
}
