package upload

import (
	"testing"

	conf "github.com/YasinDoyle/e-mall/config"
	"github.com/YasinDoyle/e-mall/consts"
)

func TestProductImageURLBuildsLocalStaticURL(t *testing.T) {
	oldConfig := conf.Config
	defer func() { conf.Config = oldConfig }()
	conf.Config = &conf.Conf{
		System: &conf.System{UploadModel: consts.UploadModelLocal, HttpPort: ":5001"},
		PhotoPath: &conf.LocalPhotoPath{
			PhotoHost:   "http://127.0.0.1",
			ProductPath: "/static/imgs/product/",
			AvatarPath:  "/static/imgs/avatar/",
		},
	}

	got := ProductImageURL("boss2/demo.jpg")
	want := "http://127.0.0.1:5001/static/imgs/product/boss2/demo.jpg"
	if got != want {
		t.Fatalf("expected %q, got %q", want, got)
	}
}

func TestProductImageURLKeepsRemoteURL(t *testing.T) {
	oldConfig := conf.Config
	defer func() { conf.Config = oldConfig }()
	conf.Config = &conf.Conf{
		System:    &conf.System{UploadModel: consts.UploadModelLocal, HttpPort: ":5001"},
		PhotoPath: &conf.LocalPhotoPath{PhotoHost: "http://127.0.0.1", ProductPath: "/static/imgs/product/"},
	}

	got := ProductImageURL("https://cdn.example.com/demo.jpg")
	if got != "https://cdn.example.com/demo.jpg" {
		t.Fatalf("expected remote URL unchanged, got %q", got)
	}
}
