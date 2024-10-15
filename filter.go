package mysql

import (
	"fmt"
	"github.com/go-estar/types/fieldUtil"
	"gorm.io/gorm"
	"reflect"
	"strings"
)

const (
	SymbolEquals            = "eq"
	SymbolNotEquals         = "ne"
	SymbolGreatThanOrEquals = "gte"
	SymbolGreatThan         = "gt"
	SymbolLessThanOrEquals  = "lte"
	SymbolLessThan          = "lt"
	SymbolLike              = "like"
	SymbolNotLike           = "notLike"
	SymbolIn                = "in"
	SymbolNotIn             = "notIn"
	SymbolFunc              = "func"
)

var Symbol = map[string]string{
	SymbolEquals:            "=",
	SymbolNotEquals:         "!=",
	SymbolGreatThanOrEquals: ">=",
	SymbolGreatThan:         ">",
	SymbolLessThanOrEquals:  "<=",
	SymbolLessThan:          "<",
	SymbolLike:              "like",
	SymbolNotLike:           "not like",
	SymbolIn:                "in",
	SymbolNotIn:             "not in",
}

type FilterKey struct {
	Column          string
	Operator        string
	IgnoreZeroValue bool // Prefix"?"
}

func keyFormat(key string) *FilterKey {
	var filterKey FilterKey

	if strings.HasPrefix(key, "?") {
		filterKey.IgnoreZeroValue = true
		key = strings.Replace(key, "?", "", 1)
	}

	nIndex := strings.Index(key, "#")
	if nIndex != -1 {
		key = key[:nIndex]
	}

	oIndex := strings.Index(key, "$")
	if oIndex == -1 {
		filterKey.Column = key
	} else {
		filterKey.Column = key[:oIndex]
		filterKey.Operator = key[oIndex+1:]
	}
	return &filterKey
}

func Filters(query *gorm.DB, filters map[string]interface{}) *gorm.DB {
	for key, val := range filters {
		if val == nil {
			continue
		}

		filterKey := keyFormat(key)
		if filterKey.IgnoreZeroValue && fieldUtil.IsEmpty(val) {
			continue
		}

		switch filterKey.Operator {
		case SymbolEquals:
			query = query.Where(filterKey.Column + " = ?", val)
		case SymbolNotEquals:
			query = query.Where(filterKey.Column + " != ?", val)
		case SymbolGreatThanOrEquals:
			query = query.Where(filterKey.Column + " >= ?", val)
		case SymbolGreatThan:
			query = query.Where(filterKey.Column + " > ?", val)
		case SymbolLessThanOrEquals:
			query = query.Where(filterKey.Column + " <= ?", val)
		case SymbolLessThan:
			query = query.Where(filterKey.Column + " < ?", val)
		case SymbolLike:
			query = query.Where(filterKey.Column + " like ?", fmt.Sprintf("%%%s%%", val))
		case SymbolNotLike:
			query = query.Where(filterKey.Column + " not like ?", fmt.Sprintf("%%%s%%", val))
		case SymbolIn:
			query = query.Where(filterKey.Column + " in (?)", val)
		case SymbolNotIn:
			query = query.Where(filterKey.Column + " not in (?)", val)
		case SymbolFunc:
			fn, ok := val.(func(db2 *gorm.DB))
			if ok {
				fn(query)
			}
		default:
			if reflect.ValueOf(val).Kind() == reflect.Slice {
				query = query.Where(filterKey.Column + " in (?)", val)
			} else {
				query = query.Where(filterKey.Column + " = ?", val)
			}
		}
	}
	return query
}
