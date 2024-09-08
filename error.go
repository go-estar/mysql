package mysql

import (
	stderrors "errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"reflect"
)

var (
	ErrorModel                     = stderrors.New("model is not ptr or mismatch")
	ErrorPrimaryKeyUnset           = stderrors.New("model primary key is undefined")
	ErrorPrimaryKeyInvalid         = stderrors.New("model primary key is invalid")
	ErrorPrimaryKeyEmpty           = stderrors.New("model primary key is empty")
	ErrorUpdateValuePtrUseMap      = stderrors.New("value is ptr type can't be a map")
	ErrorUpdateValuePtrNotStruct   = stderrors.New("value is ptr type must be a struct")
	ErrorUpdateValueNotStructOrMap = stderrors.New("value isn't ptr type must be a struct or map")
	ErrorRecordNotUnique           = stderrors.New("find duplicate record")
	ErrorRecordNotFound            = stderrors.New("record not found")
	ErrorRecordNotAffected         = stderrors.New("record for update not found")
	ErrorPluck                     = stderrors.New("pluck not supplied")
	ErrorUniqueIndexUnset          = stderrors.New("data duplicate(01)")
	ErrorUniqueIndexType           = stderrors.New("data duplicate(02)")
	ErrorUniqueIndexEmpty          = stderrors.New("data duplicate(03)")
	ErrorUniqueIndexMisMatch       = stderrors.New("data duplicate(04)")
	ErrorUniqueIndexNameEmpty      = stderrors.New("data duplicate(05)")
	ErrorUniqueIndexMessageUnset   = stderrors.New("data duplicate(06)")
	ErrorUniqueIndexColumnUnset    = stderrors.New("data duplicate(07)")
)

func (db *DB) IsUniqueIndexError(err error) bool {
	return IsUniqueIndexError(err)
}

func (db *DB) IsNotSingleError(err error) bool {
	return IsNotSingleError(err)
}

func (db *DB) IsRecordNotFoundError(err error) bool {
	return IsRecordNotFoundError(err)
}

func (db *DB) IsRecordNotAffectedError(err error) bool {
	return IsRecordNotAffectedError(err)
}

func IsUniqueIndexError(err error) bool {
	errType := reflect.TypeOf(err).String()
	if errType == "*mysql.MySQLError" && err.(*mysql.MySQLError).Number == 1062 {
		return true
	}
	return false
}

func IsNotSingleError(err error) bool {
	return stderrors.Is(err, ErrorRecordNotUnique)
}

func IsRecordNotFoundError(err error) bool {
	if stderrors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}
	return stderrors.Is(err, ErrorRecordNotFound)
}

func IsRecordNotAffectedError(err error) bool {
	return stderrors.Is(err, ErrorRecordNotAffected)
}
