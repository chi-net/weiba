package types

type DBTableRanking struct {
	ChatId int64
	UserId int64
	Count  int64
	Name   string
	Id     uint `gorm:"primaryKey"`
}
type DBTableUser struct {
	UserId int64
	Name   string
	Id     uint `gorm:"primaryKey"`
}

type DBTableKVNumbers struct {
	Id    uint `gorm:"primaryKey"`
	Key   string
	Value int64
}
