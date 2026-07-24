package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/mail"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/YasinDoyle/e-mall/repository/cache"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/email"
)

const registerEmailCodeExpire = 5 * time.Minute

func normalizeRegisterEmail(raw string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(raw))
	if normalized == "" {
		return "", errors.New("请输入邮箱")
	}
	addr, err := mail.ParseAddress(normalized)
	if err != nil || addr.Address != normalized {
		return "", errors.New("邮箱格式不正确")
	}
	return normalized, nil
}

func buildRegisterEmailCodeKey(email string) string {
	return fmt.Sprintf("register_email_code:%s", email)
}

func validateRegisterEmailCode(submitted string, stored string) error {
	code := strings.TrimSpace(submitted)
	if code == "" {
		return errors.New("请输入邮箱验证码")
	}
	if code != stored {
		return errors.New("邮箱验证码错误")
	}
	return nil
}

func ensureRegisterEmailCodeCanBeSent(exists bool) error {
	if exists {
		return errors.New("验证码已发送，请5分钟后再试")
	}
	return nil
}

func generateRegisterEmailCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

func buildRegisterEmailCodeHTML(code string) string {
	return fmt.Sprintf(`<!doctype html>
<html>
  <body style="margin:0;padding:0;background:#f4f7fb;font-family:Arial,'Microsoft YaHei',sans-serif;color:#1f2937;">
    <div style="max-width:520px;margin:0 auto;padding:32px 16px;">
      <div style="background:#ffffff;border:1px solid #e5e7eb;border-radius:8px;overflow:hidden;">
        <div style="background:#2563eb;padding:22px 28px;color:#ffffff;">
          <div style="font-size:20px;font-weight:700;line-height:1.4;">E-Mall 注册验证码</div>
          <div style="font-size:13px;opacity:.9;margin-top:4px;">完成邮箱验证后即可创建账号</div>
        </div>
        <div style="padding:28px;">
          <p style="margin:0 0 14px;font-size:15px;line-height:1.7;">您好，您的注册验证码是：</p>
          <div style="margin:18px 0;padding:18px 20px;background:#eef6ff;border-radius:8px;text-align:center;">
            <span style="font-size:34px;line-height:1;font-weight:700;letter-spacing:8px;color:#1d4ed8;">%s</span>
          </div>
          <p style="margin:0 0 12px;font-size:14px;line-height:1.7;color:#4b5563;">验证码5分钟内有效，请勿泄露给他人。</p>
          <p style="margin:0;font-size:13px;line-height:1.7;color:#9ca3af;">如果不是您本人操作，可以忽略这封邮件。</p>
        </div>
      </div>
    </div>
  </body>
</html>`, code)
}

func (s *UserSrv) SendRegisterEmailCode(ctx context.Context, req *types.RegisterEmailCodeReq) (resp any, err error) {
	normalizedEmail, err := normalizeRegisterEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if _, exist, daoErr := dao.NewUserDao(ctx).ExistOrNotByEmail(normalizedEmail); daoErr != nil {
		return nil, daoErr
	} else if exist {
		return nil, errors.New("邮箱已被注册")
	}

	code, err := generateRegisterEmailCode()
	if err != nil {
		return nil, err
	}
	if cache.RedisClient == nil {
		return nil, errors.New("验证码服务未初始化")
	}
	codeKey := buildRegisterEmailCodeKey(normalizedEmail)
	exists, err := cache.RedisClient.Exists(ctx, codeKey).Result()
	if err != nil {
		return nil, err
	}
	if err = ensureRegisterEmailCodeCanBeSent(exists > 0); err != nil {
		return nil, err
	}

	if err = cache.RedisClient.Set(ctx, codeKey, code, registerEmailCodeExpire).Err(); err != nil {
		return nil, err
	}

	mailText := buildRegisterEmailCodeHTML(code)
	if err = email.NewEmailSender().Send(mailText, normalizedEmail, "E-Mall 注册验证码"); err != nil {
		return nil, err
	}
	return "验证码已发送", nil
}

func consumeRegisterEmailCode(ctx context.Context, email string, submittedCode string) error {
	if cache.RedisClient == nil {
		return errors.New("验证码服务未初始化")
	}
	key := buildRegisterEmailCodeKey(email)
	storedCode, err := cache.RedisClient.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return errors.New("邮箱验证码已过期")
	}
	if err != nil {
		return err
	}
	if err = validateRegisterEmailCode(submittedCode, storedCode); err != nil {
		return err
	}
	return cache.RedisClient.Del(ctx, key).Err()
}
