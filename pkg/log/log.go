package log

import log "github.com/sirupsen/logrus"

func Info(msg ...interface{}) {
	log.Info(msg...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Error(msg ...interface{}) {
	log.Error(msg...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatal(msg ...interface{}) {
	log.Fatal(msg...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}

func Warn(msg ...interface{}) {
	log.Warn(msg...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Debug(msg ...interface{}) {
	log.Debug(msg...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}
