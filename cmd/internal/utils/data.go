package utils

type Lesson struct {
	Text         string
	CallbackData string
}

var DefaultLessons = []Lesson{
	{Text: "08:50 - 10:20 - 1️⃣ пара", CallbackData: "sub_creation:lesson:1"},
	{Text: "10:35 - 12:05 - 2️⃣ пара", CallbackData: "sub_creation:lesson:2"},
	{Text: "12:35 - 14:05 - 3️⃣ пара", CallbackData: "sub_creation:lesson:3"},
	{Text: "14:15 - 15:45 - 4️⃣ пара", CallbackData: "sub_creation:lesson:4"},
	{Text: "15:55 - 17:20 - 5️⃣ пара", CallbackData: "sub_creation:lesson:5"},
	{Text: "17:30 - 19:00 - 6️⃣ пара", CallbackData: "sub_creation:lesson:6"},
	{Text: "19:10 - 20:30 - 7️⃣ пара", CallbackData: "sub_creation:lesson:7"},
	{Text: "20:40 - 22:00 - 8️⃣ пара", CallbackData: "sub_creation:lesson:8"},
}
