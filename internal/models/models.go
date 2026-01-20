package models

type ValentineMsg struct {
	FromID int64
	Text   string
	ToUser int64 // Юзернейм или ID
}

type UserState struct {
	CurrentStep string
	Valentine *ValentineMsg
}

func NewUserState(userID int64) *UserState {
	return &UserState{
		CurrentStep: "idle",
		Valentine:   &ValentineMsg{
			FromID: userID,
		},
	}
}