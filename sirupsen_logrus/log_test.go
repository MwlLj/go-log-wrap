package logwrap

import (
    "testing"
)

func TestLog(t *testing.T) {
    l := NewLog("test")
    l.Infof("hello, %s\n", "world")
}
