package entity

// Image represents data about an image.
type Image struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	ImagePath string `json:"image_path"`
	ProductID int    `json:"product_id"`
}
