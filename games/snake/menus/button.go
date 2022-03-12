package menus

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type Button interface {
	draw(*pixelgl.Window, bool)
	handle()
}

type RectButton struct {
	imd     *imdraw.IMDraw
	rect    pixel.Rect
	msg     string
	handler func()
}

func newRectButton(rect pixel.Rect, msg string, handler func()) *RectButton {
	imd := imdraw.New(nil)
	imd.Push(rect.Min)
	imd.Push(rect.Max)
	imd.Rectangle(2)

	return &RectButton{imd: imd, rect: rect, msg: msg, handler: handler}
}

// draw draws the button on win and highlights it if it's chosen
func (b *RectButton) draw(win *pixelgl.Window, highlight bool) {
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	// align text to the center
	txt := text.New(b.rect.Center(), atlas)
	txt.Dot.X -= txt.BoundsOf(b.msg).W() / 2
	txt.Dot.Y -= txt.BoundsOf(b.msg).H() / 4
	// highlight the text if it's chosen
	if highlight {
		txt.Color = colornames.White
	} else {
		txt.Color = colornames.Grey
	}
	fmt.Fprint(txt, b.msg)

	txt.Draw(win, pixel.IM)
	b.imd.Draw(win)
}

func (b *RectButton) handle() {
	b.handler()
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
)

/* ================ callbacks for buttons ================ */
func newGameHandler() {
	fmt.Println("new game handler")
	
}

func leaderboardHandler() {
	fmt.Println("leaderboard handler")
}

func optionsHandler() {
	fmt.Println("option handler")
}

func exitHandler() {
	fmt.Println("exit handler")
}

func resumeHandler() {
	fmt.Println("resume handler")
}

func restartHandler() {
	fmt.Println("restart handler")
}

func retryHandler() {
	fmt.Println("retry handler")
}

func backHandler() {
	fmt.Println("back handler")
}
