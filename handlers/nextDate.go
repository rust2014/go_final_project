package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DefaultFormatDate = "20060102" // дефолтный формат

func isLeapYear(year int) bool { // проверка високосного
	if year%4 == 0 {
		if year%100 == 0 {
			return year%400 == 0
		}
		return true
	}
	return false
}

func NextDate(now time.Time, dateStr string, repeat string) (string, error) {
	date, err := time.Parse(DefaultFormatDate, dateStr) // парсинг исходной даты
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err) //???
	}

	if repeat == "" { // если правило повторения пустое то вернем ошибку
		return "", errors.New("repeat is required")
	}

	repeatParts := strings.Split(repeat, " ") // разделил правило повторения на части
	if len(repeatParts) == 0 {
		return "", errors.New("invalid repeat rule format")
	}

	switch repeatParts[0] {
	case "d":
		if len(repeatParts) != 2 {
			return "", errors.New("invalid repeat rule format for d")
		}
		days, err := strconv.Atoi(repeatParts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("invalid number of days")
		}
		date = date.AddDate(0, 0, days)

	case "y":
		if len(repeatParts) != 1 {
			return "", errors.New("invalid repeat rule format for y")
		}
		for date.Before(now) || date.Equal(now) {
			date = date.AddDate(1, 0, 0)
			if date.Month() == time.February && date.Day() == 29 && !isLeapYear(date.Year()) {
				date = date.AddDate(0, 0, 1)
			}
		}
	case "m":
		if len(repeatParts) != 2 {
			return "", errors.New("invalid repeat rule format for m")
		}
		months, err := strconv.Atoi(repeatParts[1]) // Преобразование месяцев в число
		if err != nil || months == 0 {
			return "", errors.New("invalid repeat rule format for m")
		}
		date = date.AddDate(0, months, 0)
	case "w":
		if len(repeatParts) != 2 {
			return "", errors.New("invalid repeat rule format for w")
		}
		weeks, err := strconv.Atoi(repeatParts[1]) // Преобразование недель в число
		if err != nil || weeks == 0 {
			return "", errors.New("invalid repeat rule format for w")
		}
		date = date.AddDate(0, 0, weeks*7)
	default:
		return "", errors.New("unsupported repeat rule") // ошибка
	}
	return date.Format(DefaultFormatDate), nil
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	timeNow, err := time.Parse(DefaultFormatDate, now)
	if err != nil {
		http.Error(w, "Invalid 'now' date format", http.StatusBadRequest) // вернет Invalid 'now' date format если не верный формат
		return
	}
	nextDate, err := NextDate(timeNow, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // вернет "key" is required
		return
	}
	fmt.Fprintf(w, nextDate) // ответ на запрос
}
