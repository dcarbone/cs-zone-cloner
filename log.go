package main

import (
	"github.com/dcarbone/cs-zone-cloner/definition"
	stdlog "log"
	"os"
)

type logger struct {
	log *stdlog.Logger
}

var log *logger

func init() {
	log = &logger{
		log: stdlog.New(os.Stderr, "", stdlog.LstdFlags),
	}
	definition.SetPackageLogger(log)
}

func (l *logger) Print(v ...interface{}) {
	if output != "" {
		l.log.Print(v...)
	}
}

func (l *logger) Printf(f string, v ...interface{}) {
	if output != "" {
		l.log.Printf(f, v...)
	}
}

func (l *logger) Println(v ...interface{}) {
	if output != "" {
		l.log.Println(v...)
	}
}
