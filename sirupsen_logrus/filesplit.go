package logwrap

import (
    "io/ioutil"
    "time"
    "strings"
    "bytes"
    pa "path"
    "path/filepath"
    "errors"
    "fmt"
    "os"
    "strconv"
)

var _ = fmt.Println
var _ = errors.New

const (
    dateLayoutFormat string = "2006-01-02"
)

/*
** 实现 io.Write 接口
*/
type Output struct {
    root string
    config *Config
    curDir string
    curFile *os.File
    maxSize int64
}

func (self *Output) init() {
    self.maxSize = self.config.MaxSizeKb * 1024
}

func (self *Output) SetRoot(root string) {
    self.root = root
}

func (self *Output) Write(data []byte) (int, error) {
    /*
    ** 获取日期
    */
    t := self.nowDate()
    /*
    ** 拼接目录
    */
    var fullDir bytes.Buffer
    fullDir.WriteString(self.root)
    if !strings.HasSuffix(self.root, "/") {
        /*
        ** 结尾 不存在 /
        */
        fullDir.WriteRune('/')
    }
    fullDir.WriteString(t)
    full := fullDir.String()
    /*
    ** 打开编号最大的文件
    */
    err := self.openFile(&full)
    if err != nil {
        return 0, err
    }
    /*
    ** 写入文件
    */
    return self.curFile.Write(data)
}

func (self *Output) Close() {
    if self.curFile != nil {
        self.curFile.Close()
    }
}

func (self *Output) isExist(dir *string) bool {
    return isExist(*dir)
}

func (self *Output) checkAndDelDateDir() {
    files, err := ioutil.ReadDir(self.root)
    if err != nil {
        return
    }
    for _, file := range files {
        if file.IsDir() {
            old, err := time.Parse(dateLayoutFormat, file.Name())
            if err != nil {
                return
            }
            timeUnix := time.Now().Unix()
            sub := time.Unix(timeUnix, 0).Sub(old)
            days := int(sub.Hours()/24)
            if (days < (self.config.RemainDays+1)) && (days > 0) {
            } else if days == 0 {
            } else {
                os.RemoveAll(pa.Join(self.root, file.Name()))
            }
        } else {
        }
    }
}

func (self *Output) openFile(dir *string) error {
    if self.maxSize > 0 {
        if self.curFile != nil {
            s, _ := self.curFile.Stat()
            if s.Size() > self.maxSize {
                /*
                ** 同一天, 正在写入的文件超过限制
                */
                n := filepath.Base(self.curFile.Name())
                ext := filepath.Ext(n)
                na := strings.TrimSuffix(n, ext)
                num, err := strconv.ParseInt(na, 10, 64)
                if err != nil {
                    return err
                }
                path := pa.Join(*dir, strings.Join([]string{
                    strconv.FormatInt(int64(num+1), 10), ext,
                }, ""))
                err = self.curFile.Close()
                if err != nil {
                    return err
                }
                self.curFile, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
                if err != nil {
                    return err
                }
            }
        }
    }
    if self.curDir == *dir {
        return nil
    }
    /*
    ** 目录发生变更 (日期变更), 需要重新查找
    ** 并更新缓存的当前目录
    */
    /*
    ** 当日期变更的时候, 检测一下是否需要删除存在已久的目录
    */
    if self.config.RemainDays > 0 {
        go self.checkAndDelDateDir()
    }
    /*
    ** 当日期变更的时候, 检测目录是否存在, 如果不存在就创建
    */
    if !self.isExist(dir) {
        os.MkdirAll(*dir, os.ModePerm)
    }
    self.curDir = *dir
    files, err := ioutil.ReadDir(*dir)
    if err != nil {
        return err
    }
    var max int64 = 1
    maxExt := ".log"
    for _, file := range files {
        if file.IsDir() {
        } else {
            n := file.Name()
            ext := filepath.Ext(n)
            na := strings.TrimSuffix(n, ext)
            num, err := strconv.ParseInt(na, 10, 64)
            if err != nil {
                continue
            }
            if num > max {
                max = num
                maxExt = ext
            }
        }
    }
    path := pa.Join(*dir, strings.Join([]string{
        strconv.FormatInt(int64(max), 10), maxExt,
    }, ""))
    /*
    ** 如果当前文件不是空, 需要先关闭
    */
    if self.curFile != nil {
        if err = self.curFile.Close(); err != nil {
            return err
        }
    }
    /*
    ** 判断计算出的最大的文件的大小是否超过限制
    ** 如果超过限制, 需要加1
    */
    if self.maxSize > 0 {
        size := self.fileSize(path)
        if size > self.maxSize {
            path = pa.Join(*dir, strings.Join([]string{
                strconv.FormatInt(int64(max+1), 10), maxExt,
            }, ""))
        }
    }
    self.curFile, err = os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
    if err != nil {
        return err
    }
    return nil
}

func (self *Output) fileSize(name string) int64 {
    f, err := os.Stat(name)
    if err != nil {
        return 0
    }
    return f.Size()
}

func (self *Output) nowDate() string {
    timeUnix := time.Now().Unix()
    return time.Unix(timeUnix, 0).Format(dateLayoutFormat)
}

func NewOutput(root string, config *Config) *Output {
    o := Output{
        root: root,
        config: config,
    }
    o.init()
    return &o
}
