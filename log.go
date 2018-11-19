package expvar

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/zxfonline/fileutil"
	"github.com/zxfonline/timefix"
)

var (
	expvarLog      *log.Logger
	logFile        *os.File
	fileNameFormat string
	TimePeriod     = time.Hour
)

//InitTraceLog 初始化跟踪日志
func InitExpvarLog(filePath, fileNamePrefix string) {
	appName := fileutil.ExeName
	fileNameFormat = fileutil.TransPath(filepath.Join(filePath, fileNamePrefix+"_"+appName+"_"+"%v"+".log"))

	fileName := fmt.Sprintf(fileNameFormat, time.Now().Format("2006-01-02"))
	var err error
	logFile, err = fileutil.OpenFile(fileName, fileutil.DefaultFileFlag, fileutil.DefaultFileMode)
	if err != nil {
		log.Printf("open file err:%v\n", err)
		return
	}

	expvarLog = log.New(logFile, "", 0)
	go writeloop()
}

func writeloop() {
	pm := time.NewTimer(TimePeriod)
	baset := time.Now()
	pm1 := time.NewTimer(time.Duration(timefix.NextMidnight(baset, 1).Unix()-baset.Unix()) * time.Second)
	for {
		select {
		case <-pm.C:
			pm.Reset(TimePeriod)
			SaveExpvarLog()
		case <-pm1.C:
			now := time.Now()
			fileName := fmt.Sprintf(fileNameFormat, time.Now().Format("2006-01-02"))
			if logFile1, err := fileutil.OpenFile(fileName, fileutil.DefaultFileFlag, fileutil.DefaultFileMode); err != nil {
				log.Printf("[ERROR] "+"open file err:%v\n", err)
			} else {
				logFile.Close()
				expvarLog.SetOutput(logFile1)
				logFile = logFile1
				pm1.Reset(time.Duration(timefix.NextMidnight(now, 1).Unix()-now.Unix()) * time.Second)
			}
		}
	}
}

func SaveExpvarLog() {
	defer func() {
		if x := recover(); x != nil {
			fmt.Printf("Recovered %v\n.", x)
		}
	}()
	expvarLog.Println(GetExpvarString())
}
