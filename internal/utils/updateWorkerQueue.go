package utils

import (
	"log/slog"
	"os"
	"strconv"
)

func UpdateWorkerQueueSize(newCount int) {
	newCountStr := strconv.Itoa(newCount)
	os.Setenv("WORKERS_QUEUE_SIZE", newCountStr)
	slog.Info("worker queue size updated to :", "size", os.Getenv("WORKERS_QUEUE_SIZE"))
}

func AddWorkerQueueSize(newCount int) {
	currentCount := GetWorkerQueueSize()
	currentCount += newCount
	currentCountStr := strconv.Itoa(currentCount)
	os.Setenv("WORKERS_QUEUE_SIZE", currentCountStr)
	slog.Info("worker queue size updated to :", "size", os.Getenv("WORKERS_QUEUE_SIZE"))
}

func GetWorkerQueueSize() int {
	fallback := 100
	v, err := strconv.Atoi("WORKERS_QUEUE_SIZE")
	if err != nil {
		return fallback
	}
	return v
}
