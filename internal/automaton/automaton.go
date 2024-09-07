package automaton

import (
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"runtime"
	"sync"
)

type CellularAutomaton struct {
	size         int
	cells        [][]bool
	rule         func(bool, int) bool
	neighborhood string
}

func New(size int, rule func(bool, int) bool, neighborhood string) *CellularAutomaton {
	cells := make([][]bool, size)
	for i := range cells {
		cells[i] = make([]bool, size)
		for j := range cells[i] {
			cells[i][j] = rand.Intn(2) == 1
		}
	}
	return &CellularAutomaton{size, cells, rule, neighborhood}
}

func (ca *CellularAutomaton) getNeighbors(x, y int) []bool {
	switch ca.neighborhood {
	case "moore":
		return []bool{
			ca.cells[(x-1+ca.size)%ca.size][(y-1+ca.size)%ca.size],
			ca.cells[(x-1+ca.size)%ca.size][y],
			ca.cells[(x-1+ca.size)%ca.size][(y+1)%ca.size],
			ca.cells[x][(y-1+ca.size)%ca.size],
			ca.cells[x][(y+1)%ca.size],
			ca.cells[(x+1)%ca.size][(y-1+ca.size)%ca.size],
			ca.cells[(x+1)%ca.size][y],
			ca.cells[(x+1)%ca.size][(y+1)%ca.size],
		}
	case "von_neumann":
		return []bool{
			ca.cells[(x-1+ca.size)%ca.size][y],
			ca.cells[x][(y-1+ca.size)%ca.size],
			ca.cells[x][(y+1)%ca.size],
			ca.cells[(x+1)%ca.size][y],
		}
	default:
		panic("Unsupported neighborhood type")
	}
}

func (ca *CellularAutomaton) applyRule(cell bool, neighbors []bool) bool {
	count := 0
	for _, n := range neighbors {
		if n {
			count++
		}
	}
	return ca.rule(cell, count)
}

func (ca *CellularAutomaton) Step() {
	newCells := make([][]bool, ca.size)
	for i := range newCells {
		newCells[i] = make([]bool, ca.size)
	}

	numWorkers := runtime.NumCPU()
	var wg sync.WaitGroup
	rowChan := make(chan int, ca.size)

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for row := range rowChan {
				for col := 0; col < ca.size; col++ {
					neighbors := ca.getNeighbors(row, col)
					newCells[row][col] = ca.applyRule(ca.cells[row][col], neighbors)
				}
			}
		}()
	}

	for row := 0; row < ca.size; row++ {
		rowChan <- row
	}
	close(rowChan)
	wg.Wait()

	ca.cells = newCells
}

func (ca *CellularAutomaton) Run(steps int) ([]image.Image, error) {
	frames := make([]image.Image, steps)
	for i := 0; i < steps; i++ {
		ca.Step()
		img := image.NewRGBA(image.Rect(0, 0, ca.size, ca.size))
		for y := 0; y < ca.size; y++ {
			for x := 0; x < ca.size; x++ {
				if ca.cells[y][x] {
					img.Set(x, y, color.Black)
				} else {
					img.Set(x, y, color.White)
				}
			}
		}
		frames[i] = img
	}
	return frames, nil
}

func SaveImage(img image.Image, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

func GameOfLifeRule(cell bool, count int) bool {
	return count == 3 || (cell && count == 2)
}
