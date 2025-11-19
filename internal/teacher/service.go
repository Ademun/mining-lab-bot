package teacher

import (
	"context"
	"log/slog"
	"time"

	"github.com/Ademun/mining-lab-bot/pkg/config"
	"github.com/Ademun/mining-lab-bot/pkg/logger"
	"github.com/robfig/cron/v3"
)

type Service interface {
	FindTeachersForTime(ctx context.Context, targetTime time.Time, auditorium int) []Teacher
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
		slog.Info("Updating week number", "service", logger.ServiceTeacher)
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

func (s *teacherService) FindTeachersForTime(ctx context.Context, targetTime time.Time, auditorium int) []Teacher {
	teachers, err := s.teacherRepo.FindByWeekNumberWeekdayAuditorium(ctx, calculateWeekNumber(s.weekNumber, time.Now(), targetTime), targetTime.Weekday(), auditorium)
	if err != nil {
		slog.Error("Failed to find teachers", "error", err, "service", logger.ServiceTeacher)
	}

	res := make([]Teacher, 0)
	for _, teacher := range teachers {
		normalized := targetTime.Format("15:01")
		if (teacher.TimeStart <= normalized) && (normalized <= teacher.TimeEnd) {
			res = append(res, teacher)
		}
	}

	return res
}

func calculateWeekNumber(currentWeek int, currentTime, targetTime time.Time) int {
	currentWeekStart := getWeekMonday(currentTime)
	targetWeekStart := getWeekMonday(targetTime)

	weekDiff := int(targetWeekStart.Sub(currentWeekStart).Hours() / (24 * 7))

	if weekDiff%2 == 0 {
		return currentWeek
	}

	if currentWeek == 1 {
		return 2
	}
	return 1
}

func getWeekMonday(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}

	daysToSubtract := weekday - 1
	monday := t.AddDate(0, 0, -daysToSubtract)

	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())
}
