package repository

import (
	"github.com/example/be-panganlink-data-handler/internal/model"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll(page, limit int) ([]model.Product, error)
	FindAllByUserID(userID string) ([]model.Product, error)
	FindByID(id string) (*model.Product, error)
	Create(p *model.Product) error
	Update(p *model.Product) error
	Delete(id string) error
}

type productRepository struct{ db *gorm.DB }

func NewProductRepository(db *gorm.DB) ProductRepository { return &productRepository{db} }

func (r *productRepository) FindAll(page, limit int) ([]model.Product, error) {
	var products []model.Product
	offset := (page - 1) * limit
	err := r.db.Preload("Komoditas").Preload("User").Offset(offset).Limit(limit).Find(&products).Error
	return products, err
}

func (r *productRepository) FindAllByUserID(userID string) ([]model.Product, error) {
	var p []model.Product
	err := r.db.Preload("Komoditas").Where("user_id = ?", userID).Find(&p).Error
	return p, err
}

func (r *productRepository) FindByID(id string) (*model.Product, error) {
	var p model.Product
	err := r.db.Where("id = ?", id).First(&p).Error
	return &p, err
}

func (r *productRepository) Create(p *model.Product) error { return r.db.Create(p).Error }
func (r *productRepository) Update(p *model.Product) error { return r.db.Save(p).Error }
func (r *productRepository) Delete(id string) error { return r.db.Where("id = ?", id).Delete(&model.Product{}).Error }
