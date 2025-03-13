package main

import (
	"fmt"
	"one_dead/internal/datastore"
	"one_dead/internal/game"
	"time"

	"github.com/gdamore/tcell/v2"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func (ui *ChatUI) Run() {
	// Create a ticker for updating the top bar (updates every second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	// Update the screen at regular intervals
	go func() {
		for range ticker.C {
			ui.draw()
			ui.screen.Show()
		}
	}()

	// Subscribe to chat messages
	go func() {
		gameC := ui.gameSession.Events.Subscribe()
		defer ui.gameSession.Events.Close(gameC)

		for message := range gameC {
			if message.Type == game.JOIN && message.Player.Name != ui.player.Name {
				onPlayerJoined(ui, message.Player)
			} else if message.Type == game.TRY {
				onPlaysTry(ui, &message)
			} else if message.Type == game.START {
				onGameStarted(ui)
			} else if message.Type == game.END {
				onGameEnd(ui, &message)
			}
		}
	}()

	ui.draw()
	ui.screen.Show()
	go drawMOTD(ui)

	for {
		ui.draw()
		ui.screen.Show()

		ev := ui.screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlC:
				ui.screen.Fini()
			case tcell.KeyEnter:
				handleEnter(ui)
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				handleBackSpace(ui)
			case tcell.KeyPgUp:
				handleScrollUp(ui)
			case tcell.KeyPgDn:
				handleScrollDown(ui)
			case tcell.KeyRune:
				ui.inputBuffer += string(ev.Rune())
			}
		case *tcell.EventResize:
			handleResize(ui)
		}
	}
}

func onPlayerJoined(ui *ChatUI, player *game.Player) {
	ui.addServer(Message{
		timestamp: time.Now(),
		text:      fmt.Sprintf("%s has joined this game", player.Name),
	})
}

func onPlaysTry(ui *ChatUI, message *game.SessionMessage) {
	ui.addServer(Message{
		text: fmt.Sprintf("[TRY][%s] %s => %d dead, %d injured",
			message.Player.Name,
			message.Code,
			message.Result.Dead,
			message.Result.Injured,
		),
		timestamp: time.Now(),
	})
}

func onGameStarted(ui *ChatUI) {
	ui.addServer(Message{
		text:      "Game started",
		timestamp: time.Now(),
	})
}

func onGameEnd(ui *ChatUI, message *game.SessionMessage) {
	ui.addServer(Message{
		text:      fmt.Sprintf("Game is won by %s", message.Player.Name),
		timestamp: time.Now(),
	})

	id, err := gonanoid.New()
	if err != nil {
		fmt.Println("Failed to generate game ID")
	}

	hasWon := ui.gameSession.Winner.Name == ui.player.Name

	players := []*datastore.Player{}
	for _, p := range ui.gameSession.Players {
		players = append(players, &datastore.Player{
			Name: p.Name,
		})
	}

	server.Datastore.SaveGame(&datastore.Game{
		Duration: 5,
		Id:       id,
		Won:      hasWon,
		Players:  players,
		Date:     time.Now(),
		Tries:    len(ui.gameSession.Tries),
	})

	server.RemoveSession(ui.gameSession.Id)
}

func handleScrollUp(ui *ChatUI) {
	_, height := ui.screen.Size()
	maxScroll := len(ui.messages) - (height - 2)
	if maxScroll < 0 {
		maxScroll = 0
	}
	ui.scrollOffset = min(maxScroll, ui.scrollOffset+3)
}

func handleScrollDown(ui *ChatUI) {
	ui.scrollOffset = max(0, ui.scrollOffset-3)
}

func handleResize(ui *ChatUI) {
	ui.screen.Sync()
	_, height := ui.screen.Size()
	ui.inputStartY = height - 1
}

func handleBackSpace(ui *ChatUI) {
	if len(ui.inputBuffer) > 0 {
		ui.inputBuffer = ui.inputBuffer[:len(ui.inputBuffer)-1]
	}
}

func handleEnter(ui *ChatUI) {
	if ui.inputBuffer != "" {
		timestamp := time.Now().UTC().Format("15:04:05")
		ui.C.Publish(Message{
			timestamp: time.Now(),
			text:      ui.inputBuffer,
			color:     tcell.ColorWhite,
		})
		ui.addMessage(fmt.Sprintf("[%s] <%s> %s", timestamp, ui.player.Name, ui.inputBuffer), tcell.ColorWhite)

		if ui.gameSession.Status == game.ACTIVE {
			ui.gameSession.AddTrial(ui.player.Name, game.Code(ui.inputBuffer))
		}

		ui.inputBuffer = ""
		ui.scrollOffset = 0
	}
}
