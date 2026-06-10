package service

import (
	"github.com/example/be-panganlink-data-handler/internal/model"
	"github.com/example/be-panganlink-data-handler/internal/repository"
)

type KomoditasService interface {
	GetAll() ([]model.Komoditas, error)
	Create(k *model.Komoditas) error
	Update(id string, k *model.Komoditas) error
	Delete(id string) error
}

type komoditasService struct{ repo repository.KomoditasRepository }

func NewKomoditasService(repo repository.KomoditasRepository) KomoditasService { return &komoditasService{repo} }

func (s *komoditasService) GetAll() ([]model.Komoditas, error) { return s.repo.FindAll() }
func (s *komoditasService) Create(k *model.Komoditas) error { return s.repo.Create(k) }
func (s *komoditasService) Update(id string, k *model.Komoditas) error {
	k.ID = id
	return s.repo.Update(k)
}
func (s *komoditasService) Delete(id string) error { return s.repo.Delete(id) }
