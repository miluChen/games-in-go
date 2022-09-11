package snake

import (
	"fmt"
	"strings"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/miluchen/games-in-go/games/snake/db"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

/* ================ button definition ================ */
type Button interface {
	draw(*pixelgl.Window, bool)
	handle()
	contains(pixel.Vec) bool
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
		txt.Color = colornames.Silver
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

func (b *RectButton) contains(cursor pixel.Vec) bool {
	return b.rect.Contains(cursor)
}

func (b *RectButton) handle() {
	if !b.disabled {
		b.handler()
	}
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
	playAgainButtonName   = "Play Again"
	backButtonName        = "Back"
	pausedButtonName      = "Paused"
	mainMenuButtonName    = "Main Menu"
	cancelButtonName      = "Cancel"
	confirmButtonName     = "Confirm"
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

func playAgainHandler() {
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

func cancelHandler() {
	// reset input box
	for _, inputBox := range inputNameMenu.inputBoxes {
		inputBox.reset()
	}
	// pop menu
	menuStack = menuStack[0 : len(menuStack)-1]
}

func confirmHandler() {
	// write data into database
	name := strings.TrimSpace(inputNameMenu.inputBoxes[0].input)
	if len(name) > 0 {
		db.Insert(name)
	}
	// reset input box
	inputNameMenu.inputBoxes[0].reset()
	// pop menu
	menuStack = menuStack[0 : len(menuStack)-1]
}
