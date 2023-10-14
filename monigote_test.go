package monigote

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type MyMonigote struct {
	*Monigote
}

type MyStruct struct {
}

func NewMyMonigote(t *testing.T) *MyMonigote {
	return &MyMonigote{
		NewMonigote("MyMonigoteMock", t),
	}
}

func (m *MyMonigote) MyMethod(num int) bool {
	return m.Call("MyMethod", num)[0].(bool)
}

func (m *MyMonigote) MyMethodB(b float32, num int) bool {
	return m.Call("MyMethodB", b, num)[0].(bool)
}

func TestSomething(t *testing.T) {
	Convey("MyMonigote", t, func() {
		Convey("Strict Setup", func() {
			// init
			myMonigote := NewMyMonigote(t)

			// methods setup
			myMonigote.
				Setup("MyMethod").
				WhenCalledWith(1).
				ReplyWith(true /*, ... */)

			myMonigote.
				Setup("MyMethod").
				WhenCalledWith(0).
				ReplyWith(false)

			Convey("Mocks Methods", func() {
				// test
				So(myMonigote.MyMethod(1), ShouldBeTrue)
				So(myMonigote.IsDone(), ShouldBeFalse)
				So(myMonigote.MyMethod(0), ShouldBeFalse)
				So(myMonigote.IsDone(), ShouldBeTrue)
				So(len(myMonigote.Calls["MyMethod"]), ShouldEqual, 2)
				So(myMonigote.Calls["MyMethod"][0].Args[0], ShouldEqual, 1)
				So(myMonigote.Calls["MyMethod"][1].Args[0], ShouldEqual, 0)
			})
		})

		Convey("Loose Setup", func() {
			// init
			myMonigote := NewMyMonigote(t)

			// methods setup
			myMonigote.
				Setup("MyMethod").
				WhenCalled(func(args []interface{}) bool {
					num := args[0].(int)
					return num > 10
				}).
				ReplyWith(false /*, ... */).
				Persist()

			myMonigote.
				Setup("MyMethod").
				WhenCalledWith(1).
				ReplyWith(true)

			myMonigote.
				Setup("MyMethodB").
				WhenCalled().
				ReplyWith(false).
				ReplyTimes(2)

			Convey("Mocks Methods", func() {
				// test
				So(myMonigote.MyMethod(1), ShouldBeTrue)
				So(myMonigote.IsDone(), ShouldBeFalse)
				So(myMonigote.MyMethod(11), ShouldBeFalse)
				So(myMonigote.MyMethodB(0.3, 1), ShouldBeFalse)
				So(myMonigote.MyMethodB(3, 2), ShouldBeFalse)
				So(myMonigote.IsDone(), ShouldBeTrue)
				So(len(myMonigote.Calls["MyMethod"]), ShouldEqual, 2)
				So(myMonigote.Calls["MyMethod"][0].Args[0], ShouldEqual, 1)
				So(myMonigote.Calls["MyMethod"][1].Args[0], ShouldEqual, 11)
			})
		})
	})
}
