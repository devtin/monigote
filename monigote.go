package monigote

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

type MonigoteSetup struct {
	Method    string
	Args      []interface{}
	Reply     []interface{}
	Matchers  []func(args []interface{}) bool
	Times     int
	WasCalled uint
}

type MonigoteCall struct {
	Method     string
	Args       []interface{}
	Reply      []interface{}
	Found      bool
	Persistent bool
	CalledAt   time.Time
}

func NewMonigoteSetup() *MonigoteSetup {
	return &MonigoteSetup{
		Times:     1,
		WasCalled: 0,
	}
}

// Sets up the arguments to match for the method call
func (m *MonigoteSetup) WhenCalledWith(args ...interface{}) *MonigoteSetup {
	m.Args = args
	return m
}

func (m *MonigoteSetup) WhenCalled(matchers ...func(args []interface{}) bool) *MonigoteSetup {
	m.Args = nil
	m.Matchers = matchers
	return m
}

func (m *MonigoteSetup) ReplyWith(reply ...interface{}) *MonigoteSetup {
	m.Reply = reply
	return m
}

func (m *MonigoteSetup) ReplyTimes(times int) *MonigoteSetup {
	m.Times = times
	return m
}

func (m *MonigoteSetup) Persist() *MonigoteSetup {
	m.Times = -1 // will mock forever
	return m
}

// usage: monigote.Setup("Method").WhenCalledWith("arg1", "arg2").ReplyWith("").Times(6)
type Monigote struct {
	Name   string
	Setups map[string][]*MonigoteSetup
	Calls  map[string][]*MonigoteCall
	mutex  sync.Mutex
	T      *testing.T
}

func NewMonigote(name string, t *testing.T) *Monigote {
	return &Monigote{
		Name:   name,
		Setups: map[string][]*MonigoteSetup{},
		Calls:  map[string][]*MonigoteCall{},
		mutex:  sync.Mutex{},
		T:      t,
	}
}

// Creates a new setup for a method call
func (m *Monigote) Setup(method string) *MonigoteSetup {
	setup := NewMonigoteSetup()
	setup.Method = method
	m.Setups[method] = append(m.Setups[method], setup)
	return setup
}

func (m *Monigote) removeSetup(setup *MonigoteSetup) {
	for i, s := range m.Setups[setup.Method] {
		if s == setup {
			m.Setups[setup.Method] = append(m.Setups[setup.Method][:i], m.Setups[setup.Method][i+1:]...)
			return
		}
	}
}

// Calls the method and returns the setup reply
func (m *Monigote) Call(method string, args ...interface{}) []interface{} {
	calledAt := time.Now()

	m.mutex.Lock()
	defer m.mutex.Unlock()

	monigoteCall := &MonigoteCall{
		Method:   method,
		Args:     args,
		Found:    false,
		CalledAt: calledAt,
	}

	for _, setup := range m.Setups[method] {
		if setup.Method == method {
			if setup.Args != nil && !reflect.DeepEqual(setup.Args, args) {
				continue
			}

			// check matchers
			if setup.Matchers != nil {
				matched := true
				for _, matcher := range setup.Matchers {
					if !matcher(args) {
						matched = false
						break
					}
				}
				if !matched {
					continue
				}
			}

			setup.WasCalled += 1

			if setup.Times < 0 {
				// persistent setup
				monigoteCall.Persistent = true
			}

			if setup.Times > 0 {
				setup.Times -= 1
				if setup.Times == 0 {
					m.removeSetup(setup)
				}
			}

			monigoteCall.Found = true
			monigoteCall.Reply = setup.Reply
			m.Calls[method] = append(m.Calls[method], monigoteCall)
			return setup.Reply
		}
	}

	// Fails due to unexpected call
	m.T.Fatalf("No setup found for method %s with args %v", method, args)

	return nil
}

// Resets the mock state
func (m *Monigote) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.Setups = map[string][]*MonigoteSetup{}
	m.Calls = map[string][]*MonigoteCall{}
}

// Verifies that all setups were called
func (m *Monigote) IsDone() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	err := false
	for _, setups := range m.Setups {
		for _, setup := range setups {
			if setup.Times == -1 && setup.WasCalled > 0 {
				continue
			}
			err = true
			if setup.WasCalled == 0 {
				m.T.Logf("%s NOT DONE! %s.%s(%s) was not called!", m.Name, m.Name, setup.Method, strings.TrimRight(strings.TrimLeft(fmt.Sprintf("%v", setup.Args), "["), "]"))
			} else {
				m.T.Logf("%s NOT DONE! %s.%s(%s) was called %d times but have %d pending calls!", m.Name, m.Name, setup.Method, strings.TrimRight(strings.TrimLeft(fmt.Sprintf("%v", setup.Args), "["), "]"), setup.WasCalled, setup.Times)
			}
		}
	}
	return !err
}
