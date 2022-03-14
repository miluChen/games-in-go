package snake

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

/* ================ button definition ================ */
type Button interface {
	draw(*pixelgl.Window, bool)
	handle()
	disable()
	enable()
}

type RectButton struct {
	imd      *imdraw.IMDraw
	rect     pixel.Rect
	msg      string
	disabled bool
	handler  func()
}

func newRectButton(rect pixel.Rect, msg string, disabled bool, handler func()) *RectButton {
	imd := imdraw.New(nil)
	imd.Push(rect.Min)
	imd.Push(rect.Max)
	imd.Rectangle(2)

	return &RectButton{imd: imd, rect: rect, msg: msg, disabled: disabled, handler: handler}
}

// draw draws the button on win and highlights it if it's chosen
func (b *RectButton) draw(win *pixelgl.Window, highlight bool) {
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	// align text to the center
	txt := text.New(b.rect.Center(), atlas)
	txt.Dot.X -= txt.BoundsOf(b.msg).W() / 2
	txt.Dot.Y -= txt.BoundsOf(b.msg).H() / 4
	if b.disabled {
		txt.Color = colornames.Gray
	} else {
		// highlight the text if it's chosen
		if highlight {
			txt.Color = colornames.Blue
		} else {
			txt.Color = colornames.White
		}
	}
	fmt.Fprint(txt, b.msg)

	b.imd.Draw(win)
	txt.Draw(win, pixel.IM)
}

func (b *RectButton) handle() {
	if !b.disabled {
		b.handler()
	}
}

func (b *RectButton) disable() {
	b.disabled = true
}

func (b *RectButton) enable() {
	b.disabled = false
}

/* ================ button names ================ */
const (
	newGameButtonName     = "New Game"
	leaderBoardButtonName = "Leaderboard"
	optionsButtonName     = "Options"
	exitButtonName        = "Exit"
	resumeButtonName      = "Resume"
	restartButtonName     = "Restart"
	retryButtonName       = "Retry"
	backButtonName        = "Back"
	pausedButtonName      = "Paused"
	mainMenuButtonName    = "Main Menu"
)

/* ================ callbacks for buttons ================ */
func newGameHandler() {
	startGame()
}

func leaderboardHandler() {
	menuStack = append(menuStack, leaderboardMenu)
}

func optionsHandler() {
	menuStack = append(menuStack, optionsMenu)
}

func exitHandler() {
	gameState = Exit
}

func resumeHandler() {
	menuStack = nil
	currentScene.resume()
}

func restartHandler() {
	startGame()
}

func retryHandler() {
	startGame()
}

func backHandler() {
	menuStack[len(menuStack)-1].reset()
	menuStack = menuStack[0 : len(menuStack)-1]
}

func mainMenuHandler() {
	// user can not go back after you go to main menu, so we can clear the menu stack here
	clearMenuStack()
	menuStack = append(menuStack, mainMenu)
	// remove current scene
	currentScene = nil
}
