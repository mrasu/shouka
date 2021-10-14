package log

import (
	"fmt"
	stdLog "log"
	"os"
	"time"
)

type Level int

const (
	Minimal Level = iota + 1
	Debug
)

type logger struct {
	level    Level
	instance *stdLog.Logger
}

var l = logger{
	level:    Minimal,
	instance: stdLog.New(os.Stdout, "", 0),
}

func SetLevel(level Level) {
	l.level = level
}

func Printf(format string, v ...interface{}) {
	l.printf(format, v...)
}

func Println(v ...interface{}) {
	l.println(v...)
}

func Debugf(format string, v ...interface{}) {
	l.writefWithInfo(Debug, format, v...)
}

func Debugln(v ...interface{}) {
	l.writelnWithInfo(Debug, v...)
}

func Fatalf(format string, v ...interface{}) {
	l.writefWithInfo(Debug, format, v...)
	os.Exit(1)
}

func (l *logger) println(v ...interface{}) {
	msg := fmt.Sprintln(v...)
	l.write(msg)
}

func (l *logger) printf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.write(msg)
}

func (l *logger) writefWithInfo(level Level, format string, v ...interface{}) {
	if !l.shouldLog(level) {
		return
	}
	msg := fmt.Sprintf(format, v...)
	l.writeWithFormat(msg)
}

func (l *logger) writelnWithInfo(level Level, v ...interface{}) {
	if !l.shouldLog(level) {
		return
	}
	msg := fmt.Sprintln(v...)
	l.writeWithFormat(msg)
}

func (l *logger) shouldLog(level Level) bool {
	return l.level >= level
}

func (l *logger) writeWithFormat(msg string) {
	t := time.Now().Format("2006/01/02 15:04:05")
	l.write("[DEBUG] " + t + ": " + msg)
}

func (l *logger) write(msg string) {
	l.instance.Print(msg)
}
