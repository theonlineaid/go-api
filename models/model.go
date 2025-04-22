// models/models.go
package models

import "time"

type Brand struct {
	ID                int       `json:"id"`
	BrandName         string    `json:"brand_name" binding:"required"`
	Image             string    `json:"image"`
	Status            int       `json:"status"`
	IsFeature         bool      `json:"is_feature"`
	IsPublish         bool      `json:"is_publish"`
	IsSpecial         bool      `json:"is_special"`
	IsApprovedByAdmin bool      `json:"is_approved_by_admin"`
	IsVisibleToGuest  bool      `json:"is_visible_to_guest"`
	CreatedBy         string    `json:"created_by" binding:"required"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type Category struct {
	ID              int    `json:"id"`
	CategoryName    string `json:"category_name"`
	Image           string `json:"image"`
	IsSpecial       int    `json:"is_special"`
	PriceVisibility int    `json:"price_visibility"`
	Status          int    `json:"status"`
}

type Subcategory struct {
	ID              int    `json:"id"`
	CategoryID      int    `json:"category_id"`
	SubcategoryName string `json:"subcategory_name"`
	Image           string `json:"image"`
	Status          int    `json:"status"`
}

type Attribute struct {
	ID            int    `json:"id"`
	AttributeName string `json:"attribute_name"`
	Status        int    `json:"status"`
}

type AttributeValue struct {
	ID          int    `json:"id"`
	AttributeID int    `json:"attribute_id"`
	Value       string `json:"value"`
	Status      int    `json:"status"`
}

type Product struct {
	ID              int      `json:"id"`
	BrandID         int      `json:"brand_id"`
	CategoryID      int      `json:"category_id"`
	SubCategoryID   *int     `json:"sub_category_id"`
	PCode           *string  `json:"p_code"`
	Weight          *string  `json:"weight"`
	ProductName     string   `json:"product_name"`
	ProductCode     string   `json:"product_code"`
	Price           *float64 `json:"price"`
	MTotalPrice     *float64 `json:"m_total_price"`
	Unit            *string  `json:"unit"`
	Discount        *float64 `json:"discount"`
	Tax             *float64 `json:"tax"`
	TaxType         string   `json:"tax_type"`
	SerialNo        *string  `json:"serial_no"`
	ProductVat      *float64 `json:"product_vat"`
	ProductModel    *string  `json:"product_model"`
	Warranty        *string  `json:"warranty"`
	MinimumQtyAlert *int     `json:"minimum_qty_alert"`
	Image           string   `json:"image"`
	IsMulti         int      `json:"is_multi"`
	SerialNumber    int      `json:"serial_number"`
	Tax0            *float64 `json:"tax0"`
	Tax1            *float64 `json:"tax1"`
	HsnCode         *string  `json:"hsn_code"`
	IsSaleable      int      `json:"is_saleable"`
	IsBarcode       int      `json:"is_barcode"`
	IsExpirable     int      `json:"is_expirable"`
	IsWarranty      int      `json:"is_warranty"`
	IsServiceable   int      `json:"is_serviceable"`
	IsVariation     int      `json:"is_variation"`
	Status          int      `json:"status"`
	UOMID           *int     `json:"UOM_id"`
	Rating          int      `json:"rating"`
	CurrentStock    int      `json:"current_stock"`
}

type ProductAttribute struct {
	ID               int `json:"id"`
	ProductID        int `json:"product_id"`
	AttributeID      int `json:"attribute_id"`
	AttributeValueID int `json:"attribute_value_id"`
}

type VariationProduct struct {
	ID               int      `json:"id"`
	ProductID        int      `json:"product_id"`
	SKU              string   `json:"sku"`
	SalePrice        *float64 `json:"sale_price"`
	DefaultSellPrice *float64 `json:"default_sell_price"`
	Discount         *float64 `json:"discount"`
	Image            string   `json:"image"`
	CurrentStock     int      `json:"current_stock"`
}
