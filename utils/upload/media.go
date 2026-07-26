package upload

import (
	"mime/multipart"
	"strings"

	conf "github.com/YasinDoyle/e-mall/config"
	"github.com/YasinDoyle/e-mall/consts"
)

func UploadProductImage(file multipart.File, fileSize int64, bossID uint, productName string) (string, error) {
	if conf.Config.System.UploadModel == consts.UploadModelLocal {
		return ProductUploadToLocalStatic(file, bossID, productName)
	}
	return UploadToQiNiu(file, fileSize)
}

func UploadAvatarImage(file multipart.File, fileSize int64, userID uint, userName string) (string, error) {
	if conf.Config.System.UploadModel == consts.UploadModelLocal {
		return AvatarUploadToLocalStatic(file, userID, userName)
	}
	return UploadToQiNiu(file, fileSize)
}

func UploadCarouselImage(file multipart.File, fileSize int64, fileName string) (string, error) {
	if conf.Config.System.UploadModel == consts.UploadModelLocal {
		return CarouselUploadToLocalStatic(file, fileName)
	}
	return UploadToQiNiu(file, fileSize)
}

func ProductImageURL(path string) string {
	return localStaticURL(path, conf.Config.PhotoPath.ProductPath)
}

func AvatarURL(path string) string {
	return localStaticURL(path, conf.Config.PhotoPath.AvatarPath)
}

func localStaticURL(path string, staticPath string) string {
	if path == "" || isRemoteURL(path) || conf.Config.System.UploadModel != consts.UploadModelLocal {
		return path
	}
	return conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + staticPath + strings.TrimPrefix(path, "/")
}

func isRemoteURL(path string) bool {
	return strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://")
}
