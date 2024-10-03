package mysql

import (
	baseError "github.com/go-estar/base-error"
	"regexp"
	"strings"
)

type UniqueIndexError struct {
	Index string
	Error *baseError.Error
}

func NewUniqueIndexError(index string, err *baseError.Error) *UniqueIndexError {
	return &UniqueIndexError{Index: index, Error: err}
}

func GetUniqueIndex(model interface{}) ([]*UniqueIndexError, error) {
	result, err := ModelMethod(model, "UniqueIndexErrors")
	if err != nil {
		return nil, ErrorUniqueIndexUnset
	}
	errs, ok := result.([]*UniqueIndexError)
	if !ok {
		return nil, ErrorUniqueIndexTypeMismatch
	}
	return errs, nil
}

func GetIndexName(msg string) (string, error) {
	exp := regexp.MustCompile(`for key '(.*?)'`)
	result := exp.FindAllStringSubmatch(msg, 1)
	if !(len(result) == 1 && len(result[0]) == 2) {
		return "", ErrorUniqueIndexNameEmpty
	}
	return result[0][1], nil
}

func GetUniqueIndexError(model interface{}, uniqueErr error) error {
	uniqueIndexErrors, err := GetUniqueIndex(model)
	if err != nil {
		return uniqueErr
	}

	if len(uniqueIndexErrors) == 0 {
		return uniqueErr
	}

	getIndexNameFull, err := GetIndexName(uniqueErr.Error())
	if err != nil {
		return uniqueErr
	}
	getIndexNameFull = strings.ToLower(getIndexNameFull)
	getIndexName := getIndexNameFull
	if i := strings.LastIndex(getIndexNameFull, "."); i != -1 {
		getIndexName = getIndexNameFull[i+1:]
	}

	for _, indexError := range uniqueIndexErrors {
		indexNameFull := indexError.Index
		indexNameFull = strings.ToLower(indexNameFull)
		indexName := indexNameFull
		if i := strings.LastIndex(indexNameFull, "."); i != -1 {
			indexName = indexNameFull[i+1:]
		}

		if getIndexNameFull == indexNameFull || getIndexName == indexName {
			if indexError.Error != nil {
				return indexError.Error.SetCause(uniqueErr)
			} else {
				return uniqueErr
			}
		}
	}
	return uniqueErr
}
