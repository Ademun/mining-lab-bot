package cmd

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/Ademun/mining-lab-bot/internal/polling"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/go-telegram/bot/models"
	"github.com/google/uuid"
)

func extractLabType(update *models.Update) polling.LabType {
	labTypeStr := update.CallbackQuery.Data
	labTypeStr = strings.TrimPrefix(labTypeStr, "type:")

	var labType polling.LabType
	switch labTypeStr {
	case "performance":
		labType = polling.LabTypePerformance
	case "defence":
		labType = polling.LabTypeDefence
	}
	return labType
}

func extractLabDomain(update *models.Update) *polling.LabDomain {
	labDomainStr := update.CallbackQuery.Data
	labDomainStr = strings.TrimPrefix(labDomainStr, "domain:")

	var labDomain polling.LabDomain
	switch labDomainStr {
	case "mechanics":
		labDomain = polling.LabDomainMechanics
	case "virtual":
		labDomain = polling.LabDomainVirtual
	case "electricity":
		labDomain = polling.LabDomainElectricity
	}
	return &labDomain
}

func extractWeekday(update *models.Update) *int {
	labWeekdayStr := update.CallbackQuery.Data
	labWeekdayStr = strings.TrimPrefix(labWeekdayStr, "weekday:")

	if labWeekdayStr == "skip" {
		return nil
	}
	labWeekdayInt, _ := strconv.Atoi(labWeekdayStr)
	return &labWeekdayInt
}

// extractListingData returns new sub index if the selected action was "move:idx", and sub uuid if it was "delete"
func extractListingData(update *models.Update) (*int, *uuid.UUID) {
	dataFields := strings.Split(update.CallbackQuery.Data, ":")
	switch dataFields[0] {
	case "move":
		newIndex, err := strconv.Atoi(dataFields[1])
		if err != nil {
			slog.Error("Failed to parse new sub index",
				"index", dataFields[1],
				"error", err,
				"service", logger.TelegramBot)
		}
		return &newIndex, nil
	case "delete":
		subUUID, err := uuid.Parse(dataFields[1])
		if err != nil {
			slog.Error("Failed to parse sub uuid",
				"uuid", dataFields[1],
				"error", err,
				"service", logger.TelegramBot)
		}
		return nil, &subUUID
	}
	return nil, nil
}

func extractLesson(update *models.Update) *int {
	labLessonStr := update.CallbackQuery.Data
	labLessonStr = strings.TrimPrefix(labLessonStr, "lesson:")

	if labLessonStr == "skip" {
		return nil
	}
	labLessonInt, _ := strconv.Atoi(labLessonStr)
	return &labLessonInt
}

func extractWeekParity(update *models.Update) string {
	weekParityStr := update.CallbackQuery.Data
	weekParityStr = strings.TrimPrefix(weekParityStr, "parity:")
	return weekParityStr
}
