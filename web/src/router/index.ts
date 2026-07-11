import { createRouter, createWebHistory } from "vue-router";
import { useUserStore } from "@/stores/user";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    // 认证
    {
      path: "/login",
      component: () => import("@/views/auth/LoginView.vue"),
      meta: { guest: true },
    },
    {
      path: "/register",
      component: () => import("@/views/auth/RegisterView.vue"),
      meta: { guest: true },
    },
    {
      path: "/valid-email",
      component: () => import("@/views/auth/ValidEmailView.vue"),
    },

    // 主布局（含 NavBar）
    {
      path: "/",
      component: () => import("@/components/layout/DefaultLayout.vue"),
      children: [
        { path: "", component: () => import("@/views/home/HomeView.vue") },
        {
          path: "products",
          component: () => import("@/views/product/ProductListView.vue"),
        },
        {
          path: "product/:id",
          component: () => import("@/views/product/ProductDetailView.vue"),
        },
        {
          path: "search",
          component: () => import("@/views/product/SearchView.vue"),
        },
        {
          path: "flash-sale",
          component: () => import("@/views/flash-sale/FlashSaleListView.vue"),
        },
        {
          path: "flash-sale/:id",
          component: () => import("@/views/flash-sale/FlashSaleDetailView.vue"),
        },

        // 需要登录
        {
          path: "cart",
          component: () => import("@/views/cart/CartView.vue"),
          meta: { auth: true },
        },
        {
          path: "checkout",
          component: () => import("@/views/checkout/CheckoutView.vue"),
          meta: { auth: true },
        },
        {
          path: "payment",
          component: () => import("@/views/checkout/PaymentView.vue"),
          meta: { auth: true },
        },
        {
          path: "order/success",
          component: () => import("@/views/checkout/OrderSuccessView.vue"),
          meta: { auth: true },
        },

        // 用户中心
        {
          path: "user",
          component: () => import("@/views/user/UserLayout.vue"),
          meta: { auth: true },
          children: [
            { path: "", redirect: "/user/profile" },
            {
              path: "profile",
              component: () => import("@/views/user/ProfileView.vue"),
            },
            {
              path: "orders",
              component: () => import("@/views/user/OrderListView.vue"),
            },
            {
              path: "orders/:id",
              component: () => import("@/views/user/OrderDetailView.vue"),
            },
            {
              path: "addresses",
              component: () => import("@/views/user/AddressView.vue"),
            },
            {
              path: "favorites",
              component: () => import("@/views/user/FavoriteView.vue"),
            },
            {
              path: "coupons",
              component: () => import("@/views/user/CouponView.vue"),
            },
            {
              path: "wallet",
              component: () => import("@/views/user/WalletView.vue"),
            },
          ],
        },
      ],
    },

    {
      path: "/:pathMatch(.*)*",
      component: () => import("@/views/NotFoundView.vue"),
    },
  ],
});

// 全局路由守卫
router.beforeEach((to) => {
  const userStore = useUserStore();
  if (to.meta.auth && !userStore.isLoggedIn) {
    return { path: "/login", query: { redirect: to.fullPath } };
  }
  if (to.meta.guest && userStore.isLoggedIn) {
    return { path: "/" };
  }
});

export default router;
