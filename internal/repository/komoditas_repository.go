package repository

import (
	"github.com/example/be-panganlink-data-handler/internal/model"
	"gorm.io/gorm"
)

type KomoditasRepository interface {
	FindAll() ([]model.Komoditas, error)
	FindByID(id string) (*model.Komoditas, error)
	Create(k *model.Komoditas) error
	Update(k *model.Komoditas) error
	Delete(id string) error
}

type komoditasRepository struct{ db *gorm.DB }

func NewKomoditasRepository(db *gorm.DB) KomoditasRepository { return &komoditasRepository{db} }

func (r *komoditasRepository) FindAll() ([]model.Komoditas, error) {
	var k []model.Komoditas
	err := r.db.Find(&k).Error
	return k, err
}

func (r *komoditasRepository) FindByID(id string) (*model.Komoditas, error) {
	var k model.Komoditas
	err := r.db.Where("id = ?", id).First(&k).Error
	return &k, err
}

func (r *komoditasRepository) Create(k *model.Komoditas) error { return r.db.Create(k).Error }
func (r *komoditasRepository) Update(k *model.Komoditas) error { return r.db.Save(k).Error }
func (r *komoditasRepository) Delete(id string) error { return r.db.Where("id = ?", id).Delete(&model.Komoditas{}).Error }
