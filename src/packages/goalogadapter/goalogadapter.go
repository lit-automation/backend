package goalogadapter

import (
	"context"
	"fmt"

	"github.com/goadesign/goa"
	"github.com/sirupsen/logrus"
)

// This goalogadapter is an extension on goalogrus
// It adds functionality to use a dummyAdapter which doesn't log

type (
	// adapter is the logrus goa logger adapter.
	adapter struct {
		*logrus.Entry
	}

	dummyAdapter adapter
)

func (a *dummyAdapter) New(data ...interface{}) goa.LogAdapter { return a }
func (a *dummyAdapter) Info(msg string, data ...interface{})   {}
func (a *dummyAdapter) Error(msg string, data ...interface{})  {}

// New creates a new logger given a context.
func (a *adapter) New(data ...interface{}) goa.LogAdapter {
	if len(data) > 1 {
		switch fmt.Sprintf("%v", data[0]) {
		case "ignore":
			a2 := dummyAdapter(*a)
			return &a2
		}
	}
	return &adapter{Entry: a.Entry.WithFields(data2rus(data))}
}

// New wraps a logrus logger into a goa logger.
func New(logger *logrus.Logger) goa.LogAdapter {
	return FromEntry(logrus.NewEntry(logger))
}

// FromEntry wraps a logrus log entry into a goa logger.
func FromEntry(entry *logrus.Entry) goa.LogAdapter {
	return &adapter{Entry: entry}
}

// Entry returns the logrus log entry stored in the given context if any, nil otherwise.
func Entry(ctx context.Context) *logrus.Entry {
	logger := goa.ContextLogger(ctx)
	if a, ok := logger.(*adapter); ok {
		return a.Entry
	}
	return nil
}

// Info logs messages using logrus.
func (a *adapter) Info(msg string, data ...interface{}) {
	a.Entry.WithFields(data2rus(data)).Info(msg)
}

// Error logs errors using logrus.
func (a *adapter) Error(msg string, data ...interface{}) {
	a.Entry.WithFields(data2rus(data)).Error(msg)
}

func data2rus(keyvals []interface{}) logrus.Fields {
	n := (len(keyvals) + 1) / 2
	res := make(logrus.Fields, n)
	for i := 0; i < len(keyvals); i += 2 {
		k := keyvals[i]
		var v interface{} = goa.ErrMissingLogValue
		if i+1 < len(keyvals) {
			v = keyvals[i+1]
		}
		res[fmt.Sprintf("%v", k)] = v
	}
	return res
}
