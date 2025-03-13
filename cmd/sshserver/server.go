package main

import (
	"one_dead/internal/datastore"
	"one_dead/internal/game"
	"sync"
)

type Server struct {
	ActiveSessions map[uint16]*game.Session
	Datastore      *datastore.Datastore
	Lobby          *Lobby
	userCount      int
	sync.RWMutex
}

func NewServer(db *datastore.Datastore, lobby *Lobby) *Server {
	return &Server{
		Datastore:      db,
		Lobby:          lobby,
		ActiveSessions: make(map[uint16]*game.Session),
	}
}

func (s *Server) AddSession(session *game.Session) {
	s.ActiveSessions[session.Id] = session
}

func (s *Server) RemoveSession(sessionID uint16) {
	delete(s.ActiveSessions, sessionID)
}

func (s *Server) GetUsersCount() int {
	s.RLock()
	defer s.RUnlock()
	return s.userCount
}

func (s *Server) IncrementUserCount() {
	s.Lock()
	defer s.Unlock()
	s.userCount++
}

func (s *Server) DecrementUserCount() {
	s.Lock()
	defer s.Unlock()
	if s.userCount > 0 {
		s.userCount--
	}
}
