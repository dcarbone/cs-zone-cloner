package definition

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

type (
	Formatter func(*ZoneDefinition) ([]byte, error)
)

var (
	formattersMu sync.Mutex
	formatters   map[string]Formatter
)

func init() {
	formatters = map[string]Formatter{
		"json":        FormatJSON,
		"json-indent": FormatJSONIndent,
	}
}

func FormatJSON(zd *ZoneDefinition) ([]byte, error) {
	if zd == nil {
		return nil, errors.New("zone definition cannot be empty")
	}
	return json.Marshal(zd)
}

func FormatJSONIndent(zd *ZoneDefinition) ([]byte, error) {
	if zd == nil {
		return nil, errors.New("zone definition cannot be empty")
	}
	return json.MarshalIndent(zd, "", "\t")
}

func SetFormatter(name string, fn Formatter) {
	formattersMu.Lock()
	formatters[name] = fn
	formattersMu.Unlock()
}

func Format(zd *ZoneDefinition, format string) ([]byte, error) {
	formattersMu.Lock()
	fn, ok := formatters[format]
	if !ok {
		formattersMu.Unlock()
		return nil, fmt.Errorf("no formatter named \"%s\" found", format)
	}
	formattersMu.Unlock()
	return fn(zd)
}
