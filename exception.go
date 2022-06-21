package main

type impossible struct {
	err string
}

func (e impossible) Error() string {
	return e.err
}

type QuitWithoutSaving struct {
	err string
}

func (e QuitWithoutSaving) Error() string {
	return e.err
}
