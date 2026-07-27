package routes

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	api "github.com/YasinDoyle/e-mall/api/v1"
	"github.com/YasinDoyle/e-mall/middleware"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("something-very-secret"))
	r.Use(middleware.CORS(), middleware.Jaeger())
	r.Use(sessions.Sessions("mysession", store))
	r.StaticFS("/static", http.Dir("./static"))
	v1 := r.Group("/api/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "OK",
			})
		})

		//用户操作
		v1.POST("user/register/email/code", api.RegisterEmailCodeHandler())
		v1.POST("user/register", api.UserRegisterHandler())
		v1.POST("user/login", api.UserLoginHandler())

		// 商品操作
		v1.GET("product/list", api.ListProductsHandler())
		v1.GET("product/show", api.ShowProductHandler())
		v1.POST("product/search", api.SearchProductsHandler())
		v1.GET("product/imgs/list", api.ListProductImgHandler()) // 商品图片
		v1.GET("category/list", api.ListCategoryHandler())       // 商品分类
		v1.GET("carousel/list", api.ListCarouselsHandler())      // 轮播图
		v1.GET("carousels", api.ListCarouselsHandler())          // 前端兼容别名
		// 商品评价（公开）
		v1.GET("product/reviews", api.ListReviewsHandler())

		// 优惠券（公开可查）
		v1.GET("coupons", api.CouponListHandler())

		// 支付宝/微信回调（无需 JWT）
		v1.POST("pay/wechat/notify", api.WechatNotifyHandler())
		v1.POST("pay/alipay/notify", api.AlipayNotifyHandler())

		authed := v1.Group("/") //需要用户保护

		authed.Use(middleware.AuthMiddleware())
		{
			// 用户操作
			authed.POST("user/update", api.UserUpdateHandler())
			authed.GET("user/show_info", api.ShowUserInfoHandler())
			authed.POST("user/following", api.UserFollowingHandler())
			authed.POST("user/unfollowing", api.UserUnFollowingHandler())
			authed.POST("user/avatar", api.UploadAvatarHandler()) // 上传头像

			// 商品操作
			authed.POST("product/create", api.CreateProductHandler())
			authed.POST("product/update", api.UpdateProductHandler())
			authed.POST("product/delete", api.DeleteProductHandler())

			// 卖家中心：查看自己的商品、上架/下架
			authed.GET("boss/product/list", api.BossProductListHandler())
			authed.POST("boss/product/on_sale", api.BossProductOnSaleHandler())
			authed.GET("boss/order/list", api.SellerListOrdersHandler())
			authed.GET("boss/settlement/summary", api.SellerSettlementSummaryHandler())
			authed.GET("seller/account/summary", api.SellerAccountSummaryHandler())
			authed.GET("seller/withdraw/list", api.SellerWithdrawListHandler())
			authed.POST("seller/withdraw/apply", api.SellerWithdrawApplyHandler())

			// 商品评价（需登录）
			authed.POST("reviews/create", api.CreateReviewHandler())
			authed.POST("reviews/upload", api.UploadReviewImageHandler())

			// 优惠券（需登录）
			authed.POST("coupon/claim", api.CouponClaimHandler())
			authed.GET("coupon/list", api.UserCouponListHandler())
			// 收藏夹
			authed.GET("favorites/list", api.ListFavoritesHandler())
			authed.POST("favorites/create", api.CreateFavoriteHandler())
			authed.POST("favorites/delete", api.DeleteFavoriteHandler())

			// 订单操作
			authed.POST("orders/create", api.CreateOrderHandler())
			authed.GET("orders/list", api.ListOrdersHandler())
			authed.GET("orders/show", api.ShowOrderHandler())
			authed.POST("orders/delete", api.DeleteOrderHandler())
			authed.POST("orders/ship", api.ShipOrderHandler())
			authed.POST("orders/receive", api.ReceiveOrderHandler())
			authed.POST("orders/refund/request", api.RefundRequestOrderHandler())

			// 购物车
			authed.POST("carts/create", api.CreateCartHandler())
			authed.GET("carts/list", api.ListCartHandler())
			authed.POST("carts/update", api.UpdateCartHandler()) // 购物车id
			authed.POST("carts/delete", api.DeleteCartHandler())

			// 收获地址操作
			authed.POST("addresses/create", api.CreateAddressHandler())
			authed.GET("addresses/show", api.ShowAddressHandler())
			authed.GET("addresses/list", api.ListAddressHandler())
			authed.POST("addresses/update", api.UpdateAddressHandler())
			authed.POST("addresses/delete", api.DeleteAddressHandler())

			// 支付功能
			authed.POST("paydown", api.OrderPaymentHandler())

			// 显示金额
			authed.POST("money", api.ShowMoneyHandler())
			authed.POST("money/pay-key", api.SetPayKeyHandler())

			// 商家入驻
			authed.POST("seller/apply", api.SellerApplyHandler())
			authed.GET("seller/profile", api.SellerProfileHandler())

			// 秒杀专场
			authed.POST("flash_sale/init", api.InitFlashSaleHandler())
			authed.GET("flash_sale/list", api.ListFlashSaleHandler())
			authed.GET("flash_sale/show", api.GetFlashSaleHandler())
			authed.POST("flash_sale/skill", api.FlashSaleHandler())

			// 充值（需登录）
			authed.POST("recharge/wechat", api.WechatRechargeHandler())
			authed.POST("recharge/alipay", api.AlipayRechargeHandler())
			authed.GET("recharge/status", api.RechargeStatusHandler())
			authed.GET("recharge/pending", api.GetPendingCreditHandler())
			authed.POST("recharge/apply", api.ApplyPendingCreditHandler())
		}

		// 管理员路由（需登录 + IsAdmin）
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(), middleware.AdminAuthMiddleware())
		{
			// 分类管理
			admin.GET("category/list", api.ListCategoryHandler())
			admin.POST("category/create", api.AdminCategoryCreateHandler())
			admin.POST("category/update", api.AdminCategoryUpdateHandler())
			admin.POST("category/delete", api.AdminCategoryDeleteHandler())

			// 轮播图管理
			admin.GET("carousel/list", api.AdminCarouselListHandler())
			admin.POST("carousel/create", api.AdminCarouselCreateHandler())
			admin.POST("carousel/upload", api.AdminCarouselUploadHandler())
			admin.POST("carousel/update", api.AdminCarouselUpdateHandler())
			admin.POST("carousel/delete", api.AdminCarouselDeleteHandler())

			// 公告管理
			admin.GET("notice/list", api.AdminNoticeListHandler())
			admin.POST("notice/create", api.AdminNoticeCreateHandler())
			admin.POST("notice/update", api.AdminNoticeUpdateHandler())
			admin.POST("notice/delete", api.AdminNoticeDeleteHandler())

			// 用户管理
			admin.GET("user/list", api.AdminUserListHandler())
			admin.GET("users", api.AdminUserListHandler())
			admin.POST("user/ban", api.AdminUserBanHandler())
			admin.POST("users/ban", api.AdminUserBanHandler())

			// 商家管理
			admin.GET("seller/list", api.AdminSellerListHandler())
			admin.POST("seller/audit", api.AdminSellerAuditHandler())
			admin.GET("seller/withdraw/list", api.AdminSellerWithdrawListHandler())
			admin.POST("seller/withdraw/audit", api.AdminSellerWithdrawAuditHandler())
			admin.POST("seller/withdraw/paid", api.AdminSellerWithdrawPaidHandler())
			admin.GET("seller/withdraw/detail", api.AdminSellerWithdrawDetailHandler())

			// 商品审核
			admin.GET("product/list", api.AdminProductListHandler())
			admin.POST("product/audit", api.AdminProductAuditHandler())
			admin.POST("product/delete", api.AdminProductDeleteHandler())

			// 统计
			admin.GET("stats/overview", api.AdminStatsOverviewHandler())
			admin.GET("stats/orders", api.AdminStatsOrdersHandler())

			// 评价管理
			admin.POST("review/delete", api.AdminDeleteReviewHandler())

			// 优惠券管理
			admin.GET("coupon/list", api.AdminCouponListHandler())
			admin.POST("coupon/create", api.AdminCouponCreateHandler())
			admin.POST("coupon/offline", api.AdminCouponOfflineHandler())

			// 订单管理/售后
			admin.GET("orders/list", api.AdminListOrdersHandler())
			admin.POST("orders/refund/approve", api.AdminRefundApproveOrderHandler())

			// 结算管理
			admin.GET("settlement/list", api.AdminSettlementListHandler())
			admin.POST("settlement/generate", api.AdminSettlementGenerateHandler())
			admin.POST("settlement/generate_one", api.AdminSettlementGenerateOneHandler())
			admin.POST("settlement/paid", api.AdminSettlementPaidHandler())
			admin.GET("settlement/detail", api.AdminSettlementDetailHandler())
			admin.POST("settlement/backfill", api.AdminSettlementBackfillHandler())

			// 秒杀管理
			admin.GET("flash-sale/list", api.AdminListFlashSaleHandler())
			admin.POST("flash-sale/create", api.AdminCreateFlashSaleHandler())
			admin.POST("flash-sale/update", api.AdminUpdateFlashSaleHandler())
			admin.POST("flash-sale/delete", api.AdminDeleteFlashSaleHandler())

			// 充值退款
			admin.POST("recharge/wechat/refund", api.WechatRefundHandler())
			admin.POST("recharge/alipay/refund", api.AlipayRefundHandler())
		}
	}
	return r
}
