package service

import (
	"context"
	"errors"
	"mime/multipart"
	"sync"

	conf "github.com/YasinDoyle/e-mall/config"
	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/jwt"
	"github.com/YasinDoyle/e-mall/utils/log"
	util "github.com/YasinDoyle/e-mall/utils/upload"
)

var UserSrvIns *UserSrv
var UserSrvOnce sync.Once

type UserSrv struct {
}

func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

func (s *UserSrv) UserRegister(ctx context.Context, req *types.UserRegisterReq) (resp any, err error) {
	registerEmail, err := normalizeRegisterEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if req.Password == "" || len(req.Password) < 6 {
		return nil, errors.New("密码至少6位")
	}
	if req.Password != req.PasswordConfirm {
		return nil, errors.New("两次密码输入不一致")
	}

	userDao := dao.NewUserDao(ctx)
	_, exist, err := userDao.ExistOrNotByUserName(req.UserName)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	if exist {
		err = errors.New("用户已经存在了")
		return
	}
	if _, exist, err = userDao.ExistOrNotByEmail(registerEmail); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	if exist {
		err = errors.New("邮箱已被注册")
		return
	}
	if err = consumeRegisterEmailCode(ctx, registerEmail, req.EmailCode); err != nil {
		return nil, err
	}
	user := &model.User{
		NickName: req.NickName,
		UserName: req.UserName,
		Email:    registerEmail,
		Status:   model.Active,
	}

	//加密密码
	if err = user.SetPassword(req.Password); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	user.Avatar = consts.UserDefaultAvatarLocal
	if conf.Config.System.UploadModel == consts.UploadModelOss {
		//如果是oss上传模式，默认头像设置为oss上的地址
		user.Avatar = consts.UserDefaultAvatarOss
	}

	//创建用户
	err = userDao.CreateUser(user)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	return
}

// UserLogin 用户登陆函数
func (s *UserSrv) UserLogin(ctx context.Context, req *types.UserServiceReq) (resp interface{}, err error) {
	var user *model.User
	userDao := dao.NewUserDao(ctx)
	user, exist, err := userDao.ExistOrNotByUserName(req.UserName)
	if !exist { // 如果查询不到，返回相应的错误
		log.LogrusObj.Error(err)
		return nil, errors.New("用户不存在")
	}

	if !user.CheckPassword(req.Password) {
		return nil, errors.New("账号/密码不正确")
	}
	if user.Status == "banned" {
		return nil, errors.New("账号已被封禁")
	}

	accessToken, refreshToken, err := jwt.GenerateToken(user.ID, req.UserName)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	userResp := &types.UserInfoResp{
		ID:        user.ID,
		UserName:  user.UserName,
		NickName:  user.NickName,
		Email:     user.Email,
		Status:    user.Status,
		Avatar:    user.AvatarURL(),
		PayKeySet: user.HasPayKey(),
		CreateAt:  user.CreatedAt.Unix(),
	}

	resp = &types.UserTokenData{
		User:         userResp,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return
}

// UserInfoUpdate 用户修改信息
func (s *UserSrv) UserInfoUpdate(ctx context.Context, req *types.UserInfoUpdateReq) (resp interface{}, err error) {
	// 找到用户
	u, _ := ctl.GetUserInfo(ctx)
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(u.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	if req.NickName != "" {
		user.NickName = req.NickName
	}

	err = userDao.UpdateUserById(u.Id, user)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	return
}

// UserAvatarUpload 更新头像
func (s *UserSrv) UserAvatarUpload(ctx context.Context, file multipart.File, fileSize int64, req *types.UserServiceReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	uId := u.Id
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(uId)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	var path string
	if conf.Config.System.UploadModel == consts.UploadModelLocal { // 兼容两种存储方式
		path, err = util.AvatarUploadToLocalStatic(file, uId, user.UserName)
	} else {
		path, err = util.UploadToQiNiu(file, fileSize)
	}
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	user.Avatar = path
	err = userDao.UpdateUserById(uId, user)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	return
}

// UserInfoShow 用户信息展示
func (s *UserSrv) UserInfoShow(ctx context.Context, req *types.UserInfoShowReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	user, err := dao.NewUserDao(ctx).GetUserById(u.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = &types.UserInfoResp{
		ID:        user.ID,
		UserName:  user.UserName,
		NickName:  user.NickName,
		Email:     user.Email,
		Status:    user.Status,
		Avatar:    user.AvatarURL(),
		PayKeySet: user.HasPayKey(),
		CreateAt:  user.CreatedAt.Unix(),
	}

	return
}

func (s *UserSrv) UserFollow(ctx context.Context, req *types.UserFollowingReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	err = dao.NewUserDao(ctx).FollowUser(u.Id, req.Id)

	return
}

func (s *UserSrv) UserUnFollow(ctx context.Context, req *types.UserUnFollowingReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	err = dao.NewUserDao(ctx).UnFollowUser(u.Id, req.Id)

	return
}
