package service

import (
	"github.com/budisetionugroho123/be-go-invoice/internal/models"
	"github.com/budisetionugroho123/be-go-invoice/internal/repository"
)

type ItemService struct {
	itemRepo *repository.ItemRepository
}

func NewItemService(itemRepo *repository.ItemRepository) *ItemService {
	return &ItemService{itemRepo: itemRepo}
}

func (s *ItemService) SearchByCode(code string) ([]models.Item, error) {
	return s.itemRepo.SearchByCode(code)
}
