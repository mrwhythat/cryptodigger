package main

import (
	"math/rand"
)

type BlockType = int

const (
	BlockTypes = 4

	SimpleBlock = iota
	SmallBTCBlock
	BigBTCBlock
	SurprizeMFKBlock
)

// Block is an integral part of the world.
type Block struct {
	Type      BlockType
	Integrity int
}

// Cell is a pair of coordinates in a block grid
type Cell struct {
	X, Y int
}

// Digger is a main character of the game.
type Digger struct {
}

func NewDigger() Digger {
	return Digger{}
}

// World contains game state.
type World struct {
	Digger Digger
	Grid   map[Cell]*Block
}

func NewWorld() World {
	return World{
		Digger: NewDigger(),
		Grid:   make(map[Cell]*Block),
	}
}

func (world World) GridView(min, max Cell) [][]*Block {
	height, width := max.Y-min.Y, max.X-min.X
	grid := make([][]*Block, height)
	for i := 0; i < height; i++ {
		grid[i] = make([]*Block, width)
		// Initially digger should stand above the top block level
		if i+min.Y < 1 {
			continue
		}
		for j := 0; j < width; j++ {
			cell := Cell{Y: i + min.Y, X: j + min.X}
			block, ok := world.Grid[cell]
			if !ok {
				block = &Block{
					Type:      rand.Intn(BlockTypes),
					Integrity: 4,
				}
				world.Grid[cell] = block
			}
			grid[i][j] = block
		}
	}
	return grid
}

// Check whether there is a block at a given cell.
func (world World) ContainsBlock(cell Cell) bool {
	return !(world.Grid[cell] == nil)
}

// Kick a block with a hammer, decrementing its integrity.
// When integrity falls down to 0, block dissapears.
func (world World) HammerBlock(cell Cell) {
	if world.Grid[cell].Integrity <= 0 {
		world.Grid[cell] = nil
	} else {
		world.Grid[cell].Integrity--
	}
}
