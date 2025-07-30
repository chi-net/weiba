package handlers

import (
	"context"
	"fmt"
	"github.com/chi-net/weiba/core/types"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"log"
	"runtime"
	"strconv"
	"time"
)

func InfoHandler(ctx context.Context, b *bot.Bot, update *models.Update, config types.YmlConfigurationData) {
	if update.Message.From.ID != config.AdminUID && config.AdminUID != -1 {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "weibabot@1.1.1 - A bot which can process authentication of private groups and channels\nOpenSource: https://github.com/chi-net/weiba",
		})
	} else {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		// Get CPU Info
		cpuInfo, err := cpu.Info()
		if err != nil {
			fmt.Println("Error getting CPU info:", err)
			return
		}
		// Get CPU usage percentage
		cpuPercent, err := cpu.Percent(2*time.Second, false)
		if err != nil {
			fmt.Println("Error getting CPU percent:", err)
			return
		}
		// Get Memory Info
		virtualMem, err := mem.VirtualMemory()
		if err != nil {
			fmt.Println("Error getting memory info:", err)
			return
		}
		// Get Disk Info
		diskInfo, err := disk.Usage("/")
		if err != nil {
			fmt.Println("Error getting disk info:", err)
			return
		}
		hostInfo, err := host.Info()
		if err != nil {
			log.Fatalf("Error getting host info: %v", err)
		}

		msg := "weibabot@canary - A bot which can process authentication of private groups and channels\nSystem Info:\n"
		msg += "Host:" + hostInfo.Platform + "(Kernel:" + hostInfo.OS + " " + hostInfo.KernelVersion + ")\n"
		msg += "CPU: " + cpuInfo[0].ModelName + "(" + strconv.Itoa(len(cpuInfo)) + ") " + strconv.FormatFloat(cpuInfo[0].Mhz/1000, 'f', 2, 64) + "GHz " + strconv.FormatFloat(cpuPercent[0], 'f', 2, 64) + "% Used\n"
		msg += "Memory: " + strconv.FormatFloat(float64(virtualMem.Total/1024/1024), 'f', 2, 64) + "MB Total;" + strconv.FormatFloat(float64(virtualMem.Used/1024/1024), 'f', 2, 64) + "MB Used\n"
		msg += "Disk: " + strconv.FormatFloat(float64(diskInfo.Total/1024/1024/1024), 'f', 2, 64) + "GB Total;" + strconv.FormatFloat(float64(diskInfo.Used/1024/1024/1024), 'f', 2, 64) + "GB Used\n"
		msg += "Application:\n"
		msg += "Current active goroutines:" + strconv.Itoa(runtime.NumGoroutine()) + "\n"
		msg += "Current occupied Memory:" + strconv.FormatFloat(float64(memStats.Alloc/1024/1024), 'f', 2, 64) + "MB"
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   msg,
		})
	}
}
