package httpx

type Error struct {
	Ok      bool
	Status  int
	Message string
}

func NewError()
