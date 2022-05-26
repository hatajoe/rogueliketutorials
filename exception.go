package main

type impossible struct {
	err string
}

func (e impossible) Error() string {
	return e.err
}
