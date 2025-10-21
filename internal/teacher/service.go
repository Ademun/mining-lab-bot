package teacher

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/Ademun/mining-lab-bot/pkg/model"
)

type Service interface {
	FindTeachersForTime(ctx context.Context, targetTime time.Time, auditorium int) []model.Teacher
}

type teacherService struct {
	teacherRepo Repo
	weekNumber  int
}

func New(repo Repo) Service {
	return &teacherService{
		teacherRepo: repo,
	}
}

func (s *teacherService) FindTeachersForTime(ctx context.Context, targetTime time.Time, auditorium int) []model.Teacher {
	fmt.Println(targetTime.Weekday())
	teachers, err := s.teacherRepo.FindByWeekNumberWeekdayAuditorium(ctx, 2, targetTime.Weekday(), auditorium)
	if err != nil {
		slog.Error("Failed to find teachers", "error", err, "service", logger.ServiceTeacher)
	}

	normalized := timeToMinutes(targetTime.Hour(), targetTime.Minute())
	res := make([]model.Teacher, 0)
	for _, teacher := range teachers {
		start, _ := time.Parse("15:04", teacher.TimeStart)
		startMins := timeToMinutes(start.Hour(), start.Minute())
		end, _ := time.Parse("15:04", teacher.TimeEnd)
		endMins := timeToMinutes(end.Hour(), end.Minute())
		if normalized >= startMins && normalized < endMins {
			res = append(res, teacher)
		}
	}

	return res
}

func timeToMinutes(h, m int) int {
	return h*60 + m
}
