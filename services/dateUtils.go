package services

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
		return "", fmt.Errorf("invalid date format: %v", err) // возврат ошибки если формат даты неверный
	}

	if repeat == "" { // если правило повторения пустое то вернем ошибку
		return "", errors.New("repeat is required")
	}

	repeatParts := strings.Split(repeat, " ") // разделил правило повторения на части
	if len(repeatParts) == 0 {
		return "", errors.New("invalid repeat rule format")
	}
	rule := repeatParts[0]
	switch rule {
	case "d": // кол-во дней
		if len(repeatParts) != 2 {
			return "", errors.New("invalid repeat rule format for d")
		}
		days, err := strconv.Atoi(repeatParts[1]) // конвертация строки в число дней
		if err != nil || days < 1 || days > 400 { // проверка соответсвия
			return "", errors.New("invalid number of days")
		}
		if dateStr == now.Format(DefaultFormatDate) {
			date = now
		}
		for {
			date = date.AddDate(0, 0, days) // добавление дней к дате now
			if date.After(now) {            // если новая дата после текущей, возвращаем её
				return date.Format(DefaultFormatDate), nil
			}
		}
	case "y": // ежегодно
		if len(repeatParts) != 1 {
			return "", errors.New("invalid repeat rule format for y")
		}
		for {
			date = date.AddDate(1, 0, 0) // добавление гола
			if date.After(now) {
				if date.Month() == time.February && date.Day() == 29 && !isLeapYear(date.Year()) {
					date = date.AddDate(0, 0, 1) // если следующий год не високосный то переходим на следующий день
				}
				return date.Format(DefaultFormatDate), nil
			}
		}
	case "": // пустое
		if len(repeatParts) != 2 {
		}
		return "", errors.New("invalid repeat rule format")
	default:
		return "", errors.New("unsupported repeat rule") // ошибка
	}
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) { // GET-обработчик api/nextdate
	now := r.FormValue("now")
	date := r.FormValue("date")
	repeat := r.FormValue("repeat")

	timeNow, err := time.Parse(DefaultFormatDate, now)
	if err != nil {
		http.Error(w, "invalid 'now' date format", http.StatusBadRequest) // вернет Invalid 'now' date format если не верный формат
		return
	}
	nextDate, err := NextDate(timeNow, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // вернет "key" is required
		return
	}
	fmt.Fprintf(w, nextDate) // ответ на запрос
}
