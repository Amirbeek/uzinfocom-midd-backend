package models

type Product struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Price         int     `json:"price"`
	StockQuantity float64 `json:"stock_quantity"`
}
