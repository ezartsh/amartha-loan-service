package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var JakartaTimeLocation *time.Location

func init() {
	if jktLocation, err := time.LoadLocation("Asia/Jakarta"); err == nil {
		JakartaTimeLocation = jktLocation
	} else {
		JakartaTimeLocation = time.Now().Location()
	}
}

type LocalTime struct {
	time.Time
}

func (localTime *LocalTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.ParseInLocation(time.DateOnly, s, time.Local)
	if err != nil {
		t2, err2 := time.ParseInLocation(time.RFC3339, s, time.Local)
		if err2 != nil {
			return errors.New(fmt.Sprint("Invalid Time Format (expected in YYYY-MM-DD format or RFC3339), ", err.Error(), ", ", err2.Error()))
		}
		t = t2
	}
	localTime.Time = t
	return nil
}

func (localTime *LocalTime) Value() (*time.Time, error) {
	if localTime == nil {
		return nil, nil
	}
	return &localTime.Time, nil
}

func (localTime *LocalTime) Scan(src interface{}) (err error) {
	var ok bool
	localTime.Time, ok = src.(time.Time)
	if !ok {
		return errors.New("Incompatible type src to time.Time")
	}

	return nil
}

func (localTime *LocalTime) ScanString(src string) (err error) {
	localTime.Time, err = time.Parse(time.RFC3339, src)
	if err != nil {
		localTime.Time, err = time.Parse("2006-01-02", src)
		if err != nil {
			return err
		}
	}
	return nil
}

type LocalDate struct {
	time.Time
}

func (localDate *LocalDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	s = strings.Split(s, "T")[0]
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return errors.New(fmt.Sprint("Invalid Time Format (expected in YYYY-MM-DD format), value:", s, ", ", err.Error()))
	}
	localDate.Time = t
	return nil
}

func (localDate *LocalDate) Value() (*time.Time, error) {
	if localDate == nil {
		return nil, nil
	}
	return &localDate.Time, nil
}

func (localDate *LocalDate) Scan(src interface{}) (err error) {
	var ok bool
	localDate.Time, ok = src.(time.Time)
	if !ok {
		return errors.New("incompatible type src to time.Time")
	}

	return nil
}

func (localDate *LocalDate) ScanString(src string) (err error) {
	localDate.Time, err = time.Parse(time.DateOnly, src)
	if err != nil {
		return err
	}
	return nil
}

func Delay(seconds uint) {
	time.Sleep(time.Duration(seconds) * time.Second)
}
