package main

import (
	"math"
)

type baseAI struct {
	action
	baseComponent
	Entity *actor
}

type node struct {
	Parent   *node
	Position [2]int
	G        int
	H        int
	F        int
}

func (n node) Eq(t *node) bool {
	return n.Position[0] == t.Position[0] && n.Position[1] == t.Position[1]
}

// https://medium.com/@nicholas.w.swift/easy-a-star-pathfinding-7e6689c7f7b2
func aster(tiles [][]*tile, start, end [2]int) [][2]int {
	startNode := &node{Parent: nil, Position: start, G: 0, H: 0, F: 0}
	endNode := &node{Parent: nil, Position: end, G: 0, H: 0, F: 0}

	openList := []*node{}
	closedList := []*node{}

	var currentNode *node

	openList = append(openList, startNode)
	for len(openList) > 0 {

		// Get the current node
		goal := false
		currentNode = openList[0]
		currentIndex := 0
		for i, item := range openList {
			// Found the goal
			if item.Eq(endNode) {
				goal = true
				currentNode = item
				currentIndex = i
				break
			}
			if item.F < currentNode.F {
				currentNode = item
				currentIndex = i
			}
		}
		if goal {
			break
		}

		// Pop current off open list, add to closed list
		// https://zenn.dev/mattn/articles/31dfed3c89956d#copy-%E3%82%92%E4%BD%BF%E3%81%86%E6%96%B9%E6%B3%95
		openList = openList[:currentIndex+copy(openList[currentIndex:], openList[currentIndex+1:])]
		closedList = append(closedList, currentNode)

		// // Found the goal
		// if currentNode.Eq(endNode) {
		// 	break
		// }

		// Generate children
		children := []*node{}
		for _, newPosition := range [][2]int{
			{0, -1},
			{0, 1},
			{-1, 0},
			{1, 0},
			{-1, -1},
			{-1, 1},
			{1, -1},
			{1, 1},
		} {
			// Get node position
			nodePosition := [2]int{currentNode.Position[0] + newPosition[0], currentNode.Position[1] + newPosition[1]}

			// Make sure within range
			if nodePosition[0] > len(tiles)-1 || nodePosition[0] < 0 || nodePosition[1] > len(tiles[len(tiles)-1])-1 || nodePosition[1] < 0 {
				continue
			}

			// Make sure walkable terrain
			if !tiles[nodePosition[0]][nodePosition[1]].Walkable {
				continue
			}

			// Create new node
			newNode := &node{Parent: currentNode, Position: nodePosition}

			// Append
			children = append(children, newNode)
		}

		// Loop through children
		for _, child := range children {

			// Child is on the closed list
			closed := false
			for _, closedChild := range closedList {
				if child.Eq(closedChild) {
					closed = true
					break
				}
			}
			if closed {
				continue
			}

			// Create the f, g, and h values
			child.G = currentNode.G + 1
			child.H = int(math.Pow(float64(child.Position[0]-endNode.Position[0]), 2) + math.Pow(float64(child.Position[1]-endNode.Position[1]), 2))
			child.F = child.G + child.H

			// Child is already in the open list
			opened := false
			for _, openNode := range openList {
				if child.Eq(openNode) && child.G > openNode.G {
					opened = true
				}
			}
			if opened {
				continue
			}

			// Add the child to the open list
			openList = append(openList, child)
		}
	}
	path := [][2]int{}
	current := currentNode
	for current != nil {
		path = append(path, current.Position)
		current = current.Parent
	}
	for i := 0; i < len(path)/2; i++ {
		path[i], path[len(path)-i-1] = path[len(path)-i-1], path[i]
	}
	return path[1:]
}

type hostileEnemy struct {
	baseAI
	Path [][2]int
}

func NewHostileEnemy(entity *actor) *hostileEnemy {
	return &hostileEnemy{
		baseAI: baseAI{
			Entity: entity,
			baseComponent: baseComponent{
				Entity: entity,
			},
		},
		Path: [][2]int{},
	}
}

func (ai hostileEnemy) GetPathTo(destX, destY int) [][2]int {
	return aster(ai.Entity.GameMap.Tiles, [2]int{ai.Entity.X, ai.Entity.Y}, [2]int{destX, destY})
}

func (ai hostileEnemy) Perform() error {
	target := ai.Engine().Player
	dx := target.X - ai.Entity.X
	dy := target.Y - ai.Entity.Y
	distance := int(math.Max(math.Abs(float64(dx)), math.Abs(float64(dy))))

	if ai.Engine().GameMap.Visible[ai.Entity.X][ai.Entity.Y] {
		if distance <= 1 {
			return meleeAction{
				actionWithDirection{
					baseAction: baseAction{
						Entity: ai.Entity,
					},
					Dx: dx,
					Dy: dy,
				},
			}.Perform()
		}

		ai.Path = ai.GetPathTo(target.X, target.Y)
	}

	if len(ai.Path) > 0 {
		destX, destY := ai.Path[0][0], ai.Path[0][1]
		return movementAction{
			actionWithDirection{
				baseAction: baseAction{
					Entity: ai.Entity,
				},
				Dx: destX - ai.Entity.X,
				Dy: destY - ai.Entity.Y,
			},
		}.Perform()
	}

	return waitAction{}.Perform()
}
