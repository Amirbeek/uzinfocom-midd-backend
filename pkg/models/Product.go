package models

type Product struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	Price         int64  `json:"price"`
	StockQuantity int    `json:"stock_quantity"`
}
