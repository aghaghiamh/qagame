package entity

type Game struct {
	ID          uint
	PlayerIDs   [2]uint
	QuestionIDs []uint
	CategoryID  uint
}

type Player struct {
	ID              uint
	UserID          uint
	GameID          uint
	Score           uint
	PlayerAnswerIDs []uint
}

type PlayerAnswer struct {
	ID         uint
	PlayerID   uint
	QuestionID uint
	Choice     PossibleAnswerChoice
}
