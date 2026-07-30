import { t } from "@/locales";

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

export const ApiErrorMessageKey: Record<number, string> = {
  [ApiErrorCode.SUCCESS]: "common.ok",
  [ApiErrorCode.UPDATE_PASSWORD_SUCCESS]: "user.password_updated",
  [ApiErrorCode.NOT_EXIST_IDENTIFIER]: "auth.third_party_not_bound",
  [ApiErrorCode.ERROR]: "common.error",
  [ApiErrorCode.INVALID_PARAMS]: "common.invalid_params",
  [ApiErrorCode.ERROR_EXIST_NICK]: "user.nickname_exists",
  [ApiErrorCode.ERROR_EXIST_USER]: "user.username_exists",
  [ApiErrorCode.ERROR_NOT_EXIST_USER]: "user.not_found",
  [ApiErrorCode.ERROR_NOT_COMPARE]: "auth.account_password_invalid",
  [ApiErrorCode.ERROR_NOT_COMPARE_PASSWORD]: "auth.password_confirm_mismatch",
  [ApiErrorCode.ERROR_FAIL_ENCRYPTION]: "common.encrypt_failed",
  [ApiErrorCode.ERROR_NOT_EXIST_PRODUCT]: "product.not_found",
  [ApiErrorCode.ERROR_NOT_EXIST_ADDRESS]: "address.not_found",
  [ApiErrorCode.ERROR_EXIST_FAVORITE]: "favorite.exists",
  [ApiErrorCode.ERROR_USER_NOT_FOUND]: "user.not_found",
  [ApiErrorCode.ERROR_BOSS_CHECK_TOKEN_FAIL]: "seller.token_check_failed",
  [ApiErrorCode.ERROR_BOSS_CHECK_TOKEN_TIMEOUT]: "seller.token_timeout",
  [ApiErrorCode.ERROR_BOSS_TOKEN]: "seller.token_create_failed",
  [ApiErrorCode.ERROR_BOSS]: "seller.token_invalid",
  [ApiErrorCode.ERROR_BOSS_INSUFFICIENT_AUTHORITY]: "seller.insufficient_authority",
  [ApiErrorCode.ERROR_BOSS_PRODUCT]: "seller.product_file_read_failed",
  [ApiErrorCode.ERROR_PRODUCT_EXIST_CART]: "cart.product_exists",
  [ApiErrorCode.ERROR_PRODUCT_MORE_CART]: "cart.quantity_limit_exceeded",
  [ApiErrorCode.ERROR_AUTH_CHECK_TOKEN_FAIL]: "auth.token_check_failed",
  [ApiErrorCode.ERROR_AUTH_CHECK_TOKEN_TIMEOUT]: "auth.token_timeout",
  [ApiErrorCode.ERROR_AUTH_TOKEN]: "auth.token_create_failed",
  [ApiErrorCode.ERROR_AUTH]: "auth.token_invalid",
  [ApiErrorCode.ERROR_AUTH_INSUFFICIENT_AUTHORITY]: "auth.insufficient_authority",
  [ApiErrorCode.ERROR_READ_FILE]: "file.read_failed",
  [ApiErrorCode.ERROR_SEND_EMAIL]: "common.call_api_failed",
  [ApiErrorCode.ERROR_CALL_API]: "common.call_api_failed",
  [ApiErrorCode.ERROR_UNMARSHAL_JSON]: "json.unmarshal_failed",
  [ApiErrorCode.ERROR_ADMIN_FIND_USER]: "admin.user_query_failed",
  [ApiErrorCode.ERROR_DATABASE]: "database.error",
  [ApiErrorCode.ERROR_OSS]: "oss.config_error",
  [ApiErrorCode.ERROR_UPLOAD_FILE]: "file.upload_failed",
};

function firstText(...values: unknown[]): string {
  for (const value of values) {
    if (typeof value === "string" && value.trim()) {
      return value.trim();
    }
  }
  return "";
}

function apiMessage(key: string, fallback = "") {
  return t(`api.${key}`, fallback);
}

export function resolveApiErrorMessage(payload: any, fallback = t("common.requestFailed", "请求失败")) {
  const status = Number(payload?.status ?? payload?.code);
  const detail = firstText(payload?.data, payload?.error, payload?.msg);
  const msgKey = firstText(payload?.msg_key);
  if (msgKey) {
    return apiMessage(msgKey, detail || fallback);
  }
  if (status && ApiErrorMessageKey[status]) {
    return detail || apiMessage(ApiErrorMessageKey[status], fallback);
  }
  return detail || fallback;
}
