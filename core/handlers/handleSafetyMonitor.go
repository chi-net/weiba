package handlers

import (
	"context"
	"time"

	"github.com/chi-net/weiba/core/types"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// SafetyMonitorHandler handles safety monitoring functionality
func SafetyMonitorHandler(ctx context.Context, b *bot.Bot, config types.YmlConfigurationData, safetyMonitor *types.SafetyMonitor) {
	if !config.Features.SafetyMonitor || config.AdminUID == -1 {
		return
	}

	// 每2小时在8-22点(UTC+8)检查一次
	ticker := time.NewTicker(2 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// 检查当前时间是否在8-22点（UTC+8）
			now := time.Now().In(time.FixedZone("UTC+8", 8*3600))
			hour := now.Hour()
			
			if hour >= 8 && hour <= 22 {
				// 发送新的安全检查消息（内部会自动处理30分钟检查）
				sendSafetyCheckMessage(ctx, b, config.AdminUID, safetyMonitor, config)
			}
		}
	}
}

// HandleSafetyCallback handles callback from safety check buttons
func HandleSafetyCallback(ctx context.Context, b *bot.Bot, update *models.Update, config types.YmlConfigurationData, safetyMonitor *types.SafetyMonitor) {
	if update.CallbackQuery == nil || update.CallbackQuery.From.ID != config.AdminUID {
		return
	}

	callbackData := update.CallbackQuery.Data

	// 回应callback以移除加载状态
	_, _ = b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: update.CallbackQuery.ID,
	})

	switch callbackData {
	case "safety_checkin":
		updateCheckIn(safetyMonitor)
		// 编辑消息显示已签到
		if update.CallbackQuery.Message.Message != nil {
			_, _ = b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
				MessageID: update.CallbackQuery.Message.Message.ID,
				Text:      "✅ 签到成功！管理员状态正常。",
			})
		}

	case "safety_alert":
		// 立即发送警告到监控群组
		sendAlertToGroups(ctx, b, config)
		// 手动报警后重置等待状态但不影响连续未回应计数
		safetyMonitor.IsWaitingReply = false
		// 编辑消息显示已报警
		if update.CallbackQuery.Message.Message != nil {
			_, _ = b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
				MessageID: update.CallbackQuery.Message.Message.ID,
				Text:      "",
			})
		}
	}
}

func sendSafetyCheckMessage(ctx context.Context, b *bot.Bot, adminUID int64, safetyMonitor *types.SafetyMonitor, config types.YmlConfigurationData) {
	keyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "❌ 没学", CallbackData: "safety_checkin"},
				{Text: "✅ 学了", CallbackData: "safety_alert"},
			},
		},
	}

	message := "📖 考勤检查\n\n请确认您的状态：\n• 点击「没学」表示未学习数学分析、高等代数和解析几何\n• 点击「学了」表示已学习上述内容\n\n⏰ 请在30分钟内响应，否则系统将自动警告"

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      adminUID,
		Text:        message,
		ReplyMarkup: keyboard,
	})

	if err == nil {
		markMessageSent(safetyMonitor)
		
		// 启动30分钟后的检查
		go func() {
			time.Sleep(30 * time.Minute)
			
			// 检查是否仍在等待回复（即30分钟内未回应）
			if safetyMonitor.IsWaitingReply && time.Since(safetyMonitor.LastSentTime) >= 30*time.Minute {
				safetyMonitor.ConsecutiveNoResponse++
				safetyMonitor.IsWaitingReply = false
				
				// 如果连续2次未回应，立即发送警报
				if safetyMonitor.ConsecutiveNoResponse >= 2 {
					sendAlertToGroups(ctx, b, config)
					// 发送警报后重置连续计数
					safetyMonitor.ConsecutiveNoResponse = 0
				}
			}
		}()
	}
}

func sendAlertToGroups(ctx context.Context, b *bot.Bot, config types.YmlConfigurationData) {
	// 使用UTC+8时区格式化时间
	utc8 := time.FixedZone("UTC+8", 8*3600)
	alertTime := time.Now().In(utc8).Format("2006-01-02 15:04:05")
	alertMessage := "🚨 紧急警报 🚨\n\n管理员可能处于危险状态！\n\n请立即确认管理员安全状况，必要时联系相关部门。\n\n⏰ " + alertTime + " (UTC+8)"

	for _, groupID := range config.Whitelists.MonitorStatus {
		_, _ = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: groupID,
			Text:   alertMessage,
		})
	}
}

// Helper functions for SafetyMonitor operations
func updateCheckIn(sm *types.SafetyMonitor) {
	now := time.Now()
	sm.LastCheckIn = now
	sm.IsWaitingReply = false
	// 签到成功，重置连续未回应计数
	sm.ConsecutiveNoResponse = 0
}

func resetCheckIn(sm *types.SafetyMonitor) {
	sm.CheckInCount = 0
	sm.LastCheckIn = time.Time{}
	sm.IsWaitingReply = false
	sm.ConsecutiveNoResponse = 0
}

func markMessageSent(sm *types.SafetyMonitor) {
	sm.LastSentTime = time.Now()
	sm.IsWaitingReply = true
}
