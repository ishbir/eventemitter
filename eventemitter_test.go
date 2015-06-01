package eventemitter

import (
	"fmt"
	"testing"
)

// Struct for testing Embedding of EventEmitters
type Server struct {
	EventEmitter
}

func TestEmbedding(t *testing.T) {
	s := new(Server)

	// Don't forget to allocate the memory when
	// used as sub type.
	s.EventEmitter.Init()

	s.On("recv", func(msg string) string {
		return msg
	}, 10)

	resp := <-s.Emit("recv", "Hello World")

	expected := "Hello World"

	if res := resp.Ret[0].(string); res != expected {
		t.Errorf("Expected %s, got %s", expected, res)
	}

	// Remove non-existing listener
	s.EventEmitter.RemoveListener("recv", 5)
	if l := len(s.EventEmitter.Listeners("recv")); l != 1 {
		t.Errorf("Expected length 1, got %d", l)
	}

	// Remove existing listener.
	s.EventEmitter.RemoveListener("recv", 10)
	if l := len(s.EventEmitter.Listeners("recv")); l != 0 {
		t.Errorf("Expected length 0, got %d", l)
	}

	// Test RemoveListeners
	s.On("recv", func(msg string) string {
		return msg
	}, 1)
	s.On("recv", func(msg string) string {
		return msg
	}, 2)
	s.On("recv", func(msg string) string {
		return msg
	}, 3)
	if l := len(s.EventEmitter.Listeners("recv")); l != 3 {
		t.Errorf("Expected length 3, got %d", l)
	}

	s.EventEmitter.RemoveListeners("recv")
	if l := len(s.EventEmitter.Listeners("recv")); l != 0 {
		t.Errorf("Expected length 0, got %d", l)
	}
}

func ExampleEmitReturnsEventOnChan() {
	emitter := New()

	emitter.On("hello", func(name string) string {
		return "Hello World " + name
	}, 5)

	e := <-emitter.Emit("hello", "John")

	fmt.Println(e.EventName)
	// Output:
	// hello
}

func BenchmarkEmit(b *testing.B) {
	b.StopTimer()
	emitter := New()
	nListeners := 100

	for i := 0; i < nListeners; i++ {
		emitter.On("hello", func(name string) string {
			return "Hello World " + name
		}, i)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		<-emitter.Emit("hello", "John")
	}
}
