package blub

import "sync"

type Observer[T any] func(oldV, newV T)

type box[T any] struct {
	v          T
	observers  map[int]Observer[T]
	nextHandle int
	mutex      sync.Mutex
}

func (b *box[T]) add(o Observer[T]) int {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.nextHandle++
	b.observers[b.nextHandle] = o
	return b.nextHandle
}

func (b *box[T]) remove(hnd int) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	delete(b.observers, hnd)
}

func (b *box[T]) set(v T) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	old := b.v
	b.v = v
	// note: this is prone to deadlocks for cascades etc. and we should post to an event loop for time slicing anyway
	for _, observer := range b.observers {
		observer(old, v)
	}
}

func (b *box[T]) get() T {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	return b.v
}

// A State is a stack-allocated holder which contains a heap-box and
// represents a reactive data holder between a platform renderer and a cross-platform implementation.
// Especially useful to share or override states of subcomponents. See also Share method.
type State[T any] struct {
	box   *box[T]
	mutex sync.Mutex
}

func NewState[T any](t T) State[T] {
	s := State[T]{}
	s.SetValue(t)
	return s
}

// Share mutates this state to ensure that this and the returned State shares the same heap box.
func (s *State[T]) Share() State[T] {
	s.allocBox()
	return *s
}

func (s *State[T]) allocBox() *box[T] {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.box == nil {
		s.box = &box[T]{}
	}

	return s.box
}

func (s *State[T]) Value() T {
	return s.allocBox().get()
}

func (s *State[T]) SetValue(v T) {
	s.allocBox().set(v)
}

func (s *State[T]) AddObserver(obs Observer[T]) int {
	return s.allocBox().add(obs)
}

func (s *State[T]) RemoveObserver(h int) {
	s.allocBox().remove(h)
}
