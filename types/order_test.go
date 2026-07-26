package types

import (
	"encoding/json"
	"testing"
)

func TestOrderCreateReqAcceptsDecimalMoney(t *testing.T) {
	var req OrderCreateReq
	err := json.Unmarshal([]byte(`{"money":0.01}`), &req)
	if err != nil {
		t.Fatalf("expected decimal money to unmarshal, got %v", err)
	}

	if req.Money != 0.01 {
		t.Fatalf("expected money 0.01, got %.2f", req.Money)
	}
}
