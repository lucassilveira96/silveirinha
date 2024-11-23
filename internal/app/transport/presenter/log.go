package presenter

import (
	"fmt"
	"time"
)

func LogDebbuger(message string, status string) {
	now := time.Now()

	formattedDate := now.Format("02-Jan-2006 15:04:05")
	logMessage := fmt.Sprintf("%s %s - %s", formattedDate, status, message)

	fmt.Println(logMessage)
}
