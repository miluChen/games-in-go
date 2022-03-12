package menus

import "github.com/faiface/pixel"

var MainMenu *Menu
var LeaderboardMenu *Menu
var OptionsMenu *Menu
var ExitMenu *Menu
var PauseMenu *Menu
var GameOverMenu *Menu

func InitMenus() {
	MainMenu = createMainMenu()
	LeaderboardMenu = createLeaderBoardMenu()
	OptionsMenu = createOptionsMenu()
	ExitMenu = createExitMenu()
	PauseMenu = createPauseMenu()
	GameOverMenu = createGameOverMenu()
}

func createMainMenu() *Menu {
	menu := newMenu()
	// add buttons for main menu
	rect := pixel.Rect{Min: pixel.V(300, 500), Max: pixel.V(400, 530)}
	menu.addButton(newRectButton(rect, newGameButtonName, newGameHandler))
	rect = pixel.Rect{Min: pixel.V(300, 460), Max: pixel.V(400, 490)}
	menu.addButton(newRectButton(rect, leaderBoardButtonName, leaderboardHandler))
	rect = pixel.Rect{Min: pixel.V(300, 420), Max: pixel.V(400, 450)}
	menu.addButton(newRectButton(rect, optionsButtonName, optionsHandler))
	rect = pixel.Rect{Min: pixel.V(300, 380), Max: pixel.V(400, 410)}
	menu.addButton(newRectButton(rect, exitButtonName, exitHandler))
	return menu
}

func createLeaderBoardMenu() *Menu {
	return nil
}

func createOptionsMenu() *Menu {
	return nil
}

func createExitMenu() *Menu {
	return nil
}

func createPauseMenu() *Menu {
	return nil
}

func createGameOverMenu() *Menu {
	return nil
}
