package e

const (
	SUCCESS               = 200
	UpdatePasswordSuccess = 201
	NotExistInentifier    = 202
	ERROR                 = 500
	InvalidParams         = 400

	//成员错误
	ErrorExistNick          = 10001
	ErrorExistUser          = 10002
	ErrorNotExistUser       = 10003
	ErrorNotCompare         = 10004
	ErrorNotComparePassword = 10005
	ErrorFailEncryption     = 10006
	ErrorNotExistProduct    = 10007
	ErrorNotExistAddress    = 10008
	ErrorExistFavorite      = 10009
	ErrorUserNotFound       = 10010

	//店家错误
	ErrorBossCheckTokenFail        = 20001
	ErrorBossCheckTokenTimeout     = 20002
	ErrorBossToken                 = 20003
	ErrorBoss                      = 20004
	ErrorBossInsufficientAuthority = 20005
	ErrorBossProduct               = 20006

	// 购物车
	ErrorProductExistCart = 20007
	ErrorProductMoreCart  = 20008

	//管理员错误
	ErrorAuthCheckTokenFail        = 30001 //token 错误
	ErrorAuthCheckTokenTimeout     = 30002 //token 过期
	ErrorAuthToken                 = 30003
	ErrorAuth                      = 30004
	ErrorAuthInsufficientAuthority = 30005
	ErrorReadFile                  = 30006
	ErrorCallApi                   = 30008
	ErrorUnmarshalJson             = 30009
	ErrorAdminFindUser             = 30010
	//数据库错误
	ErrorDatabase = 40001

	//对象存储错误
	ErrorOss        = 50001
	ErrorUploadFile = 50002

	//商家入驻和平台商品错误
	ErrorSellerNotApplied          = 60001
	ErrorSellerAuditPending        = 60002
	ErrorSellerAlreadyApproved     = 60003
	ErrorSellerBanned              = 60004
	ErrorSellerNotApproved         = 60005
	ErrorSellerInvalidStatus       = 60006
	ErrorSellerInvalidApplication  = 60007
	ErrorSellerShopNameRequired    = 60008
	ErrorSellerShopNameTooLong     = 60009
	ErrorSellerDescriptionTooLong  = 60010
	ErrorSellerRejectReasonMissing = 60011
	ErrorSellerAuditStatusInvalid  = 60012
	ErrorSellerPayKeyRequired      = 60013
	ErrorProductSellerNotApproved  = 60014
	ErrorCarouselProductRequired   = 60015
	ErrorCarouselProductNotExist   = 60016
)
