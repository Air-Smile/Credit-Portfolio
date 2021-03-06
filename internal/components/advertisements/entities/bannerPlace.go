package entities

import "time"

type BannerPlace struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	PricePerView float64   `json:"pricePerView"`
	Description  string    `json:"description"`
}
