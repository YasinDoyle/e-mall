package e

var MsgFlags = map[int]string{
	SUCCESS:               "ok",
	UpdatePasswordSuccess: "修改密码成功",
	NotExistInentifier:    "该第三方账号未绑定",
	ERROR:                 "fail",
	InvalidParams:         "请求参数错误",

	ErrorExistNick:          "已存在该昵称",
	ErrorExistUser:          "已存在该用户名",
	ErrorNotExistUser:       "该用户不存在",
	ErrorNotCompare:         "账号密码错误",
	ErrorNotComparePassword: "两次密码输入不一致",
	ErrorFailEncryption:     "加密失败",
	ErrorNotExistProduct:    "该商品不存在",
	ErrorNotExistAddress:    "该收获地址不存在",
	ErrorExistFavorite:      "已收藏该商品",
	ErrorUserNotFound:       "用户不存在",

	ErrorBossCheckTokenFail:        "商家的Token鉴权失败",
	ErrorBossCheckTokenTimeout:     "商家Token已超时",
	ErrorBossToken:                 "商家的Token生成失败",
	ErrorBoss:                      "商家Token错误",
	ErrorBossInsufficientAuthority: "商家权限不足",
	ErrorBossProduct:               "商家读文件错误",

	ErrorProductExistCart: "商品已经在购物车了，数量+1",
	ErrorProductMoreCart:  "超过最大上限",

	ErrorAuthCheckTokenFail:        "Token鉴权失败",
	ErrorAuthCheckTokenTimeout:     "Token已超时",
	ErrorAuthToken:                 "Token生成失败",
	ErrorAuth:                      "Token错误",
	ErrorAuthInsufficientAuthority: "权限不足",
	ErrorReadFile:                  "读文件失败",
	ErrorCallApi:                   "调用接口失败",
	ErrorUnmarshalJson:             "解码JSON失败",

	ErrorUploadFile:    "上传失败",
	ErrorAdminFindUser: "管理员查询用户失败",

	ErrorDatabase: "数据库操作出错,请重试",

	ErrorOss: "OSS配置错误",

	ErrorSellerNotApplied:           "尚未申请商家入驻",
	ErrorSellerAuditPending:         "商家入驻申请正在审核中",
	ErrorSellerAlreadyApproved:      "商家入驻已通过，无需重复申请",
	ErrorSellerBanned:               "商家账号已被封禁，无法重新申请",
	ErrorSellerNotApproved:          "请先完成商家入驻并通过审核",
	ErrorSellerInvalidStatus:        "商家状态异常",
	ErrorSellerInvalidApplication:   "参数错误",
	ErrorSellerShopNameRequired:     "店铺名称不能为空",
	ErrorSellerShopNameTooLong:      "店铺名称不能超过80个字符",
	ErrorSellerDescriptionTooLong:   "店铺描述不能超过500个字符",
	ErrorSellerRejectReasonMissing:  "拒绝原因不能为空",
	ErrorSellerAuditStatusInvalid:   "商家审核状态不正确",
	ErrorSellerPayKeyRequired:       "请先设置支付密码再上架商品",
	ErrorProductSellerNotApproved:   "商品卖家尚未完成商家入驻审核",
	ErrorCarouselProductRequired:    "请选择关联商品",
	ErrorCarouselProductNotExist:    "关联商品不存在",
	ErrorSettlementInvalidAmount:    "结算金额不合法",
	ErrorSettlementInvalidRate:      "佣金比例不合法",
	ErrorSettlementSellerInvalid:    "商家不存在或未通过审核",
	ErrorSettlementStatusInvalid:    "结算单状态不允许操作",
	ErrorOrderPayStatusInvalid:      "订单已支付或状态不允许支付",
	ErrorPaymentPayKeyRequired:      "请先设置支付密码",
	ErrorPaymentPayKeyInvalid:       "支付密码错误",
	ErrorPaymentBalanceInsufficient: "金币不足",
	ErrorPaymentStockInsufficient:   "库存不足",
	ErrorRefundStatusInvalid:        "订单状态不允许退款审批",
	ErrorRefundAmountInvalid:        "退款金额不合法",
	ErrorRefundNotFound:             "退款申请不存在",
}

// GetMsg 获取状态码对应信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
