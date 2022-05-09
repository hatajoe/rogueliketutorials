package main

type Action interface {}

type NoneAction struct {}

type MovementAction struct {
	Dx int
	Dy int
}

type EscapeAction struct {}
