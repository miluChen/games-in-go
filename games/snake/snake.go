package snake

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

type snakeState int

const (
	Idle snakeState = iota
	Moving
)

const (
	Unit     = 20 // size of a square
	Width    = 15 // width of the grid as in number of units
	Height   = 15 // height of the grid as in number of units
	MaxLevel = 10 // max game level
	AppleCnt = 3  // number of apples to advance to next level
)

var unitV = pixel.V(Unit, Unit)

type SnakeGame struct {
	alive bool        // whether the snake is still alive
	dir   Direction   // snake moving direction
	body  []pixel.Vec // the coordinates of the whole snake
	state snakeState  // state indicates whether the snake should is moving

	apple          pixel.Vec // position of the apple
	score          int       // game score, as in number apples eaten
	level          int       // game level
	freq           int64     // the number of moves the snake can make per second
	action         Direction // action for changing direction
	repeatedAction bool      // whether action is repeatedly pressed, if so, snake moves at max speed
	lastMoveTime   time.Time // last timestamp the snake moved
}

var directions = map[Direction][]float64{
	North: {0, 1},
	East:  {1, 0},
	South: {0, -1},
	West:  {-1, 0},
}

// mapping from game level to freq (index 0 is not used)
var frequencies = []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

func newSnakeGame() *SnakeGame {
	snakeGame := &SnakeGame{
		alive:        true,
		action:       East,
		lastMoveTime: time.Now(),
	}
	snakeGame.setLevel(1)
	snakeGame.resetSnake()
	snakeGame.generateApple()
	return snakeGame
}

// draw the snake and apple in window
func (s *SnakeGame) draw(win *pixelgl.Window) {
	// draw level txt and display score
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	txt := text.New(pixel.V(100, 700), atlas)
	txt.Color = colornames.Black
	fmt.Fprintf(txt, "Level %d: %d", s.level, s.score)
	// compute offset, which is the lower left boundary of the allowed area
	offsetX := (win.Bounds().W() - (Width+1)*Unit) / 2
	offsetY := (win.Bounds().H() - txt.Bounds().H() - (Height)*Unit) / 2
	offset := pixel.V(offsetX, offsetY)
	// position level txt in top center
	txt.Draw(win, pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center()).Add(pixel.V(0, win.Bounds().H()/2-txt.Bounds().H()/2))))
	// draw wall
	imd := imdraw.New(nil)
	imd.Color = colornames.Coral
	for i := -1; i < Width+1; i++ {
		v := pixel.V(float64(i*Unit), -Unit)
		imd.Push(v.Add(offset))
		imd.Push(v.Add(unitV).Add(offset))
		imd.Rectangle(0)
		v = pixel.V(float64(i*Unit), float64((Height)*Unit))
		imd.Push(v.Add(offset))
		imd.Push(v.Add(unitV).Add(offset))
		imd.Rectangle(0)
	}
	for i := 0; i < Height; i++ {
		v := pixel.V(-Unit, float64(i*Unit))
		imd.Push(v.Add(offset))
		imd.Push(v.Add(unitV).Add(offset))
		imd.Rectangle(0)
		v = pixel.V(float64((Width)*Unit), float64(i*Unit))
		imd.Push(v.Add(offset))
		imd.Push(v.Add(unitV).Add(offset))
		imd.Rectangle(0)
	}
	// draw snake body and head
	imd.Color = colornames.Limegreen
	for i := 0; i < len(s.body)-1; i++ {
		imd.Push(pixel.V(s.body[i].X*Unit, s.body[i].Y*Unit).Add(offset))
		imd.Push(pixel.V((s.body[i].X+1)*Unit, (s.body[i].Y+1)*Unit).Add(offset))
		imd.Rectangle(0)
	}
	imd.Color = colornames.Purple
	imd.Push(pixel.V(s.body[len(s.body)-1].X*Unit, s.body[len(s.body)-1].Y*Unit).Add(offset))
	imd.Push(pixel.V((s.body[len(s.body)-1].X+1)*Unit, (s.body[len(s.body)-1].Y+1)*Unit).Add(offset))
	imd.Rectangle(0)
	// draw apple
	imd.Color = colornames.Red
	imd.Push(pixel.V(s.apple.X*Unit, s.apple.Y*Unit).Add(offset))
	imd.Push(pixel.V((s.apple.X+1)*Unit, (s.apple.Y+1)*Unit).Add(offset))
	imd.Rectangle(0)

	imd.Draw(win)
}

// snake moves
func (s *SnakeGame) move() {
	if s.state == Idle {
		return
	}
	freq := s.freq
	if s.action == s.dir && s.repeatedAction {
		// if a key is held, win.Repeated will return true, false, false, false, false, true, ...
		// the max frequency is multipled by 5 to accommodate this
		freq = frequencies[len(frequencies)-1] * 5
	}
	if time.Since(s.lastMoveTime).Milliseconds() > time.Second.Milliseconds()/freq {
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
	if s.passLevel() {
		s.advanceLevel()
	}
}

// check whether it should advance to next level
func (s *SnakeGame) passLevel() bool {
	return s.score == s.level*AppleCnt
}

// advance to next level
func (s *SnakeGame) advanceLevel() {
	// TODO: check whether player wins

	s.resetSnake()
	s.setLevel(s.level + 1)
}

// reset state of snake
func (s *SnakeGame) resetSnake() {
	s.dir = East
	s.state = Idle // snake starts as idle, waiting from command
	s.body = []pixel.Vec{pixel.V(0, Height/2), pixel.V(1, Height/2), pixel.V(2, Height/2)}
}

// set game level
func (s *SnakeGame) setLevel(level int) {
	s.level = level
	s.freq = frequencies[s.level]
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
