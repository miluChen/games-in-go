package main

import (
	"flag"
	"fmt"

	"github.com/miluChen/games-in-go/games/snake"
)

const (
	snakeGame = "snake"
)

var game = flag.String("game", "", fmt.Sprintf("game: %s", snakeGame))

func main() {
	flag.Parse()
	switch *game {
	case snakeGame:
		snake.Run()
	default:
		flag.Usage()
	}
}
