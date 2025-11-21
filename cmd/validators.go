package cmd

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

func validateLabNumber(labNumberStr string) (int, string) {
	labNumber, err := strconv.Atoi(labNumberStr)
	if err != nil {
		return 0, "Номер лабораторной работы должен быть числом"
	}
	if labNumber < 1 || labNumber > 100 {
		return 0, "Номер лабораторной работы должен быть в диапазоне 1-100"
	}
	return labNumber, ""
}

func validateLabAuditorium(labAuditoriumStr string) (int, string) {
	labAuditorium, err := strconv.Atoi(labAuditoriumStr)
	if err != nil {
		return 0, "Номер аудитории должен быть числом"
	}
	if labAuditorium < 1 || labAuditorium > 1000 {
		return 0, "Номер аудитории должен быть в диапазоне 1-1000"
	}
	return labAuditorium, ""
}

func validateTeacherSurname(surname string) (string, string) {
	surname = strings.TrimSpace(surname)

	if surname == "" {
		return "", "Фамилия преподавателя не может быть пустой"
	}

	if utf8.RuneCountInString(surname) > 100 {
		return "", "Фамилия преподавателя не может быть длиннее 100 символов"
	}

	return surname, ""
}
