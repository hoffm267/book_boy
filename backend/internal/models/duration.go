package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CustomDuration struct {
	time.Duration
}

func (d *CustomDuration) MarshalJSON() ([]byte, error) {
	if d == nil {
		return []byte("null"), nil
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return []byte(fmt.Sprintf("\"%02d:%02d:%02d\"", h, m, s)), nil
}

func (d *CustomDuration) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}
	parsed, err := time.ParseDuration(parseHMS(str))
	if err != nil {
		return err
	}
	d.Duration = parsed
	return nil
}

func (d *CustomDuration) Scan(value interface{}) error {
	if value == nil {
		d.Duration = 0
		return nil
	}
	parsed, err := time.ParseDuration(parseHMS(string(value.([]uint8))))
	if err != nil {
		return err
	}
	d.Duration = parsed
	return nil
}

func (d CustomDuration) Value() (driver.Value, error) {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s), nil
}

func parseHMS(hms string) string {
	hms = strings.TrimSpace(hms)
	parts := strings.Split(hms, ":")

	var h, m, s int
	var err error

	switch len(parts) {
	case 2:
		h, err = strconv.Atoi(parts[0])
		if err != nil {
			return "0s"
		}
		m, err = strconv.Atoi(parts[1])
		if err != nil {
			return "0s"
		}
		s = 0
	case 3:
		h, err = strconv.Atoi(parts[0])
		if err != nil {
			return "0s"
		}
		m, err = strconv.Atoi(parts[1])
		if err != nil {
			return "0s"
		}
		s, err = strconv.Atoi(parts[2])
		if err != nil {
			return "0s"
		}
	default:
		return "0s"
	}

	if h < 0 || m < 0 || m >= 60 || s < 0 || s >= 60 {
		return "0s"
	}

	return fmt.Sprintf("%dh%dm%ds", h, m, s)
}
