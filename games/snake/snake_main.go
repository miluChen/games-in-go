package snake

import (
	"fmt"
	"math"
	"math/rand"
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

const Unit = 10   // size of a square
const Width = 10  // width of the grid as in number of units
const Height = 10 // height of the grid as in number of units

type SnakeGame struct {
	alive bool        // whether the snake is still alive
	dir   Direction   // snake moving direction
	body  []pixel.Vec // the coordinates of the whole snake

	apple        pixel.Vec // position of the apple
	score        int       // game score, as in number apples eaten
	freq         int64     // the number of moves the snake can make per second
	action       Direction // actions for changing direction
	lastMoveTime time.Time // last timestamp the snake moved
}

var directions map[Direction][]float64

func init() {
	directions = make(map[Direction][]float64)
	directions[North] = []float64{0, 1}
	directions[East] = []float64{1, 0}
	directions[South] = []float64{0, -1}
	directions[West] = []float64{-1, 0}
}

// draw the snake and apple in window
func (s *SnakeGame) draw(win *pixelgl.Window) {
	// draw snake body and head
	imd := imdraw.New(nil)
	imd.Color = colornames.Limegreen
	for i := 0; i < len(s.body)-1; i++ {
		imd.Push(pixel.V(s.body[i].X*Unit, s.body[i].Y*Unit))
		imd.Push(pixel.V((s.body[i].X+1)*Unit, (s.body[i].Y+1)*Unit))
		imd.Rectangle(0)
	}
	imd.Color = colornames.Purple
	imd.Push(pixel.V(s.body[len(s.body)-1].X*Unit, s.body[len(s.body)-1].Y*Unit))
	imd.Push(pixel.V((s.body[len(s.body)-1].X+1)*Unit, (s.body[len(s.body)-1].Y+1)*Unit))
	imd.Rectangle(0)
	// draw apple
	imd.Color = colornames.Red
	imd.Push(pixel.V(s.apple.X*Unit, s.apple.Y*Unit))
	imd.Push(pixel.V((s.apple.X+1)*Unit, (s.apple.Y+1)*Unit))
	imd.Rectangle(0)

	win.Clear(colornames.Aliceblue)
	imd.Draw(win)
	win.Update()
}

// snake moves
func (s *SnakeGame) move() {
	if time.Since(s.lastMoveTime).Milliseconds() > time.Second.Milliseconds()/s.freq {
		// change direction if needed
		s.dir = changeDirection(s.dir, s.action)
		// advance head
		head := s.body[len(s.body)-1]
		nx, ny := head.X+directions[s.dir][0], head.Y+directions[s.dir][1]
		// check the snake is not out of bound
		if nx < 0 || nx >= Width || ny < 0 || ny >= Height {
			s.alive = false
			return
		}
		// check the snake is not colliding with itself
		for _, pos := range s.body {
			if pos.X == nx && pos.Y == ny {
				s.alive = false
				return
			}
		}
		s.body = append(s.body, pixel.V(nx, ny))
		// if apple is eaten, generate a new apple
		if nx != s.apple.X || ny != s.apple.Y {
			s.body = s.body[1:]
		} else {
			s.generateApple()
			s.score += 1
		}
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

// generate apple randomly
func (s *SnakeGame) generateApple() {
	for {
		x := rand.Intn(Width)
		y := rand.Intn(Height)
		// check collision
		hit := false
		for _, pos := range s.body {
			if int(pos.X) == x && int(pos.Y) == y {
				hit = true
				break
			}
		}
		if !hit {
			s.apple = pixel.V(float64(x), float64(y))
			break
		}
	}
}

// check whether the snake is in an invalid position
func (s *SnakeGame) dead(win *pixelgl.Window) bool {
	return !s.alive
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

	snakeGame := SnakeGame{
		alive:        true,
		dir:          East,
		body:         []pixel.Vec{pixel.V(0, 0)},
		action:       East,
		lastMoveTime: time.Now(),
		freq:         5,
	}
	snakeGame.generateApple()

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
				snakeGame.action = West
			} else if win.Pressed(pixelgl.KeyRight) {
				snakeGame.action = East
			} else if win.Pressed(pixelgl.KeyDown) {
				snakeGame.action = South
			} else if win.Pressed(pixelgl.KeyUp) {
				snakeGame.action = North
			}
			snakeGame.move()
			if snakeGame.dead(win) {
				fmt.Println("snake is dead")
			}
			snakeGame.draw(win)
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
