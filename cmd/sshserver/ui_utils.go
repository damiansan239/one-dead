package main

import (
	"fmt"
	"one_dead/internal/game"
	"time"

	"github.com/gdamore/tcell/v2"
	_ "github.com/gdamore/tcell/v2/encoding"
)

var logo []string = []string{
	"  .oooooo.",
	" d8P'  `Y8b",
	"888      888 ooo. .oo.    .ooooo.",
	"888      888 `888P'Y88b  d88' `88b",
	"888      888  888   888  888ooo888",
	"`88b    d88'  888   888  888    .o",
	" `Y8bood8P'  o888o o888o `Y8bod8P'",
	"",
	"oooooooooo.                             .o8",
	"`888'   `Y8b                           '888",
	" 888      888  .ooooo.   .oooo.    .oooo888",
	" 888      888 d88' `88b `P  )88b  d88' `888",
	" 888      888 888ooo888  .oP'888  888   888",
	" 888     d88' 888    .o d8(  888  888   888",
	"o888bood8P'   `Y8bod8P' `Y888''8o `Y8bod88P'",
}

func (ui *ChatUI) drawTopBar() {
	width, _ := ui.screen.Size()
	topBarStyle := ui.currentStyle.Background(tcell.ColorNavy).Foreground(tcell.ColorWhite)

	// Clear the top bar
	for x := 0; x < width; x++ {
		ui.screen.SetContent(x, 0, ' ', nil, topBarStyle)
	}

	tries := 4
	currentDuration := time.Now().UTC().Format("15:04:05")

	topBarText := fmt.Sprintf(
		"One Dead: A strategic guessing game. Current Tries: %d. Current Duration: %s. Play at https://one-dead.web.app",
		tries,
		currentDuration,
	)

	// Draw first line of the top bar
	for x, ch := range topBarText {
		if x >= width {
			break
		}
		ui.screen.SetContent(x, 0, ch, nil, topBarStyle)
	}
}

func (ui *ChatUI) drawPromptBar() {
	width, _ := ui.screen.Size()

	// Draw input area with colored background and prompt
	inputStyle := tcell.StyleDefault.
		Background(tcell.ColorDarkBlue).
		Foreground(tcell.ColorWhite)

	promptStyle := tcell.StyleDefault.
		Background(tcell.ColorDarkBlue).
		Foreground(tcell.ColorLightGreen).
		Bold(true)

	// Draw input background
	for x := 0; x < width; x++ {
		ui.screen.SetContent(x, ui.inputStartY, ' ', nil, inputStyle)
	}

	// Draw prompt
	prompt := ">> "
	for x, ch := range prompt {
		ui.screen.SetContent(x, ui.inputStartY, ch, nil, promptStyle)
	}

	// Draw input text
	for x, ch := range ui.inputBuffer {
		if x+len(prompt) >= width {
			break
		}
		ui.screen.SetContent(x+len(prompt), ui.inputStartY, ch, nil, inputStyle)
	}

	// Position cursor after the prompt
	ui.screen.ShowCursor(len(prompt)+len(ui.inputBuffer), ui.inputStartY)
}

func drawMOTD(ui *ChatUI) {
	ui.addSystem(Message{text: "", timestamp: time.Now()})

	for _, line := range logo {
		ui.addWarning(Message{
			text:      line,
			timestamp: time.Now(),
		})
		time.Sleep(50 * time.Millisecond)
	}
	time.Sleep(700 * time.Millisecond)

	ui.addSystem(Message{text: "", timestamp: time.Now()})

	ui.addServer(Message{
		timestamp: time.Now(),
		parts: []TextPart{
			{text: "Welcome to One Dead: The strategic guessing game", bold: true},
		},
	})
	time.Sleep(300 * time.Millisecond)

	ui.addSystem(Message{text: "", timestamp: time.Now()})

	ui.addServer(Message{
		text:      "How to Play:",
		timestamp: time.Now(),
	})
	ui.addServer(Message{
		text:      "-------------",
		timestamp: time.Now(),
	})
	ui.addServer(Message{
		timestamp: time.Now(),
		text:      "You and your opponent will each set a secret 4-digit code",
	})
	ui.addServer(Message{
		timestamp: time.Now(),
		text:      "Your goal is to correctly guess the secret code of your opponent with the fewest number of guesses",
	})
	ui.addServer(Message{
		timestamp: time.Now(),
		text:      "In the shortest amount of time",
	})
	ui.addServer(Message{
		timestamp: time.Now(),
		text:      "You will be given two feedback values for each guess:",
	})
	ui.addServer(Message{
		timestamp: time.Now(),
		text:      "1. The number of digits that are correct and in the correct position",
	})
	ui.addServer(Message{
		timestamp: time.Now(),
		text:      "2. The number of digits that are correct but in the wrong position",
	})
	ui.addServer(Message{
		timestamp: time.Now(),
		text:      "Good luck!",
	})
	ui.addServer(Message{
		text:      "-------------",
		timestamp: time.Now(),
	})
	time.Sleep(200 * time.Millisecond)

	ui.addSystem(Message{text: "", timestamp: time.Now()})

	ui.addServer(Message{
		timestamp: time.Now(),
		text:      "This server is 2 months 3 weeks 6 days 5 hours 12 minutes 3 seconds old",
	})
	time.Sleep(300 * time.Millisecond)

	activeSessions := len(server.ActiveSessions)
	activeConnections := server.GetUsersCount()

	ui.addServer(Message{
		timestamp: time.Now(),
		text:      fmt.Sprintf("%d active sessions", activeSessions),
	})

	ui.addServer(Message{
		timestamp: time.Now(),
		text:      fmt.Sprintf("%d active connections", activeConnections),
	})
	time.Sleep(300 * time.Millisecond)

	ui.addSystem(Message{text: "", timestamp: time.Now()})

	createNewGameForPlayer(ui)
}

func createNewGameForPlayer(ui *ChatUI) {
	ui.addStyledMessage([]TextPart{
		{text: "Logging in as ", bold: false},
		{text: ui.player.Name, bold: true},
		{text: "...", bold: false},
	},
		tcell.ColorLightGreen, time.Now(),
	)
	time.Sleep(600 * time.Millisecond)

	// ui.addServer(Message{
	// 	timestamp: time.Now(),
	// 	text:      "Your Stats:",
	// })
	// ui.addServer(Message{
	// 	text:      "-------------",
	// 	timestamp: time.Now(),
	// })
	// ui.addServer(Message{
	// 	text:      "Games Played: 0",
	// 	timestamp: time.Now(),
	// })
	// ui.addServer(Message{
	// 	text:      "Games Won: 0",
	// 	timestamp: time.Now(),
	// })
	// ui.addServer(Message{
	// 	text:      "Games Lost: 0",
	// 	timestamp: time.Now(),
	// })
	// ui.addServer(Message{
	// 	text:      "-------------",
	// 	timestamp: time.Now(),
	// })
	// time.Sleep(200 * time.Millisecond)

	ui.addSystem(Message{text: "", timestamp: time.Now()})

	ui.addServer(Message{
		text:      "Finding a game...",
		timestamp: time.Now(),
	})
	time.Sleep(300 * time.Millisecond)

	ui.addServer(Message{
		text:      fmt.Sprintf("Adding you to a game %d...", ui.gameSession.Id),
		timestamp: time.Now(),
	})
	time.Sleep(200 * time.Millisecond)

	ui.addSystem(Message{text: "", timestamp: time.Now()})
	time.Sleep(300 * time.Millisecond)

	ui.addServer(Message{
		parts: []TextPart{
			{text: "READY!", bold: true},
		},
		timestamp: time.Now(),
	})

	ui.addWarning(Message{
		text:      "Enter your test code",
		timestamp: time.Now(),
	})

	uiC := ui.C.Subscribe()
	defer ui.C.Close(uiC)

	data := <-uiC

	ui.gameSession.AddPlayer(&game.Player{
		Name: ui.player.Name,
		Code: game.Code(data.text),
	})

	ui.addWarning(Message{
		text:      "Now waiting for opponent to join...",
		timestamp: time.Now(),
	})

	gameC := ui.gameSession.Events.Subscribe()
	defer ui.gameSession.Events.Close(gameC)

	for {
		gameEvent := <-gameC
		if gameEvent.Type == game.JOIN {
			ui.addServer(Message{
				parts: []TextPart{
					{text: "Opponent ", bold: false},
					{text: gameEvent.Player.Name, bold: true},
					{text: " has joined the game!", bold: false},
				},
				timestamp: time.Now(),
			})
			break
		} else if gameEvent.Type == game.START {
			ui.addServer(Message{
				parts: []TextPart{
					{text: "Game has started!", bold: true},
				},
				timestamp: time.Now(),
			})
			break
		}
	}
}
