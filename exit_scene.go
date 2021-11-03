package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font/basicfont"
)

type ExitScene struct {
	txt     *text.Text
	rects   []pixel.Rect
	actions []GameState
	scale   float64
	imds    []*imdraw.IMDraw
}

func newExitScene(orig pixel.Vec, scale float64) *ExitScene {
	scene := ExitScene{actions: []GameState{Play, Exit}, scale: scale}
	msgs := []string{"RESUME", "EXIT"}
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)

	scene.txt = text.New(orig, atlas)
	scene.txt.LineHeight = 3 * scale

	for _, msg := range msgs {
		rect := scene.txt.BoundsOf(msg)
		rectResized := rect.Resized(orig, pixel.V(rect.W(), rect.H()).Scaled(scale))
		scene.rects = append(scene.rects, rectResized)
		imd := imdraw.New(nil)
		imd.Push(rectResized.Min)
		imd.Push(rectResized.Max)
		imd.Rectangle(2)
		scene.imds = append(scene.imds, imd)

		fmt.Fprintln(scene.txt, msg)
	}
	return &scene
}

func (scene *ExitScene) draw(win *pixelgl.Window) {
	scene.txt.Draw(win, pixel.IM.Scaled(scene.txt.Orig, scene.scale))
	// draw bounding box for "start" and "exit" messages
	for _, imd := range scene.imds {
		imd.Draw(win)
	}
}

func (scene *ExitScene) action(win *pixelgl.Window) GameState {
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		// retrieve cursor position and check which action to take
		cursor := win.MousePosition()
		for i, rect := range scene.rects {
			if rect.Contains(cursor) {
				return scene.actions[i]
			}
		}
	}
	return NOOP
}
