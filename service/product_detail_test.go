package service

import (
	"testing"

	"github.com/YasinDoyle/e-mall/repository/db/model"
)

func TestBuildProductRespIncludesAuditEnrichmentFields(t *testing.T) {
	product := &model.Product{
		Name:              "有机苹果",
		CategoryID:        1,
		Title:             "产地直发",
		Info:              "脆甜多汁",
		ImgPath:           "boss2/apple.jpg",
		Price:             "19.90",
		DiscountPrice:     "15.90",
		Num:               20,
		BossID:            2,
		BossName:          "果园店",
		Brand:             "山谷果园",
		Origin:            "山东烟台",
		Specification:     "5斤装",
		ProductionDate:    "2026-07-01",
		ShelfLife:         "30天",
		ServiceGuarantees: "坏果包赔",
		CertificateMeta:   "绿色食品认证",
	}

	resp := buildProductResp(product)

	if resp.Brand != product.Brand {
		t.Fatalf("expected brand %q, got %q", product.Brand, resp.Brand)
	}
	if resp.Origin != product.Origin {
		t.Fatalf("expected origin %q, got %q", product.Origin, resp.Origin)
	}
	if resp.Specification != product.Specification {
		t.Fatalf("expected specification %q, got %q", product.Specification, resp.Specification)
	}
	if resp.ProductionDate != product.ProductionDate {
		t.Fatalf("expected production date %q, got %q", product.ProductionDate, resp.ProductionDate)
	}
	if resp.ShelfLife != product.ShelfLife {
		t.Fatalf("expected shelf life %q, got %q", product.ShelfLife, resp.ShelfLife)
	}
	if resp.ServiceGuarantees != product.ServiceGuarantees {
		t.Fatalf("expected service guarantees %q, got %q", product.ServiceGuarantees, resp.ServiceGuarantees)
	}
	if resp.CertificateMeta != product.CertificateMeta {
		t.Fatalf("expected certificate meta %q, got %q", product.CertificateMeta, resp.CertificateMeta)
	}
}

func TestBuildProductCertificateRespKeepsCertificateMetadata(t *testing.T) {
	certificate := &model.ProductCertificate{
		ProductID:       10,
		CertificateType: "quality_inspection",
		Name:            "质检报告",
		FilePath:        "boss2/cert-quality.jpg",
	}

	resp := buildProductCertificateResp(certificate)

	if resp.ProductID != certificate.ProductID {
		t.Fatalf("expected product id %d, got %d", certificate.ProductID, resp.ProductID)
	}
	if resp.CertificateType != certificate.CertificateType {
		t.Fatalf("expected certificate type %q, got %q", certificate.CertificateType, resp.CertificateType)
	}
	if resp.Name != certificate.Name {
		t.Fatalf("expected certificate name %q, got %q", certificate.Name, resp.Name)
	}
	if resp.FilePath == "" {
		t.Fatal("expected certificate file path to be exposed")
	}
}
