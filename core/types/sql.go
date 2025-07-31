package types

type DBTableRanking struct {
	ChatId int64
	UserId int64
	Count  int64
}
type DBTableUser struct {
	UserId   int64
	Username string
}
