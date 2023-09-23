package entity

// Product represents data about an product.
type Product struct {
	ID          int     `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Categories  string  `json:"categories"`
	Size        string  `json:"size"`
	Price       float64 `json:"price"`
	Tax         float64 `json:"tax"`
	Available   bool    `json:"available"`
	StoreID     int     `json:"store_id"`
	Images      []Image `json:"images"`
}
