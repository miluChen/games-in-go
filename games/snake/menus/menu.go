package menus

import (
	"fmt"

	"github.com/faiface/pixel/pixelgl"
)

type Menu struct {
	buttonIndex int
	buttons     []Button
}

func newMenu() *Menu {
	return &Menu{buttonIndex: -1}
}

func (m *Menu) Draw(win *pixelgl.Window) {
	for i, button := range m.buttons {
		highlight := i == m.buttonIndex
		button.draw(win, highlight)
	}
}

func (m *Menu) addButton(button Button) {
	m.buttons = append(m.buttons, button)
}

func (m *Menu) HandleEvent(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyEnter) {
		if m.buttonIndex >= 0 && m.buttonIndex < len(m.buttons) {
			m.buttons[m.buttonIndex].handle()
		}
		return
	}
	if win.JustPressed(pixelgl.KeyDown) {
		m.buttonIndex = min(m.buttonIndex+1, len(m.buttons)-1)
		fmt.Println("hehe")
		fmt.Println(m.buttonIndex, len(m.buttons))
	} else if win.JustPressed(pixelgl.KeyUp) {
		m.buttonIndex = max(0, m.buttonIndex-1)
	}
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
