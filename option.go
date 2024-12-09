package mysql

import (
	"gorm.io/gorm"
	"reflect"
)

type Option func(*QueryOption)

type Options struct {
	list *[]Option
}

func NewOptions(opts ...Option) Options {
	return Options{
		list: &opts,
	}
}

func (o Options) Append(opts ...Option) []Option {
	*o.list = append(*o.list,opts...)
	return *o.list
}

func (o Options) Appends(opts []Option) []Option {
	*o.list = append(*o.list,opts...)
	return *o.list
}

func (o Options) List() []Option {
	return *o.list
}


type Select struct {
	Query string
	Args  []interface{}
}
type QueryOption struct {
	DB               *gorm.DB
	Table            string
	PrimaryKey       string
	Updates          map[string]interface{}
	Select           Select
	Omit             []string
	IgnoreOmit       bool
	Attend           []string
	Join             [][]interface{}
	Where            [][]interface{}
	Or               [][]interface{}
	Filters          map[string]interface{}
	Group            string
	Limit            int
	Offset           int
	Pageable         *Pageable
	Sort             []string
	Pluck            []interface{}
	Dest             interface{}
	First            bool
	Last             bool
	WithDeleted      bool
	IgnoreNotFound   bool
	MustAffected     bool
	ErrorNotFound    error
	ErrorNotAffected error
	ErrorNotSingle   error
}

func WithDB(val *gorm.DB) Option {
	return func(opts *QueryOption) {
		opts.DB = val
	}
}
func WithTable(val string) Option {
	return func(opts *QueryOption) {
		opts.Table = val
	}
}
func WithPrimaryKey(val string) Option {
	return func(opts *QueryOption) {
		opts.PrimaryKey = val
	}
}
func WithSelect(query string, args ...interface{}) Option {
	return func(opts *QueryOption) {
		if opts.Select.Query != "" {
			opts.Select.Query += ","
		}
		opts.Select.Query += query
		opts.Select.Args = append(opts.Select.Args, args...)
	}
}
func WithSelectCondition(condition bool, query string, args ...interface{}) Option {
	return func(opts *QueryOption) {
		if condition {
			if opts.Select.Query != "" {
				opts.Select.Query += ","
			}
			opts.Select.Query += query
			opts.Select.Args = append(opts.Select.Args, args...)
		}
	}
}
func WithIgnoreOmit() Option {
	return func(opts *QueryOption) {
		opts.IgnoreOmit = true
	}
}
func WithOmit(val ...string) Option {
	return func(opts *QueryOption) {
		opts.Omit = append(opts.Omit, val...)
	}
}
func WithAttend(val ...string) Option {
	return func(opts *QueryOption) {
		opts.Attend = append(opts.Attend, val...)
	}
}
func WithUpdates(val map[string]interface{}) Option {
	return func(opts *QueryOption) {
		if opts.Updates == nil {
			opts.Updates = map[string]interface{}{}
		}
		for k, v := range val {
			opts.Updates[k] = v
		}
	}
}
func WithJoin(query string, val ...interface{}) Option {
	return func(opts *QueryOption) {
		opts.Join = append(opts.Join, []interface{}{query, val})
	}
}
func WithJoinCondition(condition bool, query string, val ...interface{}) Option {
	return func(opts *QueryOption) {
		if condition {
			opts.Join = append(opts.Join, []interface{}{query, val})
		}
	}
}
func WithWhere(val ...interface{}) Option {
	return func(opts *QueryOption) {
		opts.Where = append(opts.Where, val)
	}
}
func WithWhereCondition(condition bool, val ...interface{}) Option {
	return func(opts *QueryOption) {
		if condition {
			opts.Where = append(opts.Where, val)
		}
	}
}
func WithWheres(val ...[]interface{}) Option {
	return func(opts *QueryOption) {
		for _, val := range val {
			opts.Where = append(opts.Where, val)
		}
	}
}
func WithWheresCondition(condition bool, val ...[]interface{}) Option {
	return func(opts *QueryOption) {
		if condition {
			for _, val := range val {
				opts.Where = append(opts.Where, val)
			}
		}
	}
}
func WithOr(val ...interface{}) Option {
	return func(opts *QueryOption) {
		opts.Or = append(opts.Or, val)
	}
}
func WithOrCondition(condition bool, val ...interface{}) Option {
	return func(opts *QueryOption) {
		if condition {
			opts.Or = append(opts.Or, val)
		}
	}
}
func WithFilters(val ...map[string]interface{}) Option {
	return func(opts *QueryOption) {
		if opts.Filters == nil {
			opts.Filters = map[string]interface{}{}
		}
		for _, m := range val {
			for k, v := range m {
				opts.Filters[k] = v
			}
		}
	}
}
func WithGroup(val string) Option {
	return func(opts *QueryOption) {
		opts.Group = val
	}
}
func WithLimit(val int) Option {
	return func(opts *QueryOption) {
		opts.Limit = val
	}
}
func WithOffset(val int) Option {
	return func(opts *QueryOption) {
		opts.Offset = val
	}
}
func WithPage(page int, size int, sort string) Option {
	return func(opts *QueryOption) {
		opts.Pageable = &Pageable{
			Page: page, Size: size, Sort: sort,
		}
	}
}
func WithPageable(val *Pageable) Option {
	return func(opts *QueryOption) {
		opts.Pageable = val
	}
}
func WithSort(val ...string) Option {
	return func(opts *QueryOption) {
		opts.Sort = append(opts.Sort, val...)
	}
}
func WithPluck(column string, val interface{}) Option {
	return func(opts *QueryOption) {
		opts.Pluck = []interface{}{column, val}
	}
}
func WithDest(dest interface{}) Option {
	return func(opts *QueryOption) {
		opts.Dest = dest
	}
}
func WithFirst() Option {
	return func(opts *QueryOption) {
		opts.First = true
	}
}
func WithLast() Option {
	return func(opts *QueryOption) {
		opts.Last = true
	}
}
func WithDeleted() Option {
	return func(opts *QueryOption) {
		opts.WithDeleted = true
	}
}
func WithIgnoreNotFound() Option {
	return func(opts *QueryOption) {
		opts.IgnoreNotFound = true
	}
}
func WithMustAffected() Option {
	return func(opts *QueryOption) {
		opts.MustAffected = true
	}
}
func WithErrorNotFound(val error) Option {
	return func(opts *QueryOption) {
		opts.ErrorNotFound = val
	}
}
func WithErrorNotAffected(val error) Option {
	return func(opts *QueryOption) {
		opts.ErrorNotAffected = val
	}
}
func WithErrorNotSingle(val error) Option {
	return func(opts *QueryOption) {
		opts.ErrorNotSingle = val
	}
}

func (db *DB) queryBuilder(model interface{}, opts ...Option) (*gorm.DB, *QueryOption) {
	queryOption := &QueryOption{}
	for _, apply := range opts {
		if apply != nil {
			apply(queryOption)
		}
	}

	var query *gorm.DB
	if queryOption.DB != nil {
		query = queryOption.DB
	} else {
		query = db.DB
	}

	if queryOption.Table != "" {
		query = query.Table(queryOption.Table)
	} else if model != nil {
		query = query.Model(model)
	}

	if len(queryOption.Select.Query) > 0 {
		query = query.Select(queryOption.Select.Query, queryOption.Select.Args...)
	}

	if !queryOption.IgnoreOmit && len(queryOption.Omit) > 0 {
		query = query.Omit(queryOption.Omit...)
	}

	if len(queryOption.Join) > 0 {
		for _, joins := range queryOption.Join {
			if joins == nil || joins[0] == nil {
				continue
			}
			if len(joins) == 1 {
				query = query.Joins(joins[0].(string))
			} else {
				query = query.Joins(joins[0].(string), joins[1:]...)
			}
		}
	}

	if queryOption.Group != "" {
		query = query.Group(queryOption.Group)
	}

	if len(queryOption.Where) > 0 {
		for _, where := range queryOption.Where {
			if where == nil || where[0] == nil {
				continue
			}
			if len(where) == 1 {
				//query = query.Where(wheres[0])
				if (reflect.TypeOf(where[0]).Kind() == reflect.Ptr && reflect.TypeOf(where[0]).Elem().Kind() == reflect.Struct) ||
					reflect.TypeOf(where[0]).Kind() == reflect.Struct {
					query = query.Where(structToMap(query.Config,where[0]))
				} else {
					query = query.Where(where[0])
				}
			} else {
				query = query.Where(where[0], where[1:]...)
			}
		}
	}

	if len(queryOption.Or) > 0 {
		for _, or := range queryOption.Or {
			if or == nil || or[0] == nil {
				continue
			}
			if len(or) == 1 {
				//query = query.Or(wheres[0])
				if (reflect.TypeOf(or[0]).Kind() == reflect.Ptr && reflect.TypeOf(or[0]).Elem().Kind() == reflect.Struct) ||
					reflect.TypeOf(or[0]).Kind() == reflect.Struct {
					query = query.Or(structToMap(query.Config,or[0]))
				} else {
					query = query.Or(or[0])
				}
			} else {
				query = query.Or(or[0], or[1:]...)
			}
		}
	}

	if len(queryOption.Sort) != 0 {
		for _, sort := range queryOption.Sort {
			if queryOption.Sort[0] != "" {
				query = query.Order(sort)
			}
		}
	}

	if queryOption.Pageable != nil && queryOption.Pageable.Size > 0 {
		query = query.Limit(queryOption.Pageable.Size).Offset((queryOption.Pageable.Page - 1) * queryOption.Pageable.Size)
		if queryOption.Pageable.Sort != "" {
			query = query.Order(queryOption.Pageable.Sort)
		}
	}
	if queryOption.Limit != 0 {
		query = query.Limit(queryOption.Limit)
	}
	if queryOption.Offset != 0 {
		query = query.Offset(queryOption.Offset)
	}

	if !queryOption.WithDeleted && model != nil {
		modelT := reflect.TypeOf(model)
		if modelT.Kind() == reflect.Ptr {
			modelT = modelT.Elem()
		}
		if _, ok := modelT.FieldByName("Deleted"); ok {
			tableName := db.Config.NamingStrategy.TableName(modelT.Name())
			if v := getTableName(model); v != "" {
				tableName = v
			}
			query = query.Where(tableName + ".deleted = 0")
		}
	}

	if queryOption.Filters != nil {
		query = Filters(query, queryOption.Filters)
	}

	return query, queryOption
}

func countBuilder(query *gorm.DB) *gorm.DB {
	return query.Select("*").Limit(-1).Offset(-1)
}

func defaultSort(model interface{}, query *gorm.DB, queryOption *QueryOption) *gorm.DB {
	if (queryOption.Pageable == nil || queryOption.Pageable.Sort == "") && queryOption.Sort == nil {
		if primaryKey := getPKName(query.Config,model); primaryKey != "" {
			query = query.Order(primaryKey + " desc")
		}
	}
	return query
}
