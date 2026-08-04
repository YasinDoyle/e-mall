package service

import (
	"strings"

	"github.com/YasinDoyle/e-mall/types"
	"github.com/YasinDoyle/e-mall/utils/e"
)

const shipmentInfoMaxLen = 64

func NormalizeShipmentInfo(company, trackingNo string) (types.ShipmentInfo, error) {
	info := types.ShipmentInfo{
		LogisticsCompany: strings.TrimSpace(company),
		TrackingNo:       strings.TrimSpace(trackingNo),
	}
	if info.LogisticsCompany == "" || info.TrackingNo == "" {
		return types.ShipmentInfo{}, e.NewBusinessError(e.InvalidParams)
	}
	if len(info.LogisticsCompany) > shipmentInfoMaxLen || len(info.TrackingNo) > shipmentInfoMaxLen {
		return types.ShipmentInfo{}, e.NewBusinessError(e.InvalidParams)
	}
	return info, nil
}
