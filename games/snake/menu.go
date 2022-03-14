package snake

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

/* ========== menu definition ========== */
type Menu struct {
	buttonIndex int
	buttons     []Button
}

func newMenu() *Menu {
	return &Menu{buttonIndex: 0}
}

func (m *Menu) draw(win *pixelgl.Window) {
	for i, button := range m.buttons {
		highlight := i == m.buttonIndex
		button.draw(win, highlight)
	}
}

func (m *Menu) addButton(button Button) {
	m.buttons = append(m.buttons, button)
}

// handleEvent handles user input, it should be called before Draw
func (m *Menu) handleEvent(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.KeyEnter) {
		if m.buttonIndex >= 0 && m.buttonIndex < len(m.buttons) {
			m.buttons[m.buttonIndex].handle()
		}
		return
	}
	if win.JustPressed(pixelgl.KeyDown) {
		m.buttonIndex = min(m.buttonIndex+1, len(m.buttons)-1)
	} else if win.JustPressed(pixelgl.KeyUp) {
		m.buttonIndex = max(0, m.buttonIndex-1)
	}
}

func (m *Menu) update(win *pixelgl.Window) {
	m.handleEvent(win)
	m.draw(win)
}

func (m *Menu) reset() {
	m.buttonIndex = 0
}

/* ========== menu handle functions ========== */

var mainMenu *Menu
var leaderboardMenu *Menu
var optionsMenu *Menu
var pauseMenu *Menu
var gameOverMenu *Menu

var menuStack []*Menu

func initMenus() {
	mainMenu = createMainMenu()
	leaderboardMenu = createLeaderBoardMenu()
	optionsMenu = createOptionsMenu()
	pauseMenu = createPauseMenu()
	gameOverMenu = createGameOverMenu()

	menuStack = append(menuStack, mainMenu)
}

func createMainMenu() *Menu {
	menu := newMenu()
	// add buttons for main menu
	rect := pixel.Rect{Min: pixel.V(300, 500), Max: pixel.V(400, 530)}
	menu.addButton(newRectButton(rect, newGameButtonName, false, newGameHandler))
	rect = pixel.Rect{Min: pixel.V(300, 460), Max: pixel.V(400, 490)}
	menu.addButton(newRectButton(rect, leaderBoardButtonName, false, leaderboardHandler))
	rect = pixel.Rect{Min: pixel.V(300, 420), Max: pixel.V(400, 450)}
	menu.addButton(newRectButton(rect, optionsButtonName, false, optionsHandler))
	rect = pixel.Rect{Min: pixel.V(300, 380), Max: pixel.V(400, 410)}
	menu.addButton(newRectButton(rect, exitButtonName, false, exitHandler))
	return menu
}

func createLeaderBoardMenu() *Menu {
	menu := newMenu()
	// add buttons for leaderboard menu
	rect := pixel.Rect{Min: pixel.V(300, 340), Max: pixel.V(400, 370)}
	menu.addButton(newRectButton(rect, backButtonName, false, backHandler))
	return menu
}

func createOptionsMenu() *Menu {
	menu := newMenu()
	// add buttons for options menu
	rect := pixel.Rect{Min: pixel.V(300, 340), Max: pixel.V(400, 370)}
	menu.addButton(newRectButton(rect, backButtonName, false, backHandler))
	return menu
}

func createPauseMenu() *Menu {
	menu := newMenu()
	// add buttons for pause menu
	rect := pixel.Rect{Min: pixel.V(300, 540), Max: pixel.V(400, 570)}
	menu.addButton(newRectButton(rect, pausedButtonName, true, nil))
	rect = pixel.Rect{Min: pixel.V(300, 500), Max: pixel.V(400, 530)}
	menu.addButton(newRectButton(rect, resumeButtonName, false, resumeHandler))
	rect = pixel.Rect{Min: pixel.V(300, 460), Max: pixel.V(400, 490)}
	menu.addButton(newRectButton(rect, restartButtonName, false, restartHandler))
	rect = pixel.Rect{Min: pixel.V(300, 420), Max: pixel.V(400, 450)}
	menu.addButton(newRectButton(rect, optionsButtonName, false, optionsHandler))
	rect = pixel.Rect{Min: pixel.V(300, 380), Max: pixel.V(400, 410)}
	menu.addButton(newRectButton(rect, exitButtonName, false, exitHandler))
	return menu
}

func createGameOverMenu() *Menu {
	menu := newMenu()
	// add buttons for pause menu
	rect := pixel.Rect{Min: pixel.V(300, 500), Max: pixel.V(400, 530)}
	menu.addButton(newRectButton(rect, retryButtonName, false, retryHandler))
	rect = pixel.Rect{Min: pixel.V(300, 460), Max: pixel.V(400, 490)}
	menu.addButton(newRectButton(rect, mainMenuButtonName, false, mainMenuHandler))
	rect = pixel.Rect{Min: pixel.V(300, 380), Max: pixel.V(400, 410)}
	menu.addButton(newRectButton(rect, exitButtonName, false, exitHandler))
	return menu
}

func clearMenuStack() {
	for _, menu := range menuStack {
		menu.reset()
	}
	menuStack = nil
}
