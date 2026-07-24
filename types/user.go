package types

type UserServiceReq struct {
	NickName string `form:"nick_name" json:"nick_name"`
	UserName string `form:"user_name" json:"user_name"`
	Password string `form:"password" json:"password"`
	Key      string `form:"key" json:"key"`
}

type UserRegisterReq struct {
	NickName        string `form:"nick_name" json:"nick_name"`
	UserName        string `form:"user_name" json:"user_name"`
	Email           string `form:"email" json:"email"`
	EmailCode       string `form:"email_code" json:"email_code"`
	Password        string `form:"password" json:"password"`
	PasswordConfirm string `form:"password_confirm" json:"password_confirm"`
}

type RegisterEmailCodeReq struct {
	Email string `form:"email" json:"email"`
}

type UserTokenData struct {
	User         any    `json:"user"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserLoginReq struct {
	UserName string `form:"user_name" json:"user_name"`
	Password string `form:"password" json:"password"`
}

type UserInfoUpdateReq struct {
	NickName string `form:"nick_name" json:"nick_name"`
}

type UserInfoShowReq struct {
}

type UserFollowingReq struct {
	Id uint `json:"id" form:"id"`
}

type UserUnFollowingReq struct {
	Id uint `json:"id" form:"id"`
}

type UserInfoResp struct {
	ID        uint   `json:"id"`
	UserName  string `json:"user_name"`
	NickName  string `json:"nickname"`
	Type      int    `json:"type"`
	Email     string `json:"email"`
	Status    string `json:"status"`
	Avatar    string `json:"avatar"`
	PayKeySet bool   `json:"pay_key_set"`
	CreateAt  int64  `json:"create_at"`
}
