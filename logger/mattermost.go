package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type Transporter struct {
	logs           []LogEntry
	allowedLevels  []string
	batchSize      int
	interval       time.Duration
	URL            string
	botToken       string
	logChannelID   string
	errorChannelID string
	mutex          sync.Mutex
}

type LogEntry struct {
	Message string
	Level   string
}

type TransportOption struct {
	BatchSize      int
	Interval       time.Duration
	URL            string
	BotToken       string
	LogChannelID   string
	ErrorChannelID string
}

func NewTransport(opts TransportOption) *Transporter {
	mt := &Transporter{
		logs:           []LogEntry{},
		allowedLevels:  []string{"notify", "error"},
		batchSize:      opts.BatchSize,
		interval:       opts.Interval,
		URL:            opts.URL,
		botToken:       opts.BotToken,
		logChannelID:   opts.LogChannelID,
		errorChannelID: opts.ErrorChannelID,
	}
	go mt.batchLogs() // Run batching in a separate goroutine
	return mt
}

func (mt *Transporter) Log(level string, message string) error {
	if contains(mt.allowedLevels, level) {
		return nil
	}
	mt.mutex.Lock()
	defer mt.mutex.Unlock()
	mt.logs = append(mt.logs, LogEntry{
		Message: message,
		Level:   level,
	})
	return nil
}

func (mt *Transporter) batchLogs() {
	ticker := time.NewTicker(mt.interval)
	for range ticker.C {
		mt.mutex.Lock()
		if len(mt.logs) == 0 {
			mt.mutex.Unlock()
			continue
		}

		// Create batch for processing
		var logEntries []LogEntry
		if len(mt.logs) <= mt.batchSize {
			logEntries = mt.logs
			mt.logs = []LogEntry{}
		} else {
			logEntries = mt.logs[:mt.batchSize]
			mt.logs = mt.logs[mt.batchSize:]
		}
		mt.mutex.Unlock()

		// Process the batch
		mt.processBatch(logEntries)
	}
}

func (mt *Transporter) processBatch(entries []LogEntry) {
	var logMessage, errorMessage string
	separator := "\n\n---\n\n"
	for _, entry := range entries {
		fullMessage := entry.Message + separator
		if entry.Level == "error" {
			errorMessage += fullMessage
		} else {
			logMessage += fullMessage
		}
	}

	if logMessage != "" {
		mt.sendNotification(logMessage, true)
	}
	if errorMessage != "" {
		mt.sendNotification(errorMessage, false)
	}
}

func (mt *Transporter) sendNotification(message string, isLog bool) {
	channelID := mt.logChannelID
	tag := "@all"
	if !isLog {
		channelID = mt.errorChannelID
	}
	postData := map[string]string{
		"channel_id": channelID,
		"message":    tag + "\n\n" + message,
	}
	jsonData, err := json.Marshal(postData)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return
	}

	req, err := http.NewRequest("POST", mt.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", mt.botToken))
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Notification failed with status: %v", resp.Status)
	}
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
