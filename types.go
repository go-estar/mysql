package mysql

import (
	"database/sql/driver"
	"errors"
	"strconv"
	"strings"
)

type IntArray []int

func (a *IntArray) Scan(src any) error {
	bytes, ok := src.([]byte)
	if !ok {
		return errors.New("src value cannot cast to []byte")
	}
	var arr = make([]int, 0)
	if string(bytes) == "" {
		*a = arr
		return nil
	}
	for _, val := range strings.Split(string(bytes), ",") {
		i, err := strconv.Atoi(val)
		if err != nil {
			return err
		}
		arr = append(arr, i)
	}
	*a = arr
	return nil
}

func (a IntArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	val := ""
	for _, v := range a {
		if val != "" {
			val += ","
		}
		val += strconv.Itoa(v)
	}
	return val, nil
}

type StringArray []string

func (a *StringArray) Scan(src any) error {
	bytes, ok := src.([]byte)
	if !ok {
		return errors.New("src value cannot cast to []byte")
	}
	if string(bytes) == "" {
		*a = make([]string, 0)
		return nil
	}
	*a = strings.Split(string(bytes), ",")
	return nil
}

func (a StringArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return strings.Join(a, ","), nil
}
