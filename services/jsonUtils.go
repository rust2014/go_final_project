package services

import (
	"unicode/utf8"
)

func IsValidJSON(s string) bool { // обработчик JSON, чтобы избежать ошибок с некорректными символами
	return utf8.ValidString(s)
}
