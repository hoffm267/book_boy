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
	parts := strings.Split(hms, ":")
	if len(parts) != 3 {
		return "0s"
	}
	h, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	s, _ := strconv.Atoi(parts[2])
	return fmt.Sprintf("%dh%dm%ds", h, m, s)
}
