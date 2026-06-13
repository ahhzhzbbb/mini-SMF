package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger(lv zerolog.Level) *zerolog.Logger {
	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	consoleWriter := zerolog.ConsoleWriter{
		Out:          os.Stdout,
		NoColor:      false,
		TimeFormat:   time.RFC3339,
		TimeLocation: loc,
		FormatLevel: func(i interface{}) string {
			return fmt.Sprintf("[%s]", i)
		},
		FormatCaller: func(i interface{}) string {
			return fmt.Sprintf("[%s]", i)
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf(">> %s <<", i)
		},
		FormatFieldName: func(i interface{}) string {
			return fmt.Sprintf("\"%s\"", i)
		},
		// FormatFieldValue:      nil,
		FormatErrFieldName: func(i interface{}) string {
			return fmt.Sprintf("(%s)", i)
		},
		// FormatErrFieldValue:   nil,
		// FormatPartValueByName: nil,
		// FormatExtra: func(map[string]interface{}, *bytes.Buffer) error {
		// 	panic("TODO")
		// },
		// FormatPrepare: func(map[string]interface{}) error {
		// 	panic("TODO")
		// },
	}

	logger := zerolog.New(consoleWriter).Level(lv).With().Caller().Timestamp().Logger()
	return &logger
}
