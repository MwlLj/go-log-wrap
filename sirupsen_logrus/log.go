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

func (self *CLog) init(root string) error {
    if !isExist(root) {
        err := os.MkdirAll(root, os.ModePerm)
        if err != nil {
            return err
        }
    }
    // self.log.SetFormatter(&logrus.JSONFormatter{
    //     PrettyPrint: false,
    // })
    self.log.SetFormatter(&logrus.TextFormatter{
        TimestampFormat: "2006-01-02 15:04:05",
    })
    self.log.SetOutput(self.outputer)
    return nil
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

func NewLog(root string, config *Config) (*CLog, error) {
    l := CLog{
        log: logrus.New(),
        outputer: NewOutput(root, config),
    }
    err := l.init(root)
    if err != nil {
        return nil, err
    }
    return &l, nil
}
