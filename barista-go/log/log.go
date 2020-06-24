package log

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"runtime/debug"

	"github.com/alecthomas/repr"
	"github.com/logrusorgru/aurora"
)

type ExitCode int

const (
	UnknownReason ExitCode = iota
	ConfigFailure
	BackendFailure
)

func loc() aurora.Value {
	_, file, line, _ := runtime.Caller(2)
	return aurora.Gray(12, fmt.Sprintf("[%s:%d]", path.Base(file), line))
}

// CanPanic calls a function that can panic
func CanPanic(f func()) {
	defer func() {
		if r := recover(); r != nil {
			CatchPanic(debug.Stack())
		}
	}()
	f()
}

// CatchPanic pretty prints a caught panic
func CatchPanic(stack []byte) {
	defer recover()
	fmt.Println(aurora.Sprintf(
		"%s %s %s\n%s",
		loc(),
		aurora.Bold(aurora.Blink(aurora.Red("PANIC"))),
		aurora.Bold(aurora.White("==>")),
		string(stack),
	))
}

// DebugValue pretty prints a value for debugging purposes
func DebugValue(message string, value interface{}) {
	defer func() {
		if r := recover(); r != nil {
			CatchPanic(debug.Stack())
		}
	}()
	fmt.Println(aurora.Sprintf(
		"%s %s %s %s:\n%s",
		loc(),
		aurora.Bold(aurora.Yellow("Debug Value")),
		aurora.Bold(aurora.White("==>")),
		message,
		repr.String(value, repr.Indent("\t")),
	))
}

// Info outputs an info message
func Info(format string, v ...interface{}) {
	var predefined []interface{}
	predefined = []interface{}{loc(), aurora.Bold(aurora.Cyan("Info")), aurora.Bold(aurora.White("==>"))}
	fmt.Println(aurora.Sprintf(
		"%s %s %s "+format,
		append(predefined, v...)...,
	))
}

// Error outputs an error message
func Error(format string, v ...interface{}) {
	var predefined []interface{}
	predefined = []interface{}{loc(), aurora.Bold(aurora.Red("Error")), aurora.Bold(aurora.White("==>"))}
	fmt.Println(aurora.Sprintf(
		"%s %s %s "+format,
		append(predefined, v...)...,
	))
}

// Fatal outputs an error message
func Fatal(code ExitCode, format string, v ...interface{}) {
	var predefined []interface{}
	predefined = []interface{}{loc(), aurora.Bold(aurora.Red("Fatal")), aurora.Bold(aurora.White("==>"))}
	fmt.Println(aurora.Sprintf(
		"%s %s %s "+format,
		append(predefined, v...)...,
	))
	os.Exit(int(code))
}
