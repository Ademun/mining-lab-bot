package teacher

type Teacher struct {
	Name       string `db:"name"`
	Auditorium int    `db:"auditorium"`
	WeekNumber int    `db:"week_number"`
	Weekday    int    `db:"weekday"`
	TimeStart  string `db:"time_start"`
	TimeEnd    string `db:"time_end"`
}
