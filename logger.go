package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	DefaultTimeFormat = "15:04:05.000"

	FLAT = "%[1]s[%[2]s]%[3]s%[4]s%[3]s%[5]s%[3]sgoroutine:%[6]d%[3]s%[7]s\u001B[0m\n"
	JSON = "%[1]s{\"level\":\"%[2]s\",\"time\":\"%[4]s\",\"location\":\"%[5]s\",\"goroutine\":%[6]d,\"message\":\"%[7]s\"}\u001B[0m\n"

	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
	OFF // disable logging
)

// Level describes available levels of logging
type Level uint8

// Label describes options of label (name and color)
type Label struct {
	name  string
	color string
}

// Config describes structure of logger configuration
type Config struct {
	level      Level
	timeFormat string
	logFormat  string
	output     io.Writer
	delimiter  string
}

// Message describes structure of logger message
type Message struct {
	label       Label
	time        time.Time
	location    string
	goroutineId uint64
	text        string
}

var (

	// using for synchronization of log message printing
	mutex = &sync.Mutex{}

	config = Config{
		level:      INFO,
		timeFormat: DefaultTimeFormat,
		logFormat:  FLAT,
		output:     os.Stdout,
		delimiter:  " ",
	}

	// mapping of log levels to printing names
	levelStore = map[Level]Label{
		DEBUG: {name: "DEBUG", color: "\u001B[37m"},
		INFO:  {name: "INFO ", color: "\u001B[32m"},
		WARN:  {name: "WARN ", color: "\u001B[33m"},
		ERROR: {name: "ERROR", color: "\u001B[31m"},
		FATAL: {name: "FATAL", color: "\u001B[35m"},
	}
)

// newMessage creates log message with current time
func newMessage(level Level, v ...interface{}) Message {

	fileName, line := getFileNameAndLine()

	var textBuilder strings.Builder
	for idx, msg := range v {
		if idx != 0 {
			textBuilder.WriteString(" ")
		}
		textBuilder.WriteString(fmt.Sprintf("%v", msg))
	}

	return Message{
		label:       levelStore[level],
		time:        time.Now(),
		location:    fmt.Sprintf("%s:%d", fileName, line),
		goroutineId: getGoroutineId(),
		text:        textBuilder.String(),
	}
}

// String converts log message to string in accordance with log format
func (message Message) String() string {
	return fmt.Sprintf(
		config.logFormat,
		message.label.color,
		message.label.name,
		config.delimiter,
		message.time.Format(config.timeFormat),
		message.location,
		message.goroutineId,
		message.text,
	)
}

// SetLevel sets level of logger.
// Supported levels are DEBUG, INFO, WARN, ERROR, FATAL
// By default using INFO.
func SetLevel(level Level) {
	config.level = level
}

// SetLevelStr sets level of logger from string value.
// Supported levels are DEBUG, INFO, WARN, ERROR, FATAL.
// If you set not supported value argument then will used level INFO.
func SetLevelStr(level string) {
	switch strings.ToUpper(level) {
	case "DEBUG":
		SetLevel(DEBUG)
	case "ERROR":
		SetLevel(ERROR)
	case "WARN":
		SetLevel(WARN)
	case "FATAL":
		SetLevel(FATAL)
	case "OFF":
		SetLevel(OFF)
	default:
		SetLevel(INFO)
	}
}

// SetTimeFormat sets format of time in log message.
// Use constants from package 'time'
// By default using constant 'DefaultTimeFormat' from this package.
func SetTimeFormat(format string) {
	config.timeFormat = format
}

// SetLogFormat sets format of logger.
// By default using constant 'SimpleFlatFormat' from this package.
func SetLogFormat(format string) {
	config.logFormat = format
}

// SetDelimiter sets delimiter between fields of log message.
// By default using " ".
func SetDelimiter(delimiter string) {
	config.delimiter = delimiter
}

// SetOutput sets writer for log messages.
// By default using console (os.Stdout).
func SetOutput(output io.Writer) {
	config.output = output
}

// getFileNameAndLine returns short file name and line where log function has been called
func getFileNameAndLine() (fileName string, line int) {
	_, fileName, line, _ = runtime.Caller(3)
	return filepath.Base(fileName), line
}

// getGoroutineId returns identifier of goroutine where log function has been called
func getGoroutineId() uint64 {
	buf := make([]byte, 64)
	buf = buf[:runtime.Stack(buf, false)]
	buf = bytes.TrimPrefix(buf, []byte("goroutine "))
	buf = buf[:bytes.IndexByte(buf, ' ')]
	id, _ := strconv.ParseUint(string(buf), 10, 64)
	return id
}

// writeLog writes log message with using selected format and level to output.
// See: SetOutput
func writeLog(message Message) {

	// synchronised printing of message
	mutex.Lock()
	_, _ = fmt.Fprint(config.output, message)
	mutex.Unlock()

}

// Fatal prints log message with level FATAL if the level is allow to printing
func Fatal(v ...interface{}) {
	if config.level <= FATAL {
		writeLog(newMessage(FATAL, v...))
	}
}

// Error prints log message with level ERROR if the level is allow to printing
func Error(v ...interface{}) {
	if config.level <= ERROR {
		writeLog(newMessage(ERROR, v...))
	}
}

// Warn prints log message with level WARN if the level is allow to printing
func Warn(v ...interface{}) {
	if config.level <= WARN {
		writeLog(newMessage(WARN, v...))
	}
}

// Info prints log message with level INFO if the level is allow to printing
func Info(v ...interface{}) {
	if config.level <= INFO {
		writeLog(newMessage(INFO, v...))
	}
}

// Debug prints log message with level DEBUG if the level is allow to printing
func Debug(v ...interface{}) {
	if config.level <= DEBUG {
		writeLog(newMessage(DEBUG, v...))
	}
}

// FatalF prints log message by format with level FATAL if the level is allow to printing
func FatalF(format string, v ...interface{}) {
	if config.level <= FATAL {
		writeLog(newMessage(FATAL, fmt.Sprintf(format, v...)))
	}
}

// ErrorF prints log message by format with level ERROR if the level is allow to printing
func ErrorF(format string, v ...interface{}) {
	if config.level <= ERROR {
		writeLog(newMessage(ERROR, fmt.Sprintf(format, v...)))
	}
}

// WarnF prints log message by format with level WARN if the level is allow to printing
func WarnF(format string, v ...interface{}) {
	if config.level <= WARN {
		writeLog(newMessage(WARN, fmt.Sprintf(format, v...)))
	}
}

// InfoF prints log message by format with level INFO if the level is allow to printing
func InfoF(format string, v ...interface{}) {
	if config.level <= INFO {
		writeLog(newMessage(INFO, fmt.Sprintf(format, v...)))
	}
}

// DebugF prints log message with level DEBUG if the level is allow to printing
func DebugF(format string, v ...interface{}) {
	if config.level <= DEBUG {
		writeLog(newMessage(DEBUG, fmt.Sprintf(format, v...)))
	}
}
