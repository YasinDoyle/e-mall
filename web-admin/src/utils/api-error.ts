export const ApiErrorCode = {
  SUCCESS: 200,
  UPDATE_PASSWORD_SUCCESS: 201,
  NOT_EXIST_IDENTIFIER: 202,
  ERROR: 500,
  INVALID_PARAMS: 400,
  ERROR_EXIST_NICK: 10001,
  ERROR_EXIST_USER: 10002,
  ERROR_NOT_EXIST_USER: 10003,
  ERROR_NOT_COMPARE: 10004,
  ERROR_NOT_COMPARE_PASSWORD: 10005,
  ERROR_FAIL_ENCRYPTION: 10006,
  ERROR_NOT_EXIST_PRODUCT: 10007,
  ERROR_NOT_EXIST_ADDRESS: 10008,
  ERROR_EXIST_FAVORITE: 10009,
  ERROR_USER_NOT_FOUND: 10010,
  ERROR_BOSS_CHECK_TOKEN_FAIL: 20001,
  ERROR_BOSS_CHECK_TOKEN_TIMEOUT: 20002,
  ERROR_BOSS_TOKEN: 20003,
  ERROR_BOSS: 20004,
  ERROR_BOSS_INSUFFICIENT_AUTHORITY: 20005,
  ERROR_BOSS_PRODUCT: 20006,
  ERROR_PRODUCT_EXIST_CART: 20007,
  ERROR_PRODUCT_MORE_CART: 20008,
  ERROR_AUTH_CHECK_TOKEN_FAIL: 30001,
  ERROR_AUTH_CHECK_TOKEN_TIMEOUT: 30002,
  ERROR_AUTH_TOKEN: 30003,
  ERROR_AUTH: 30004,
  ERROR_AUTH_INSUFFICIENT_AUTHORITY: 30005,
  ERROR_READ_FILE: 30006,
  ERROR_SEND_EMAIL: 30007,
  ERROR_CALL_API: 30008,
  ERROR_UNMARSHAL_JSON: 30009,
  ERROR_ADMIN_FIND_USER: 30010,
  ERROR_DATABASE: 40001,
  ERROR_OSS: 50001,
  ERROR_UPLOAD_FILE: 50002,
} as const;

export const ApiErrorMessage: Record<number, string> = {
  [ApiErrorCode.SUCCESS]: "ok",
  [ApiErrorCode.UPDATE_PASSWORD_SUCCESS]: "修改密码成功",
  [ApiErrorCode.NOT_EXIST_IDENTIFIER]: "该第三方账号未绑定",
  [ApiErrorCode.ERROR]: "fail",
  [ApiErrorCode.INVALID_PARAMS]: "请求参数错误",
  [ApiErrorCode.ERROR_EXIST_NICK]: "已存在该昵称",
  [ApiErrorCode.ERROR_EXIST_USER]: "已存在该用户名",
  [ApiErrorCode.ERROR_NOT_EXIST_USER]: "该用户不存在",
  [ApiErrorCode.ERROR_NOT_COMPARE]: "账号密码错误",
  [ApiErrorCode.ERROR_NOT_COMPARE_PASSWORD]: "两次密码输入不一致",
  [ApiErrorCode.ERROR_FAIL_ENCRYPTION]: "加密失败",
  [ApiErrorCode.ERROR_NOT_EXIST_PRODUCT]: "该商品不存在",
  [ApiErrorCode.ERROR_NOT_EXIST_ADDRESS]: "该收获地址不存在",
  [ApiErrorCode.ERROR_EXIST_FAVORITE]: "已收藏该商品",
  [ApiErrorCode.ERROR_USER_NOT_FOUND]: "用户不存在",
  [ApiErrorCode.ERROR_BOSS_CHECK_TOKEN_FAIL]: "商家的Token鉴权失败",
  [ApiErrorCode.ERROR_BOSS_CHECK_TOKEN_TIMEOUT]: "商家的Token已超时",
  [ApiErrorCode.ERROR_BOSS_TOKEN]: "商家的Token生成失败",
  [ApiErrorCode.ERROR_BOSS]: "商家Token错误",
  [ApiErrorCode.ERROR_BOSS_INSUFFICIENT_AUTHORITY]: "商家权限不足",
  [ApiErrorCode.ERROR_BOSS_PRODUCT]: "商家读文件错误",
  [ApiErrorCode.ERROR_PRODUCT_EXIST_CART]: "商品已经在购物车了，数量+1",
  [ApiErrorCode.ERROR_PRODUCT_MORE_CART]: "超过最大上限",
  [ApiErrorCode.ERROR_AUTH_CHECK_TOKEN_FAIL]: "Token鉴权失败",
  [ApiErrorCode.ERROR_AUTH_CHECK_TOKEN_TIMEOUT]: "Token已超时",
  [ApiErrorCode.ERROR_AUTH_TOKEN]: "Token生成失败",
  [ApiErrorCode.ERROR_AUTH]: "Token错误",
  [ApiErrorCode.ERROR_AUTH_INSUFFICIENT_AUTHORITY]: "权限不足",
  [ApiErrorCode.ERROR_READ_FILE]: "读文件失败",
  [ApiErrorCode.ERROR_SEND_EMAIL]: "发送邮件失败",
  [ApiErrorCode.ERROR_CALL_API]: "调用接口失败",
  [ApiErrorCode.ERROR_UNMARSHAL_JSON]: "解码JSON失败",
  [ApiErrorCode.ERROR_ADMIN_FIND_USER]: "管理员查询用户失败",
  [ApiErrorCode.ERROR_DATABASE]: "数据库操作出错,请重试",
  [ApiErrorCode.ERROR_OSS]: "OSS配置错误",
  [ApiErrorCode.ERROR_UPLOAD_FILE]: "上传失败",
};

function firstText(...values: unknown[]): string {
  for (const value of values) {
    if (typeof value === "string" && value.trim()) {
      return value.trim();
    }
  }
  return "";
}

export function resolveApiErrorMessage(payload: any, fallback = "请求失败") {
  const status = Number(payload?.status ?? payload?.code);
  const detail = firstText(payload?.data, payload?.error, payload?.msg);
  if (status && ApiErrorMessage[status]) {
    return detail || ApiErrorMessage[status];
  }
  return detail || fallback;
}
