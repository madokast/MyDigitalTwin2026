package optional

type Optional[T any] struct {
	value T
	has   bool
}

func Some[T any](value T) Optional[T] {
	return Optional[T]{
		value: value,
		has:   true,
	}
}

func None[T any]() Optional[T] {
	var obj Optional[T] = Optional[T]{
		has: false,
	}
	return obj
}

func (o Optional[T]) Get() (T, bool) {
	return o.value, o.has
}

func (o Optional[T]) Or(other T) T {
	if o.has {
		return o.value
	}
	return other
}

func (o Optional[T]) Exists() bool {
	return o.has
}

func (o Optional[T]) Absent() bool {
	return !o.has
}

func (o *Optional[T]) SetIfAbsent(value T) bool {
	if o.has {
		return false
	}
	o.value = value
	o.has = true
	return true
}
