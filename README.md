# monigote

## Manifesto

Creating interface mocks with ability of spying for unit tests in go is crucial. Monigote is a tool that auto generate mocks given a interface. These mocks come with batteries providing an easy to use interface to create assertions and perform unit tests.

## Installation

```sh
go install github.com/devtin/monigote@latest
```

## At a glance

Given 

> file: `internal/io/x-api-service/models/models.go`

```go
package xapiservicemodels

//go:generate monigote x-api-service.go -i XApiService

interface XApiService struct {
    GetTemperature() (float32, error)
}
```

When calling:

```sh
go generate
```

**Then generates:**

> file: `internal/io/x-api-service/testing/mock/mock.go`

```go
package xapiservicemock

import (
    "github.com/devtin/monigote"
)

type XApiServiceMock struct {
    *monigote.Monigote
}

func NewXApiServiceMock(t *testing.T) *XApiServiceMock {
    return &XApiServiceMock{
        Monigote: monigote.NewMonigote("XApiServiceMock", t)
    }
}
```

**And:**

> file: `internal/io/testing/mock/GetTemperature.go`

```go
package xapiservicemock

func (m *XApiServiceMock) GetTemperature(id string) (float32, error) {
    reply := m.Monigote.Call("GetTemperature", id)

    var r1 float32
    var r2 error

    if reply[0] != nil {
        r1 = reply[0].(float32)
    }

    if reply[1] != nil {
        r2 = reply[1].(error)
    }

    return r1, r2
}
```

**Which can be used as...**

> file: internal/core/service/models/models.go

```go
package test

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSomething(t *testing.T) {
	Convey("My test", t, func() {
		// setup
		Convey("Expectation", func() {
			// test
			So(true, ShouldBeTrue)
		})
	})
}

```
