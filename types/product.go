package types

type ProductSearchReq struct {
	ID            uint   `form:"id" json:"id"`
	Name          string `form:"name" json:"name"`
	CategoryID    int    `form:"category_id" json:"category_id"`
	Title         string `form:"title" json:"title" `
	Info          string `form:"info" json:"info" `
	Price         string `form:"price" json:"price"`
	DiscountPrice string `form:"discount_price" json:"discount_price"`
	OnSale        bool   `form:"on_sale" json:"on_sale"`
	BasePage
}

// BossProductListReq 卖家查看自己的商品列表
type BossProductListReq struct {
	BasePage
}

// BossProductOnSaleReq 卖家上架/下架商品
type BossProductOnSaleReq struct {
	ID     uint `json:"id" binding:"required"`
	OnSale bool `json:"on_sale"`
}

type ProductCreateReq struct {
	ID            uint   `form:"id" json:"id"`
	Name          string `form:"name" json:"name"`
	CategoryID    uint   `form:"category_id" json:"category_id"`
	Title         string `form:"title" json:"title" `
	Info          string `form:"info" json:"info" `
	ImgPath       string `form:"img_path" json:"img_path"`
	Price         string `form:"price" json:"price"`
	DiscountPrice string `form:"discount_price" json:"discount_price"`
	OnSale        bool   `form:"on_sale" json:"on_sale"`
	Num           int    `form:"num" json:"num"`

	Brand             string   `form:"brand" json:"brand"`
	Origin            string   `form:"origin" json:"origin"`
	Specification     string   `form:"specification" json:"specification"`
	ProductionDate    string   `form:"production_date" json:"production_date"`
	ShelfLife         string   `form:"shelf_life" json:"shelf_life"`
	ServiceGuarantees string   `form:"service_guarantees" json:"service_guarantees"`
	CertificateMeta   string   `form:"certificate_meta" json:"certificate_meta"`
	CertificateTypes  []string `form:"certificate_type" json:"certificate_type"`
	CertificateNames  []string `form:"certificate_name" json:"certificate_name"`
}

type ProductListReq struct {
	CategoryID uint `form:"category_id" json:"category_id"`
	BasePage
}

type ProductDeleteReq struct {
	ID uint `form:"id" json:"id"`
	BasePage
}

type ProductShowReq struct {
	ID uint `form:"id" json:"id"`
}

type ProductUpdateReq struct {
	ID            uint   `form:"id" json:"id"`
	Name          string `form:"name" json:"name"`
	CategoryID    uint   `form:"category_id" json:"category_id"`
	Title         string `form:"title" json:"title" `
	Info          string `form:"info" json:"info" `
	ImgPath       string `form:"img_path" json:"img_path"`
	Price         string `form:"price" json:"price"`
	DiscountPrice string `form:"discount_price" json:"discount_price"`
	OnSale        bool   `form:"on_sale" json:"on_sale"`
	Num           int    `form:"num" json:"num"`

	Brand               string   `form:"brand" json:"brand"`
	Origin              string   `form:"origin" json:"origin"`
	Specification       string   `form:"specification" json:"specification"`
	ProductionDate      string   `form:"production_date" json:"production_date"`
	ShelfLife           string   `form:"shelf_life" json:"shelf_life"`
	ServiceGuarantees   string   `form:"service_guarantees" json:"service_guarantees"`
	CertificateMeta     string   `form:"certificate_meta" json:"certificate_meta"`
	ReplaceCertificates bool     `form:"replace_certificates" json:"replace_certificates"`
	CertificateTypes    []string `form:"certificate_type" json:"certificate_type"`
	CertificateNames    []string `form:"certificate_name" json:"certificate_name"`
}

type ListProductImgReq struct {
	ID uint `json:"id" form:"id"`
}

type ProductResp struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	CategoryID    uint   `json:"category_id"`
	Title         string `json:"title"`
	Info          string `json:"info"`
	ImgPath       string `json:"img_path"`
	Price         string `json:"price"`
	DiscountPrice string `json:"discount_price"`
	View          uint64 `json:"view"`
	CreatedAt     int64  `json:"created_at"`
	Num           int    `json:"num"`
	OnSale        bool   `json:"on_sale"`
	BossID        uint   `json:"boss_id"`
	BossName      string `json:"boss_name"`
	BossAvatar    string `json:"boss_avatar"`
	AuditStatus   uint   `json:"audit_status"`

	Brand             string                   `json:"brand"`
	Origin            string                   `json:"origin"`
	Specification     string                   `json:"specification"`
	ProductionDate    string                   `json:"production_date"`
	ShelfLife         string                   `json:"shelf_life"`
	ServiceGuarantees string                   `json:"service_guarantees"`
	CertificateMeta   string                   `json:"certificate_meta"`
	Certificates      []ProductCertificateResp `json:"certificates"`
}

type ProductCertificateResp struct {
	ID              uint   `json:"id"`
	ProductID       uint   `json:"product_id"`
	CertificateType string `json:"certificate_type"`
	Name            string `json:"name"`
	FilePath        string `json:"file_path"`
	CreatedAt       int64  `json:"created_at"`
}
