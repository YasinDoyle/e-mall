package service

import (
	"strings"
	"testing"
)

func TestNormalizeRegisterEmail(t *testing.T) {
	email, err := normalizeRegisterEmail("  USER@Example.COM  ")
	if err != nil {
		t.Fatalf("expected valid email, got %v", err)
	}
	if email != "user@example.com" {
		t.Fatalf("expected normalized email, got %q", email)
	}
}

func TestBuildRegisterEmailCodeKey(t *testing.T) {
	key := buildRegisterEmailCodeKey("user@example.com")
	if key != "register_email_code:user@example.com" {
		t.Fatalf("unexpected redis key %q", key)
	}
}

func TestValidateRegisterEmailCode(t *testing.T) {
	if err := validateRegisterEmailCode("123456", "123456"); err != nil {
		t.Fatalf("expected matching code to pass, got %v", err)
	}
	if err := validateRegisterEmailCode("123456", "654321"); err == nil {
		t.Fatal("expected mismatched code to fail")
	}
	if err := validateRegisterEmailCode("", "123456"); err == nil {
		t.Fatal("expected empty submitted code to fail")
	}
}

func TestEnsureRegisterEmailCodeCanBeSent(t *testing.T) {
	if err := ensureRegisterEmailCodeCanBeSent(false); err != nil {
		t.Fatalf("expected no existing code to pass, got %v", err)
	}
	if err := ensureRegisterEmailCodeCanBeSent(true); err == nil {
		t.Fatal("expected existing code to block repeat sends")
	}
}

func TestBuildRegisterEmailCodeHTML(t *testing.T) {
	html := buildRegisterEmailCodeHTML("123456")
	for _, want := range []string{
		"<!doctype html>",
		"E-Mall 注册验证码",
		"123456",
		"5分钟内有效",
		"font-family",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("expected email html to contain %q, got %s", want, html)
		}
	}
}
