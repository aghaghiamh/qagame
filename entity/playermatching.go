package entity

type WaitingMember struct {
	UserID uint
}

type MatchedUsers struct {
	Category Category
	UserIDs  []uint
}
