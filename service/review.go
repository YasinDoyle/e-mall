package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"sync"
	"time"

	conf "github.com/YasinDoyle/e-mall/config"
	"github.com/YasinDoyle/e-mall/consts"
	"github.com/YasinDoyle/e-mall/repository/db/dao"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/ctl"
	"github.com/YasinDoyle/e-mall/utils/log"
	util "github.com/YasinDoyle/e-mall/utils/upload"
)

var ReviewSrvIns *ReviewSrv
var ReviewSrvOnce sync.Once

type ReviewSrv struct{}

func GetReviewSrv() *ReviewSrv {
	ReviewSrvOnce.Do(func() { ReviewSrvIns = &ReviewSrv{} })
	return ReviewSrvIns
}

func (s *ReviewSrv) ReviewImageUpload(ctx context.Context, file multipart.File, fileSize int64) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return nil, err
	}

	fileName := fmt.Sprintf("review_%d", time.Now().UnixNano())
	var path string
	if conf.Config.System.UploadModel == consts.UploadModelLocal {
		path, err = util.ReviewUploadToLocalStatic(file, u.Id, fileName)
		if err == nil {
			path = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.ProductPath + path
		}
	} else {
		path, err = util.UploadToQiNiu(file, fileSize)
	}
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}

	resp = &types.ReviewImageUploadResp{URL: path}
	return
}

// ReviewCreate 用户创建评价（仅订单状态为已收货时可评价）
func (s *ReviewSrv) ReviewCreate(ctx context.Context, req *types.ReviewCreateReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		return
	}

	// 校验订单归属且状态为已收货
	order, err := dao.NewOrderDao(ctx).GetOrderById(req.OrderID, u.Id)
	if err != nil {
		return nil, errors.New("订单不存在")
	}
	if order.Type != consts.OrderTypeReceipt {
		return nil, errors.New("仅已收货的订单才能评价")
	}
	if order.ProductID != req.ProductID {
		return nil, errors.New("订单商品与评价商品不匹配")
	}

	// 防重复评价
	if dao.NewReviewDao(ctx).HasReviewed(u.Id, req.OrderID) {
		return nil, errors.New("该订单已评价")
	}

	// 查询用户信息（用于存储昵称和头像快照）
	user, _ := dao.NewUserDao(ctx).GetUserById(u.Id)

	review := &model.Review{
		UserID:     u.Id,
		ProductID:  req.ProductID,
		OrderID:    req.OrderID,
		Rating:     req.Rating,
		Content:    req.Content,
		Images:     req.Images,
		UserName:   user.NickName,
		UserAvatar: user.Avatar,
	}
	if err = dao.NewReviewDao(ctx).CreateReview(review); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "评价成功"
	return
}

// ReviewList 获取商品评价列表
func (s *ReviewSrv) ReviewList(ctx context.Context, req *types.ReviewListReq) (resp interface{}, err error) {
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	reviews, total, err := dao.NewReviewDao(ctx).ListByProduct(req.ProductID, req.PageNum, req.PageSize)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	list := make([]*types.ReviewResp, 0, len(reviews))
	for _, r := range reviews {
		avatar := r.UserAvatar
		if conf.Config.System.UploadModel == consts.UploadModelLocal {
			avatar = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.AvatarPath + avatar
		}
		list = append(list, &types.ReviewResp{
			ID:         r.ID,
			UserName:   r.UserName,
			UserAvatar: avatar,
			Rating:     r.Rating,
			Content:    r.Content,
			Images:     r.Images,
			CreatedAt:  r.CreatedAt.Unix(),
		})
	}
	resp = &types.DataListResp{Item: list, Total: total}
	return
}

// ReviewAdminDelete 管理员删除评价
func (s *ReviewSrv) ReviewAdminDelete(ctx context.Context, req *types.ReviewAdminDeleteReq) (resp interface{}, err error) {
	if err = dao.NewReviewDao(ctx).DeleteReview(req.ID); err != nil {
		log.LogrusObj.Error(err)
		return
	}
	resp = "删除成功"
	return
}
