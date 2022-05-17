package main

type RenderOrder int

const (
	RenderOrder_Corpse RenderOrder = iota
	RenderOrder_Item
	RenderOrder_Actor
)
