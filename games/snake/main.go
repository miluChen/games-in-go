package snake

import (
	"log"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/miluChen/games-in-go/games/snake/db"
	"golang.org/x/image/colornames"
)

type GameState int

const (
	InGame GameState = iota
	Exit
)

var gameState GameState
var currentScene *Scene

func initialize(win *pixelgl.Window) {
	gameState = InGame
	currentScene = nil
	// application starts, push main menu to menuStack
	initMenus(win)
	// open DB
	err := db.Open()
	if err != nil {
		log.Printf("open db failed: %v\n", err)
		return
	}
}

func close() {
	if err := db.Close(); err != nil {
		log.Printf("close db failed: %v\n", err)
	}
}

func startGame() {
	clearMenuStack()
	// initialize snake game
	currentScene = &Scene{active: true, snakeGame: newSnakeGame()}
}

func run() {
	// initialize window
	cfg := pixelgl.WindowConfig{
		Title:  "snake",
		Bounds: pixel.R(0, 0, 500, 400),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	initialize(win)
	// game loop
	for !win.Closed() && gameState != Exit {
		win.Clear(colornames.Black)
		// if there is current scene, render it
		if currentScene != nil {
			currentScene.update(win)
		}
		// update menu if needed
		if len(menuStack) > 0 {
			menuStack[len(menuStack)-1].update(win)
		}

		win.Update()
	}
	close()
}

func Run() {
	pixelgl.Run(run)
}
