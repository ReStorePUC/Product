package repository

import (
	"context"
	"github.com/restore/product/entity"
	"gorm.io/gorm"
	"strings"
)

type Product struct {
	db *gorm.DB
}

func NewProduct(db *gorm.DB) *Product {
	return &Product{
		db: db,
	}
}

func (p *Product) CreateProduct(ctx context.Context, product *entity.Product) (int, error) {
	result := p.db.Create(product)
	if result.Error != nil {
		return 0, result.Error
	}
	return product.ID, nil
}

func (p *Product) GetProduct(ctx context.Context, id int) (*entity.Product, error) {
	result := entity.Product{ID: id}
	res := p.db.First(&result)
	if res.Error != nil {
		return nil, res.Error
	}

	var images []entity.Image
	query := p.db.Where("product_id = ?", id)
	res = query.Find(&images)
	if res.Error != nil {
		return nil, res.Error
	}

	result.Images = images
	return &result, nil
}

func (p *Product) Unavailable(ctx context.Context, id int) error {
	result := entity.Product{ID: id}
	res := p.db.First(&result)
	if res.Error != nil {
		return res.Error
	}

	result.Available = false

	res = p.db.Save(result)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

func (p *Product) ListProduct(ctx context.Context, id int, unavailable bool, name string) ([]entity.Product, error) {
	var result []entity.Product

	query := p.db.Where("store_id = ? AND available = ?", id, !unavailable)
	if name != "" {
		query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(name)+"%")
	}
	res := query.Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}

	for i, prod := range result {
		var images []entity.Image
		query = p.db.Where("product_id = ?", prod.ID)
		res = query.Find(&images)
		if res.Error != nil {
			return nil, res.Error
		}

		result[i].Images = images
	}

	return result, nil
}

func (p *Product) ListRecent(ctx context.Context) ([]entity.Product, error) {
	var result []entity.Product
	res := p.db.Where("available = TRUE").Limit(10).Order("id desc").Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}

	for i, prod := range result {
		var images []entity.Image
		res = p.db.Where("product_id = ?", prod.ID).Find(&images)
		if res.Error != nil {
			return nil, res.Error
		}

		result[i].Images = images
	}

	return result, nil
}

func (p *Product) Search(ctx context.Context, name string, categories []string) ([]entity.Product, error) {
	var result []entity.Product
	query := p.db.Where("available = TRUE")

	if name != "" {
		query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(name)+"%")
	}

	if categories != nil {
		for _, cat := range categories {
			query.Where("LOWER(categories) LIKE ?", "%"+strings.ToLower(cat)+"%")
		}
	}

	res := query.Find(&result)
	if res.Error != nil {
		return nil, res.Error
	}

	for i, prod := range result {
		var images []entity.Image
		res = p.db.Where("product_id = ?", prod.ID).Find(&images)
		if res.Error != nil {
			return nil, res.Error
		}

		result[i].Images = images
	}

	return result, nil
}
