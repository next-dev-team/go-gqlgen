package model

import (
	"gorm.io/gorm"
)

type Base[T any] interface {
	Create(t *T) (*T, error)
	Update(t *T) (*T, error)
	Delete(t *T) error
	FindById(id any) (*T, error)
	FindAllByIds(id []any) ([]*T, error)
}

type base[T any] struct {
	db *gorm.DB
}

func (b *base[T]) Create(t *T) (*T, error) {
	if err := b.db.Create(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (b *base[T]) Update(t *T) (*T, error) {
	if err := b.db.Save(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (b *base[T]) Delete(t *T) error {
	if err := b.db.Delete(t).Error; err != nil {
		return err
	}
	return nil
}
func (b *base[T]) FindById(id any) (*T, error) {
	var t T
	if err := b.db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (b *base[T]) FindAllByIds(id []any) ([]*T, error) {
	var t []*T
	if len(id) == 0 {
		if err := b.db.Find(&t).Error; err != nil {
			return nil, err
		}
		return t, nil
	}
	if err := b.db.Where("id in (?)", id).Find(&t).Error; err != nil {
		return nil, err
	}
	return t, nil
}
