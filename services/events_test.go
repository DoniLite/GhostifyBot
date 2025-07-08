package services_test

import (
	"sync"
	"testing"
	"time"

	"github.com/DoniLite/GhostifyBot/services"
)

func TestCreateEvent_ShouldRegisterAndReturnSameEvent(t *testing.T) {
	bus := services.EventBus

	ev1 := bus.CreateEvent("test:event")
	ev2 := bus.CreateEvent("test:event")

	if ev1 != ev2 {
		t.Errorf("Expected same event instance, got different ones")
	}
}

func TestOn_ShouldRegisterHandler(t *testing.T) {
	bus := services.EventBus
	event := bus.CreateEvent("on:event")

	called := false
	handler := func(data *services.EventData, args ...string) {
		called = true
	}

	bus.On(event, handler)
	bus.Emit(event, &services.EventData{Message: "Hello"})
	bus.Wait()

	if !called {
		t.Errorf("Handler was not called after Emit")
	}
}

func TestEmit_ShouldCallAllHandlers(t *testing.T) {
	bus := services.EventBus
	event := bus.CreateEvent("multi:handler")

	var mu sync.Mutex
	calls := 0

	handler1 := func(data *services.EventData, args ...string) {
		mu.Lock()
		calls++
		mu.Unlock()
	}
	handler2 := func(data *services.EventData, args ...string) {
		mu.Lock()
		calls++
		mu.Unlock()
	}

	bus.On(event, handler1)
	bus.On(event, handler2)
	bus.Emit(event, &services.EventData{Message: "Event triggered"})
	bus.Wait()

	if calls != 2 {
		t.Errorf("Expected 2 handlers to be called, got %d", calls)
	}
}

func TestOff_ShouldRemoveHandler(t *testing.T) {
	bus := services.EventBus
	event := bus.CreateEvent("off:event")

	called := false
	handler := func(data *services.EventData, args ...string) {
		called = true
	}

	bus.On(event, handler)
	bus.Off(event, handler)
	bus.Emit(event, &services.EventData{Message: "Should not be called"})
	bus.Wait()

	if called {
		t.Errorf("Handler was called after being removed")
	}
}

func TestEmit_ShouldPassArguments(t *testing.T) {
	bus := services.EventBus
	event := bus.CreateEvent("args:event")

	var received string
	handler := func(data *services.EventData, args ...string) {
		if len(args) > 0 {
			received = args[0]
		}
	}

	bus.On(event, handler)
	bus.Emit(event, &services.EventData{}, "Doni")
	bus.Wait()

	if received != "Doni" {
		t.Errorf("Expected argument 'Doni', got '%s'", received)
	}
}

func TestEmit_ShouldBeAsynchronous(t *testing.T) {
	bus := services.EventBus
	event := bus.CreateEvent("async:event")

	start := time.Now()
	handler := func(data *services.EventData, args ...string) {
		time.Sleep(50 * time.Millisecond)
	}

	bus.On(event, handler)
	bus.Emit(event, &services.EventData{})
	bus.Wait()

	elapsed := time.Since(start)
	if elapsed > 100*time.Millisecond {
		t.Errorf("Emit took too long, expected async behavior")
	}
}
