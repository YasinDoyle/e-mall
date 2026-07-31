package e

const DefaultLocale = "zh-CN"
const EnglishLocale = "en-US"

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

	ErrorSellerNotApplied:                  "尚未申请商家入驻",
	ErrorSellerAuditPending:                "商家入驻申请正在审核中",
	ErrorSellerAlreadyApproved:             "商家入驻已通过，无需重复申请",
	ErrorSellerBanned:                      "商家账号已被封禁，无法重新申请",
	ErrorSellerNotApproved:                 "请先完成商家入驻并通过审核",
	ErrorSellerInvalidStatus:               "商家状态异常",
	ErrorSellerInvalidApplication:          "参数错误",
	ErrorSellerShopNameRequired:            "店铺名称不能为空",
	ErrorSellerShopNameTooLong:             "店铺名称不能超过80个字符",
	ErrorSellerDescriptionTooLong:          "店铺描述不能超过500个字符",
	ErrorSellerRejectReasonMissing:         "拒绝原因不能为空",
	ErrorSellerAuditStatusInvalid:          "商家审核状态不正确",
	ErrorSellerPayKeyRequired:              "请先设置支付密码再上架商品",
	ErrorProductSellerNotApproved:          "商品卖家尚未完成商家入驻审核",
	ErrorCarouselProductRequired:           "请选择关联商品",
	ErrorCarouselProductNotExist:           "关联商品不存在",
	ErrorSettlementInvalidAmount:           "结算金额不合法",
	ErrorSettlementInvalidRate:             "佣金比例不合法",
	ErrorSettlementSellerInvalid:           "商家不存在或未通过审核",
	ErrorSettlementStatusInvalid:           "结算单状态不允许操作",
	ErrorOrderPayStatusInvalid:             "订单已支付或状态不允许支付",
	ErrorPaymentPayKeyRequired:             "请先设置支付密码",
	ErrorPaymentPayKeyInvalid:              "支付密码错误",
	ErrorPaymentBalanceInsufficient:        "金币不足",
	ErrorPaymentStockInsufficient:          "库存不足",
	ErrorRefundStatusInvalid:               "订单状态不允许退款审批",
	ErrorRefundAmountInvalid:               "退款金额不合法",
	ErrorRefundNotFound:                    "退款申请不存在",
	ErrorOrderSelfPurchaseForbidden:        "不能购买自己发布的商品",
	ErrorSellerWithdrawAmountInvalid:       "提现金额不合法",
	ErrorSellerWithdrawPayeeRequired:       "提现收款信息不能为空",
	ErrorSellerWithdrawInsufficientBalance: "可提现余额不足",
	ErrorSellerWithdrawStatusInvalid:       "提现单状态不允许操作",
	ErrorSellerWithdrawReasonMissing:       "拒绝原因不能为空",
	ErrorSellerWithdrawNotFound:            "提现单不存在",
	ErrorOrderStatusTransitionInvalid:      "订单状态不允许该操作",
}

var LocaleMsgFlags = map[string]map[int]string{
	DefaultLocale: MsgFlags,
	EnglishLocale: {
		SUCCESS:               "ok",
		UpdatePasswordSuccess: "Password updated successfully",
		NotExistInentifier:    "This third-party account is not linked",
		ERROR:                 "fail",
		InvalidParams:         "Invalid request parameters",

		ErrorExistNick:          "Nickname already exists",
		ErrorExistUser:          "Username already exists",
		ErrorNotExistUser:       "User does not exist",
		ErrorNotCompare:         "Invalid account or password",
		ErrorNotComparePassword: "The two passwords do not match",
		ErrorFailEncryption:     "Encryption failed",
		ErrorNotExistProduct:    "Product does not exist",
		ErrorNotExistAddress:    "Address does not exist",
		ErrorExistFavorite:      "Product already favorited",
		ErrorUserNotFound:       "User not found",

		ErrorBossCheckTokenFail:        "Seller token authentication failed",
		ErrorBossCheckTokenTimeout:     "Seller token has expired",
		ErrorBossToken:                 "Seller token generation failed",
		ErrorBoss:                      "Invalid seller token",
		ErrorBossInsufficientAuthority: "Insufficient seller permissions",
		ErrorBossProduct:               "Seller file read failed",

		ErrorProductExistCart: "Product is already in cart; quantity increased by 1",
		ErrorProductMoreCart:  "Quantity exceeds the limit",

		ErrorAuthCheckTokenFail:        "Token authentication failed",
		ErrorAuthCheckTokenTimeout:     "Token has expired",
		ErrorAuthToken:                 "Token generation failed",
		ErrorAuth:                      "Invalid token",
		ErrorAuthInsufficientAuthority: "Insufficient permissions",
		ErrorReadFile:                  "Failed to read file",
		ErrorCallApi:                   "API call failed",
		ErrorUnmarshalJson:             "Failed to decode JSON",

		ErrorUploadFile:    "Upload failed",
		ErrorAdminFindUser: "Admin user query failed",

		ErrorDatabase: "Database operation failed, please try again",

		ErrorOss: "OSS configuration error",

		ErrorSellerNotApplied:                  "Seller onboarding has not been submitted",
		ErrorSellerAuditPending:                "Seller onboarding application is under review",
		ErrorSellerAlreadyApproved:             "Seller onboarding has already been approved",
		ErrorSellerBanned:                      "Seller account is banned and cannot reapply",
		ErrorSellerNotApproved:                 "Please complete seller onboarding and pass review first",
		ErrorSellerInvalidStatus:               "Invalid seller status",
		ErrorSellerInvalidApplication:          "Invalid parameters",
		ErrorSellerShopNameRequired:            "Shop name is required",
		ErrorSellerShopNameTooLong:             "Shop name cannot exceed 80 characters",
		ErrorSellerDescriptionTooLong:          "Shop description cannot exceed 500 characters",
		ErrorSellerRejectReasonMissing:         "Rejection reason is required",
		ErrorSellerAuditStatusInvalid:          "Invalid seller audit status",
		ErrorSellerPayKeyRequired:              "Please set a payment password before listing products",
		ErrorProductSellerNotApproved:          "Product seller has not passed seller onboarding review",
		ErrorCarouselProductRequired:           "Please select a related product",
		ErrorCarouselProductNotExist:           "Related product does not exist",
		ErrorSettlementInvalidAmount:           "Invalid settlement amount",
		ErrorSettlementInvalidRate:             "Invalid commission rate",
		ErrorSettlementSellerInvalid:           "Seller does not exist or has not been approved",
		ErrorSettlementStatusInvalid:           "Settlement status does not allow this operation",
		ErrorOrderPayStatusInvalid:             "Order has been paid or cannot be paid in its current status",
		ErrorPaymentPayKeyRequired:             "Please set a payment password first",
		ErrorPaymentPayKeyInvalid:              "Invalid payment password",
		ErrorPaymentBalanceInsufficient:        "Insufficient coin balance",
		ErrorPaymentStockInsufficient:          "Insufficient stock",
		ErrorRefundStatusInvalid:               "Order status does not allow refund approval",
		ErrorRefundAmountInvalid:               "Invalid refund amount",
		ErrorRefundNotFound:                    "Refund request does not exist",
		ErrorOrderSelfPurchaseForbidden:        "You cannot buy your own product",
		ErrorSellerWithdrawAmountInvalid:       "Invalid withdrawal amount",
		ErrorSellerWithdrawPayeeRequired:       "Withdrawal payee information is required",
		ErrorSellerWithdrawInsufficientBalance: "Insufficient withdrawable balance",
		ErrorSellerWithdrawStatusInvalid:       "Withdrawal status does not allow this operation",
		ErrorSellerWithdrawReasonMissing:       "Rejection reason is required",
		ErrorSellerWithdrawNotFound:            "Withdrawal request does not exist",
		ErrorOrderStatusTransitionInvalid:      "Order status does not allow this operation",
	},
}

var MsgKeys = map[int]string{
	SUCCESS:               "common.ok",
	UpdatePasswordSuccess: "user.password_updated",
	NotExistInentifier:    "auth.third_party_not_bound",
	ERROR:                 "common.error",
	InvalidParams:         "common.invalid_params",

	ErrorExistNick:          "user.nickname_exists",
	ErrorExistUser:          "user.username_exists",
	ErrorNotExistUser:       "user.not_found",
	ErrorNotCompare:         "auth.account_password_invalid",
	ErrorNotComparePassword: "auth.password_confirm_mismatch",
	ErrorFailEncryption:     "common.encrypt_failed",
	ErrorNotExistProduct:    "product.not_found",
	ErrorNotExistAddress:    "address.not_found",
	ErrorExistFavorite:      "favorite.exists",
	ErrorUserNotFound:       "user.not_found",

	ErrorBossCheckTokenFail:        "seller.token_check_failed",
	ErrorBossCheckTokenTimeout:     "seller.token_timeout",
	ErrorBossToken:                 "seller.token_create_failed",
	ErrorBoss:                      "seller.token_invalid",
	ErrorBossInsufficientAuthority: "seller.insufficient_authority",
	ErrorBossProduct:               "seller.product_file_read_failed",

	ErrorProductExistCart: "cart.product_exists",
	ErrorProductMoreCart:  "cart.quantity_limit_exceeded",

	ErrorAuthCheckTokenFail:        "auth.token_check_failed",
	ErrorAuthCheckTokenTimeout:     "auth.token_timeout",
	ErrorAuthToken:                 "auth.token_create_failed",
	ErrorAuth:                      "auth.token_invalid",
	ErrorAuthInsufficientAuthority: "auth.insufficient_authority",
	ErrorReadFile:                  "file.read_failed",
	ErrorCallApi:                   "common.call_api_failed",
	ErrorUnmarshalJson:             "json.unmarshal_failed",
	ErrorAdminFindUser:             "admin.user_query_failed",

	ErrorDatabase:   "database.error",
	ErrorOss:        "oss.config_error",
	ErrorUploadFile: "file.upload_failed",

	ErrorSellerNotApplied:                  "seller.not_applied",
	ErrorSellerAuditPending:                "seller.audit_pending",
	ErrorSellerAlreadyApproved:             "seller.already_approved",
	ErrorSellerBanned:                      "seller.banned",
	ErrorSellerNotApproved:                 "seller.not_approved",
	ErrorSellerInvalidStatus:               "seller.invalid_status",
	ErrorSellerInvalidApplication:          "seller.invalid_application",
	ErrorSellerShopNameRequired:            "seller.shop_name_required",
	ErrorSellerShopNameTooLong:             "seller.shop_name_too_long",
	ErrorSellerDescriptionTooLong:          "seller.description_too_long",
	ErrorSellerRejectReasonMissing:         "seller.reject_reason_missing",
	ErrorSellerAuditStatusInvalid:          "seller.audit_status_invalid",
	ErrorSellerPayKeyRequired:              "seller.pay_key_required",
	ErrorProductSellerNotApproved:          "product.seller_not_approved",
	ErrorCarouselProductRequired:           "carousel.product_required",
	ErrorCarouselProductNotExist:           "carousel.product_not_exist",
	ErrorSettlementInvalidAmount:           "settlement.invalid_amount",
	ErrorSettlementInvalidRate:             "settlement.invalid_rate",
	ErrorSettlementSellerInvalid:           "settlement.seller_invalid",
	ErrorSettlementStatusInvalid:           "settlement.status_invalid",
	ErrorOrderPayStatusInvalid:             "order.pay_status_invalid",
	ErrorPaymentPayKeyRequired:             "payment.pay_key_required",
	ErrorPaymentPayKeyInvalid:              "payment.pay_key_invalid",
	ErrorPaymentBalanceInsufficient:        "payment.balance_insufficient",
	ErrorPaymentStockInsufficient:          "payment.stock_insufficient",
	ErrorRefundStatusInvalid:               "refund.status_invalid",
	ErrorRefundAmountInvalid:               "refund.amount_invalid",
	ErrorRefundNotFound:                    "refund.not_found",
	ErrorOrderSelfPurchaseForbidden:        "order.self_purchase_forbidden",
	ErrorSellerWithdrawAmountInvalid:       "seller_withdraw.amount_invalid",
	ErrorSellerWithdrawPayeeRequired:       "seller_withdraw.payee_required",
	ErrorSellerWithdrawInsufficientBalance: "seller_withdraw.insufficient_balance",
	ErrorSellerWithdrawStatusInvalid:       "seller_withdraw.status_invalid",
	ErrorSellerWithdrawReasonMissing:       "seller_withdraw.reason_missing",
	ErrorSellerWithdrawNotFound:            "seller_withdraw.not_found",
	ErrorOrderStatusTransitionInvalid:      "order.status_transition_invalid",
}

// GetMsg 获取状态码对应信息
func GetMsg(code int) string {
	return GetMsgByLocale(code, DefaultLocale)
}

func GetMsgKey(code int) string {
	key, ok := MsgKeys[code]
	if ok {
		return key
	}
	return MsgKeys[ERROR]
}

func GetMsgByLocale(code int, locale string) string {
	messages, ok := LocaleMsgFlags[NormalizeLocale(locale)]
	if !ok {
		messages = MsgFlags
	}
	msg, ok := messages[code]
	if ok {
		return msg
	}
	return messages[ERROR]
}

func NormalizeLocale(locale string) string {
	switch locale {
	case "zh", "zh-CN", "zh-Hans", "zh_CN":
		return DefaultLocale
	case "en", "en-US", "en_US":
		return EnglishLocale
	default:
		return DefaultLocale
	}
}
