package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/restore/product/config"
	"github.com/restore/product/entity"
	"net/http"
)

type controller interface {
	Create(ctx context.Context, product *entity.Product) (int, error)
	Unavailable(ctx context.Context, id string) error
	UnavailablePayment(ctx context.Context, id string) error
	GetProduct(ctx context.Context, id string) (*entity.Product, error)
	ListProduct(ctx context.Context, id string, unavailable bool, name string) ([]entity.Product, error)
	ListRecent(ctx context.Context) ([]entity.Product, error)
	Search(ctx context.Context, name, categories string) ([]entity.Product, error)
}

type Product struct {
	controller controller
}

func NewProduct(c controller) *Product {
	return &Product{
		controller: c,
	}
}

// Create creates a new Product.
func (p *Product) Create(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), config.EmailHeader, c.GetHeader(config.EmailHeader))

	var product entity.Product
	if err := c.BindJSON(&product); err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	id, err := p.controller.Create(ctx, &product)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, struct {
		ID int
	}{
		id,
	})
}

// GetProduct returns the product information.
func (p *Product) GetProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			"invalid ID",
		})
		return
	}

	result, err := p.controller.GetProduct(c.Request.Context(), id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

// ListProduct returns the products from a store.
func (p *Product) ListProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			"invalid ID",
		})
		return
	}
	unavailable := false
	if c.Query("unavailable") == "true" {
		unavailable = true
	}

	result, err := p.controller.ListProduct(c.Request.Context(), id, unavailable, c.Query("name"))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

// ListRecent returns the recent products.
func (p *Product) ListRecent(c *gin.Context) {
	result, err := p.controller.ListRecent(c.Request.Context())
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, result)
}

// Unavailable changes the product status to unavailable.
func (p *Product) Unavailable(c *gin.Context) {
	ctx := context.WithValue(c.Request.Context(), config.EmailHeader, c.GetHeader(config.EmailHeader))

	id := c.Param("id")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			"invalid ID",
		})
		return
	}

	err := p.controller.Unavailable(ctx, id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, struct {
	}{})
}

// Search searches for products.
func (p *Product) Search(c *gin.Context) {
	categories := c.Query("categories")
	name := c.Query("name")
	if categories == "" && name == "" {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			"invalid parameters",
		})
		return
	}

	result, err := p.controller.Search(c.Request.Context(), name, categories)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusCreated, result)
}
