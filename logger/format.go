package logger

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/youngpto/funs_tool/times"
	"runtime"
	"strings"
	"sync"
)

type Format interface {
	logrus.Formatter
}

func TextFormat() Format {
	return new(logrus.TextFormatter)
}

func JSONFormat() Format {
	return new(logrus.JSONFormatter)
}

var (
	loggerPkg       string
	miniCallerDepth = 1
	callerInitOnce  sync.Once
)

const (
	toolPackage            = "github.com/youngpto/funs_tool"
	logrusPackage          = "github.com/sirupsen/logrus"
	mqantPackage           = "github.com/youngpto/mqant/log"
	maximumCallerDepth int = 30
	knownFrames        int = 4
)

type funFormat struct {
	black []string
}

func FunFormat(black []string) *funFormat {
	formatter := &funFormat{
		black: []string{
			toolPackage,
			logrusPackage,
			mqantPackage,
		},
	}

	for _, s := range black {
		formatter.black = append(formatter.black, s)
	}

	return formatter
}

func (formatter *funFormat) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := times.Time2Loc(entry.Time).Format("2006-01-02 15:04:05.000")

	var newLog strings.Builder
	newLog.WriteString("[")

	if name, ok := entry.Data["name"]; ok {
		if str, convert := name.(string); convert {
			newLog.WriteString(fmt.Sprintf("%s ", str))
		}
	}

	newLog.WriteString(fmt.Sprintf("%s %s", timestamp, formatter.ParseLevel(entry.Level.String())))

	if caller := formatter.getCaller(); caller != nil {
		//newLog.WriteString(fmt.Sprintf(" %s:%d", filepath.Base(caller.File), caller.Line))
		newLog.WriteString(fmt.Sprintf(" %s:%d", caller.File, caller.Line))
	}
	newLog.WriteString("] ")

	newLog.WriteString(fmt.Sprintf("%s\n", entry.Message))

	b.WriteString(newLog.String())
	return b.Bytes(), nil
}

func (formatter *funFormat) getCaller() *runtime.Frame {
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(miniCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	callerInitOnce.Do(func() {
		loggerPkg = getPackageName(runtime.FuncForPC(pcs[0]).Name())

		miniCallerDepth = knownFrames
	})

NEXT:
	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		for _, blackPkg := range formatter.black {
			if strings.Contains(pkg, blackPkg) {
				continue NEXT
			}
		}
		return &f
	}

	return nil
}

func (formatter *funFormat) ParseLevel(level string) string {
	switch level {
	case "trace":
		return "TRACE"
	case "debug":
		return "DEBUG"
	case "info":
		return "INFO"
	case "warning":
		return "WARN"
	case "error":
		return "ERROR"
	case "fatal":
		return "FATAL"
	case "panic":
		return "PANIC"
	default:
		return ""
	}
}

func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}
