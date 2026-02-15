package types

import "time"

// SafetyMonitor holds the safety monitoring state in memory for admin user
type SafetyMonitor struct {
	LastCheckIn           time.Time
	CheckInCount          int
	LastSentTime          time.Time
	IsWaitingReply        bool
	ConsecutiveNoResponse int // 连续未回应次数
}