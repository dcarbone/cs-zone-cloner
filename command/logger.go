package command

import "sync"

type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
}

type MutableLogger interface {
	Mute()
	UnMute()
}

type mutableLogger struct {
	mu   sync.RWMutex
	mute bool
	log  Logger
}

func NewMutableLogger(l Logger) Logger {
	return &mutableLogger{log: l}
}

func (ml *mutableLogger) Mute() {
	ml.mu.Lock()
	ml.mute = true
	ml.mu.Unlock()
}

func (ml *mutableLogger) UnMute() {
	ml.mu.Lock()
	ml.mute = false
	ml.mu.Unlock()
}

func (ml *mutableLogger) Print(v ...interface{}) {
	if !ml.muted() {
		ml.log.Print(v...)
	}
}

func (ml *mutableLogger) Printf(f string, v ...interface{}) {
	if !ml.muted() {
		ml.log.Printf(f, v...)
	}
}

func (ml *mutableLogger) Println(v ...interface{}) {
	if !ml.muted() {
		ml.log.Println(v...)
	}
}

func (ml *mutableLogger) muted() bool {
	ml.mu.RLock()
	m := ml.mute
	ml.mu.RUnlock()
	return m
}
