package logger

import (
	"fmt"
	"fs_monitor/goesl"
)

type LoggerFiller interface {
	Type() string
	Detail() string
}

func logf(l LoggerFiller, dolog func(fmt string, v ...any), format string, v ...any) {
	if l != nil {
		format = fmt.Sprintf(`{detail:%s,type:%s} %s`, l.Detail(), l.Type(), format)
	}
	dolog(format, v...)
}

func Debugf(l LoggerFiller, format string, v ...any) {
	logf(l, goesl.Debugf, format, v...)
}

func Infof(l LoggerFiller, format string, v ...any) {
	logf(l, goesl.Infof, format, v...)
}

func Noticef(l LoggerFiller, format string, v ...any) {
	logf(l, goesl.Noticef, format, v...)
}

func Warnf(l LoggerFiller, format string, v ...any) {
	logf(l, goesl.Warnf, format, v...)
}

func Errorf(l LoggerFiller, format string, v ...any) {
	logf(l, goesl.Errorf, format, v...)
}

func Fatalf(l LoggerFiller, format string, v ...any) {
	logf(l, goesl.Fatalf, format, v...)
}
