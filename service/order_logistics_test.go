package service

import "testing"

func TestNormalizeShipmentInfoRequiresCompanyAndTrackingNo(t *testing.T) {
	_, err := NormalizeShipmentInfo(" SF ", " SF123456789 ")
	if err != nil {
		t.Fatal(err)
	}
	if _, err = NormalizeShipmentInfo("", "SF123"); err == nil {
		t.Fatal("expected missing logistics company to fail")
	}
	if _, err = NormalizeShipmentInfo("SF", ""); err == nil {
		t.Fatal("expected missing tracking number to fail")
	}
}
