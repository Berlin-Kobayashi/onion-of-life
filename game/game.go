package game

import (
	"github.com/DanShu93/onion-of-life/oled"
	"encoding/csv"
	"bufio"
	"os"
)

type world [][]bool

type Controller struct {
	view                      oled.Screen
	world                     world
	aliveAmounts, bornAmounts []int
}

func NewController(configPath string, aliveAmounts, bornAmounts []int, ) Controller {
	config := readCsv(configPath)

	world := newWorldFromConfig(config)

	return Controller{
		view:         oled.NewScreen(),
		world:        world,
		aliveAmounts: aliveAmounts,
		bornAmounts:  bornAmounts,
	}
}

func readCsv(csvPath string) [][]string {
	file, err := os.Open(csvPath)
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(bufio.NewReader(file))

	content, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	return content
}

func newWorldFromConfig(config [][]string) world {
	world := newWorld()

	rowOffset := len(world)/2 - len(config)/2
	columnOffset := len(world[0])/2 - len(config[0])/2

	for x := 0; x < len(config); x++ {
		for y := 0; y < len(config[0]); y++ {
			if config[x][y] == "1" {
				world[x+rowOffset][y+columnOffset] = true
			}
		}
	}

	return world
}

func newWorld() world {
	world := make(world, oled.PixelsPerColumn)
	for i := range world {
		world[i] = make([]bool, oled.PixelsPerRow)
	}

	return world
}

func (c *Controller) Play() {
	for {
		c.render()

		nextWorld := newWorld()
		for x, row := range c.world {
			for y := range row {
				nextWorld[x][y] = c.isCellAlive(x, y)
			}
		}
		c.world = nextWorld
	}
}

func (c *Controller) render() {
	for y := 0; y < len(c.world[0]); y++ {
		c.view.MoveCursorToPixel(0, y)

		var characterValue byte
		for x := 0; x < len(c.world); x++ {
			characterX := x % oled.PixelsPerCharacter

			if characterX == 0 && x != 0 {
				c.view.WriteByte(characterValue)
				c.view.MoveCursorToPixel(x/oled.PixelsPerCharacter, y)

				characterValue = 0
			}
			if c.world[x][y] {
				characterValue += 1 << byte(characterX)

			}
		}
		c.view.WriteByte(characterValue)
	}
	c.view.Commit()
}

func (c *Controller) isCellAlive(x, y int) bool {
	wasCellAlive := c.world[x][y]
	neighboursAlive := 0
	for neighbourX := x - 1; neighbourX <= x+1; neighbourX++ {
		for neighbourY := y - 1; neighbourY <= y+1; neighbourY++ {
			if neighbourX > 0 && neighbourX < len(c.world) {
				if neighbourY > 0 && neighbourY < len(c.world[0]) {
					if !(neighbourX == x && neighbourY == y) {
						if c.world[neighbourX][neighbourY] {
							neighboursAlive++
						}
					}
				}
			}
		}
	}

	var amounts []int
	if wasCellAlive {
		amounts = c.aliveAmounts
	} else {
		amounts = c.bornAmounts
	}

	for _, amount := range amounts {
		if neighboursAlive == amount {
			return true
		}
	}

	return false
}
