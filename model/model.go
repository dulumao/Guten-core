package model

import (
	"github.com/dulumao/Guten-utils/paginater"
	"github.com/jinzhu/gorm"
	"github.com/dulumao/Guten-core/database"
)

type Model struct {
	db *gorm.DB
}

type TransactionFunc func(m *Model) error
type Func func(m *Model) error
type NotFoundCallback func(m *Model)
type Where []func(db *gorm.DB) *gorm.DB

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
	var wheres Where
	var itemCount = 0
	var field = "id"

	if len(fields) > 0 {
		field = fields[0]
	}

	wheres = append(wheres, func(db *gorm.DB) *gorm.DB {
		return db.Where("`?` = ?", gorm.Expr(field), id)
	})

	itemCount = m.Count(model, wheres)

	return itemCount > 0
}

func (m *Model) Find(model interface{}, id interface{}, fields ...string) error {
	var field = "id"

	if len(fields) > 0 {
		field = fields[0]
	}

	if err := m.DB().Where("`?` = ?", gorm.Expr(field), id).Find(model).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) Get(model interface{}, wheres ...Where) error {
	var query = m.DB().Model(model)

	if len(wheres) > 0 {
		for _, scope := range wheres[0] {
			query = query.Scopes(scope)
		}
	}

	if err := query.Find(model).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) First(model interface{}, withTrasheds ...bool) error {
	var db = m.DB()

	if len(withTrasheds) > 0 && withTrasheds[0] {
		db = db.Unscoped()
	}

	if err := db.First(model).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) FirstForUpdate(model interface{}, withTrasheds ...bool) error {
	var db = m.DB().Set("gorm:query_option", "FOR UPDATE")

	if len(withTrasheds) > 0 && withTrasheds[0] {
		db = db.Unscoped()
	}

	if err := db.First(model).Error; err != nil {
		return err
	}

	return nil
}

func (m *Model) Update(model interface{}, attrs ...interface{}) error {
	return m.DB().Model(model).Update(attrs...).Error
}

func (m *Model) Updates(model interface{}, attrs interface{}, ignoreProtectedAttrs ...bool) error {
	return m.DB().Model(model).Updates(attrs, ignoreProtectedAttrs...).Error
}

func (m *Model) Delete(model interface{}, where ...interface{}) error {
	return m.DB().Delete(model, where...).Error
}

func (m *Model) Count(model interface{}, wheres ...Where) int {
	var count int
	var query = m.DB().Model(model)

	if len(wheres) > 0 {
		for _, scope := range wheres[0] {
			query = query.Scopes(scope)
		}
	}

	query.Count(&count)

	return count
}

func (m *Model) Paginate(model interface{}, page, pageCount, numPages int, wheres ...Where) (*paginater.Paginater, error) {
	var total = m.Count(model, wheres ...)

	var _wheres Where

	if len(wheres) > 0 {
		_wheres = append(_wheres, wheres[0]...)
	}

	_wheres = append(_wheres, func(db *gorm.DB) *gorm.DB {
		return db.Offset(pageCount * (page - 1)).Limit(pageCount)
	})

	if err := m.Get(model, _wheres); err != nil {
		return nil, err
	}

	return paginater.New(total, pageCount, page, numPages), nil
}
