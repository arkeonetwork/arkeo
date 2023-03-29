// Structured logging with field annotation support.
// Details about where the log was called from are included.

package logging

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Fields is an alias to logrus.Fields
type Fields = logrus.Fields

type CustomLogFormatter struct {
	TimestampFormat string
	LevelDesc       []string
}

// Logger wraps a logrus.FieldLogger to provide all standard logging functionality
type Logger interface {
	logrus.FieldLogger
}

const timestampFormat = "2006-01-02T15:04:05.999Z07:00"

var logger Logger

func init() {
	logger = logrus.StandardLogger()

	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetReportCaller(true)
	// logrus.SetFormatter(&logrus.TextFormatter{TimestampFormat: timestampFormat})
	logrus.SetFormatter(NewCustomLogFormatter())
	if os.Getenv("ARKEO_DIR_JSON_LOGS") == "true" {
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: timestampFormat,
			PrettyPrint:     false,
		})
	}
}

// WithFields adds field annotations to the logger instance
func WithFields(fields Fields) Logger {
	return logger.WithFields(fields)
}

// WithoutFields uses the default logger with no extra field annotations
func WithoutFields() Logger {
	return logger
}

func NewCustomLogFormatter() *CustomLogFormatter {
	return &CustomLogFormatter{
		TimestampFormat: "2006-01-02T15:04:05.999Z07:00",
		LevelDesc:       []string{"PANC", "FATL", "ERRO", "WARN", "INFO", "DEBG"},
	}
}

func (f *CustomLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	fArr := strings.Split(entry.Caller.File, "/")
	if len(fArr) > 2 {
		fArr = fArr[len(fArr)-2:]
	}
	fileName := strings.Join(fArr, "/")

	timestamp := fmt.Sprint(entry.Time.Format(f.TimestampFormat))

	fields := ""
	if len(entry.Data) > 0 {
		sb := strings.Builder{}
		fmt.Fprintf(&sb, " {")
		for k, v := range entry.Data {
			fmt.Fprintf(&sb, "(%s: %s)", k, v)
		}
		fmt.Fprintf(&sb, "}")
		fields = sb.String()
	}

	return []byte(
		fmt.Sprintf("%s %s \"%s%s\" [%s:%d]\n", timestamp, f.LevelDesc[entry.Level], entry.Message, fields, fileName, entry.Caller.Line),
	), nil
}
