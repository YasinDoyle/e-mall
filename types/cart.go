package types

type CartServiceReq struct {
	Id        uint `form:"id" json:"id"`
	BossID    uint `form:"boss_id" json:"boss_id"`
	ProductId uint `form:"product_id" json:"product_id"`
	Num       uint `form:"num" json:"num"`
	UserId    uint `form:"user_id" json:"user_id"`
}

type CartCreateReq struct {
	BossID    uint `form:"boss_id" json:"boss_id"`
	ProductId uint `form:"product_id" json:"product_id"`
}

type CartDeleteReq struct {
	Id uint `form:"id" json:"id"`
}

type UpdateCartServiceReq struct {
	Id  uint `form:"id" json:"id"`
	Num uint `form:"num" json:"num"`
}

type CartListReq struct {
	BasePage
}

// 购物车
type CartResp struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"user_id"`
	ProductID     uint   `json:"product_id"`
	CreatedAt     int64  `json:"created_at"`
	Num           uint   `json:"num"`
	MaxNum        uint   `json:"max_num"`
	Check_        bool   `json:"check"`
	Name          string `json:"name"`
	ImgPath       string `json:"img_path"`
	DiscountPrice string `json:"discount_price"`
	BossId        uint   `json:"boss_id"`
	BossName      string `json:"boss_name"`
	Info          string `json:"info"`
}
