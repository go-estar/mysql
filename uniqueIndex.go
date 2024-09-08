package mysql

import (
	"github.com/pkg/errors"
	"reflect"
	"regexp"
	"strings"
)

type UniqueIndexError struct {
	IndexName string
	Error     error
}

func GetUniqueIndex(model interface{}) ([]UniqueIndexError, error) {
	result, err := ModelMethod(model, "UniqueIndexErrors")
	if err != nil {
		return nil, errors.WithStack(ErrorUniqueIndexUnset)
	}
	errs, ok := result.([]UniqueIndexError)
	if !ok {
		return nil, errors.WithStack(ErrorUniqueIndexType)
	}
	return errs, nil
}

func GetIndexName(msg string) (string, error) {
	exp := regexp.MustCompile(`for key '(.*?)'`)
	result := exp.FindAllStringSubmatch(msg, 1)
	if !(len(result) == 1 && len(result[0]) == 2) {
		return "", errors.WithStack(ErrorUniqueIndexNameEmpty)
	}
	return result[0][1], nil
}

func GetFieldValue(model interface{}, fieldName string) (interface{}, error) {
	modelV := reflect.ValueOf(model)
	if modelV.Kind() == reflect.Ptr {
		modelV = modelV.Elem()
	}
	field := modelV.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, errors.WithStack(ErrorUniqueIndexColumnUnset)
	}
	return field.Interface(), nil
}

func GetUniqueIndexError(model interface{}, error string) error {

	uniqueIndexErrors, err := GetUniqueIndex(model)
	if err != nil {
		return err
	}

	if len(uniqueIndexErrors) == 0 {
		return errors.WithStack(ErrorUniqueIndexEmpty)
	}

	indexName, err := GetIndexName(error)
	if err != nil {
		return err
	}

	for _, uniqueIndex := range uniqueIndexErrors {
		indexNameArr := strings.Split(uniqueIndex.IndexName, ".")
		indexNameWithoutTable := ""
		if len(indexNameArr) == 2 {
			indexNameWithoutTable = indexNameArr[1]
		}

		if uniqueIndex.IndexName == indexName || indexNameWithoutTable == indexName {
			uniqueIndexErr := uniqueIndex.Error
			if uniqueIndexErr != nil {
				return uniqueIndexErr
			} else {
				return errors.WithStack(ErrorUniqueIndexMessageUnset)
			}
		}
	}
	return errors.WithStack(ErrorUniqueIndexMisMatch)
}
