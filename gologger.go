package gologger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/zan8in/gologger/formatter"
	"github.com/zan8in/gologger/levels"
	"github.com/zan8in/gologger/writer"
)

var (
	// labels2 = map[levels.Level]string{
	// 	levels.LevelFatal:   "FTL",
	// 	levels.LevelError:   "ERR",
	// 	levels.LevelInfo:    "INF",
	// 	levels.LevelWarning: "WRN",
	// 	levels.LevelDebug:   "DBG",
	// 	levels.LevelVerbose: "VER",
	// }

	// 默认 Unicode 符号
	unicodeSymbols = map[levels.Level]string{
		levels.LevelFatal:   "✖", // 红色大写X
		levels.LevelError:   "✖", // 红色小写x
		levels.LevelInfo:    "✓", // 蓝色小写i
		levels.LevelWarning: "⚠", // 黄色叹号
		levels.LevelDebug:   "#", // 灰色井号
		levels.LevelVerbose: "~", // 浅灰色波浪号
	}
	// ASCII 兼容符号
	asciiSymbols = map[levels.Level]string{
		levels.LevelFatal:   "X", // 红色大写X
		levels.LevelError:   "x", // 红色小写x
		levels.LevelInfo:    "√", // 蓝色小写i
		levels.LevelWarning: "!", // 黄色叹号
		levels.LevelDebug:   "!", // 灰色井号
		levels.LevelVerbose: "~", // 浅灰色波浪号
	}

	labels = map[levels.Level]string{}

	unicode bool

	// DefaultLogger is the default logging instance
	DefaultLogger *Logger
)

func init() {
	detectTerminalCapabilities()
	DefaultLogger = &Logger{}
	DefaultLogger.SetMaxLevel(levels.LevelInfo)
	DefaultLogger.SetFormatter(formatter.NewCLI(false))
	DefaultLogger.SetWriter(writer.NewCLI())
}

// Logger is a logger for logging structured data in a beautfiul and fast manner.
type Logger struct {
	writer    writer.Writer
	maxLevel  levels.Level
	formatter formatter.Formatter
}

// Log logs a message to a logger instance
func (l *Logger) Log(event *Event) {
	if event.level > l.maxLevel {
		return
	}
	event.message = strings.TrimSuffix(event.message, "\n")
	data, err := l.formatter.Format(&formatter.LogEvent{
		Message:  event.message,
		Level:    event.level,
		Metadata: event.metadata,
	})
	if err != nil {
		return
	}
	l.writer.Write(data, event.level)

	if event.level == levels.LevelFatal {
		os.Exit(1)
	}
}

// SetMaxLevel sets the max logging level for logger
func (l *Logger) SetMaxLevel(level levels.Level) {
	l.maxLevel = level
}

// SetFormatter sets the formatter instance for a logger
func (l *Logger) SetFormatter(formatter formatter.Formatter) {
	l.formatter = formatter
}

// SetWriter sets the writer instance for a logger
func (l *Logger) SetWriter(writer writer.Writer) {
	l.writer = writer
}

// Event is a log event to be written with data
type Event struct {
	logger   *Logger
	level    levels.Level
	message  string
	metadata map[string]string
}

// Label applies a custom label on the log event
func (e *Event) Label(label string) *Event {
	e.metadata["label"] = label
	return e
}

// Str adds a string metadata item to the log
func (e *Event) Str(key, value string) *Event {
	e.metadata[key] = value
	return e
}

// Msg logs a message to the logger
func (e *Event) Msg(format string) {
	e.message = format
	e.logger.Log(e)
}

// Msgf logs a printf style message to the logger
func (e *Event) Msgf(format string, args ...interface{}) {
	e.message = fmt.Sprintf(format, args...)
	e.logger.Log(e)
}

// Info writes a info message on the screen with the default label
func Info() *Event {
	level := levels.LevelInfo
	event := &Event{
		logger:   DefaultLogger,
		level:    level,
		metadata: make(map[string]string),
	}
	event.metadata["label"] = labels[level]
	return event
}

// Warning writes a warning message on the screen with the default label
func Warning() *Event {
	level := levels.LevelWarning
	event := &Event{
		logger:   DefaultLogger,
		level:    level,
		metadata: make(map[string]string),
	}
	event.metadata["label"] = labels[level]
	return event
}

// Error writes a error message on the screen with the default label
func Error() *Event {
	level := levels.LevelError
	event := &Event{
		logger:   DefaultLogger,
		level:    level,
		metadata: make(map[string]string),
	}
	event.metadata["label"] = labels[level]
	return event
}

// Debug writes an error message on the screen with the default label
func Debug() *Event {
	level := levels.LevelDebug
	event := &Event{
		logger:   DefaultLogger,
		level:    level,
		metadata: make(map[string]string),
	}
	event.metadata["label"] = labels[level]
	return event
}

// Fatal exits the program if we encounter a fatal error
func Fatal() *Event {
	level := levels.LevelFatal
	event := &Event{
		logger:   DefaultLogger,
		level:    level,
		metadata: make(map[string]string),
	}
	event.metadata["label"] = labels[level]
	return event
}

// Silent prints a string on stdout without any extra labels.
func Silent() *Event {
	level := levels.LevelSilent
	event := &Event{
		logger:   DefaultLogger,
		level:    level,
		metadata: make(map[string]string),
	}
	return event
}

// Print prints a string on stderr without any extra labels.
func Print() *Event {
	level := levels.LevelInfo
	event := &Event{
		logger:   DefaultLogger,
		level:    level,
		metadata: make(map[string]string),
	}
	return event
}

// Verbose prints a string only in verbose output mode.
func Verbose() *Event {
	level := levels.LevelVerbose
	event := &Event{
		logger:   DefaultLogger,
		level:    level,
		metadata: make(map[string]string),
	}
	event.metadata["label"] = labels[level]
	return event
}

// Info writes a info message on the screen with the default label
func (l *Logger) Info() *Event {
	level := levels.LevelInfo
	event := &Event{
		logger:   l,
		level:    level,
		metadata: make(map[string]string),
	}
	event.metadata["label"] = labels[level]
	return event
}

// Warning writes a warning message on the screen with the default label
func (l *Logger) Warning() *Event {
	level := levels.LevelWarning
	event := &Event{
		logger:   l,
		level:    level,
		metadata: make(map[string]string),
	}
	event.metadata["label"] = labels[level]
	return event
}

// Error writes a error message on the screen with the default label
func (l *Logger) Error() *Event {
	level := levels.LevelError
	event := &Event{
		logger:   l,
		level:    level,
		metadata: make(map[string]string),
	}
	event.metadata["label"] = labels[level]
	return event
}

// Debug writes an error message on the screen with the default label
func (l *Logger) Debug() *Event {
	level := levels.LevelDebug
	event := &Event{
		logger:   l,
		level:    level,
		metadata: make(map[string]string),
	}
	event.metadata["label"] = labels[level]
	return event
}

// Fatal exits the program if we encounter a fatal error
func (l *Logger) Fatal() *Event {
	level := levels.LevelFatal
	event := &Event{
		logger:   l,
		level:    level,
		metadata: make(map[string]string),
	}
	event.metadata["label"] = labels[level]
	return event
}

// Print prints a string on screen without any extra labels.
func (l *Logger) Print() *Event {
	level := levels.LevelSilent
	event := &Event{
		logger:   l,
		level:    level,
		metadata: make(map[string]string),
	}
	return event
}

// Verbose prints a string only in verbose output mode.
func (l *Logger) Verbose() *Event {
	level := levels.LevelVerbose
	event := &Event{
		logger:   l,
		level:    level,
		metadata: make(map[string]string),
	}
	event.metadata["label"] = labels[level]
	return event
}

// 检测终端能力
func detectTerminalCapabilities() {
	// 检查 Unicode 支持
	switch runtime.GOOS {
	case "windows":
		unicode = detectWindowsUnicodeSupport()
	default:
		unicode = detectUnixUnicodeSupport()
	}
	if unicode {
		labels = unicodeSymbols
	} else {
		labels = asciiSymbols
	}
}

func detectWindowsUnicodeSupport() bool {
	// 检测 Windows Terminal 或配置了 UTF-8 代码页
	return os.Getenv("WT_SESSION") != "" ||
		strings.Contains(os.Getenv("PROMPT"), "$E") || // ANSI 转义支持
		os.Getenv("PYCHARM_HOSTED") == "1" // IDE 终端
}

func detectUnixUnicodeSupport() bool {
	return strings.Contains(strings.ToLower(os.Getenv("LANG")), "utf-8")
}
