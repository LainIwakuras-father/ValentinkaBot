package db

import (
	"sync"

	"github.com/LainIwakuras-father/ValentinkaBot/internal/models"
)

type Db struct {
	userStates map[int64]*models.UserState
	mu         sync.RWMutex
}

func NewDb() *Db {
	return &Db{
		userStates: make(map[int64]*models.UserState),
	}
}

func (d *Db) GetUserState(userID int64) *models.UserState {
	d.mu.RLock()
	defer d.mu.RUnlock()

	state, exists := d.userStates[userID]
	if !exists {
		state = models.NewUserState(userID)
	}
	return state
}

func (d *Db) SetUserState(userID int64, state *models.UserState) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.userStates[userID] = state
}

// ResetState сбрасывает состояние пользователя
func (d *Db) ResetState(userID int64) {
	d.mu.Lock()
	defer d.mu.Unlock()

	delete(d.userStates, userID)
}