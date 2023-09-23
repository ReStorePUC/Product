package controller

import (
	"context"
	"github.com/restore/product/entity"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

type repository interface {
	CreateProduct(ctx context.Context, product *entity.Product) (int, error)
	GetProduct(ctx context.Context, id int) (*entity.Product, error)
	ListProduct(ctx context.Context, id int, unavailable bool) ([]entity.Product, error)
	ListRecent(ctx context.Context) ([]entity.Product, error)
	Unavailable(ctx context.Context, id int) error
	Search(ctx context.Context, name string, categories []string) ([]entity.Product, error)
}

type Product struct {
	repo repository
}

func NewProduct(r repository) *Product {
	return &Product{
		repo: r,
	}
}

func (p *Product) Create(ctx context.Context, product *entity.Product) (int, error) {
	log := zap.NewNop()

	id, err := p.repo.CreateProduct(ctx, product)
	if err != nil {
		log.Error(
			"error to create product",
			zap.Error(err),
		)
		return 0, err
	}

	return id, nil
}

func (p *Product) GetProduct(ctx context.Context, id string) (*entity.Product, error) {
	log := zap.NewNop()

	productID, err := strconv.Atoi(id)
	if err != nil {
		log.Error(
			"error validating id",
			zap.Error(err),
		)
		return nil, err
	}

	result, err := p.repo.GetProduct(ctx, productID)
	if err != nil {
		log.Error(
			"error to get product",
			zap.Error(err),
		)
		return nil, err
	}

	return result, nil
}

func (p *Product) ListProduct(ctx context.Context, id string, unavailable bool) ([]entity.Product, error) {
	log := zap.NewNop()

	storeID, err := strconv.Atoi(id)
	if err != nil {
		log.Error(
			"error validating id",
			zap.Error(err),
		)
		return nil, err
	}

	result, err := p.repo.ListProduct(ctx, storeID, unavailable)
	if err != nil {
		log.Error(
			"error to list products",
			zap.Error(err),
		)
		return nil, err
	}

	return result, nil
}

func (p *Product) ListRecent(ctx context.Context) ([]entity.Product, error) {
	log := zap.NewNop()

	result, err := p.repo.ListRecent(ctx)
	if err != nil {
		log.Error(
			"error to list products",
			zap.Error(err),
		)
		return nil, err
	}

	return result, nil
}

func (p *Product) Unavailable(ctx context.Context, id string) error {
	log := zap.NewNop()

	productID, err := strconv.Atoi(id)
	if err != nil {
		log.Error(
			"error validating id",
			zap.Error(err),
		)
		return err
	}

	err = p.repo.Unavailable(ctx, productID)
	if err != nil {
		log.Error(
			"error to change product status",
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (p *Product) Search(ctx context.Context, name, categories string) ([]entity.Product, error) {
	log := zap.NewNop()

	cat := strings.Split(categories, ",")
	result, err := p.repo.Search(ctx, name, cat)
	if err != nil {
		log.Error(
			"error to search for products",
			zap.Error(err),
		)
		return nil, err
	}

	return result, nil
}
