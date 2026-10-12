package utils

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/google/uuid"
)

func ValidationRequired(fieldName, value string) error {
	if value == "" {
		return fmt.Errorf("%s is required", fieldName)
	}

	return nil
}

func ValidationStringLength(fieldName, value string, min, max int) error {
	l := len(value)

	if l < min || l > max {
		return fmt.Errorf("%s Must be between %d and %d characters", fieldName, min, max)
	}

	return nil
}

func ValidationRegex(fieldName, value string, re *regexp.Regexp, errMessage string) error {
	if !re.MatchString(value) {
		return fmt.Errorf("%s %s", fieldName, errMessage)
	}

	return nil
}

func ValidationPositiveInt(fieldName, value string) (int, error) {
	v, err := strconv.Atoi(value)

	if err != nil {
		return 0, fmt.Errorf("%s must be a number", fieldName)
	}

	if v <= 0 {
		return 0, fmt.Errorf("%s must be a number", fieldName)
	}

	return v, nil
}

func ValidationUUID(fieldName, value string) (uuid.UUID, error) {
	uid, err := uuid.Parse(value)

	if err != nil {
		return uuid.Nil, fmt.Errorf("%s must be a UUID", fieldName)
	}

	return uid, nil
}

func ValidationInList(fieldName, value string, allowed map[string]bool) error {
	if !allowed[value] {
		return fmt.Errorf("%s must be one of the category: %v", fieldName, keys(allowed))
	}
	return nil
}

func keys(mp map[string]bool) []string {
	var k []string

	for v := range mp {
		k = append(k, v)
	}

	return k
}
