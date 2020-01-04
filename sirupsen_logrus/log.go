package logwrap

import (
    "github.com/sirupsen/logrus"
)

type CLog struct {
    log *logrus.Logger
    outputer *Output
}

func (self *CLog) init() {
    self.log.SetFormatter(&logrus.JSONFormatter{
        PrettyPrint: false,
    })
    self.log.SetOutput(self.outputer)
}

func (self *CLog) Close() {
    self.outputer.Close()
}

func (self *CLog) Infof(format string, args ...interface{}) {
    self.log.Infof(format, args...)
}

func NewLog(name string) *CLog {
    l := CLog{
        log: logrus.New(),
        outputer: NewOutput(name),
    }
    l.init()
    return &l
}
