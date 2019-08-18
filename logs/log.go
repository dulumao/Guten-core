package logs

import (
	"fmt"
	"github.com/dulumao/Guten-core/env"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type logger log.Logger

const saveLogTime int64 = 60 * 60 * 24 * 7

var mu sync.RWMutex
var fd *os.File
var logHandle *log.Logger

func New() error {
	mu.Lock()
	defer mu.Unlock()

	if err := Mkdir(env.Value.Server.LogDir); err != nil {
		fmt.Printf("mkdir error %s\n", err)
		return err
	}

	if fd != nil {
		fd.Close()
	}

	fileName := fmt.Sprintf("%s_%s.log", env.Value.Server.LogDir+string(os.PathSeparator)+filepath.Base(os.Args[0]), time.Now().Format("20060102"))
	fd, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("create file error %s\n", err)
		return err
	}

	if env.Value.Server.Debug {
		logHandle = log.New(os.Stdout, "", log.LstdFlags)
	} else {
		// logHandle = log.New(fd, "", log.Ltime|log.Lshortfile)
		logHandle = log.New(fd, "", log.LstdFlags)
	}

	return nil
}

func LogInfo() *log.Logger {
	logHandle.SetPrefix("[INFO] ")
	return logHandle
}

func LogError() *log.Logger {
	logHandle.SetPrefix("[ERROR] ")
	return logHandle
}

func LogWarning() *log.Logger {
	logHandle.SetPrefix("[WARNING] ")
	return logHandle
}

func LogDebug() *log.Logger {
	logHandle.SetPrefix("[DEBUG] ")
	return logHandle
}

func (l *logger) Println(v ...interface{}) {
	mu.RLock()
	defer mu.RUnlock()

	logHandle.Println(v...)
}

func (l *logger) Printf(format string, v ...interface{}) {
	mu.RLock()
	defer mu.RUnlock()

	logHandle.Printf(format, v...)
}

func Errorf(format string, v ...interface{}) error {
	return fmt.Errorf(format, v...)
}

func Println(v ...interface{}) {
	fmt.Println(v...)
}

func Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

// 生成文件夹
func Mkdir(path string) error {
	if IsNotExist(path) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

// 判断文件是否存在
func IsNotExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return true
	}

	return false
}
