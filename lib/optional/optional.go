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

func (o Optional[T]) Map[U any](f func(T) U) Optional[U] {
	if !o.has {
		return None[U]()
	}
	return Some(f(o.value))
}

func (o Optional[T]) Map2[U any, V any](f func(T) (U, V)) (Optional[U], V) {
	if !o.has {
		var v V
		return None[U](), v
	}
	var u, v = f(o.value)
	return Some(u), v
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

func (o *Optional[T]) SetIfAbsent(value T) bool {
	if o.has {
		return false
	}
	o.value = value
	o.has = true
	return true
}
