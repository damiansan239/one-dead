package main

import (
	"fmt"
	"log"
	"one_dead/internal/datastore"
	"one_dead/internal/game"

	_ "github.com/gdamore/tcell/v2/encoding"
	"github.com/gliderlabs/ssh"
)

var server *Server

func init() {
	db, err := datastore.NewDatastore("players.db")
	if err != nil {
		fmt.Printf("Failed to create player database: %v", err)
	}

	lobby := NewLobby()
	server = NewServer(db, lobby)
}

func main() {
	ssh.Handle(func(sess ssh.Session) {
		server.IncrementUserCount()

		// fetch player info
		player, err := server.Datastore.GetByName(sess.User())

		if err != nil || player == nil {
			fmt.Println("Creating new player")

			player, err = server.Datastore.CreateNewPlayer(sess.User())
			if err != nil {
				fmt.Printf("Failed to create new player: %v", err)
			}
		}

		gameSession := server.Lobby.GetFreeSession(player)
		server.AddSession(gameSession)

		go func() {
			<-sess.Context().Done()
			server.DecrementUserCount()
			if gameSession.Status == game.PENDING {
				server.RemoveSession(gameSession.Id)
				server.Lobby.AddSession(gameSession)
			}

			log.Println("connection closed")
		}()

		ui, err := NewChatUI(sess, player, gameSession)
		if err != nil {
			fmt.Printf("Failed to create UI: %v", err)
		}

		ui.Run()
	})
	fmt.Println("Starting server on port 2222...")

	opts := []ssh.Option{}
	// if true {
	// 	opts = append(opts, ssh.HostKeyFile("/home/yungwarlock/.ssh/id_rsa"))
	// }

	log.Fatal(ssh.ListenAndServe(":2222", nil, opts...))
}
