package logger

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

const lineFatalLog = 20

// printAllLevels is helper function for printing messages for all available log levels
func printAllLevels(msg string) {
	Debug("-->", msg)
	Info("-->", msg)
	Warn("-->", msg)
	Error("-->", msg)
	Fatal("-->", msg)
	DebugF("--> %s", msg)
	InfoF("--> %s", msg)
	WarnF("--> %s", msg)
	ErrorF("--> %s", msg)
	FatalF("--> %s", msg)
}

// Test of function Debug()
func TestDebug(t *testing.T) {
	out := strings.Builder{}
	SetOutput(&out)
	SetLevelStr("DEBUG")
	t.Run("test debug function", func(t *testing.T) {
		printAllLevels("test debug message")
		if strings.Count(out.String(), "DEBUG") != 2 {
			t.Errorf("Message with status DEBUG should printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "INFO") != 2 {
			t.Errorf("Message with status INFO should printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "WARN") != 2 {
			t.Errorf("Message with status WARN should printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "ERROR") != 2 {
			t.Errorf("Message with status ERROR should printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "FATAL") != 2 {
			t.Errorf("Message with status FATAL should printed:\n%v", out.String())
		}
	})
}

// Test of function Info()
func TestInfo(t *testing.T) {
	out := strings.Builder{}
	SetOutput(&out)
	SetLevelStr("INFO")
	t.Run("test info function", func(t *testing.T) {
		printAllLevels("test info message")
		if strings.Count(out.String(), "DEBUG") == 2 {
			t.Errorf("Message with status DEBUG should not printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "INFO") != 2 {
			t.Errorf("Message with status INFO should printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "WARN") != 2 {
			t.Errorf("Message with status WARN should printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "ERROR") != 2 {
			t.Errorf("Message with status ERROR should printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "FATAL") != 2 {
			t.Errorf("Message with status FATAL should printed:\n%v", out.String())
		}
	})
}

// Test of function Warn()
func TestWarn(t *testing.T) {
	out := strings.Builder{}
	SetOutput(&out)
	SetLevelStr("WARN")
	t.Run("test warn function", func(t *testing.T) {
		printAllLevels("test warm message")
		if strings.Count(out.String(), "DEBUG") == 2 {
			t.Errorf("Message with status DEBUG should not printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "INFO") == 2 {
			t.Errorf("Message with status INFO should not printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "WARN") != 2 {
			t.Errorf("Message with status WARN should printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "ERROR") != 2 {
			t.Errorf("Message with status ERROR should printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "FATAL") != 2 {
			t.Errorf("Message with status FATAL should printed:\n%v", out.String())
		}
	})
}

// Test of function Error()
func TestError(t *testing.T) {
	out := strings.Builder{}
	SetOutput(&out)
	SetLevelStr("ERROR")
	t.Run("test error function", func(t *testing.T) {
		printAllLevels("test error message")
		if strings.Count(out.String(), "DEBUG") == 1 {
			t.Errorf("Message with status DEBUG should not printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "INFO") == 2 {
			t.Errorf("Message with status INFO should not printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "WARN") == 2 {
			t.Errorf("Message with status WARN should not printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "ERROR") != 2 {
			t.Errorf("Message with status ERROR should printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "FATAL") != 2 {
			t.Errorf("Message with status FATAL should printed:\n%v", out.String())
		}
	})
}

// Test of function Fatal()
func TestFatal(t *testing.T) {
	out := strings.Builder{}
	SetOutput(&out)
	SetLevelStr("FATAL")
	t.Run("test fatal function", func(t *testing.T) {
		printAllLevels("test fatal message")
		if strings.Count(out.String(), "DEBUG") == 2 {
			t.Errorf("Message with status DEBUG should not printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "INFO") == 2 {
			t.Errorf("Message with status INFO should not printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "WARN") == 2 {
			t.Errorf("Message with status WARN should not printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "ERROR") == 2 {
			t.Errorf("Message with status ERROR should printed:\n%v", out.String())
		}
		if strings.Count(out.String(), "FATAL") != 2 {
			t.Errorf("Message with status FATAL should printed:\n%v", out.String())
		}
	})
}

// Test disable of logger
func TestOff(t *testing.T) {
	out := strings.Builder{}
	SetOutput(&out)
	SetLevelStr("OFF")
	t.Run("test disable of logger", func(t *testing.T) {
		printAllLevels("test off message")
		if out.Len() != 0 {
			t.Errorf("Log has been empty:\n%v", out.String())
		}
	})
}

// Test of invocation location
func TestInvocationLocation(t *testing.T) {
	out := strings.Builder{}
	SetOutput(&out)
	SetLevelStr("FATAL")
	t.Run("test of invocation location", func(t *testing.T) {
		printAllLevels("test invocation location")
		_, file, _, _ := runtime.Caller(0)
		location := fmt.Sprintf("%s:%d", filepath.Base(file), lineFatalLog)
		if !strings.Contains(out.String(), location) {
			t.Errorf("Log should contain locatiion of invocation '%s':\n%v", location, out.String())
		}
	})
}

// Test of setting field delimiter
func TestSetDelimiter(t *testing.T) {
	uniqueDelimiter := "<~>"
	out := strings.Builder{}
	SetOutput(&out)
	SetLevelStr("FATAL")
	SetDelimiter(uniqueDelimiter)
	t.Run("test of setting field delimiter", func(t *testing.T) {
		printAllLevels("test message")
		if !strings.Contains(out.String(), uniqueDelimiter) {
			t.Errorf("Log should contain delimiter '%s':\n%v", uniqueDelimiter, out.String())
		}
	})
	SetDelimiter(" ")
}

// Test setting of time format
func TestSetFormat(t *testing.T) {
	testFormat := "15:04"
	out := strings.Builder{}
	SetOutput(&out)
	SetLevelStr("FATAL")
	SetTimeFormat(testFormat)
	t.Run("test setting of time format", func(t *testing.T) {
		currentTime := time.Now().Format(testFormat)
		printAllLevels("test format")
		if strings.Count(out.String(), " "+currentTime+" ") != 2 {
			t.Errorf("Log should contain time in format '%s':\n%v", testFormat, out.String())
		}
	})
	SetTimeFormat(DefaultTimeFormat)
}

// Test setting of JSON log format
func TestJsonFormat(t *testing.T) {
	out := strings.Builder{}
	SetOutput(&out)
	SetLevelStr("DEBUG")
	SetLogFormat(JSON)
	t.Run("test setting of JSON log format", func(t *testing.T) {
		printAllLevels("test json format")
		if !(strings.Count(out.String(), "{\"level\":\"DEBUG\",\"time\":\"") == 2 &&
			strings.Count(out.String(), "\",\"location\":\"logger_test.go:") == 10 &&
			strings.Count(out.String(), "\"goroutine\":") == 10 &&
			strings.Count(out.String(), ",\"message\":\"--> test json format\"}") == 10) {
			t.Errorf("Log should be in JSON format:\n%v", out.String())
		}
	})
}

// Test setting of FLAT log format
func TestFlatFormat(t *testing.T) {
	out := strings.Builder{}
	SetOutput(&out)
	SetLevelStr("DEBUG")
	SetLogFormat(FLAT)
	t.Run("test setting of JSON log format", func(t *testing.T) {
		printAllLevels("test json format")
		if !(strings.Count(out.String(), "[DEBUG] ") == 2 &&
			strings.Count(out.String(), " logger_test.go:") == 10 &&
			strings.Count(out.String(), " goroutine:") == 10 &&
			strings.Count(out.String(), " --> test json format") == 10) {
			t.Errorf("Log should be in JSON format:\n%v", out.String())
		}
	})
}
