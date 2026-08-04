package dto

type CreateServiceTypeRequest struct {
	Name           string  `json:"name" binding:"required,min=2,max=50"`
	Description    string  `json:"description"`
	BasePrice      float64 `json:"base_price" binding:"gte=0"`
	PricePerWeight float64 `json:"price_per_weight" binding:"gte=0"`
	PricePerPiece  float64 `json:"price_per_piece" binding:"gte=0"`
}
