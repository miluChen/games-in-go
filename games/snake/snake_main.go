package snake

import (
	"fmt"
	"math"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type GameState int
type Direction int

const (
	Start GameState = iota
	Play
	TryExit
	Exit
	NOOP
)

const (
	North Direction = iota
	East
	South
	West
)

const Unit = 10 // size of a square

type Snake struct {
	body         []pixel.Vec // the coordinates of the whole snake
	dir          Direction   // move direction
	action       Direction   // actions for changing direction
	lastMoveTime time.Time   // last timestamp the snake moved
	freq         int64       // the number of moves the snake can make per second
}

var directions map[Direction][]float64

func init() {
	directions = make(map[Direction][]float64)
	directions[North] = []float64{0, Unit}
	directions[East] = []float64{Unit, 0}
	directions[South] = []float64{0, -Unit}
	directions[West] = []float64{-Unit, 0}
}

// draw the snake in window
func (s *Snake) draw(win *pixelgl.Window) {
	imd := imdraw.New(nil)
	imd.Color = colornames.Limegreen
	for _, point := range s.body {
		imd.Push(point)
		imd.Push(pixel.V(point.X+Unit, point.Y+Unit))
		imd.Rectangle(0)
	}
	win.Clear(colornames.Aliceblue)
	imd.Draw(win)
	win.Update()
}

// snake moves
func (s *Snake) move() {
	if time.Since(s.lastMoveTime).Milliseconds() > time.Second.Milliseconds()/s.freq {
		// change direction if needed
		s.dir = changeDirection(s.dir, s.action)
		// advance head
		head := s.body[len(s.body)-1]
		s.body = append(s.body, pixel.V(head.X+directions[s.dir][0], head.Y+directions[s.dir][1]))
		// remove tail
		s.body = s.body[1:]
		// update last move timestamp
		s.lastMoveTime = time.Now()
	}
}

// change direction
func changeDirection(dir Direction, action Direction) Direction {
	if dir == action {
		return dir
	}
	if math.Abs(float64(dir)-float64(action)) == 2 {
		return dir
	}
	return action
}

// check whether the snake is in an invalid position
func (s *Snake) dead(win *pixelgl.Window) bool {
	return false
}

var gameState GameState
var startScene, exitScene ChoiceScene

func initialize(win *pixelgl.Window) {
	gameState = Start
	startScene = *newChoiceScene(pixel.V(300, 500), 7, []string{"START", "EXIT"}, []GameState{Play, Exit})
	exitScene = *newChoiceScene(pixel.V(300, 500), 7, []string{"RESUME", "EXIT"}, []GameState{Play, Exit})
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "snake",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	initialize(win)

	snake := Snake{
		body:         []pixel.Vec{pixel.V(0, 0)},
		dir:          East,
		action:       East,
		lastMoveTime: time.Now(),
		freq:         5,
	}
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
			if win.Pressed(pixelgl.KeyLeft) {
				snake.action = West
			} else if win.Pressed(pixelgl.KeyRight) {
				snake.action = East
			} else if win.Pressed(pixelgl.KeyDown) {
				snake.action = South
			} else if win.Pressed(pixelgl.KeyUp) {
				snake.action = North
			}
			snake.move()
			if snake.dead(win) {
				fmt.Println("snake is dead")
			}
			snake.draw(win)
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

func Run() {
	pixelgl.Run(run)
}
