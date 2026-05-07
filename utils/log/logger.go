package log

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/YasinDoyle/e-mall/repository/es"

	"github.com/sirupsen/logrus"
)

var LogrusObj *logrus.Logger

func InitLogger() {
	if LogrusObj != nil {
		src, _ := setOutputFile()
		LogrusObj.Out = src
		return
	}

	// 实例化
	logger := logrus.New()
	src, _ := setOutputFile()
	// 设置输出
	logger.Out = src
	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)
	// 设置日志格式
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	hook := es.EsHookLog()
	logger.AddHook(hook)
	LogrusObj = logger
}

func setOutputFile() (*os.File, error) {
	now := time.Now()
	logFilePath := ""
	if dir, err := os.Getwd(); err == nil {
		logFilePath = dir + "/logs/"
	}
	_, err := os.Stat(logFilePath)

	if os.IsNotExist(err) {
		if err := os.MkdirAll(logFilePath, 0777); err != nil {
			log.Println(err.Error())
			return nil, err
		}
	}

	logFileName := now.Format("2006-01-02") + ".log"
	fileName := path.Join(logFilePath, logFileName)
	src, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Println("err", err.Error())
		return nil, err
	}
	return src, nil

}
