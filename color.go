package golog

import (
	"io/ioutil"
	"strings"
)

type Color int

const (
	NoColor Color = iota
	Black
	Red
	Green
	Yellow
	Blue
	Purple
	DarkGreen
	White
)

var logColorPrefix = []string{
	"",
	"\x1b[030m",
	"\x1b[031m",
	"\x1b[032m",
	"\x1b[033m",
	"\x1b[034m",
	"\x1b[035m",
	"\x1b[036m",
	"\x1b[037m",
}

type colorData struct {
	name string
	c    Color
}

var colorByName = []colorData{
	{"none", NoColor},
	{"black", Black},
	{"red", Red},
	{"green", Green},
	{"yellow", Yellow},
	{"blue", Blue},
	{"purple", Purple},
	{"darkgreen", DarkGreen},
	{"white", White},
}

func matchColor(name string) Color {
	lower := strings.ToLower(name)
	for _, d := range colorByName {
		if d.name == lower {
			return d.c
		}
	}
	return NoColor
}

func colorFromLevel(l Level) Color {
	switch l {
	case Level_Warn:
		return Yellow
	case Level_Error:
		return Red
	}
	return NoColor
}

var logColorSuffix = "\x1b[0m"

func SetColorDefine(loggerName string, jsonFormat string) error {
	cf := NewColorFile()

	if err := cf.Load(jsonFormat); err != nil {
		return err
	}

	return VisitLogger(loggerName, func(l *Logger) bool {
		l.SetColorFile(cf)
		return true
	})
}

func EnableColorLogger(loggerName string, enable bool) error {
	return VisitLogger(loggerName, func(l *Logger) bool {
		l.EnableColor(enable)
		return true
	})
}

func SetColorFile(loggerName string, colorFileName string) error {
	data, err := ioutil.ReadFile(colorFileName)
	if err != nil {
		return err
	}
	return SetColorDefine(loggerName, string(data))
}

//设置当条颜色
func (slf *Logger) SetColor(name string) *Logger {
	slf.mu.Lock()
	slf.currColor = matchColor(name)
	slf.mu.Unlock()
	return slf
}

func (slf *Logger) SetVColor(v Color) *Logger {
	slf.mu.Lock()
	slf.currColor = v
	slf.mu.Unlock()
	return slf
}

func (slf *Logger) ColBlack() *Logger {
	return slf.SetVColor(Black)
}
func (slf *Logger) ColRed() *Logger {
	return slf.SetVColor(Red)
}
func (slf *Logger) ColGreen() *Logger {
	return slf.SetVColor(Green)
}
func (slf *Logger) ColYellow() *Logger {
	return slf.SetVColor(Yellow)
}
func (slf *Logger) ColBlue() *Logger {
	return slf.SetVColor(Blue)
}
func (slf *Logger) ColPurple() *Logger {
	return slf.SetVColor(Purple)
}
func (slf *Logger) ColDarkGreen() *Logger {
	return slf.SetVColor(DarkGreen)
}
func (slf *Logger) ColWhite() *Logger {
	return slf.SetVColor(White)
}

// 注意, 加色只能在Gogland的main方式启用, Test方式无法加色
func (slf *Logger) SetColorFile(file *ColorFile) {
	slf.colorFile = file
}

func (slf *Logger) selectColorByLevel() {
	if levelColor := colorFromLevel(slf.currLevel); levelColor != NoColor {
		slf.currColor = levelColor
	}
}

func (slf *Logger) selectColorByText() {
	if slf.enableColor && globalColorable && slf.colorFile != nil && slf.currColor == NoColor {
		slf.currColor = slf.colorFile.ColorFromText(slf.currText)
	}
}
