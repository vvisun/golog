package golog

import (
	"fmt"
	"os"
)

type Level int

const (
	Level_Debug Level = iota
	Level_Info
	Level_Warn
	Level_Error
	Level_Fatal
)

var levelString = [...]string{
	"[DEBU]",
	"[INFO]",
	"[WARN]",
	"[ERRO]",
	"[FATA]",
}

var globalLevel = Level_Debug

func SetGlobalLevel(lvl Level) {
	globalLevel = lvl
}

func str2loglevel(level string) Level {

	switch level {
	case "debug":
		return Level_Debug
	case "info":
		return Level_Info
	case "warn":
		return Level_Warn
	case "error", "err":
		return Level_Error
	case "fatal":
		return Level_Fatal
	}

	return Level_Debug
}

// 通过字符串设置某一类日志的级别
func SetLevelByString(loggerName string, level string) error {
	return VisitLogger(loggerName, func(l *Logger) bool {
		l.SetLevelByString(level)
		return true
	})
}

func (slf *Logger) Debugf(format string, v ...interface{}) {
	if slf.level < globalLevel {
		return
	}
	slf.LogText(Level_Debug, fmt.Sprintf(format, v...), nil)
}

func (slf *Logger) Infof(format string, v ...interface{}) {
	if slf.level < globalLevel {
		return
	}
	slf.LogText(Level_Info, fmt.Sprintf(format, v...), nil)
}

func (slf *Logger) Warnf(format string, v ...interface{}) {
	if slf.level < globalLevel {
		return
	}
	slf.LogText(Level_Warn, fmt.Sprintf(format, v...), nil)
}

func (slf *Logger) Errorf(format string, v ...interface{}) {
	slf.LogText(Level_Error, fmt.Sprintf(format, v...), nil)
}

func (slf *Logger) Fatalf(format string, v ...interface{}) {
	slf.LogText(Level_Fatal, fmt.Sprintf(format, v...), nil)
	os.Exit(1)
}

func (slf *Logger) Debugln(v ...interface{}) {
	if slf.level < globalLevel {
		return
	}
	slf.LogText(Level_Debug, fmt.Sprintln(v...), nil)
}

func (slf *Logger) Infoln(v ...interface{}) {
	if slf.level < globalLevel {
		return
	}
	slf.LogText(Level_Info, fmt.Sprintln(v...), nil)
}

func (slf *Logger) Warnln(v ...interface{}) {
	if slf.level < globalLevel {
		return
	}
	slf.LogText(Level_Warn, fmt.Sprintln(v...), nil)
}

func (slf *Logger) Errorln(v ...interface{}) {
	slf.LogText(Level_Error, fmt.Sprintln(v...), nil)
}

func (slf *Logger) Fatalln(v ...interface{}) {
	slf.LogText(Level_Fatal, fmt.Sprintln(v...), nil)
	os.Exit(1)
}

func (slf *Logger) SetLevelByString(level string) {
	slf.SetLevel(str2loglevel(level))
}

func (slf *Logger) SetLevel(lv Level) *Logger {
	slf.level = lv
	return slf
}

func (slf *Logger) Level() Level {
	return slf.level
}

func (slf *Logger) CurrLevelString() string {
	return levelString[slf.currLevel]
}
