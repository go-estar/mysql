package mysql

import (
	"gorm.io/gorm"
	"reflect"
)

//WARNING when update with struct, GORM will only update those fields that with non blank value
//For below Update, nothing will be updated as "", 0, false are blank values of their types

// NOTE When query with struct, GORM will only query with those fields has non-zero value,
// that means if your field’s value is 0, ”, false or other zero values, it won’t be used to build query conditions
func (db *DB) Create(model interface{}, opts ...Option) error {
	if reflect.TypeOf(model).Kind() != reflect.Ptr || reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return WithStack(ErrorModel)
	}
	query, _ := db.queryBuilder(model, opts...)
	if err := query.Create(model).Error; err != nil {
		if db.IsUniqueIndexError(err) {
			return GetUniqueIndexError(model, err)
		}
		return WithStack(err)
	}
	return nil
}

func (db *DB) Count(model interface{}, opts ...Option) (int, error) {
	if reflect.TypeOf(model).Kind() != reflect.Ptr || reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return 0, WithStack(ErrorModel)
	}
	query, _ := db.queryBuilder(model, opts...)
	var count int64 = 0
	if err := countBuilder(query).Count(&count).Error; err != nil {
		return 0, WithStack(err)
	}
	return int(count), nil
}

func (db *DB) FindById(model interface{}, opts ...Option) error {
	if reflect.TypeOf(model).Kind() != reflect.Ptr || reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return WithStack(ErrorModel)
	}

	query, queryOpt := db.queryBuilder(model, opts...)

	if _, err := validatePK(model, queryOpt.PrimaryKey); err != nil {
		return err
	}

	var dest = model
	if queryOpt.Dest != nil {
		dest = queryOpt.Dest
	}
	if err := query.Take(dest).Error; err != nil {
		if db.IsRecordNotFoundError(err) {
			if queryOpt.IgnoreNotFound {
				return nil
			} else {
				if err := queryOpt.ErrorNotFound; err != nil {
					return err
				}
				return getRecordNotFoundError(model)
			}
		} else {
			return WithStack(err)
		}
	}
	return nil
}

func (db *DB) FindOne(model interface{}, opts ...Option) error {
	if reflect.TypeOf(model).Kind() != reflect.Ptr || reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return WithStack(ErrorModel)
	}

	query, queryOpt := db.queryBuilder(model, opts...)

	var dest = model
	if queryOpt.Dest != nil {
		dest = queryOpt.Dest
	}
	list := reflect.New(reflect.SliceOf(reflect.TypeOf(dest).Elem()))
	if err := query.Find(list.Interface()).Error; err != nil {
		return WithStack(err)
	}

	count := list.Elem().Len()
	if count == 0 {
		if queryOpt.IgnoreNotFound {
			return nil
		} else {
			if err := queryOpt.ErrorNotFound; err != nil {
				return err
			}
			return getRecordNotFoundError(model)
		}
	}
	if count > 1 {
		if err := queryOpt.ErrorNotSingle; err != nil {
			return err
		}
		return WithStack(ErrorRecordNotUnique)
	}

	reflect.ValueOf(model).Elem().Set(list.Elem().Index(0))
	return nil
}

func (db *DB) CloneById(model interface{}, opts ...Option) (interface{}, error) {
	if reflect.TypeOf(model).Kind() != reflect.Ptr || reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return nil, WithStack(ErrorModel)
	}
	clone := reflect.New(reflect.TypeOf(model).Elem()).Interface()
	pk, err := validatePK(model)
	if err != nil {
		return nil, err
	}
	if err := setPKValue(clone, pk.Name, pk.Value); err != nil {
		return nil, err
	}
	opts = append(opts, WithIgnoreOmit())
	if err := db.FindById(clone, opts...); err != nil {
		return nil, err
	}
	return clone, nil
}

func (db *DB) CloneOne(model interface{}, opts ...Option) (interface{}, error) {
	if reflect.TypeOf(model).Kind() != reflect.Ptr || reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return nil, WithStack(ErrorModel)
	}
	clone := reflect.New(reflect.TypeOf(model).Elem()).Interface()
	opts = append(opts, WithIgnoreOmit())
	if err := db.FindOne(clone, opts...); err != nil {
		return nil, err
	}
	return clone, nil
}

func (db *DB) delete(model interface{}, query *gorm.DB, queryOpt *QueryOption) error {
	query = query.Delete(model)
	if err := query.Error; err != nil {
		return WithStack(err)
	}
	if query.RowsAffected == 0 && queryOpt.MustAffected {
		if err := queryOpt.ErrorNotAffected; err != nil {
			return err
		}
		return getRecordNotAffectedError(model)
	}
	return nil
}

func (db *DB) DeleteAll(model interface{}, opts ...Option) error {
	if reflect.TypeOf(model).Kind() != reflect.Ptr || reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return WithStack(ErrorModel)
	}
	query, queryOpt := db.queryBuilder(model, opts...)
	return db.delete(model, query, queryOpt)
}

func (db *DB) DeleteById(model interface{}, opts ...Option) error {
	if reflect.TypeOf(model).Kind() != reflect.Ptr || reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return WithStack(ErrorModel)
	}
	query, queryOpt := db.queryBuilder(model, opts...)
	if _, err := validatePK(model, queryOpt.PrimaryKey); err != nil {
		return err
	}
	return db.delete(model, query, queryOpt)
}

func (db *DB) DeleteOne(model interface{}, opts ...Option) error {
	if reflect.TypeOf(model).Kind() != reflect.Ptr || reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return WithStack(ErrorModel)
	}
	query, queryOpt := db.queryBuilder(model, opts...)
	var dest = model
	if queryOpt.Dest != nil {
		dest = queryOpt.Dest
	}
	list := reflect.New(reflect.SliceOf(reflect.TypeOf(dest).Elem()))
	if err := query.Find(list.Interface()).Error; err != nil {
		return WithStack(err)
	}
	count := list.Elem().Len()
	if count == 0 {
		if queryOpt.IgnoreNotFound {
			return nil
		} else {
			if err := queryOpt.ErrorNotFound; err != nil {
				return err
			}
			return getRecordNotFoundError(model)
		}
	}
	if count > 1 {
		if err := queryOpt.ErrorNotSingle; err != nil {
			return err
		}
		return WithStack(ErrorRecordNotUnique)
	}
	return db.delete(model, query, queryOpt)
}

func (db *DB) update(model interface{}, updates interface{}, query *gorm.DB, queryOpt *QueryOption) (int, error) {
	if len(queryOpt.Attend) > 0 {
		attends := make([]interface{}, 0)
		for _, attend := range queryOpt.Attend {
			attends = append(attends, attend)
		}
		if len(attends) == 1 {
			query = query.Select(attends[0])
		} else {
			query = query.Select(attends[0], attends[1:]...)
		}
	}
	query = query.Updates(updates)
	if err := query.Error; err != nil {
		if db.IsUniqueIndexError(err) {
			return 0, GetUniqueIndexError(model, err)
		}
		return 0, WithStack(err)
	}

	if query.RowsAffected == 0 && queryOpt.MustAffected {
		if err := queryOpt.ErrorNotAffected; err != nil {
			return 0, err
		}
		return 0, getRecordNotAffectedError(model)
	}
	return int(query.RowsAffected), nil
}

func (db *DB) UpdateAll(model interface{}, updates interface{}, opts ...Option) (int, error) {
	if reflect.TypeOf(model).Kind() != reflect.Ptr || reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return 0, WithStack(ErrorModel)
	}
	query, queryOpt := db.queryBuilder(model, opts...)
	return db.update(model, updates, query, queryOpt)
}

func (db *DB) UpdateById(model interface{}, values interface{}, opts ...Option) error {
	if reflect.TypeOf(model).Kind() != reflect.Ptr || reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return WithStack(ErrorModel)
	}
	if isStruct(values) {
		_, _, err := db.UpdateByIdReturnChangedValues(model, values, opts...)
		return err
	}
	query, queryOpt := db.queryBuilder(model, opts...)
	if _, err := validatePK(model, queryOpt.PrimaryKey); err != nil {
		return err
	}
	_, err := db.update(model, values, query, queryOpt)
	return err
}

func (db *DB) UpdateByIdReturnChangedValues(model interface{}, values interface{}, opts ...Option) (map[string]interface{}, interface{}, error) {
	if reflect.TypeOf(model).Kind() != reflect.Ptr || reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return nil, nil, WithStack(ErrorModel)
	}
	clone, err := db.CloneById(model, opts...)
	if err != nil {
		return nil, nil, err
	}
	updates, err := getUpdateValue(db.Config, clone, values)
	if err != nil {
		return nil, nil, err
	}
	_, err = db.UpdateAll(model, updates, opts...)
	if err != nil {
		return nil, nil, err
	}
	return updates, clone, nil
}

func (db *DB) UpdateOne(model interface{}, values interface{}, opts ...Option) error {
	if reflect.TypeOf(model).Kind() != reflect.Ptr || reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return WithStack(ErrorModel)
	}
	if isStruct(values) {
		_, _, err := db.UpdateOneReturnChangedValues(model, values, opts...)
		return err
	}
	query, queryOpt := db.queryBuilder(model, opts...)
	list := reflect.New(reflect.SliceOf(reflect.TypeOf(model).Elem()))
	if err := query.Find(list.Interface()).Error; err != nil {
		return WithStack(err)
	}
	count := list.Elem().Len()
	if count == 0 {
		if queryOpt.IgnoreNotFound {
			return nil
		} else {
			if err := queryOpt.ErrorNotFound; err != nil {
				return err
			}
			return getRecordNotFoundError(model)
		}
	}
	if count > 1 {
		if err := queryOpt.ErrorNotSingle; err != nil {
			return err
		}
		return WithStack(ErrorRecordNotUnique)
	}
	_, err := db.update(model, values, query, queryOpt)
	return err
}

func (db *DB) UpdateOneReturnChangedValues(model interface{}, values interface{}, opts ...Option) (map[string]interface{}, interface{}, error) {
	if reflect.TypeOf(model).Kind() != reflect.Ptr || reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return nil, nil, WithStack(ErrorModel)
	}
	clone, err := db.CloneOne(model, opts...)
	if err != nil {
		return nil, nil, err
	}
	updates, err := getUpdateValue(db.Config, clone, values)
	if err != nil {
		return nil, nil, err
	}
	_, err = db.UpdateAll(model, updates, opts...)
	if err != nil {
		return nil, nil, err
	}
	return updates, clone, nil
}

func (db *DB) UpdateOneOrCreate(model interface{}, opts ...Option) error {
	err := db.UpdateOne(model, model, opts...)
	if err != nil {
		if IsRecordNotFoundError(err) {
			return db.Create(model, opts...)
		}
		return err
	}
	return nil
}

func (db *DB) Find(model interface{}, opts ...Option) error {
	if reflect.TypeOf(model).Kind() != reflect.Ptr || reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return WithStack(ErrorModel)
	}
	query, queryOpt := db.queryBuilder(model, opts...)

	var dest = model
	if queryOpt.Dest != nil {
		dest = queryOpt.Dest
	}
	if queryOpt.First {
		query.First(dest)
	} else if queryOpt.Last {
		query.Last(dest)
	} else {
		query.Find(dest)
	}

	if err := query.Error; err != nil {
		if db.IsRecordNotFoundError(err) {
			if queryOpt.IgnoreNotFound {
				return nil
			} else {
				if err := queryOpt.ErrorNotFound; err != nil {
					return err
				}
				return getRecordNotFoundError(model)
			}
		} else {
			return WithStack(err)
		}
	}
	return nil
}

func (db *DB) FindAllWithModel(model interface{}, opts ...Option) (interface{}, error) {
	if reflect.TypeOf(model).Kind() != reflect.Ptr || reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return nil, WithStack(ErrorModel)
	}
	query, queryOpt := db.queryBuilder(model, opts...)
	query = defaultSort(model, query, queryOpt)
	var dest = model
	if queryOpt.Dest != nil {
		dest = queryOpt.Dest
	}
	list := reflect.New(reflect.SliceOf(reflect.TypeOf(dest)))
	if err := query.Find(list.Interface()).Error; err != nil {
		return nil, WithStack(err)
	}
	return list.Interface(), nil
}

func (db *DB) FindPageWithModel(model interface{}, opts ...Option) (interface{}, int, error) {
	if reflect.TypeOf(model).Kind() != reflect.Ptr || reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return nil, 0, WithStack(ErrorModel)
	}
	query, queryOpt := db.queryBuilder(model, opts...)
	query = defaultSort(model, query, queryOpt)
	var dest = model
	if queryOpt.Dest != nil {
		dest = queryOpt.Dest
	}
	list := reflect.New(reflect.SliceOf(reflect.TypeOf(dest)))
	if err := query.Find(list.Interface()).Error; err != nil {
		return nil, 0, WithStack(err)
	}

	var total int64 = 0
	if queryOpt.Pageable != nil {
		if err := countBuilder(query).Count(&total).Error; err != nil {
			return nil, 0, err
		}
	}
	return list.Interface(), int(total), nil
}

func (db *DB) FindAll(list interface{}, opts ...Option) error {
	listT := reflect.TypeOf(list)
	if listT.Kind() != reflect.Ptr || listT.Elem().Kind() != reflect.Slice {
		return WithStack(ErrorModel)
	}
	if !(listT.Elem().Elem().Kind() == reflect.Struct || (listT.Elem().Elem().Kind() == reflect.Ptr && listT.Elem().Elem().Elem().Kind() == reflect.Struct)) {
		return WithStack(ErrorModel)
	}

	elem := listT.Elem().Elem()
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}

	model := reflect.New(elem).Interface()
	query, queryOpt := db.queryBuilder(model, opts...)
	query = defaultSort(model, query, queryOpt)

	//TODO: cache
	//sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
	//	return tx.Find(list)
	//})
	if err := query.Find(list).Error; err != nil {
		return WithStack(err)
	}
	return nil
}

func (db *DB) FindPage(list interface{}, opts ...Option) (int, error) {
	listT := reflect.TypeOf(list)
	if listT.Kind() != reflect.Ptr || listT.Elem().Kind() != reflect.Slice {
		return 0, WithStack(ErrorModel)
	}
	if !(listT.Elem().Elem().Kind() == reflect.Struct || (listT.Elem().Elem().Kind() == reflect.Ptr && listT.Elem().Elem().Elem().Kind() == reflect.Struct)) {
		return 0, WithStack(ErrorModel)
	}

	elem := listT.Elem().Elem()
	if elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}

	model := reflect.New(elem).Interface()
	query, queryOpt := db.queryBuilder(model, opts...)
	query = defaultSort(model, query, queryOpt)

	var total int64 = 0
	if err := query.Find(list).Error; err != nil {
		return 0, WithStack(err)
	}
	if queryOpt.Pageable != nil {
		if err := countBuilder(query).Count(&total).Error; err != nil {
			return 0, err
		}
	}
	return int(total), nil
}

func (db *DB) FindPluck(model interface{}, opts ...Option) error {
	if reflect.TypeOf(model).Kind() != reflect.Ptr || reflect.TypeOf(model).Elem().Kind() != reflect.Struct {
		return WithStack(ErrorModel)
	}
	query, queryOpt := db.queryBuilder(model, opts...)
	if queryOpt.Pluck == nil {
		return WithStack(ErrorPluck)
	}
	if err := query.Pluck(queryOpt.Pluck[0].(string), queryOpt.Pluck[1]).Error; err != nil {
		return WithStack(err)
	}
	return nil
}
