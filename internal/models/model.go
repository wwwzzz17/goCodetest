package models

type ContextKey string

const TraceIDKey ContextKey = "trace_id"

type HTTPResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type ProductListQuery struct {
	Page        int     `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit       int     `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	IDs         []int64 `form:"ids[]" json:"ids"`
	Name        string  `form:"name" json:"name"`
	PriceMin    float64 `form:"price_min" binding:"omitempty,min=0" json:"price_min"`
	PriceMax    float64 `form:"price_max" binding:"omitempty,min=0" json:"price_max"`
	QuantityMin int64   `form:"quantity_min" binding:"omitempty,min=0" json:"quantity_min"`
	QuantityMax int64   `form:"quantity_max" binding:"omitempty,min=0" json:"quantity_max"`
}
