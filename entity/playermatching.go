package entity

type WaitingMember struct {
	UserID uint
}

type MatchedPlayers struct {
	Category Category
	UserIDs  []uint
}
