import assert from "node:assert/strict";
import { describe, it, beforeEach } from "node:test";
import { readFileSync, readdirSync, statSync } from "node:fs";
import { resolve } from "node:path";
import ts from "typescript";

const root = resolve(import.meta.dirname, "..");

function listVueFiles(dir) {
  const files = [];
  for (const entry of readdirSync(dir)) {
    const fullPath = resolve(dir, entry);
    if (statSync(fullPath).isDirectory()) {
      files.push(...listVueFiles(fullPath));
    } else if (fullPath.endsWith(".vue")) {
      files.push(fullPath);
    }
  }
  return files;
}

class MemoryStorage {
  constructor(initial = {}) {
    this.store = new Map(Object.entries(initial));
  }
  getItem(key) {
    return this.store.has(key) ? this.store.get(key) : null;
  }
  setItem(key, value) {
    this.store.set(key, String(value));
  }
  removeItem(key) {
    this.store.delete(key);
  }
  clear() {
    this.store.clear();
  }
}

async function loadTranspiledModule(sourcePath) {
  const source = readFileSync(sourcePath, "utf8");
  const transpiled = ts.transpileModule(source, {
    compilerOptions: {
      target: ts.ScriptTarget.ES2020,
      module: ts.ModuleKind.ESNext,
      sourceMap: false,
    },
    fileName: sourcePath,
  }).outputText;
  const moduleUrl = `data:text/javascript;base64,${Buffer.from(transpiled).toString("base64")}`;
  return import(moduleUrl);
}

async function loadUserSessionModule() {
  return loadTranspiledModule(resolve(root, "src/utils/session.ts"));
}

async function loadAdminSessionModule() {
  return loadTranspiledModule(resolve(root, "../web-admin/src/utils/session.ts"));
}

function seedUserSessions() {
  localStorage.setItem(
    "mall:user:sessions",
    JSON.stringify({
      "1": {
        id: "1",
        token: "token-a",
        refreshToken: "refresh-a",
        userInfo: { id: 1, user_name: "account_a" },
        data: { sellerProfile: "seller-a" },
      },
      "2": {
        id: "2",
        token: "token-b",
        refreshToken: "refresh-b",
        userInfo: { id: 2, user_name: "account_b" },
        data: {},
      },
    }),
  );
  localStorage.setItem("mall:user:last_active_session", "1");
}

describe("Post-P1 multi-account tab isolation", () => {
  beforeEach(() => {
    globalThis.localStorage = new MemoryStorage();
    globalThis.sessionStorage = new MemoryStorage();
  });

  it("does not auto-login a fresh tab from the last active account", async () => {
    seedUserSessions();
    const session = await loadUserSessionModule();

    assert.equal(session.getActiveUserToken(), "");
    assert.equal(session.getActiveUserInfo(), null);
  });

  it("logs out only the current tab and keeps shared account sessions", async () => {
    seedUserSessions();
    sessionStorage.setItem("mall:user:active_session", "1");
    const session = await loadUserSessionModule();

    session.clearActiveUserSession();

    const sessions = JSON.parse(localStorage.getItem("mall:user:sessions"));
    assert.equal(sessions["1"].token, "token-a");
    assert.equal(sessions["2"].token, "token-b");
    assert.equal(session.getActiveUserToken(), "");
  });

  it("keeps user request token selection isolated across two tab-scoped sessions", async () => {
    seedUserSessions();
    const tabA = new MemoryStorage({ "mall:user:active_session": "1" });
    const tabB = new MemoryStorage({ "mall:user:active_session": "2" });
    const session = await loadUserSessionModule();

    globalThis.sessionStorage = tabA;
    assert.equal(session.getActiveUserToken(), "token-a");
    assert.equal(session.getActiveUserRefreshToken(), "refresh-a");

    globalThis.sessionStorage = tabB;
    assert.equal(session.getActiveUserToken(), "token-b");
    assert.equal(session.getActiveUserRefreshToken(), "refresh-b");

    globalThis.sessionStorage = tabA;
    assert.equal(session.getActiveUserToken(), "token-a");
  });

  it("does not expose a stale pending user login token in a fresh tab", async () => {
    localStorage.setItem(
      "mall:user:sessions",
      JSON.stringify({
        pending: {
          id: "pending",
          token: "stale-login-token",
          refreshToken: "stale-refresh-token",
          userInfo: null,
          data: {},
        },
      }),
    );
    const session = await loadUserSessionModule();

    assert.equal(session.getActiveUserToken(), "");
    assert.equal(session.getActiveUserInfo(), null);
    assert.deepEqual(JSON.parse(localStorage.getItem("mall:user:sessions")), {});
  });

  it("does not expose stale pending user token after the tab already initialized pending", async () => {
    const session = await loadUserSessionModule();
    assert.equal(session.getActiveUserToken(), "");
    assert.equal(sessionStorage.getItem("mall:user:active_session"), "pending");
    localStorage.setItem(
      "mall:user:sessions",
      JSON.stringify({
        pending: {
          id: "pending",
          token: "late-stale-login-token",
          refreshToken: "late-stale-refresh-token",
          userInfo: null,
          data: {},
        },
      }),
    );

    assert.equal(session.getActiveUserToken(), "");
    assert.deepEqual(JSON.parse(localStorage.getItem("mall:user:sessions")), {});
  });

  it("allows pending user token only during the current tab login flow", async () => {
    const session = await loadUserSessionModule();

    session.beginUserLoginSession();
    session.setActiveUserTokens("login-token", "login-refresh");

    assert.equal(session.getActiveUserToken(), "login-token");
    assert.equal(session.getActiveUserRefreshToken(), "login-refresh");
  });

  it("ignores legacy user localStorage login keys at runtime", async () => {
    localStorage.setItem("token", "legacy-token");
    localStorage.setItem("refreshToken", "legacy-refresh-token");
    localStorage.setItem(
      "userInfo",
      JSON.stringify({ id: 7, user_name: "legacy_user" }),
    );
    localStorage.setItem("sellerProfile", JSON.stringify({ shop_name: "legacy shop" }));
    const session = await loadUserSessionModule();

    assert.equal(session.getActiveUserToken(), "");
    assert.equal(session.getActiveUserInfo(), null);
    assert.equal(localStorage.getItem("mall:user:sessions"), null);
  });

  it("does not expose a stale pending admin login token in a fresh tab", async () => {
    localStorage.setItem(
      "mall:admin:sessions",
      JSON.stringify({
        pending: {
          id: "pending",
          token: "stale-admin-token",
          refreshToken: "stale-admin-refresh-token",
          adminInfo: null,
          data: {},
        },
      }),
    );
    const session = await loadAdminSessionModule();

    assert.equal(session.getActiveAdminToken(), "");
    assert.equal(session.getActiveAdminInfo(), null);
    assert.deepEqual(JSON.parse(localStorage.getItem("mall:admin:sessions")), {});
  });

  it("does not expose stale pending admin token after the tab already initialized pending", async () => {
    const session = await loadAdminSessionModule();
    assert.equal(session.getActiveAdminToken(), "");
    assert.equal(sessionStorage.getItem("mall:admin:active_session"), "pending");
    localStorage.setItem(
      "mall:admin:sessions",
      JSON.stringify({
        pending: {
          id: "pending",
          token: "late-stale-admin-token",
          refreshToken: "late-stale-admin-refresh-token",
          adminInfo: null,
          data: {},
        },
      }),
    );

    assert.equal(session.getActiveAdminToken(), "");
    assert.deepEqual(JSON.parse(localStorage.getItem("mall:admin:sessions")), {});
  });

  it("allows pending admin token only during the current tab login flow", async () => {
    const session = await loadAdminSessionModule();

    session.beginAdminLoginSession();
    session.setActiveAdminTokens("admin-login-token", "admin-login-refresh");

    assert.equal(session.getActiveAdminToken(), "admin-login-token");
    assert.equal(session.getActiveAdminRefreshToken(), "admin-login-refresh");
  });

  it("ignores legacy admin localStorage login keys at runtime", async () => {
    localStorage.setItem("admin_token", "legacy-admin-token");
    localStorage.setItem("admin_refresh_token", "legacy-admin-refresh-token");
    localStorage.setItem(
      "admin_info",
      JSON.stringify({ id: 8, user_name: "legacy_admin", nick_name: "Legacy" }),
    );
    const session = await loadAdminSessionModule();

    assert.equal(session.getActiveAdminToken(), "");
    assert.equal(session.getActiveAdminInfo(), null);
    assert.equal(localStorage.getItem("mall:admin:sessions"), null);
  });

  it("keeps admin request token selection isolated across two tab-scoped sessions", async () => {
    localStorage.setItem(
      "mall:admin:sessions",
      JSON.stringify({
        "1": {
          id: "1",
          token: "admin-token-a",
          refreshToken: "admin-refresh-a",
          adminInfo: { id: 1, user_name: "admin_a", nick_name: "Admin A" },
          data: {},
        },
        "2": {
          id: "2",
          token: "admin-token-b",
          refreshToken: "admin-refresh-b",
          adminInfo: { id: 2, user_name: "admin_b", nick_name: "Admin B" },
          data: {},
        },
      }),
    );
    const tabA = new MemoryStorage({ "mall:admin:active_session": "1" });
    const tabB = new MemoryStorage({ "mall:admin:active_session": "2" });
    const session = await loadAdminSessionModule();

    globalThis.sessionStorage = tabA;
    assert.equal(session.getActiveAdminToken(), "admin-token-a");
    assert.equal(session.getActiveAdminRefreshToken(), "admin-refresh-a");

    globalThis.sessionStorage = tabB;
    assert.equal(session.getActiveAdminToken(), "admin-token-b");
    assert.equal(session.getActiveAdminRefreshToken(), "admin-refresh-b");

    globalThis.sessionStorage = tabA;
    assert.equal(session.getActiveAdminToken(), "admin-token-a");
  });
});

describe("Runtime legacy compatibility policy", () => {
  it("keeps user and admin session utilities free of legacy migration helpers", () => {
    const userSession = readFileSync(resolve(root, "src/utils/session.ts"), "utf8");
    const adminSession = readFileSync(resolve(root, "../web-admin/src/utils/session.ts"), "utf8");

    for (const source of [userSession, adminSession]) {
      assert.doesNotMatch(source, /migrateLegacy/i);
      assert.doesNotMatch(source, /clearLegacy/i);
    }
    assert.doesNotMatch(userSession, /localStorage\.getItem\("token"\)/);
    assert.doesNotMatch(userSession, /localStorage\.getItem\("userInfo"\)/);
    assert.doesNotMatch(userSession, /localStorage\.getItem\("sellerProfile"\)/);
    assert.doesNotMatch(adminSession, /localStorage\.getItem\("admin_token"\)/);
    assert.doesNotMatch(adminSession, /localStorage\.getItem\("admin_info"\)/);
  });

  it("documents that runtime code should not carry old-version compatibility patches", () => {
    const source = readFileSync(resolve(root, "../AGENTS.md"), "utf8");

    assert.match(source, /不在运行时代码里做旧版兼容处理/);
    assert.match(source, /旧版升级.*迁移脚本/);
  });
});

describe("Post-P1 i18n acceptance coverage", () => {
  const cases = [
    ["src/views/home/HomeView.vue", ["商品分类", "秒杀专场", "限时秒杀", "好价商品限量开抢，先到先得", "热门商品", "进入专场", "立即查看", "查看全部"]],
    ["src/views/user/UserLayout.vue", ["个人资料", "我的订单", "收货地址", "我的收藏", "我的优惠券", "我的钱包", "消息通知", "卖家中心"]],
    ["src/views/user/ProfileView.vue", ["个人资料", "更换头像", "用户名", "昵称", "邮箱", "未绑定", "保存修改"]],
    ["src/views/seller/SellerLayout.vue", ["入驻状态", "资金账户", "商品管理", "订单管理", "发布商品"]],
    ["src/views/flash-sale/FlashSaleListView.vue", ["秒杀专场", "火热进行中", "距本场结束", "暂无秒杀商品", "逛逛普通商品", "剩余", "立即抢购", "已售罄"]],
    ["src/views/flash-sale/FlashSaleDetailView.vue", ["返回秒杀列表", "距本场结束", "限时特惠", "剩余库存", "收货地址", "支付密码", "登录后抢购", "活动说明", "每人限购", "抢购成功"]],
    ["src/views/cart/CartView.vue", ["购物车", "购物车是空的", "去逛逛", "单价", "数量", "小计", "全选", "删除选中", "已选", "结算"]],
    ["src/views/user/NotificationView.vue", ["消息通知", "只看未读", "全部已读", "暂无通知", "未读"]],
    ["src/views/user/OrderListView.vue", ["我的订单", "暂无订单", "订单号", "确认收货", "查看详情", "退款申请中"]],
    ["src/views/user/AddressView.vue", ["收货地址", "新增地址", "还没有收货地址", "编辑地址", "姓名", "手机号", "请填写完整信息"]],
    ["src/views/user/FavoriteView.vue", ["我的收藏", "收藏夹是空的", "取消收藏", "已取消收藏"]],
    ["src/views/user/CouponView.vue", ["我的优惠券", "可领取", "暂无优惠券", "有效期至", "无门槛", "领取成功"]],
    ["src/views/user/WalletView.vue", ["我的钱包", "账户余额", "支付密码", "充值", "刷新待入账", "钱包充值", "确认入账", "充值已入账"]],
    ["src/views/seller/SellerProductListView.vue", ["我的商品", "批量上架", "发布商品", "销售状态", "销售中", "下架", "批量上架完成"]],
    ["src/views/seller/SellerProductFormView.vue", ["编辑商品", "发布商品", "商品名称", "审核资料", "资质证书", "提交审核", "返回列表"]],
    ["src/views/checkout/CheckoutView.vue", ["确认订单", "收货地址", "商品清单", "价格汇总", "商品总价", "优惠券", "应付金额", "返回购物车", "提交订单", "选择收货地址", "选择优惠券"]],
    ["src/views/checkout/PaymentView.vue", ["支付订单", "订单信息", "支付方式", "支付密码", "余额支付", "微信充值", "支付宝充值", "立即支付", "确认入账"]],
    ["src/views/product/ProductDetailView.vue", ["卖家", "库存", "浏览量", "分类", "购买数量", "加入购物车", "收藏", "资质材料", "商品评价", "返回商品列表"]],
    ["src/views/product/ProductListView.vue", ["全部", "暂无商品"]],
    ["src/views/product/SearchView.vue", ["搜索商品", "搜索", "未找到相关商品"]],
    ["src/views/auth/RegisterView.vue", ["注册", "用户名", "昵称", "邮箱", "邮箱验证码", "密码", "确认密码", "发送验证码", "立即登录"]],
    ["src/views/NotFoundView.vue", ["抱歉，您访问的页面不存在", "返回首页", "返回上一页"]],
    ["src/views/checkout/OrderSuccessView.vue", ["支付成功", "感谢您的购买", "查看订单", "继续购物"]],
  ];

  for (const [file, texts] of cases) {
    it(`${file} moves accepted visible copy behind i18n`, () => {
      const source = readFileSync(resolve(root, file), "utf8");
      for (const text of texts) {
        assert.equal(source.includes(text), false, `${file} still contains ${text}`);
      }
    });
  }
});

describe("Project i18n working rule", () => {
  it("documents that frontend visible copy must update zh-CN and en-US together", () => {
    const source = readFileSync(resolve(root, "../AGENTS.md"), "utf8");
    assert.match(source, /前端.*文本.*zh-CN.*en-US/s);
  });
});

describe("Post-P1 admin i18n acceptance coverage", () => {
  it("keeps admin page components free of hardcoded Chinese copy", () => {
    const files = [
      ...listVueFiles(resolve(root, "../web-admin/src/views")),
      ...listVueFiles(resolve(root, "../web-admin/src/components")),
    ];

    for (const file of files) {
      const source = readFileSync(file, "utf8");
      assert.doesNotMatch(source, /\p{Script=Han}/u, file);
    }
  });
});

describe("NavBar dropdown interaction", () => {
  it("opens the avatar dropdown by click instead of hover", () => {
    const source = readFileSync(resolve(root, "src/components/common/NavBar.vue"), "utf8");
    assert.match(source, /<el-dropdown\s+trigger="click"\s+@command="handleCommand">/);
  });
});

describe("Post-P1 account switcher UI", () => {
  it("shows a saved-account switcher entry and picker in the user navbar", () => {
    const source = readFileSync(resolve(root, "src/components/common/NavBar.vue"), "utf8");
    assert.match(source, /switchAccount/);
    assert.match(source, /accountSwitch/);
    assert.match(source, /savedAccounts/);
  });

  it("shows a saved-account switcher entry and picker in the admin layout", () => {
    const source = readFileSync(resolve(root, "../web-admin/src/components/AdminLayout.vue"), "utf8");
    assert.match(source, /switchAccount/);
    assert.match(source, /accountSwitch/);
    assert.match(source, /savedAccounts/);
  });
});

describe("Post-P1 notification subscription policy", () => {
  it("uses SSE as the default primary channel", () => {
    const source = readFileSync(resolve(root, "src/stores/appConfig.ts"), "utf8");
    assert.match(source, /notification_sse:\s*true/);
  });

  it("does not start steady unread polling while starting SSE", () => {
    const source = readFileSync(resolve(root, "src/components/common/NavBar.vue"), "utf8");
    assert.doesNotMatch(source, /onMounted\(\(\)\s*=>\s*{[^}]*startNotificationPolling\(\)/s);
  });
});

describe("Post-P1 reported regression coverage", () => {
  it("orders product management actions as reject, list, delete with action-specific copy", () => {
    const source = readFileSync(resolve(root, "../web-admin/src/views/product/ProductView.vue"), "utf8");
    const rejectIndex = source.indexOf('t("page.product.rejectAction")');
    const listIndex = source.indexOf('t("page.product.listAction")');
    const deleteIndex = source.indexOf('t("common.delete")');

    assert.ok(rejectIndex > -1, "missing reject action label");
    assert.ok(listIndex > -1, "missing list action label");
    assert.ok(deleteIndex > -1, "missing delete action label");
    assert.ok(rejectIndex < listIndex, "reject action should appear before list action");
    assert.ok(listIndex < deleteIndex, "list action should appear before delete action");
    assert.match(source, /class="product-actions"/);
  });

  it("updates the in-memory cart badge after checkout payment succeeds or partially succeeds", () => {
    const source = readFileSync(resolve(root, "src/views/checkout/PaymentView.vue"), "utf8");

    assert.match(source, /useUserStore/);
    assert.match(source, /syncPaidCartCount/);
    assert.match(source, /userStore\.setCartCount/);
  });

  it("renders an admin avatar in the navbar account dropdown trigger", () => {
    const source = readFileSync(resolve(root, "../web-admin/src/components/AdminLayout.vue"), "utf8");

    assert.match(source, /<el-avatar/);
    assert.match(source, /UserFilled/);
    assert.match(source, /admin-account-trigger/);
  });

  it("keeps withdraw review action buttons spaced in a dedicated action group", () => {
    const source = readFileSync(resolve(root, "../web-admin/src/views/withdraw/WithdrawView.vue"), "utf8");

    assert.match(source, /class="withdraw-actions"/);
    assert.match(source, /\.withdraw-actions\s*{[^}]*display:\s*flex/s);
    assert.match(source, /\.withdraw-actions :deep\(\.el-button \+ \.el-button\)/);
  });

  it("shows admin menu pending badges from pending-only list totals", () => {
    const source = readFileSync(resolve(root, "../web-admin/src/components/AdminLayout.vue"), "utf8");

    assert.match(source, /pendingMenuCounts/);
    assert.match(source, /refreshPendingMenuCounts/);
    assert.match(source, /getAdminProductList\(\{\s*page_num:\s*1,\s*page_size:\s*1,\s*audit_status:\s*0/s);
    assert.match(source, /getAdminSellerList\(\{\s*page_num:\s*1,\s*page_size:\s*1,\s*status:\s*0/s);
    assert.match(source, /getAdminOrderList\(\{\s*page_num:\s*1,\s*page_size:\s*1,\s*refund_status:\s*1/s);
    assert.match(source, /getAdminSettlementList\(\{\s*page_num:\s*1,\s*page_size:\s*1,\s*status:\s*"pending"/s);
    assert.match(source, /getAdminSellerWithdrawList\(\{\s*page_num:\s*1,\s*page_size:\s*1,\s*status:\s*"pending"/s);
    assert.match(source, /pending-menu-badge/);
  });
});
