package service

import (
	"github.com/example/be-panganlink-data-handler/internal/model"
	"github.com/example/be-panganlink-data-handler/internal/repository"
	"github.com/google/uuid"
)

type ProductService interface {
	GetAll(page, limit int) ([]model.Product, error)
	GetByUserID(userID string) ([]model.Product, error)
	GetByID(id string) (*model.Product, error)
	Create(p *model.Product) error
	Update(id string, p *model.Product) error
	Delete(id string) error
}

type productService struct{ repo repository.ProductRepository }

func NewProductService(repo repository.ProductRepository) ProductService { return &productService{repo} }

func (s *productService) GetAll(page, limit int) ([]model.Product, error) { return s.repo.FindAll(page, limit) }
func (s *productService) GetByUserID(userID string) ([]model.Product, error) { return s.repo.FindAllByUserID(userID) }
func (s *productService) GetByID(id string) (*model.Product, error) { return s.repo.FindByID(id) }
func (s *productService) Create(p *model.Product) error {
	p.ID = uuid.New().String()
	p.Status = "pending"
	return s.repo.Create(p)
}
func (s *productService) Update(id string, p *model.Product) error {
	p.ID = id
	return s.repo.Update(p)
}
func (s *productService) Delete(id string) error { return s.repo.Delete(id) }
