package teacher

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/config"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/Ademun/mining-lab-bot/pkg/model"
	"github.com/robfig/cron/v3"
)

type Service interface {
	FindTeachersForTime(ctx context.Context, targetTime time.Time, auditorium int) []model.Teacher
}

type teacherService struct {
	teacherRepo   Repo
	options       config.TeacherConfig
	weekNumber    int
	cronScheduler *cron.Cron
}

func New(repo Repo, opts *config.TeacherConfig) Service {
	return &teacherService{
		teacherRepo: repo,
		options:     *opts,
		weekNumber:  opts.StartingWeek,
	}
}

func (s *teacherService) Start() error {
	slog.Info("Starting", "service", logger.ServiceTeacher)
	c := cron.New(cron.WithLocation(time.Local))
	_, err := c.AddFunc("0 0 * 0", func() {
		slog.Info("Updating week number", logger.ServiceTeacher)
		if s.weekNumber == 1 {
			s.weekNumber = 2
			return
		}
		s.weekNumber = 1
	})
	if err != nil {
		slog.Info("Cron error", "error", err, "service", logger.ServiceTeacher)
	}
	c.Start()
	s.cronScheduler = c
	slog.Info("Started", "service", logger.ServiceTeacher)
	return nil
}

func (s *teacherService) Stop(ctx context.Context) {
	<-ctx.Done()
	s.cronScheduler.Stop()
	slog.Info("Stopped", "service", logger.ServiceTeacher)
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
