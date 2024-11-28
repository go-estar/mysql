package mysql

import (
	"database/sql/driver"
	"errors"
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
)

type IdReq struct {
	Id int `json:"id" validate:"required"`
}
type Id64Req struct {
	Id int64 `json:"id,string" validate:"required"`
}

type IdListReq struct {
	IdList []int `json:"idList" validate:"required"`
}

type Id64ListReq[T any] struct {
	IdList []int64 `json:"idList" validate:"required"`
}

type FilterReq struct {
	Filters map[string]interface{} `json:"filters"`
}

type Pageable struct {
	Page int    `json:"page"`
	Size int    `json:"size"`
	Sort string `json:"sort"`
}

type PageReq struct {
	*Pageable `validate:"required"`
	Filters   map[string]interface{} `json:"filters"`
}

type PageRes[T any] struct {
	List  *[]*T `json:"list"`
	Total int   `json:"total"`
}

type TitleRes struct {
	Title string `json:"title"`
}

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
		return "", nil
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
		return "", nil
	}
	return strings.Join(a, ","), nil
}

type DecimalArray []decimal.Decimal

func (a *DecimalArray) Scan(src any) error {
	bytes, ok := src.([]byte)
	if !ok {
		return errors.New("src value cannot cast to []byte")
	}
	var arr = make([]decimal.Decimal, 0)
	if string(bytes) == "" {
		*a = arr
		return nil
	}
	for _, val := range strings.Split(string(bytes), ",") {
		i, err := decimal.NewFromString(val)
		if err != nil {
			return err
		}
		arr = append(arr, i)
	}
	*a = arr
	return nil
}

func (a DecimalArray) Value() (driver.Value, error) {
	if len(a) == 0 {
		return "", nil
	}
	val := ""
	for _, v := range a {
		if val != "" {
			val += ","
		}
		val += v.String()
	}
	return val, nil
}
