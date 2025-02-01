package models

type PricesSummary struct {
  TotalItems      int     `json:"total_items"`
  TotalCategories int     `json:"total_categories"`
  TotalPrice      float32 `json:"total_price"`
}

func NewPricesSummary() *PricesSummary {
  ps := PricesSummary{}

  ps.TotalItems = 12
  ps.TotalCategories = 10
  ps.TotalPrice = 1000.5

  return &ps
}
