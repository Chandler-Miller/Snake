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

	// grid := [][]int{
	// 	{2, 0, 0, 0, 0, 0, 1, 0, 0, 0},
	// 	{0, 0, 1, 0, 1, 0, 1, 0, 0, 0},
	// 	{0, 1, 1, 0, 1, 0, 1, 0, 0, 0},
	// 	{0, 0, 0, 0, 1, 0, 1, 0, 0, 0},
	// 	{0, 0, 1, 1, 1, 0, 1, 0, 0, 0},
	// 	{0, 0, 1, 0, 0, 0, 1, 0, 0, 0},
	// 	{0, 0, 0, 0, 0, 0, 0, 3, 0, 0},
	// 	{0, 1, 0, 0, 1, 0, 0, 0, 0, 0},
	// 	{1, 0, 0, 1, 0, 1, 0, 0, 0, 0},
	// 	{1, 0, 0, 0, 0, 1, 0, 0, 0, 0},
	// }

	// start := &path.Node{X: 0, Y: 0}
	// dest := &path.Node{X: 7, Y: 6}
	// resultPath := path.AStarSearch(start, dest, grid)

	// for _, i := range resultPath {
	// 	fmt.Println(i.X, i.Y)
	// }

	// path.PrintPathOnGrid(grid, start, dest, resultPath)

	screen, err := tcell.NewScreen()
	// width, height := screen.Size()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)

	screen.SetStyle(defStyle)

	width, height := screen.Size()
	grid := make([][]int, width)
	for i := 0; i < width; i++ {
		grid[i] = make([]int, height)
	}

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
			} else if event.Key() == tcell.KeyUp && game.snakeBody.Yspeed == 0 {
				game.snakeBody.ChangeDir(-1, 0)
			} else if event.Key() == tcell.KeyDown && game.snakeBody.Yspeed == 0 {
				game.snakeBody.ChangeDir(1, 0)
			} else if event.Key() == tcell.KeyLeft && game.snakeBody.Xspeed == 0 {
				game.snakeBody.ChangeDir(0, -1)
			} else if event.Key() == tcell.KeyRight && game.snakeBody.Xspeed == 0 {
				game.snakeBody.ChangeDir(0, 1)
			} else if event.Rune() == 'y' && game.GameOver {
				go game.Run()
			} else if event.Rune() == 'n' && game.GameOver {
				game.Screen.Fini()
				os.Exit(0)
			}
		}
	}
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
	Screen    tcell.Screen
	snakeBody SnakeBody
	FoodPos   SnakePart
	Score     int
	GameOver  bool
	Grid      [][]int
}

func drawParts(s tcell.Screen, snakeParts []SnakePart, foodPos SnakePart, snakeStyle tcell.Style, foodStyle tcell.Style, grid [][]int, checkPath bool) (*path.Node, *path.Node, [][]int) {
	s.SetContent(foodPos.X, foodPos.Y, '\u25CF', nil, foodStyle)

	for _, part := range snakeParts {
		s.SetContent(part.X, part.Y, ' ', nil, snakeStyle)
	}

	if checkPath {
		grid[foodPos.X][foodPos.Y] = 3
		grid[snakeParts[0].X][snakeParts[0].Y] = 2
		start := &path.Node{X: snakeParts[0].X, Y: snakeParts[0].Y}
		dest := &path.Node{X: foodPos.X, Y: foodPos.Y}

		for _, part := range snakeParts {
			grid[part.X][part.Y] = 1
		}

		return start, dest, grid
	} else {
		return nil, nil, nil
	}
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
	g.FoodPos.X = rand.Intn(width)
	g.FoodPos.Y = rand.Intn(height)

	if g.FoodPos.Y == 1 && g.FoodPos.X < 10 {
		g.UpdateFoodPos(width, height)
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
	snakeStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorWhite)

	i := 0

	for {
		longerSnake := false
		checkPath := false
		g.Screen.Clear()
		if checkCollision(g.snakeBody.Parts[len(g.snakeBody.Parts)-1:], g.FoodPos) {
			g.UpdateFoodPos(width, height)
			longerSnake = true
			checkPath = true
			g.Score++
		}

		if checkCollision(g.snakeBody.Parts[:len(g.snakeBody.Parts)-1], g.snakeBody.Parts[len(g.snakeBody.Parts)-1]) {
			break
		}

		g.snakeBody.Update(width, height, longerSnake)
		start, dest, grid := drawParts(g.Screen, g.snakeBody.Parts, g.FoodPos, snakeStyle, defStyle, g.Grid, checkPath)
		drawText(g.Screen, 1, 1, 8+len(strconv.Itoa(g.Score)), 1, "Score: "+strconv.Itoa(g.Score))
		if checkPath {
			res := Result{
				Start: *start,
				Dest:  *dest,
				Grid:  grid,
			}
			newPath := path.AStarSearch(start, dest, grid)
			if newPath != nil {
				panic(newPath)
			}
			i++
			if i >= 2 {
				panic(res)
			}
		}

		time.Sleep(40 * time.Millisecond)
		g.Screen.Show()
	}

	g.GameOver = true
	drawText(g.Screen, width/2-20, height/2, width/2+20, height/2, "Game Over, Score: "+strconv.Itoa(g.Score)+", Play Again? y/n")
	g.Screen.Show()
}

type Result struct {
	Start path.Node
	Dest  path.Node
	Grid  [][]int
}
