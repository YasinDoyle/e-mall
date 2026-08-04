package model

import (
	"testing"

	conf "github.com/YasinDoyle/e-mall/config"
)

func TestNewRegisteredUserHasNoPayKey(t *testing.T) {
	user := &User{}
	if user.HasPayKey() {
		t.Fatal("new registered user should not have a pay key")
	}
}

func TestSetInitialMoneyWithPayKey(t *testing.T) {
	conf.Config = &conf.Conf{
		EncryptSecret: &conf.EncryptSecret{
			MoneySecret: "1234567890123456",
		},
	}
	user := &User{}
	if err := user.SetInitialMoneyWithPayKey("123456"); err != nil {
		t.Fatalf("expected pay key setup to pass, got %v", err)
	}
	if !user.HasPayKey() {
		t.Fatal("expected user to have pay key after setup")
	}
	money, err := user.DecryptMoney()
	if err != nil {
		t.Fatalf("expected money to decrypt with platform wallet key, got %v", err)
	}
	if money != 10000 {
		t.Fatalf("expected initial money 10000, got %v", money)
	}
}

func TestPayKeyIsAuthenticationNotWalletEncryptionKey(t *testing.T) {
	conf.Config = &conf.Conf{
		EncryptSecret: &conf.EncryptSecret{
			MoneySecret: "1234567890123456",
		},
	}
	user := &User{}
	if err := user.SetInitialMoneyWithPayKey("123456"); err != nil {
		t.Fatal(err)
	}
	if !user.CheckPayKey("123456") {
		t.Fatal("expected payment password to authenticate")
	}
	if user.CheckPayKey("654321") {
		t.Fatal("expected wrong payment password to fail")
	}
	if user.PayKeyDigest == "" {
		t.Fatal("expected payment password digest to be stored separately")
	}

	money, err := user.DecryptMoney()
	if err != nil {
		t.Fatalf("expected platform wallet decrypt without buyer key, got %v", err)
	}
	if money != 10000 {
		t.Fatalf("expected initial money 10000, got %v", money)
	}
}
