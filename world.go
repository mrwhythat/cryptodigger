package main

import (
	"math"
	"math/rand"

	"github.com/faiface/pixel"
)

type BlockType = int

const BlockTypes = 4

const (
	SimpleBlock = iota
	SmallCoinBlock
	BigCoinBlock
	SurprizeBlock
)

// Block is an integral part of the world.
type Block struct {
	Type      BlockType
	Integrity int
}

func (block Block) Reward() int {
	r := 0
	switch block.Type {
	case SmallCoinBlock:
		r = 10
	case BigCoinBlock:
		r = 50
	case SurprizeBlock:
		r = 100
	}
	return r
}

// Cell is a pair of coordinates in a block grid
type Cell struct {
	X, Y int
}

func CellFromVec(pos pixel.Vec) Cell {
	j, i := pos.XY()
	var x float64
	if j < 0 {
		x = math.Floor(j/ASize + 0.5)
	} else {
		x = math.Ceil(j/ASize - 0.5)
	}
	ix, iy := int(x), int(math.Floor(i/ASize+0.5))
	return Cell{X: ix, Y: iy}
}

func (cell Cell) Right(pos pixel.Vec) Cell {
	if pos.X/ASize > float64(cell.X) {
		cell.X++
	}
	return cell
}

func (cell Cell) Left(pos pixel.Vec) Cell {
	if pos.X/ASize < float64(cell.X) {
		cell.X--
	}
	return cell
}

func (cell Cell) Up() Cell {
	return Cell{X: cell.X, Y: cell.Y - 1}
}

func (cell Cell) Down() Cell {
	return Cell{X: cell.X, Y: cell.Y + 1}
}

// Digger is a main character of the game.
type Digger struct {
	Coins int
}

func NewDigger() Digger {
	return Digger{}
}

func (digger *Digger) DigCell(world World, cell Cell) {
	digger.Coins += world.HammerBlock(cell)
}

// World contains game state.
type World struct {
	Grid Grid
}

type Grid map[Cell]*Block

func (grid Grid) Get(cell Cell) *Block {
	if _, ok := grid[cell]; !ok && cell.Y > 0 {
		grid[cell] = &Block{Type: rand.Intn(BlockTypes), Integrity: 8}
	}
	return grid[cell]
}

func (grid Grid) Del(cell Cell) {
	grid[cell] = nil
}

func NewWorld() World {
	return World{make(map[Cell]*Block)}
}

func (world World) GridView(min, max Cell) [][]*Block {
	height, width := max.Y-min.Y, max.X-min.X
	grid := make([][]*Block, height)
	for i := 0; i < height; i++ {
		grid[i] = make([]*Block, width)
		for j := 0; j < width; j++ {
			grid[i][j] = world.Grid.Get(Cell{Y: i + min.Y, X: j + min.X})
		}
	}
	return grid
}

// VisibleBlocks returns map from vectors to blocks, translating vectors into grid cell
// coordintates and building a map from a grid represented as matrix.
func (world World) VisibleBlocks(min pixel.Vec, max pixel.Vec) map[pixel.Vec]*Block {
	res := make(map[pixel.Vec]*Block)
	bMinX, bMinY := int(math.Ceil(min.X/ASize)), int(math.Ceil(min.Y/ASize))
	bMaxX := bMinX + int(math.Ceil((max.X-min.X)/ASize))
	bMaxY := bMinY + int(math.Ceil((max.Y-min.Y)/ASize))
	view := world.GridView(Cell{X: bMinX - 1, Y: bMinY - 1}, Cell{X: bMaxX + 1, Y: bMaxY + 1})
	for i := 0; i < len(view); i++ {
		for j := 0; j < len(view[i]); j++ {
			x, y := float64(j-1)*ASize+min.X, float64(i-1)*ASize+min.Y
			res[pixel.V(x, y)] = view[i][j]
		}
	}
	return res
}

// Check whether there is a block at a given cell.
func (world World) ContainsBlock(cell Cell) bool {
	return !(world.Grid.Get(cell) == nil)
}

// Kick a block with a hammer, decrementing its integrity.
// When integrity falls down to 0, block dissapears.
func (world World) HammerBlock(cell Cell) (coins int) {
	if block := world.Grid.Get(cell); block != nil {
		block.Integrity--
		if block.Integrity < 0 {
			r := block.Reward()
			coins += r
			world.Grid.Del(cell)
		}
	}
	return
}
