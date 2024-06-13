package logger

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"io"
	"os"
	"time"
)

func Stdout() io.Writer {
	return os.Stdout
}

func FileOutput(name string, count uint, aliveTime time.Duration) io.Writer {
	if count > 0 {
		writer, err := rotatelogs.New(name, rotatelogs.WithRotationCount(count), rotatelogs.WithRotationTime(aliveTime))
		if err != nil {
			fmt.Println(err.Error())
			return Stdout()
		}
		return writer
	}
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err.Error())
		return Stdout()
	}
	return file
}

func DefMulFileOut() io.Writer {
	return FileOutput("./bin/logs/%Y-%m-%d.log", 30, time.Minute*60*24)
}
