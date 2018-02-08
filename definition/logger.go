package definition

import (
	stdlog "log"
	"os"
)

type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
}

var log Logger

func init() {
	log = stdlog.New(os.Stderr, "", stdlog.LstdFlags)
}

func SetPackageLogger(l Logger) {
	log = l
}
