package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"snake/path"

	"strconv"
	"time"

	"github.com/gdamore/tcell"
)

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)

	screen.SetStyle(defStyle)

	width, height := screen.Size()
	grid := makeGrid(width, height)

	snakeParts := []SnakePart{
		{
			X: 5,
			Y: 10,
		},
		{
			X: 6,
			Y: 10,
		},
		{
			X: 7,
			Y: 10,
		},
	}

	snakeBody := SnakeBody{
		Parts:  snakeParts,
		Xspeed: 1,
		Yspeed: 0,
	}

	game := Game{
		Screen:    screen,
		snakeBody: snakeBody,
		Grid:      grid,
	}

	go game.Run()

	for {
		switch event := game.Screen.PollEvent().(type) {
		case *tcell.EventResize:
			game.Screen.Sync()
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
				game.Screen.Fini()
				os.Exit(0)
			} else if event.Rune() == 'w' || event.Rune() == 'W' && game.snakeBody.Yspeed != 1 {
				game.snakeBody.ChangeDir(-1, 0)
			} else if event.Rune() == 's' || event.Rune() == 'S' && game.snakeBody.Yspeed != -1 {
				game.snakeBody.ChangeDir(1, 0)
			} else if event.Rune() == 'a' || event.Rune() == 'A' && game.snakeBody.Xspeed != 1 {
				game.snakeBody.ChangeDir(0, -1)
			} else if event.Rune() == 'd' || event.Rune() == 'D' && game.snakeBody.Xspeed != -1 {
				game.snakeBody.ChangeDir(0, 1)
			} else if event.Rune() == 'y' || event.Rune() == 'Y' && game.GameOver {
				go game.Run()
			} else if event.Rune() == 'n' || event.Rune() == 'N' && game.GameOver {
				game.Screen.Fini()
				os.Exit(0)
			}
		}
	}
}

func makeGrid(width, height int) [][]int {
	grid := make([][]int, width)
	for i := range grid {
		grid[i] = make([]int, height)
	}

	return grid
}

func transpose(slice [][]int) [][]int {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]int, xl)
	for i := range result {
		result[i] = make([]int, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}

///////////////////////////////////////////////////////////////////////////////////////////////////

type SnakePart struct {
	X int
	Y int
}

type SnakeBody struct {
	Parts  []SnakePart
	Xspeed int
	Yspeed int
}

func (sb *SnakeBody) ChangeDir(vertical, horizontal int) {
	sb.Yspeed = vertical
	sb.Xspeed = horizontal
}

func (sb *SnakeBody) Update(width, height int, longerSnake bool) {
	sb.Parts = append(sb.Parts, sb.Parts[len(sb.Parts)-1].GetUpdatedPart(sb, width, height))
	if !longerSnake {
		sb.Parts = sb.Parts[1:]
	}
}

func (sb *SnakeBody) ResetPos(width, height int) {

	snakeParts := []SnakePart{
		{
			X: int(width / 2),
			Y: int(height / 2),
		},
		{
			X: int(width/2) + 1,
			Y: int(height / 2),
		},
		{
			X: int(width/2) + 2,
			Y: int(height / 2),
		},
	}

	sb.Parts = snakeParts
	sb.Xspeed = 1
	sb.Yspeed = 0
}

func (sp *SnakePart) GetUpdatedPart(sb *SnakeBody, width int, height int) SnakePart {
	newPart := *sp
	newPart.X = (newPart.X + sb.Xspeed) % width

	if newPart.X < 0 {
		newPart.X += width
	}

	newPart.Y = (newPart.Y + sb.Yspeed) % height

	if newPart.Y < 0 {
		newPart.Y += height
	}

	return newPart
}

///////////////////////////////////////////////////////////////////////////////////////////////////

type Game struct {
	Screen         tcell.Screen
	snakeBody      SnakeBody
	FoodPos        SnakePart
	Score          int
	GameOver       bool
	Grid           [][]int
	foodPosChanged bool
	Path           []*path.Node
}

func drawParts(s tcell.Screen, snakeParts []SnakePart, foodPos SnakePart, snakeStyle tcell.Style, foodStyle tcell.Style) (*path.Node, *path.Node, [][]int) {
	s.SetContent(foodPos.X, foodPos.Y, '\u25CF', nil, foodStyle)

	for _, part := range snakeParts {
		s.SetContent(part.X, part.Y, ' ', nil, snakeStyle)
	}

	width, height := s.Size()
	grid := makeGrid(width, height)

	for _, part := range snakeParts {
		grid[part.X][part.Y] = 1
	}
	grid[foodPos.X][foodPos.Y] = 3
	grid[snakeParts[len(snakeParts)-1].X][snakeParts[len(snakeParts)-1].Y] = 2
	start := &path.Node{X: snakeParts[len(snakeParts)-1].X, Y: snakeParts[len(snakeParts)-1].Y}
	dest := &path.Node{X: foodPos.X, Y: foodPos.Y}

	resGrid := transpose(grid)

	return start, dest, resGrid
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, text string) {
	row := y1
	col := x1
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)

	for _, r := range text {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}

		if row > y2 {
			break
		}
	}
}

func checkCollision(parts []SnakePart, otherPart SnakePart) bool {
	for _, part := range parts {
		if part.X == otherPart.X && part.Y == otherPart.Y {
			return true
		}
	}
	return false
}

func (g *Game) UpdateFoodPos(width, height int) {
	prevFoodPos := g.FoodPos

	g.FoodPos.X = rand.Intn(width)
	g.FoodPos.Y = rand.Intn(height)

	if g.FoodPos.Y == 1 && g.FoodPos.X < 10 {
		g.UpdateFoodPos(width, height)
	}

	// Check if the food position has changed
	if prevFoodPos.X != g.FoodPos.X || prevFoodPos.Y != g.FoodPos.Y {
		g.foodPosChanged = true
	} else {
		g.foodPosChanged = false
	}
}

func (g *Game) Run() {
	defer func() {
		if r := recover(); r != nil {
			g.Screen.Fini()
			fmt.Fprintf(os.Stderr, "Panic occurred: %v\n", r)
			os.Exit(1)
		}
	}()

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	g.Screen.SetStyle(defStyle)

	width, height := g.Screen.Size()

	g.snakeBody.ResetPos(width, height)
	g.UpdateFoodPos(width, height)
	g.GameOver = false
	g.Score = 0
	g.Path = nil
	snakeStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorWhite)
	f, err := os.OpenFile("output.log", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
	log.SetFlags(0)

	for {
		longerSnake := false
		g.Screen.Clear()
		if checkCollision(g.snakeBody.Parts[len(g.snakeBody.Parts)-1:], g.FoodPos) {
			g.UpdateFoodPos(width, height)
			longerSnake = true
			g.Score++

		}

		if checkCollision(g.snakeBody.Parts[:len(g.snakeBody.Parts)-1], g.snakeBody.Parts[len(g.snakeBody.Parts)-1]) {
			break
		}

		g.snakeBody.Update(width, height, longerSnake)

		start, dest, grid := drawParts(g.Screen, g.snakeBody.Parts, g.FoodPos, snakeStyle, defStyle)

		drawText(g.Screen, 1, 1, 8+len(strconv.Itoa(g.Score)), 1, "Score: "+strconv.Itoa(g.Score))

		newPath := path.AStarSearch(start, dest, grid)

		g.Path = newPath
		g.pathSnake(newPath, &g.snakeBody)

		time.Sleep(40 * time.Millisecond)
		g.Screen.Show()
	}

	g.GameOver = true
	drawText(g.Screen, width/2-20, height/2, width/2+20, height/2, "Game Over, Score: "+strconv.Itoa(g.Score)+", Play Again? y/n")
	g.Screen.Show()
}

func (g *Game) pathSnake(path []*path.Node, sb *SnakeBody) {
	for i := 0; i < len(path)-1; i++ {
		x, y := calcDifference(path[i].X, path[i+1].X, path[i].Y, path[i+1].Y)
		if (x != -sb.Xspeed || y != -sb.Yspeed) && (x != sb.Xspeed || y != sb.Yspeed) {
			sb.ChangeDir(y, x)
		}
	}
	g.Path = nil
}

// func (g *Game) pathSnake(path []*path.Node, sb *SnakeBody) {
// 	if len(path) < 2 {
// 		return
// 	}

// 	// Get the current direction of the snake
// 	currentX, currentY := sb.Xspeed, sb.Yspeed

// 	for i := 0; i < len(path)-1; i++ {
// 		nextX, nextY := calcDifference(path[i].X, path[i+1].X, path[i].Y, path[i+1].Y)

// 		// Check if the next direction is valid (not opposite or equal to the current direction)
// 		if (nextX != -currentX || nextY != -currentY) && (nextX != currentX || nextY != currentY) {
// 			sb.ChangeDir(nextY, nextX)
// 			break
// 		}
// 	}

// 	// Move one step forward in the path
// 	sb.Update(len(g.Grid), len(g.Grid[0]), false)
// }

func calcDifference(x1, x2, y1, y2 int) (int, int) {
	return x2 - x1, y2 - y1
}
