package mysql

import (
	stderrors "errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"io"
	"reflect"
	"runtime"
	"strings"
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

func WithStack(err error, depth ...int) error {
	if err == nil {
		return nil
	}
	var d = 3
	if len(depth) > 0 && depth[0] > 0 {
		d = depth[0]
	}
	return &withStack{
		err,
		callers(3, d),
	}
}

type withStack struct {
	error
	*stack
}

func (w *withStack) Cause() error { return w.error }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *withStack) Unwrap() error { return w.error }

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Cause())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

type stack []uintptr

func (s *stack) Format(st fmt.State, verb rune) {
	switch verb {
	case 'v':
		switch {
		case st.Flag('+'):
			for _, pc := range *s {
				f := errors.Frame(pc)
				fmt.Fprintf(st, "\n%+v", f)
			}
		}
	}
}

func (s *stack) StackTrace() errors.StackTrace {
	f := make([]errors.Frame, len(*s))
	for i := 0; i < len(f); i++ {
		f[i] = errors.Frame((*s)[i])
	}
	return f
}

func callers(skip int, depth int) *stack {
	var s = skip
	for i := skip; i < 15; i++ {
		_, file, _, ok := runtime.Caller(i)
		if ok && (!strings.HasPrefix(file, sourceDir) || strings.HasSuffix(file, "_test.go")) {
			s = i + 1
			break
		}
	}
	pcs := make([]uintptr, depth)
	n := runtime.Callers(s, pcs[:])
	var st stack = pcs[0:n]
	return &st
}
