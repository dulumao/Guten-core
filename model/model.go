package model

import (
	"github.com/dulumao/Guten-core/database"
	"github.com/dulumao/Guten-utils/paginater"
	"github.com/jinzhu/gorm"
)

type Model struct {
	db *gorm.DB
}

type Func func(m *Model) error
type TransactionFunc Func
type NotFoundCallback func(m *Model)
type Scope func(db *gorm.DB) *gorm.DB
type Scopes []Scope

func (m *Model) DB() *gorm.DB {
	if m.db == nil {
		return database.DB
	}

	return m.db
}

// 接受者不使用指针，防止污染其他正常model调用者
func (m Model) setTx() *Model {
	if m.db == nil {
		m.db = database.DB.Begin()
	} else {
		m.db = m.db.Begin()
	}

	return &m
}

func (m *Model) Transaction(f TransactionFunc) (error, error) {
	var mTx = m.setTx()

	if err := f(mTx); err != nil {
		return mTx.DB().Rollback().Error, err
	}

	return mTx.DB().Commit().Error, nil
}

func (m Model) Scopes(f Func, scopes ...Scope) error {
	var db = m.DB()

	if len(scopes) > 0 {
		for _, scope := range scopes {
			db = db.Scopes(scope)
		}
	}

	m.db = db

	if err := f(&m); err != nil {
		return err
	}

	return nil
}

func (m Model) Sets(f Func, kv ...map[string]interface{}) error {
	var db = m.DB()

	if len(kv) > 0 {
		for k, v := range kv[0] {
			db = db.Set(k, v)
		}
	}

	m.db = db

	if err := f(&m); err != nil {
		return err
	}

	return nil
}

func (m *Model) IsNotFound(f Func, callbacks ...NotFoundCallback) bool {
	if err := f(m); err != nil {
		if gorm.IsRecordNotFoundError(err) {
			if len(callbacks) > 0 {
				callbacks[0](m)
			}

			return true
		}

		panic(err)
	}

	return false
}

func (m *Model) Create(value interface{}) error {
	return m.DB().Create(value).Error
}

func (m *Model) Save(value interface{}) error {
	return m.DB().Save(value).Error
}

func (m *Model) Exists(model interface{}, id interface{}, fields ...string) bool {
	var scopes Scopes
	var itemCount = 0
	var field = "id"

	if len(fields) > 0 {
		field = fields[0]
	}

	scopes = append(scopes, func(db *gorm.DB) *gorm.DB {
		return db.Where("`?` = ?", gorm.Expr(field), id)
	})

	itemCount = m.Count(model, scopes)

	return itemCount > 0
}

func (m *Model) Find(model interface{}, where ...interface{}) error {
	if err := m.DB().Find(model, where...).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) FindWithUnscoped(model interface{}, where ...interface{}) error {
	if err := m.DB().Unscoped().Find(model, where...).Error; err != nil {
		return err
	}

	return nil
}

// func (m *Model) Find(model interface{}, id interface{}, fields ...string) error {
// 	var field = "id"
//
// 	if len(fields) > 0 {
// 		field = fields[0]
// 	}
//
// 	if err := m.DB().Where("`?` = ?", gorm.Expr(field), id).Find(model).Error; err != nil {
// 		return err
// 	}
//
// 	return nil
// }

// alias Get
func (m *Model) FindScopes(model interface{}, scopes ...Scopes) error {
	return m.Get(model, scopes...)
}

func (m *Model) Get(model interface{}, scopes ...Scopes) error {
	var query = m.DB().Model(model)

	if len(scopes) > 0 {
		for _, scope := range scopes[0] {
			query = query.Scopes(scope)
		}
	}

	if err := query.Find(model).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) First(model interface{}, where ...interface{}) error {
	var db = m.DB().Model(model)

	if err := db.First(model, where...).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) FirstWithUnscoped(model interface{}, where ...interface{}) error {
	var db = m.DB().Model(model).Unscoped()

	if err := db.First(model, where...).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) FirstScopes(model interface{}, scopes ...Scopes) error {
	var query = m.DB().Model(model)

	if len(scopes) > 0 {
		for _, scope := range scopes[0] {
			query = query.Scopes(scope)
		}
	}

	if err := query.First(model).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) FirstWithUnscopedForUpdate(model interface{}, where ...interface{}) error {
	var db = m.DB().Set("gorm:query_option", "FOR UPDATE").Model(model).Unscoped()

	if err := db.First(model, where...).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) FirstForUpdate(model interface{}, where ...interface{}) error {
	var db = m.DB().Set("gorm:query_option", "FOR UPDATE").Model(model)

	if err := db.First(model, where...).Error; err != nil {
		return err
	}

	return nil
}

// 激活0值保存，默认struct保存会忽略0值
// model.M().Updates(&model, structs.Map(form))
func (m *Model) Update(model interface{}, attrs ...interface{}) error {
	return m.DB().Model(model).Update(attrs...).Error
}

func (m *Model) Updates(model interface{}, attrs interface{}, ignoreProtectedAttrs ...bool) error {
	return m.DB().Model(model).Updates(attrs, ignoreProtectedAttrs...).Error
}

func (m *Model) Delete(model interface{}, where ...interface{}) error {
	return m.DB().Model(model).Delete(model, where...).Error
}

func (m *Model) Count(model interface{}, scopes ...Scopes) int {
	var count int
	var query = m.DB().Model(model)

	if len(scopes) > 0 {
		for _, scope := range scopes[0] {
			query = query.Scopes(scope)
		}
	}

	query.Count(&count)

	return count
}

func (m *Model) Paginate(model interface{}, page, pageCount, numPages int, scopes ...Scopes) (*paginater.Paginater, error) {
	var total = m.Count(model, scopes ...)

	var _scopes Scopes

	if len(scopes) > 0 {
		_scopes = append(_scopes, scopes[0]...)
	}

	_scopes = append(_scopes, func(db *gorm.DB) *gorm.DB {
		return db.Offset(pageCount * (page - 1)).Limit(pageCount)
	})

	if err := m.Get(model, _scopes); err != nil {
		return nil, err
	}

	return paginater.New(total, pageCount, page, numPages), nil
}
