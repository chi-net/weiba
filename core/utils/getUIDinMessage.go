package utils

import (
	"regexp"
	"strconv"
)

func GetUIDinMessage(content string) int64 {
	// Define the regex pattern to match UID
	re := regexp.MustCompile(`UID:(\d+)`)

	// Find the match
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		// Extract the first capture group
		uidStr := matches[1]

		// Convert the UID to int64
		uid, err := strconv.ParseInt(uidStr, 10, 64)
		if err != nil {
			return -1
		}

		return uid
	} else {
		return 0
	}
}
