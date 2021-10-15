package event

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = &logrus.Logger{
		Out:      os.Stderr,
		Level:    logrus.DebugLevel,
		ExitFunc: os.Exit,
	}
}
