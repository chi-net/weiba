package utils

import (
	"fmt"
	"github.com/chi-net/weiba/core/types"
	"github.com/go-telegram/bot"
	"strconv"
)

func GetRankingMessage(result []types.RankingList, count types.RankingData) string {
	message := ""
	message += bot.EscapeMarkdown("总消息数: " + strconv.FormatInt(count.TotalMessages, 10) + ", 统计发言用户数: " + strconv.Itoa(count.TotalUsers) + "\n")
	for i, val := range result {
		percentage := (float64(val.Count) / float64(count.TotalMessages)) * 100
		percentageOutput := fmt.Sprintf("%.2f%%", percentage)
		refmsg := bot.EscapeMarkdown(strconv.Itoa(i+1) + ". " + val.Name + ": " + strconv.FormatInt(val.Count, 10) + "条, " + percentageOutput)
		if i == 0 {
			message += "**"
		}
		message += "> " + refmsg
		if i == len(result)-1 {
			message += "||"
		}
		message += "\n"
	}
	message += "\n"
	return message
}
