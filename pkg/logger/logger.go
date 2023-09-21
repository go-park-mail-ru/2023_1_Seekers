package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type logHook struct {
	Writers   []io.Writer
	LogLevels []logrus.Level
}

func (h *logHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		return err
	}
	for _, w := range h.Writers {
		if _, err = w.Write([]byte(line)); err != nil {
			return err
		}
	}
	return nil
}

func (h *logHook) Levels() []logrus.Level {
	return h.LogLevels
}

type Logger struct {
	*logrus.Entry
}

func (l *Logger) LoggerWithFields(fields map[string]any) *Logger {
	return &Logger{l.WithFields(fields)}
}

func (l *Logger) LoggerWithField(key string, val any) *Logger {
	return &Logger{l.WithField(key, val)}
}

func Init(level logrus.Level, toStdout bool, logsFileName, timeFormat, baseDir, logsDir string) *Logger {
	l := logrus.New()
	loc, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		logrus.Fatalf("Error: cant load time location")
	}

	now := time.Now().In(loc)
	logsFName := logsDir + logsFileName + now.Format(timeFormat) + ".log"

	l.SetReportCaller(true)
	l.SetFormatter(&Formatter{
		TimeFormat:  timeFormat,
		Colors:      false,
		Caller:      true,
		FieldsSpace: true,
		LevelFirst:  false,
		CustomCallerFormatter: func(f *runtime.Frame) string {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			var dir string
			if idx := strings.Index(f.File, baseDir); idx != -1 {
				dir = f.File[idx+len(baseDir):]
			} else {
				dir = f.File
			}
			return fmt.Sprintf("[%s:%d @%s]", dir, f.Line, funcName)
		},
	})

	dirName := filepath.Dir(logsFName)
	if err = os.MkdirAll(dirName, 0777); err != nil {
		if os.IsExist(err) {
			logrus.Fatalf("Error: Such path already exists")
		}
		logrus.Fatalf("Error: create logs directory %v", err)
	}

	f, err := os.OpenFile(logsFName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logrus.Fatalf("Error: open file %v", err)
	}

	l.SetOutput(io.Discard)
	if toStdout {
		l.AddHook(&logHook{
			Writers:   []io.Writer{f, os.Stdout},
			LogLevels: logrus.AllLevels,
		})
	} else {
		l.AddHook(&logHook{
			Writers:   []io.Writer{f},
			LogLevels: logrus.AllLevels,
		})
	}

	l.SetLevel(level)

	return &Logger{logrus.NewEntry(l)}
}

type color uint8

const (
	Red    color = 31
	Yellow color = 33
	Cyan   color = 36
	Gray   color = 37
)

func (f *Formatter) printColored(b *bytes.Buffer, level logrus.Level, message string) {
	var levelColor color
	switch level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = Gray
	case logrus.WarnLevel:
		levelColor = Yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = Red
	default:
		levelColor = Cyan
	}
	_, err := fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m", levelColor, message)
	if err != nil {
		return
	}
}

type Formatter struct {
	TimeFormat            string
	Colors                bool
	Caller                bool
	LevelFirst            bool
	FieldsSpace           bool
	CustomCallerFormatter func(*runtime.Frame) string
}

func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	if f.TimeFormat == "" {
		f.TimeFormat = time.RFC3339
	}

	b := &bytes.Buffer{}

	f.writeStringBrackets(b, entry.Time.Format(f.TimeFormat))
	if f.FieldsSpace {
		b.WriteString(" ")
	}

	if !f.LevelFirst {
		if f.Caller {
			f.writeCaller(b, entry)
		}

		if f.FieldsSpace {
			b.WriteString(" ")
		}
	}

	level := strings.ToUpper(entry.Level.String())
	if f.Colors {
		f.printColored(b, entry.Level, f.getStringInBr(level))
	} else {
		f.writeStringBrackets(b, level)
	}

	if f.FieldsSpace {
		b.WriteString(" ")
	}

	f.writeFields(b, entry)

	if f.FieldsSpace {
		b.WriteString(" ")
	}

	if f.LevelFirst {
		if f.Caller {
			f.writeCaller(b, entry)
		}

		if f.FieldsSpace {
			b.WriteString(" ")
		}
	}

	b.WriteString(strings.TrimSpace(entry.Message))
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *Formatter) getStringInBr(message string) string {
	return fmt.Sprintf("[%s]", message)
}

func (f *Formatter) writeStringBrackets(b *bytes.Buffer, message string) {
	b.WriteString(f.getStringInBr(message))
}

func (f *Formatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		if f.CustomCallerFormatter != nil {
			_, _ = fmt.Fprint(b, f.CustomCallerFormatter(entry.Caller))
		} else {
			_, _ = fmt.Fprintf(
				b,
				"[%s:%d @%s]",
				entry.Caller.File,
				entry.Caller.Line,
				entry.Caller.Function,
			)
		}
	}
}

func (f *Formatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}

		sort.Strings(fields)

		for _, field := range fields {
			f.writeField(b, entry, field)
		}
	}
}

func (f *Formatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	str := fmt.Sprintf("[%s:%v]", field, entry.Data[field])

	if f.FieldsSpace {
		str += " "
	}
	if f.Colors {
		f.printColored(b, entry.Level, str)
	} else {
		b.WriteString(str)
	}
}
