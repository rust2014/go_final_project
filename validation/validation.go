package validation

import (
	"errors"
	"regexp"
)

func ValidateRepeatRule(repeat string) error { // проверяет формат правила повторения
	var (
		dayPattern  = regexp.MustCompile(`^d\s\d+$`)
		yearPattern = regexp.MustCompile(`^y$`)
	)

	if repeat == "" {
		return nil
	}

	if !dayPattern.MatchString(repeat) && !yearPattern.MatchString(repeat) { // проверка на соответствие правилу повторения
		return errors.New("the repetition rule is in the wrong format") // если оба правила не совпадают
	}
	return nil
}
