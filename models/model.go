// models/models.go
package models

import "time"

type User struct {
	ID          int       `json:"id"`
	Username    string    `json:"username" binding:"required"`
	Email       string    `json:"email" binding:"required,email"`
	Password    string    `json:"password" binding:"required"`
	Role        string    `json:"role" binding:"required"` // user, admin, manager
	PhoneNumber *string   `json:"phone_number,omitempty"`
	Image       *string   `json:"image,omitempty"`
	IsVerified  bool      `json:"is_verified"`
	IsBlocked   bool      `json:"is_blocked"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Address struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	AddressLine1 string    `json:"address_line1" binding:"required"`
	City         string    `json:"city" binding:"required"`
	Country      string    `json:"country" binding:"required"`
	PostalCode   string    `json:"postal_code" binding:"required"`
	Type         string    `json:"type" binding:"required,oneof=home office other"`
	CreatedAt    time.Time `json:"created_at"`
}

type LoginSession struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Browser   string    `json:"browser"`
	OS        string    `json:"os"`
	Device    string    `json:"device"`
	IPAddress string    `json:"ip_address"`
	LoginAt   time.Time `json:"login_at"`
}

type ShippingAddress struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	AddressLine1 string    `json:"address_line1" binding:"required"`
	City         string    `json:"city" binding:"required"`
	Country      string    `json:"country" binding:"required"`
	PostalCode   string    `json:"postal_code" binding:"required"`
	Type         string    `json:"type" binding:"required,oneof=home office other"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
}

type BillingAddress struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	AddressLine1 string    `json:"address_line1" binding:"required"`
	City         string    `json:"city" binding:"required"`
	Country      string    `json:"country" binding:"required"`
	PostalCode   string    `json:"postal_code" binding:"required"`
	Type         string    `json:"type" binding:"required,oneof=home office other"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
}

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

type SubCategory struct {
	ID              int     `json:"id"`
	CategoryID      int     `json:"category_id"`
	SubCategoryName string  `json:"subcategory_name"`
	Image           *string `json:"image,omitempty"`
	Status          int     `json:"status"`
}

type Category struct {
	ID                 int           `json:"id"`
	Code               *string       `json:"code,omitempty"`
	CategoryName       string        `json:"category_name" binding:"required"`
	CategoryImg        *string       `json:"category_img,omitempty"`
	Image              *string       `json:"image,omitempty"`
	CategoryVisibility int           `json:"category_visibility"`
	IsSpecial          int           `json:"is_special"`
	IsFeatured         int           `json:"is_featured"`
	IsApproved         bool          `json:"is_approved"`
	IsPublished        bool          `json:"is_published"`
	Position           *int          `json:"position,omitempty"`
	PriceVisibility    int           `json:"price_visibility"`
	Status             int           `json:"status"`
	CreatedBy          string        `json:"created_by" binding:"required"`
	CreatedAt          time.Time     `json:"created_at"`
	UpdatedAt          time.Time     `json:"updated_at"`
	ProductsCount      int           `json:"products_count"`
	SubCategories      []SubCategory `json:"sub_categories"`
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
