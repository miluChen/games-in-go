package snake

import (
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Scene struct {
	active    bool
	snakeGame *SnakeGame
}

func (s *Scene) draw(win *pixelgl.Window) {
	win.Clear(colornames.Aliceblue)
	s.snakeGame.draw(win)
}

func (s *Scene) update(win *pixelgl.Window) {
	if !s.active {
		return
	}
	// check whether to pause the game
	if win.JustPressed(pixelgl.KeyEscape) {
		s.active = false
		menuStack = append(menuStack, pauseMenu)
		return
	}
	// check whether player has won
	if s.snakeGame.won {
		s.active = false
		menuStack = append(menuStack, winMenu)
		menuStack = append(menuStack, inputNameMenu)
		return
	}

	if win.JustPressed(pixelgl.KeyLeft) {
		s.snakeGame.action = West
		s.snakeGame.state = Moving
	} else if win.JustPressed(pixelgl.KeyRight) {
		s.snakeGame.action = East
		s.snakeGame.state = Moving
	} else if win.JustPressed(pixelgl.KeyDown) {
		s.snakeGame.action = South
		s.snakeGame.state = Moving
	} else if win.JustPressed(pixelgl.KeyUp) {
		s.snakeGame.action = North
		s.snakeGame.state = Moving
	}

	s.snakeGame.repeatedAction = false
	if win.Repeated(pixelgl.KeyLeft) {
		s.snakeGame.repeatedAction = true
	} else if win.Repeated(pixelgl.KeyRight) {
		s.snakeGame.repeatedAction = true
	} else if win.Repeated(pixelgl.KeyDown) {
		s.snakeGame.repeatedAction = true
	} else if win.Repeated(pixelgl.KeyUp) {
		s.snakeGame.repeatedAction = true
	}

	s.snakeGame.move()
	if s.snakeGame.dead(win) {
		s.active = false
		menuStack = append(menuStack, gameOverMenu)
		return
	}
	s.draw(win)
}

func (s *Scene) resume() {
	s.active = true
}
