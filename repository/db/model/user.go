package model

import (
	"github.com/YasinDoyle/e-mall/utils/secret"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	conf "github.com/YasinDoyle/e-mall/config"
	"github.com/YasinDoyle/e-mall/consts"
)

type User struct {
	gorm.Model
	UserName       string `gorm:"unique"`
	Email          string
	PasswordDigest string
	NickName       string
	Status         string
	Avatar         string `gorm:"size:1000"`
	Money          string
	PayKeyDigest   string
	PayKeySet      bool   `gorm:"default:false"`
	IsAdmin        bool   `gorm:"default:false"`
	Relations      []User `gorm:"many2many:relation;"`
}

const (
	PasswordCost        = 12       //密码加密难度
	Active       string = "active" //激活难度
)

const walletServerKeyFallback = "mall-wallet-key"

func (u *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return err
	}
	u.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordDigest), []byte(password))
	return err == nil
}

// AvatarURL 头像地址
func (u *User) AvatarURL() string {
	if conf.Config.System.UploadModel == consts.UploadModelOss {
		return u.Avatar
	}
	pConfig := conf.Config.PhotoPath
	return pConfig.PhotoHost + conf.Config.System.HttpPort + pConfig.AvatarPath + u.Avatar
}

func walletServerKey() string {
	if conf.Config != nil && conf.Config.EncryptSecret != nil {
		if conf.Config.EncryptSecret.MoneySecret != "" {
			return conf.Config.EncryptSecret.MoneySecret
		}
		if conf.Config.EncryptSecret.JwtSecret != "" {
			return conf.Config.EncryptSecret.JwtSecret
		}
	}
	return walletServerKeyFallback
}

// EncryptMoney encrypts the wallet balance with the platform wallet key.
func (u *User) EncryptMoney() (money string, err error) {
	aesObj, err := secret.NewAesEncrypt(walletServerKey(), "wallet", "", secret.AesEncrypt128, secret.AesModeTypeCBC)
	if err != nil {
		return
	}
	money = aesObj.SecretEncrypt(u.Money)

	return
}

// DecryptMoney decrypts the wallet balance with the platform wallet key.
func (u *User) DecryptMoney() (money float64, err error) {
	aesObj, err := secret.NewAesEncrypt(walletServerKey(), "wallet", "", secret.AesEncrypt128, secret.AesModeTypeCBC)
	if err != nil {
		return
	}

	money = cast.ToFloat64(aesObj.SecretDecrypt(u.Money))
	return
}

func (u *User) HasPayKey() bool {
	return u.PayKeySet && u.PayKeyDigest != ""
}

func (u *User) CheckPayKey(key string) bool {
	if key == "" || u.PayKeyDigest == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(u.PayKeyDigest), []byte(key)) == nil
}

func (u *User) SetInitialMoneyWithPayKey(key string) error {
	digest, err := bcrypt.GenerateFromPassword([]byte(key), PasswordCost)
	if err != nil {
		return err
	}
	u.Money = consts.UserInitMoney
	money, err := u.EncryptMoney()
	if err != nil {
		return err
	}
	u.Money = money
	u.PayKeyDigest = string(digest)
	u.PayKeySet = true
	return nil
}
