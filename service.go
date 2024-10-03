package mysql

import (
	"encoding/json"
	"reflect"
	"time"
)

type IdReq struct {
	Id json.Number `json:"id" validate:"required"`
}

type FilterReq struct {
	Filters map[string]interface{}
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

type BaseService[T any] struct {
	DB       *DB
	GetTitle func(*T) string
	Pk       string
	model    *T
}

func (b *BaseService[T]) GetPk() (string, error) {
	if b.Pk != "" {
		return b.Pk, nil
	}
	t := b.GetModel()
	pkField := GetPKField(t)
	pkName := pkField.Name
	if pkName == "" {
		return "", WithStack(ErrorPrimaryKeyUnset)
	}
	b.Pk = pkName
	return b.Pk, nil
}

func (b *BaseService[T]) GetModel() *T {
	if b.model == nil {
		b.model = b.NewModel()
	}
	return b.model
}

func (b *BaseService[T]) NewModel() *T {
	return new(T)
}

func (b *BaseService[T]) NewModelList() *[]*T {
	return new([]*T)
	//model := b.GetModel()
	//list := reflect.New(reflect.SliceOf(reflect.TypeOf(model)))
}

func (b *BaseService[T]) NewModelWithId(id interface{}) (*T, error) {
	pk, err := b.GetPk()
	if err != nil {
		return nil, err
	}
	t := new(T)
	modelV := reflect.ValueOf(t).Elem()
	modelFieldV := modelV.FieldByName(pk)
	modelFieldV.Set(reflect.ValueOf(id).Convert(modelFieldV.Type()))
	return t, nil
}

func (b *BaseService[T]) FindTitle(id interface{}, filters ...map[string]interface{}) (*TitleRes, error) {
	model, err := b.FindById(id, filters...)
	if err != nil {
		return nil, err
	}
	return &TitleRes{Title: b.GetTitle(model)}, nil
}

func (b *BaseService[T]) FindById(id interface{}, filters ...map[string]interface{}) (*T, error) {
	model, err := b.NewModelWithId(id)
	if err != nil {
		return nil, err
	}
	if err := b.DB.FindById(
		model,
		b.DB.WithFilters(filters...),
	); err != nil {
		return nil, err
	}
	return model, nil
}


func (b *BaseService[T]) FindByIdWithOpts(id interface{}, opts ...Option) (*T, error) {
	model, err := b.NewModelWithId(id)
	if err != nil {
		return nil, err
	}
	if err := b.DB.FindById(
		model,
		opts...,
	); err != nil {
		return nil, err
	}
	return model, nil
}

func (b *BaseService[T]) FindOne(filters ...map[string]interface{}) (*T, error) {
	model := b.NewModel()
	if err := b.DB.FindOne(
		model,
		b.DB.WithFilters(filters...),
	); err != nil {
		return nil, err
	}
	return model, nil
}

func (b *BaseService[T]) FindOneWithOpts(opts ...Option) (*T, error) {
	model := b.NewModel()
	if err := b.DB.FindOne(
		model,
		opts...,
	); err != nil {
		return nil, err
	}
	return model, nil
}

func (b *BaseService[T]) FindAll(filters ...map[string]interface{}) (*[]*T, error) {
	list := b.NewModelList()
	err := b.DB.FindAll(
		list,
		b.DB.WithFilters(filters...),
	)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (b *BaseService[T]) FindAllWithOpts(opts ...Option) (*[]*T, error) {
	list := b.NewModelList()
	err := b.DB.FindAll(
		list,
		opts...,
	)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (b *BaseService[T]) FindPage(pageable *Pageable, filters ...map[string]interface{}) (*PageRes[T], error) {
	list := b.NewModelList()
	total, err := b.DB.FindPage(
		list,
		b.DB.WithFilters(filters...),
		b.DB.WithPageable(pageable),
	)
	if err != nil {
		return nil, err
	}
	return &PageRes[T]{Total: total, List: list}, nil
}

func (b *BaseService[T]) FindPageWithOpts(opts ...Option) (*PageRes[T], error) {
	list := b.NewModelList()
	total, err := b.DB.FindPage(
		list,
		opts...,
	)
	if err != nil {
		return nil, err
	}
	return &PageRes[T]{Total: total, List: list}, nil
}

func (b *BaseService[T]) Create(value *T) (*T, error) {
	err := b.DB.Create(value)
	return value, err
}

func (b *BaseService[T]) CreateWithUserId(value *T, userId int) (*T, error) {
	SetCreatedBy(value, userId)
	return b.Create(value)
}

func (b *BaseService[T]) Update(value *T, filters ...map[string]interface{}) (*T, error) {
	if err := b.DB.UpdateById(
		value,
		value,
		b.DB.WithOmit("createdAt", "createdBy"),
		b.DB.WithFilters(filters...),
	); err != nil {
		return nil, err
	}

	if err := b.DB.FindById(
		value,
	); err != nil {
		return nil, err
	}
	return value, nil
}

func (b *BaseService[T]) UpdateWithUserId(value *T, userId int, filters ...map[string]interface{}) (*T, error) {
	SetUpdatedBy(value, userId)
	return b.Update(value, filters...)
}

func (b *BaseService[T]) UpdateById(id interface{}, value interface{}, filters ...map[string]interface{}) (*T, error) {
	model, err := b.NewModelWithId(id)
	if err != nil {
		return nil, err
	}
	if err := b.DB.UpdateById(
		model,
		value,
		b.DB.WithOmit("createdAt", "createdBy"),
		b.DB.WithFilters(filters...),
	); err != nil {
		return nil, err
	}
	if err := b.DB.FindById(
		model,
	); err != nil {
		return nil, err
	}
	return model, err
}

func (b *BaseService[T]) UpdateByIdWithUserId(id interface{}, value interface{}, userId int, filters ...map[string]interface{}) (*T, error) {
	SetUpdatedBy(value, userId)
	return b.UpdateById(id, value, filters...)
}

func (b *BaseService[T]) UpdateOneOrCreate(value *T, filters ...map[string]interface{}) (*T, error) {
	if err := b.DB.UpdateOne(
		value,
		value,
		b.DB.WithOmit("createdAt", "createdBy"),
		b.DB.WithFilters(filters...),
	); err != nil {
		if IsRecordNotFoundError(err) {
			return b.Create(value)
		}
		return nil, err
	}
	return value, nil
}

func (b *BaseService[T]) UpdateOneOrCreateWithUserId(value *T, userId int, filters ...map[string]interface{}) (*T, error) {
	SetUpdatedBy(value, userId)
	if err := b.DB.UpdateOne(
		value,
		value,
		b.DB.WithOmit("createdAt", "createdBy"),
		b.DB.WithFilters(filters...),
	); err != nil {
		if IsRecordNotFoundError(err) {
			SetCreatedBy(value, userId)
			return b.Create(value)
		}
		return nil, err
	}
	return value, nil
}

func (b *BaseService[T]) Remove(value *T, filters ...map[string]interface{}) error {
	SetDeleted(value)
	return b.DB.UpdateById(
		value,
		value,
		b.DB.WithAttend("updatedAt", "updatedBy", "deleted"),
		b.DB.WithFilters(filters...),
	)
}

func (b *BaseService[T]) RemoveWithUserId(value *T, userId int, filters ...map[string]interface{}) error {
	SetUpdatedBy(value, userId)
	return b.Remove(value, filters...)
}

func (b *BaseService[T]) RemoveById(id interface{}, filters ...map[string]interface{}) error {
	model, err := b.NewModelWithId(id)
	if err != nil {
		return err
	}
	return b.Remove(model, filters...)
}

func (b *BaseService[T]) RemoveByIdWithUserId(id interface{}, userId int, filters ...map[string]interface{}) error {
	model, err := b.NewModelWithId(id)
	if err != nil {
		return err
	}
	SetUpdatedBy(model, userId)
	return b.Remove(model, filters...)
}

func SetCreatedBy(value interface{}, userId int) {
	valueV := reflect.ValueOf(value)
	if valueV.Kind() == reflect.Ptr {
		valueV = valueV.Elem()
	}
	createdByFieldV := valueV.FieldByName("CreatedBy")
	if createdByFieldV.IsValid() && createdByFieldV.CanSet() {
		createdByFieldV.Set(reflect.ValueOf(NewUserId(userId)))
	}
	updatedByFieldV := valueV.FieldByName("UpdatedBy")
	if updatedByFieldV.IsValid() && updatedByFieldV.CanSet() {
		updatedByFieldV.Set(reflect.ValueOf(NewUserId(userId)))
	}
}

func SetUpdatedBy(value interface{}, userId int) {
	valueV := reflect.ValueOf(value)
	if valueV.Kind() == reflect.Ptr {
		valueV = valueV.Elem()
	}

	updatedByFieldV := valueV.FieldByName("UpdatedBy")
	if updatedByFieldV.IsValid() && updatedByFieldV.CanSet() {
		updatedByFieldV.Set(reflect.ValueOf(NewUserId(userId)))
	}
}

func SetDeleted(value interface{}) {
	valueV := reflect.ValueOf(value)
	if valueV.Kind() == reflect.Ptr {
		valueV = valueV.Elem()
	}

	updatedByFieldV := valueV.FieldByName("Deleted")
	if updatedByFieldV.IsValid() && updatedByFieldV.CanSet() {
		updatedByFieldV.Set(reflect.ValueOf(time.Now().Unix()).Convert(updatedByFieldV.Type()))
	}
}
