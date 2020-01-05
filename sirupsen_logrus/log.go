package logwrap

import (
    "github.com/sirupsen/logrus"
    "os"
)

type Config struct {
    MaxSizeKb int64
    RemainDays int
}

type CLog struct {
    log *logrus.Logger
    outputer *Output
}

func (self *CLog) init() {
    // self.log.SetFormatter(&logrus.JSONFormatter{
    //     PrettyPrint: false,
    // })
    self.log.SetFormatter(&logrus.TextFormatter{
        TimestampFormat: "2006-01-02 15:04:05",
    })
    self.log.SetOutput(self.outputer)
}

func isExist(name string) bool {
    _, err := os.Stat(name)
    if err != nil {
        if os.IsExist(err) {
            return true
        }
        return false
    }
    return true
}

func (self *CLog) Close() {
    self.outputer.Close()
}

func (self *CLog) Logger() *logrus.Logger {
    return self.log
}

func (self *CLog) SetRoot(root string) {
    self.outputer.SetRoot(root)
}

func NewLog(root string, config *Config) *CLog {
    l := CLog{
        log: logrus.New(),
        outputer: NewOutput(root, config),
    }
    l.init()
    return &l
}
