// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package log implements a simple logging package. It defines a type, Logger,
// with methods for formatting output. It also has a predefined 'standard'
// Logger accessible through helper functions Print[f|ln], Fatal[f|ln], and
// Panic[f|ln], which are easier to use than creating a Logger manually.
// That logger writes to standard error and prints the date and time
// of each logged message.
// The Fatal functions call os.Exit(1) after writing the log message.
// The Panic functions call panic after writing the log message.
package golog

import (
	"io"
	"runtime"
	"strings"
	"sync"
)

type PartFunc func(*Logger)

// A Logger represents an active logging object that generates lines of
// output to an io.Writer.  Each logging operation makes a single call to
// the Writer's Write method.  A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	mu          sync.Mutex // ensures atomic writes; protects the following fields
	buf         []byte     // for accumulating text to write
	level       Level
	enableColor bool
	name        string
	pkgName     string
	userData    interface{}
	colorFile   *ColorFile

	parts []PartFunc

	output io.Writer

	currColor     Color
	currLevel     Level
	currText      string
	currCondition bool
	currContext   interface{}
}

// New creates a new Logger.   The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line.
// The flag argument defines the logging properties.

const lineBuffer = 32

func getPackageName() string {
	pc, _, _, _ := runtime.Caller(2)
	raw := runtime.FuncForPC(pc).Name()
	return strings.TrimSuffix(raw, ".init.ializers")
}

func New(name string) *Logger {

	l := &Logger{
		level:         Level_Debug,
		name:          name,
		pkgName:       getPackageName(),
		buf:           make([]byte, 0, lineBuffer),
		currCondition: true,
	}

	l.SetParts(LogPart_CurrLevel, LogPart_Name, LogPart_Time)

	add(l)

	return l
}

func (slf *Logger) EnableColor(v bool) {
	slf.mu.Lock()
	slf.enableColor = v
	slf.mu.Unlock()
}

func (slf *Logger) SetParts(f ...PartFunc) {
	slf.parts = []PartFunc{logPart_ColorBegin}
	slf.parts = append(slf.parts, f...)
	slf.parts = append(slf.parts, logPart_Text, logPart_ColorEnd, logPart_Line)
}

func (slf *Logger) SetFullParts(f ...PartFunc) {
	slf.parts = f
}

// 二次开发接口
func (slf *Logger) WriteRawString(s string) {
	slf.buf = append(slf.buf, s...)
}

func (slf *Logger) WriteRawByte(b byte) {
	slf.buf = append(slf.buf, b)
}

func (slf *Logger) WriteRawByteSlice(b []byte) {
	slf.buf = append(slf.buf, b...)
}

func (slf *Logger) Name() string {
	return slf.name
}

func (slf *Logger) SetUserData(data interface{}) {
	slf.userData = data
}

func (slf *Logger) UserData() interface{} {
	return slf.userData
}

func (slf *Logger) PkgName() string {
	return slf.pkgName
}

func (slf *Logger) Buff() []byte {
	return slf.buf
}

// 仅供LogPart访问
func (slf *Logger) Text() string {
	return slf.currText
}

// 仅供LogPart访问
func (slf *Logger) Context() interface{} {
	return slf.currContext
}

func (slf *Logger) LogText(level Level, text string, ctx interface{}) {

	// 防止日志并发打印导致的文本错位
	slf.mu.Lock()
	defer slf.mu.Unlock()

	slf.currLevel = level
	slf.currText = text
	slf.currContext = ctx

	defer slf.resetState()

	if slf.currLevel < slf.level || !slf.currCondition {
		return
	}

	slf.selectColorByText()
	slf.selectColorByLevel()

	slf.buf = slf.buf[:0]

	for _, p := range slf.parts {
		p(slf)
	}

	if slf.output != nil {
		slf.output.Write(slf.buf)
	} else {
		globalWrite(slf.buf)
	}

}

func (slf *Logger) Condition(value bool) *Logger {
	slf.mu.Lock()
	slf.currCondition = value
	slf.mu.Unlock()

	return slf
}

func (slf *Logger) resetState() {
	slf.currColor = NoColor
	slf.currCondition = true
	slf.currContext = nil
}

func (slf *Logger) IsDebugEnabled() bool {
	return slf.level == Level_Debug
}
