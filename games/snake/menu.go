package snake

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/miluChen/games-in-go/games/snake/db"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

/* ========== menu definition ========== */
type Menu struct {
	buttonIndex  int
	buttons      []Button
	texts        []*text.Text
	textMatrices []pixel.Matrix
	inputBoxes   []*InputBox

	isLeaderBoard bool // leader board menu needs to read from DB every time
}

func newMenu() *Menu {
	return &Menu{buttonIndex: 0}
}

func (m *Menu) draw(win *pixelgl.Window) {
	win.Clear(colornames.Gray)
	cursor := win.MousePosition()
	for _, button := range m.buttons {
		highlight := button.contains(cursor)
		button.draw(win, highlight)
	}
	for i, text := range m.texts {
		text.Draw(win, m.textMatrices[i])
	}
	for _, inputBox := range m.inputBoxes {
		inputBox.draw(win)
	}
}

func (m *Menu) addButton(button Button) {
	m.buttons = append(m.buttons, button)
}

func (m *Menu) addText(text *text.Text, matrix pixel.Matrix) {
	m.texts = append(m.texts, text)
	m.textMatrices = append(m.textMatrices, matrix)
}

func (m *Menu) setInputBox(inputBox *InputBox) {
	m.inputBoxes = append(m.inputBoxes, inputBox)
}

// handleEvent handles user input, it should be called before Draw.
func (m *Menu) handleEvent(win *pixelgl.Window) {
	if win.JustPressed(pixelgl.MouseButtonLeft) {
		cursor := win.MousePosition()
		// check whether button is clicked
		for _, button := range m.buttons {
			if button.contains(cursor) {
				button.handle()
				return
			}
		}
		// check whether input box is chosen
		index := -1
		for i, inputBox := range m.inputBoxes {
			if inputBox.rect.Contains(cursor) {
				index = i
				break
			}
		}
		if index != -1 {
			for i, inputBox := range m.inputBoxes {
				if i == index {
					inputBox.activate()
				} else {
					inputBox.deactivate()
				}
			}
		}
	} else {
		for _, inputBox := range m.inputBoxes {
			if inputBox.activated {
				inputBox.handle(win)
				return
			}
		}
	}
}

func (m *Menu) update(win *pixelgl.Window) {
	if m.isLeaderBoard {
		generateLeaderBoard(win, m)
	}
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
var winMenu *Menu
var inputNameMenu *Menu

var menuStack []*Menu

func initMenus(win *pixelgl.Window) {
	mainMenu = createMainMenu()
	leaderboardMenu = createLeaderBoardMenu(win)
	optionsMenu = createOptionsMenu()
	pauseMenu = createPauseMenu()
	gameOverMenu = createGameOverMenu()
	winMenu = createWinMenu(win)
	inputNameMenu = createInputNameMenu(win)

	menuStack = append(menuStack, mainMenu)
}

func createMainMenu() *Menu {
	menu := newMenu()
	// add buttons for main menu
	rect := pixel.Rect{Min: pixel.V(200, 350), Max: pixel.V(300, 380)}
	menu.addButton(newRectButton(rect, newGameButtonName, false, newGameHandler))
	rect = pixel.Rect{Min: pixel.V(200, 310), Max: pixel.V(300, 340)}
	menu.addButton(newRectButton(rect, leaderBoardButtonName, false, leaderboardHandler))
	rect = pixel.Rect{Min: pixel.V(200, 270), Max: pixel.V(300, 300)}
	menu.addButton(newRectButton(rect, optionsButtonName, false, optionsHandler))
	rect = pixel.Rect{Min: pixel.V(200, 230), Max: pixel.V(300, 260)}
	menu.addButton(newRectButton(rect, exitButtonName, false, exitHandler))
	return menu
}

func createLeaderBoardMenu(win *pixelgl.Window) *Menu {
	menu := newMenu()
	menu.isLeaderBoard = true
	generateLeaderBoard(win, menu)
	// add buttons for leaderboard menu
	rect := pixel.Rect{Min: pixel.V(200, 190), Max: pixel.V(300, 220)}
	menu.addButton(newRectButton(rect, backButtonName, false, backHandler))
	return menu
}

func generateLeaderBoard(win *pixelgl.Window, menu *Menu) {
	// read from db and draw leaderboard
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	txt := text.New(pixel.V(100, 700), atlas)
	txt.Color = colornames.Green
	fmt.Fprintln(txt, "Leaderboard")
	matrix := pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center()).Add(pixel.V(0, win.Bounds().H()/2-txt.Bounds().H()/2)))

	menu.texts = nil
	menu.textMatrices = nil
	menu.addText(txt, matrix)

	names, err := db.Read()
	if err != nil {
		txt.Color = colornames.Red
		fmt.Fprintf(txt, "err: %s\n", err.Error())
	} else {
		txt.Color = colornames.Greenyellow
		for i, name := range names {
			fmt.Fprintf(txt, "%d\t%s\n", i+1, name)
		}
	}
}

func createOptionsMenu() *Menu {
	menu := newMenu()
	// add buttons for options menu
	rect := pixel.Rect{Min: pixel.V(200, 230), Max: pixel.V(300, 260)}
	menu.addButton(newRectButton(rect, backButtonName, false, backHandler))
	return menu
}

func createPauseMenu() *Menu {
	menu := newMenu()
	// add buttons for pause menu
	rect := pixel.Rect{Min: pixel.V(200, 350), Max: pixel.V(300, 380)}
	menu.addButton(newRectButton(rect, pausedButtonName, true, nil))
	rect = pixel.Rect{Min: pixel.V(200, 310), Max: pixel.V(300, 340)}
	menu.addButton(newRectButton(rect, resumeButtonName, false, resumeHandler))
	rect = pixel.Rect{Min: pixel.V(200, 270), Max: pixel.V(300, 300)}
	menu.addButton(newRectButton(rect, restartButtonName, false, restartHandler))
	rect = pixel.Rect{Min: pixel.V(200, 230), Max: pixel.V(300, 260)}
	menu.addButton(newRectButton(rect, optionsButtonName, false, optionsHandler))
	rect = pixel.Rect{Min: pixel.V(200, 190), Max: pixel.V(300, 220)}
	menu.addButton(newRectButton(rect, exitButtonName, false, exitHandler))
	return menu
}

func createGameOverMenu() *Menu {
	menu := newMenu()
	// add buttons for pause menu
	rect := pixel.Rect{Min: pixel.V(200, 350), Max: pixel.V(300, 380)}
	menu.addButton(newRectButton(rect, retryButtonName, false, retryHandler))
	rect = pixel.Rect{Min: pixel.V(200, 310), Max: pixel.V(300, 340)}
	menu.addButton(newRectButton(rect, mainMenuButtonName, false, mainMenuHandler))
	rect = pixel.Rect{Min: pixel.V(200, 230), Max: pixel.V(300, 260)}
	menu.addButton(newRectButton(rect, exitButtonName, false, exitHandler))
	return menu
}

func createWinMenu(win *pixelgl.Window) *Menu {
	menu := newMenu()
	// add win text
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	txt := text.New(pixel.V(100, 700), atlas)
	txt.Color = colornames.Red
	fmt.Fprint(txt, "You Win!")
	matrix := pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center()).Add(pixel.V(0, win.Bounds().H()/2-txt.Bounds().H()/2)))
	menu.addText(txt, matrix)
	// add buttons for win menu
	rect := pixel.Rect{Min: pixel.V(200, 310), Max: pixel.V(300, 340)}
	menu.addButton(newRectButton(rect, playAgainButtonName, false, playAgainHandler))
	rect = pixel.Rect{Min: pixel.V(200, 270), Max: pixel.V(300, 300)}
	menu.addButton(newRectButton(rect, mainMenuButtonName, false, mainMenuHandler))
	return menu
}

func createInputNameMenu(win *pixelgl.Window) *Menu {
	menu := newMenu()
	// add win text
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	txt := text.New(pixel.V(100, 700), atlas)
	txt.Color = colornames.Red
	fmt.Fprintf(txt, "You Win! Your Name:\n")
	matrix := pixel.IM.Moved(win.Bounds().Center().Sub(txt.Bounds().Center()).Add(pixel.V(0, win.Bounds().H()/2-txt.Bounds().H()/2)))
	menu.addText(txt, matrix)
	// add input box
	rect := pixel.Rect{Min: pixel.V(200, 350), Max: pixel.V(300, 380)}
	menu.setInputBox(newInputBox(rect))
	// add other buttons
	rect = pixel.Rect{Min: pixel.V(200, 270), Max: pixel.V(300, 300)}
	menu.addButton(newRectButton(rect, cancelButtonName, false, cancelHandler))
	rect = pixel.Rect{Min: pixel.V(200, 230), Max: pixel.V(300, 260)}
	menu.addButton(newRectButton(rect, confirmButtonName, false, confirmHandler))

	// set button index to -1
	menu.buttonIndex = -1

	return menu
}

func clearMenuStack() {
	for _, menu := range menuStack {
		menu.reset()
	}
	menuStack = nil
}
