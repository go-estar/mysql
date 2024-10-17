package mysql

import (
	"encoding/json"
	"errors"
	"reflect"
	"time"
)

type IdReq struct {
	Id json.Number `json:"id" validate:"required"`
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

type Service[T any] struct {
	DB         *DB
	Pk         string
	model      *T
	TitleQuery string
}

func (b *Service[T]) GetPk() (string, error) {
	if b.Pk != "" {
		return b.Pk, nil
	}
	t := b.GetModel()
	pkField := getPKField(t)
	pkName := pkField.Name
	if pkName == "" {
		return "", WithStack(ErrorPrimaryKeyUnset)
	}
	b.Pk = pkName
	return b.Pk, nil
}

func (b *Service[T]) GetModel() *T {
	if b.model == nil {
		b.model = b.NewModel()
	}
	return b.model
}

func (b *Service[T]) NewModel() *T {
	return new(T)
}

func (b *Service[T]) NewModelList() *[]*T {
	return new([]*T)
	//model := b.GetModel()
	//list := reflect.New(reflect.SliceOf(reflect.TypeOf(model)))
}

func (b *Service[T]) NewModelWithId(id interface{}) (*T, error) {
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

func (b *Service[T]) FindTitle(id interface{}, opts ...Option) (*TitleRes, error) {
	if b.TitleQuery == "" {
		return nil, errors.New("TitleQuery not configured")
	}
	model, err := b.NewModelWithId(id)
	if err != nil {
		return nil, err
	}
	title := &TitleRes{}
	if err := b.DB.FindById(
		model,
		append(
			opts,
			WithSelect(b.TitleQuery+" as title"),
			WithWhere(getPKName(b.DB.Config, model)+"=?", id),
			WithDest(title),
		)...,
	); err != nil {
		return nil, err
	}
	return title, nil
}

func (b *Service[T]) FindById(id interface{}, opts ...Option) (*T, error) {
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

func (b *Service[T]) FindOne(opts ...Option) (*T, error) {
	model := b.NewModel()
	if err := b.DB.FindOne(
		model,
		opts...,
	); err != nil {
		return nil, err
	}
	return model, nil
}

func (b *Service[T]) FindAll(opts ...Option) (*[]*T, error) {
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

func (b *Service[T]) FindPage(opts ...Option) (*PageRes[T], error) {
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

func (b *Service[T]) Create(value *T, opts ...Option) (*T, error) {
	err := b.DB.Create(value, opts...)
	return value, err
}

func (b *Service[T]) CreateWithUserId(value *T, userId int, opts ...Option) (*T, error) {
	SetCreatedBy(value, userId)
	return b.Create(value, opts...)
}

func (b *Service[T]) Update(value *T, opts ...Option) (*T, error) {
	if err := b.DB.UpdateById(
		value,
		value,
		opts...,
	); err != nil {
		return nil, err
	}
	return value, nil
}

func (b *Service[T]) UpdateWithUserId(value *T, userId int, opts ...Option) (*T, error) {
	SetUpdatedBy(value, userId)
	return b.Update(value, opts...)
}

func (b *Service[T]) UpdateById(id interface{}, value interface{}, opts ...Option) error {
	model, err := b.NewModelWithId(id)
	if err != nil {
		return err
	}
	if err := b.DB.UpdateById(
		model,
		value,
		opts...,
	); err != nil {
		return err
	}
	return nil
}

func (b *Service[T]) UpdateByIdWithUserId(id interface{}, value interface{}, userId int, opts ...Option) error {
	SetUpdatedBy(value, userId)
	return b.UpdateById(id, value, opts...)
}

func (b *Service[T]) UpdateOneOrCreate(value *T, opts ...Option) (*T, error) {
	if err := b.DB.UpdateOne(
		value,
		value,
		opts...,
	); err != nil {
		if IsRecordNotFoundError(err) {
			return b.Create(value)
		}
		return nil, err
	}
	return value, nil
}

func (b *Service[T]) UpdateOneOrCreateWithUserId(value *T, userId int, opts ...Option) (*T, error) {
	SetUpdatedBy(value, userId)
	if err := b.DB.UpdateOne(
		value,
		value,
		opts...,
	); err != nil {
		if IsRecordNotFoundError(err) {
			SetCreatedBy(value, userId)
			return b.Create(value)
		}
		return nil, err
	}
	return value, nil
}

func (b *Service[T]) Remove(value *T, opts ...Option) error {
	SetDeleted(value)
	return b.DB.UpdateById(
		value,
		value,
		append(opts, WithAttend("updatedAt", "updatedBy", "deleted"))...,
	)
}

func (b *Service[T]) RemoveWithUserId(value *T, userId int, opts ...Option) error {
	SetUpdatedBy(value, userId)
	return b.Remove(value, opts...)
}

func (b *Service[T]) RemoveById(id interface{}, opts ...Option) error {
	model, err := b.NewModelWithId(id)
	if err != nil {
		return err
	}
	return b.Remove(model, opts...)
}

func (b *Service[T]) RemoveByIdWithUserId(id interface{}, userId int, opts ...Option) error {
	model, err := b.NewModelWithId(id)
	if err != nil {
		return err
	}
	SetUpdatedBy(model, userId)
	return b.Remove(model, opts...)
}

func (b *Service[T]) DeleteById(id interface{}, opts ...Option) error {
	model, err := b.NewModelWithId(id)
	if err != nil {
		return err
	}
	return b.DB.DeleteById(model, opts...)
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
