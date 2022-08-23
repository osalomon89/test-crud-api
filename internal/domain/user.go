package domain

type User struct {
	ID    uint
	Email string
	Items []Item
}

type UserItem struct {
	ID     uint
	UserID uint
	ItemID uint
}
