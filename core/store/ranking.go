package store

import (
	"errors"
	"github.com/chi-net/weiba/core/types"
	"gorm.io/gorm"
	"sort"
)

func RecordRanking(Userid int64, Chatid int64, Name string) {
	mu.Lock()
	defer mu.Unlock()
	var ranking types.DBTableRanking
	result := sql.First(&ranking, "user_id = ? AND chat_id = ? AND name = ?", Userid, Chatid, Name)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		sql.Create(&types.DBTableRanking{
			ChatId: Chatid,
			UserId: Userid,
			Count:  1,
			Name:   Name,
		})
	} else {
		count := ranking.Count
		sql.Model(&ranking).Update("Count", count+1)
	}
	//result = sql.First(&ranking, "user_id = ? AND chat_id = ? AND name = ?", 0, Chatid, Name)
	//if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	//	sql.Create(&types.DBTableRanking{
	//		ChatId: Chatid,
	//		UserId: 0,
	//		Count:  1,
	//		Name:   Name,
	//	})
	//} else {
	//	count := ranking.Count
	//	sql.Model(&ranking).Update("Count", count+1)
	//}
}

func UpdateUsername(Userid int64, Name string) {
	mu.Lock()
	defer mu.Unlock()
	var model types.DBTableUser
	result := sql.First(&model, "user_id = ? AND name = ?", Userid, Name)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		sql.Create(&types.DBTableUser{
			UserId: Userid,
			Name:   Name,
		})
	} else {
		if model.Name != Name {
			sql.Model(&model).Update("name", Name)
		}
	}
}

func GetRanking(Name string) []types.RankingList {
	mu.RLock()
	defer mu.RUnlock()
	var rankings []types.DBTableRanking
	result := sql.Where("name = ?", Name).Find(&rankings)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	} else {
		sort.Slice(rankings, func(i, j int) bool {
			return rankings[i].Count > rankings[j].Count
		})
		var res []types.RankingList
		for _, ranking := range rankings {
			var resultUser types.DBTableUser
			resp := sql.First(&resultUser, "user_id = ?", ranking.UserId)
			if errors.Is(resp.Error, gorm.ErrRecordNotFound) {
				continue
			} else {
				res = append(res, types.RankingList{
					Name:  resultUser.Name,
					Count: ranking.Count,
					Id:    ranking.UserId,
				})
			}
		}
		if len(res) > 10 {
			return res[:10]
		} else {
			return res
		}
	}
}

func GetGroupRecordedMessages(ChatId int64, Name string) int64 {
	mu.Lock()
	defer mu.Unlock()
	var ranking types.DBTableRanking
	result := sql.Model(&ranking).Where("name = ? AND chat_id = ?", Name, ChatId)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return -1
	} else {
		return ranking.Count
	}
}
