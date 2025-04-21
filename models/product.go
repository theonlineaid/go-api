// models/product.go
package models

type Category struct {
	ID              int    `json:"id"`
	CategoryName    string `json:"category_name"`
	IsSpecial       int    `json:"is_special"`
	PriceVisibility int    `json:"price_visibility"`
}

type VariationProduct struct {
	ID                 int      `json:"id"`
	VariationID        int      `json:"variation_id"`
	VariationDetailsID int      `json:"variation_details_id"`
	ProductID          int      `json:"product_id"`
	SKU                string   `json:"sku"`
	Value              string   `json:"value"`
	Color              string   `json:"color"`
	SalePrice          float64  `json:"sale_price"`
	DefaultSellPrice   float64  `json:"default_sell_price"`
	Discount           *float64 `json:"discount"`
	Image              string   `json:"image"`
}

type Product struct {
	ID                int                `json:"id"`
	BrandID           int                `json:"brand_id"`
	CategoryID        int                `json:"category_id"`
	SubCategoryID     *int               `json:"sub_category_id"`
	PCode             *string            `json:"p_code"`
	Weight            *string            `json:"weight"`
	ProductName       string             `json:"product_name"`
	ProductCode       string             `json:"product_code"`
	Price             *float64           `json:"price"`
	MTotalPrice       *float64           `json:"m_total_price"`
	Unit              *string            `json:"unit"`
	Discount          *float64           `json:"discount"`
	Tax               *float64           `json:"tax"`
	TaxType           string             `json:"tax_type"`
	SerialNo          *string            `json:"serial_no"`
	ProductVat        *float64           `json:"product_vat"`
	ProductModel      *string            `json:"product_model"`
	Warranty          *string            `json:"warranty"`
	MinimumQtyAlert   *int               `json:"minimum_qty_alert"`
	Image             string             `json:"image"`
	IsMulti           *int               `json:"is_multi"`
	SerialNumber      int                `json:"serial_number"`
	Tax0              *float64           `json:"tax0"`
	Tax1              *float64           `json:"tax1"`
	HsnCode           *string            `json:"hsn_code"`
	IsSaleable        int                `json:"is_saleable"`
	IsBarcode         int                `json:"is_barcode"`
	IsExpirable       *int               `json:"is_expirable"`
	IsWarranty        int                `json:"is_warranty"`
	IsServiceable     *int               `json:"is_serviceable"`
	IsVariation       int                `json:"is_variation"`
	Status            int                `json:"status"`
	UOMID             *int               `json:"UOM_id"`
	Rating            int                `json:"rating"`
	CurrentStock      int                `json:"current_stock"`
	Category          Category           `json:"category"`
	VariationProducts []VariationProduct `json:"variation_products"`
}
