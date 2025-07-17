package commands

import (
	"fmt"
	"time"

	"github.com/yamamushi/EscapingEden/logging"
)

// DefaultCommandLogger provides basic command logging
type DefaultCommandLogger struct {
	log logging.LoggerType
}

// NewDefaultCommandLogger creates a new command logger
func NewDefaultCommandLogger(log logging.LoggerType) *DefaultCommandLogger {
	return &DefaultCommandLogger{
		log: log,
	}
}

// LogCommand logs command execution
func (l *DefaultCommandLogger) LogCommand(cmd *PlayerCommand, result *CommandResult, duration time.Duration) {
	if result == nil {
		l.log.Println(logging.LogWarn, fmt.Sprintf(
			"Command executed with nil result - Type: %s, Player: %s, Duration: %v",
			cmd.Type, cmd.PlayerID, duration))
		return
	}

	if result.Success {
		l.log.Println(logging.LogInfo, fmt.Sprintf(
			"Command executed successfully - Type: %s, Player: %s, Duration: %v, Message: %s",
			cmd.Type, cmd.PlayerID, duration, result.Message))
	} else if result.Error != nil {
		l.log.Println(logging.LogWarn, fmt.Sprintf(
			"Command failed - Type: %s, Player: %s, Duration: %v, Error: %s",
			cmd.Type, cmd.PlayerID, duration, result.Error.Error()))
	} else {
		l.log.Println(logging.LogWarn, fmt.Sprintf(
			"Command failed without error - Type: %s, Player: %s, Duration: %v",
			cmd.Type, cmd.PlayerID, duration))
	}
}

// LogError logs command errors
func (l *DefaultCommandLogger) LogError(cmd *PlayerCommand, err error) {
	l.log.Println(logging.LogError, fmt.Sprintf(
		"Command error - Type: %s, Player: %s, Error: %v",
		cmd.Type, cmd.PlayerID, err))
}

// DefaultCommandMetrics provides basic command metrics
type DefaultCommandMetrics struct {
	stats map[CommandType]*commandStats
}

type commandStats struct {
	totalExecutions int64
	successCount    int64
	errorCount      int64
	totalTime       time.Duration
	errorsByCode    map[ErrorCode]int64
}

// NewDefaultCommandMetrics creates a new metrics collector
func NewDefaultCommandMetrics() *DefaultCommandMetrics {
	return &DefaultCommandMetrics{
		stats: make(map[CommandType]*commandStats),
	}
}

// RecordExecution records command execution metrics
func (m *DefaultCommandMetrics) RecordExecution(cmdType CommandType, duration time.Duration, success bool) {
	if _, exists := m.stats[cmdType]; !exists {
		m.stats[cmdType] = &commandStats{
			errorsByCode: make(map[ErrorCode]int64),
		}
	}

	m.stats[cmdType].totalExecutions++
	m.stats[cmdType].totalTime += duration

	if success {
		m.stats[cmdType].successCount++
	} else {
		m.stats[cmdType].errorCount++
	}
}

// RecordError records command error metrics
func (m *DefaultCommandMetrics) RecordError(cmdType CommandType, errorCode ErrorCode) {
	if _, exists := m.stats[cmdType]; !exists {
		m.stats[cmdType] = &commandStats{
			errorsByCode: make(map[ErrorCode]int64),
		}
	}

	m.stats[cmdType].errorsByCode[errorCode]++
}

// GetStats returns command execution statistics
func (m *DefaultCommandMetrics) GetStats() map[CommandType]CommandStats {
	result := make(map[CommandType]CommandStats)

	for cmdType, stats := range m.stats {
		var avgTime time.Duration
		if stats.totalExecutions > 0 {
			avgTime = stats.totalTime / time.Duration(stats.totalExecutions)
		}

		result[cmdType] = CommandStats{
			TotalExecutions: stats.totalExecutions,
			SuccessCount:    stats.successCount,
			ErrorCount:      stats.errorCount,
			AverageTime:     avgTime,
			ErrorsByCode:    stats.errorsByCode,
		}
	}

	return result
}
