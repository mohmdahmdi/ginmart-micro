package models

type Review struct {
	ID          int     `json:"id"`
	ProductId   string  `json:"product_id"`
	UserId      string  `json:"user_id"`
	Rating      float64 `json:"rating"`
	Comment     string  `json:"comment"`
}
