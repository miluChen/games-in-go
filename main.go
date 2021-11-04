package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type GameState int

const (
	Start GameState = iota
	Play
	TryExit
	Exit
	NOOP
)

var gameState GameState
var startScene, exitScene ChoiceScene

func initialize(win *pixelgl.Window) {
	gameState = Start
	startScene = *newChoiceScene(pixel.V(300, 500), 7, []string{"START", "EXIT"}, []GameState{Play, Exit})
	exitScene = *newChoiceScene(pixel.V(300, 500), 7, []string{"RESUME", "EXIT"}, []GameState{Play, Exit})
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Maze",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	initialize(win)

	for !win.Closed() && gameState != Exit {
		win.Clear(colornames.Black)
		if win.JustPressed(pixelgl.MouseButtonLeft) {
			win.SetTitle(fmt.Sprintf("%f %f", win.MousePosition().X, win.MousePosition().Y))
		}
		if win.JustPressed(pixelgl.KeyEscape) {
			if gameState != Start {
				gameState = TryExit
			}
		}
		state := NOOP
		switch gameState {
		case Start:
			startScene.draw(win)
			state = startScene.action(win)
		case Play:
		case TryExit:
			exitScene.draw(win)
			state = exitScene.action(win)
		default:
		}
		if state != NOOP {
			gameState = state
		}
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
