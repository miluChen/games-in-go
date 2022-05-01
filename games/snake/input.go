// input.go contains logic that lets user input its name

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

var numbers = []pixelgl.Button{
	pixelgl.Key0, pixelgl.Key1, pixelgl.Key2, pixelgl.Key3, pixelgl.Key4,
	pixelgl.Key5, pixelgl.Key6, pixelgl.Key7, pixelgl.Key8, pixelgl.Key9,
}

var chars = []pixelgl.Button{
	pixelgl.KeyA, pixelgl.KeyB, pixelgl.KeyC, pixelgl.KeyD, pixelgl.KeyE,
	pixelgl.KeyF, pixelgl.KeyG, pixelgl.KeyH, pixelgl.KeyI, pixelgl.KeyJ,
	pixelgl.KeyK, pixelgl.KeyL, pixelgl.KeyM, pixelgl.KeyN, pixelgl.KeyO,
	pixelgl.KeyP, pixelgl.KeyQ, pixelgl.KeyR, pixelgl.KeyS, pixelgl.KeyT,
	pixelgl.KeyU, pixelgl.KeyV, pixelgl.KeyW, pixelgl.KeyX, pixelgl.KeyY,
	pixelgl.KeyZ,
}

type InputBox struct {
	input    string
	curPos   int
	capsLock bool
	rect     pixel.Rect
	imd      *imdraw.IMDraw
}

func newInputBox(rect pixel.Rect) *InputBox {
	imd := imdraw.New(nil)
	imd.Push(rect.Min)
	imd.Push(rect.Max)
	imd.Rectangle(0)
	imd.Color = colornames.White

	return &InputBox{rect: rect, imd: imd}
}

func (ui *InputBox) reset() {
	ui.input = ""
	ui.curPos = 0
}

func (ui *InputBox) handle(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyCapsLock) {
		ui.capsLock = !ui.capsLock
	} else if win.JustPressed(pixelgl.KeyHome) {
		ui.curPos = 0
	} else if win.JustPressed(pixelgl.KeyEnd) {
		ui.curPos = len(ui.input)
	} else if win.JustPressed(pixelgl.KeyLeft) {
		ui.curPos = max(0, ui.curPos-1)
	} else if win.JustPressed(pixelgl.KeyRight) {
		ui.curPos = min(len(ui.input), ui.curPos+1)
	} else if win.JustPressed(pixelgl.KeyBackspace) || win.JustPressed(pixelgl.KeyDelete) {
		// delete the character before the cursor
		ui.input = ui.input[0 : len(ui.input)-1]
	} else {
		// read the input if it's valid, i.e. if it's space, a number or a letter
		input := ui.readCharInput(win)
		ui.input += input
		ui.curPos += len(input)
	}
}

func (ui *InputBox) draw(win *pixelgl.Window, activated bool) {
	// draw rectangle box
	ui.imd.Draw(win)
	// draw input text
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	txt := text.New(pixel.V(ui.rect.Min.X+1, ui.rect.Min.Y+ui.rect.H()/2), atlas)
	txt.Dot.Y -= txt.BoundsOf(ui.input).H() / 4
	txt.Color = colornames.Green
	fmt.Fprint(txt, ui.input)
	txt.Draw(win, pixel.IM)
	// draw cursor if input box is activated
	if activated {
		cursor := imdraw.New(nil)
		cursor.Color = colornames.Gray
		cursor.Push(pixel.V(txt.Dot.X, ui.rect.Min.Y))
		cursor.Push(pixel.V(txt.Dot.X+1, ui.rect.Max.Y))
		cursor.Rectangle(0)
		cursor.Draw(win)
	}
}

func (ui *InputBox) readCharInput(win *pixelgl.Window) string {
	if win.JustPressed(pixelgl.KeySpace) {
		return string(rune(int(pixelgl.KeySpace)))
	}
	for _, num := range numbers {
		if win.JustPressed(num) {
			return string(rune(int(num)))
		}
	}

	offset := 32
	if ui.capsLock != (win.Pressed(pixelgl.KeyLeftShift) || win.Pressed(pixelgl.KeyRightShift)) {
		offset = 0
	}
	for _, c := range chars {
		if win.JustPressed(c) {
			return string(rune(int(c) + offset))
		}
	}
	return ""
}
