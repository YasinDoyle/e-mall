package es

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"

	conf "github.com/YasinDoyle/e-mall/config"
	"github.com/YasinDoyle/e-mall/repository/db/model"
	"github.com/YasinDoyle/e-mall/types"
)

type ProductIndexRepo struct{}

type ProductDocument struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	CategoryID    uint   `json:"category_id"`
	Title         string `json:"title"`
	Info          string `json:"info"`
	ImgPath       string `json:"img_path"`
	Price         string `json:"price"`
	DiscountPrice string `json:"discount_price"`
	CreatedAt     int64  `json:"created_at"`
	Num           int    `json:"num"`
	OnSale        bool   `json:"on_sale"`
	BossID        uint   `json:"boss_id"`
	BossName      string `json:"boss_name"`
	BossAvatar    string `json:"boss_avatar"`
}

type productSearchResponse struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source ProductDocument `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func NewProductIndexRepo() *ProductIndexRepo {
	return &ProductIndexRepo{}
}

func (repo *ProductIndexRepo) IndexProduct(ctx context.Context, product *model.Product) error {
	if product == nil {
		return errors.New("product is nil")
	}

	body, err := json.Marshal(buildProductDocument(product))
	if err != nil {
		return err
	}

	req := esapi.IndexRequest{
		Index:      productIndexName(),
		DocumentID: strconv.FormatUint(uint64(product.ID), 10),
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}

	return doRequest(ctx, req.Do)
}

func (repo *ProductIndexRepo) DeleteProduct(ctx context.Context, productID uint) error {
	req := esapi.DeleteRequest{
		Index:      productIndexName(),
		DocumentID: strconv.FormatUint(uint64(productID), 10),
		Refresh:    "true",
	}

	res, err := req.Do(normalizeContext(ctx), EsClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil
	}

	if res.IsError() {
		return readResponseError(res)
	}

	return nil
}

func (repo *ProductIndexRepo) SearchProducts(ctx context.Context, keyword string, page types.BasePage) (products []*types.ProductResp, total int64, err error) {
	page = normalizePage(page)
	body, err := json.Marshal(buildSearchBody(keyword))
	if err != nil {
		return nil, 0, err
	}

	from := (page.PageNum - 1) * page.PageSize
	size := page.PageSize
	trackTotalHits := true
	req := esapi.SearchRequest{
		Index:          []string{productIndexName()},
		Body:           bytes.NewReader(body),
		From:           &from,
		Size:           &size,
		TrackTotalHits: trackTotalHits,
	}

	res, err := req.Do(normalizeContext(ctx), EsClient)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, 0, readResponseError(res)
	}

	var result productSearchResponse
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, 0, err
	}

	products = make([]*types.ProductResp, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		products = append(products, &types.ProductResp{
			ID:            hit.Source.ID,
			Name:          hit.Source.Name,
			CategoryID:    hit.Source.CategoryID,
			Title:         hit.Source.Title,
			Info:          hit.Source.Info,
			ImgPath:       hit.Source.ImgPath,
			Price:         hit.Source.Price,
			DiscountPrice: hit.Source.DiscountPrice,
			CreatedAt:     hit.Source.CreatedAt,
			Num:           hit.Source.Num,
			OnSale:        hit.Source.OnSale,
			BossID:        hit.Source.BossID,
			BossName:      hit.Source.BossName,
			BossAvatar:    hit.Source.BossAvatar,
		})
	}

	return products, result.Hits.Total.Value, nil
}

func buildProductDocument(product *model.Product) ProductDocument {
	return ProductDocument{
		ID:            product.ID,
		Name:          product.Name,
		CategoryID:    product.CategoryID,
		Title:         product.Title,
		Info:          product.Info,
		ImgPath:       product.ImgPath,
		Price:         product.Price,
		DiscountPrice: product.DiscountPrice,
		CreatedAt:     product.CreatedAt.Unix(),
		Num:           product.Num,
		OnSale:        product.OnSale,
		BossID:        product.BossID,
		BossName:      product.BossName,
		BossAvatar:    product.BossAvatar,
	}
}

func buildSearchBody(keyword string) map[string]interface{} {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return map[string]interface{}{
			"query": map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
			"sort": []map[string]interface{}{
				{"id": map[string]interface{}{"order": "desc"}},
			},
		}
	}

	return map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"name": map[string]interface{}{
								"query": keyword,
								"boost": 3,
							},
						},
					},
					{
						"match": map[string]interface{}{
							"title": map[string]interface{}{
								"query": keyword,
								"boost": 2,
							},
						},
					},
					{
						"match": map[string]interface{}{
							"info": keyword,
						},
					},
				},
				"minimum_should_match": 1,
			},
		},
		"sort": []map[string]interface{}{
			{"_score": map[string]interface{}{"order": "desc"}},
			{"id": map[string]interface{}{"order": "desc"}},
		},
	}
}

func normalizePage(page types.BasePage) types.BasePage {
	if page.PageNum <= 0 {
		page.PageNum = 1
	}
	if page.PageSize <= 0 {
		page.PageSize = 15
	}
	return page
}

func productIndexName() string {
	if conf.Config != nil && conf.Config.Es != nil && conf.Config.Es.EsIndex != "" {
		return conf.Config.Es.EsIndex + "_product"
	}
	return "product"
}

func normalizeContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	if EsClient == nil {
		return ctx
	}
	return ctx
}

func doRequest(ctx context.Context, fn func(context.Context, esapi.Transport) (*esapi.Response, error)) error {
	if EsClient == nil {
		return errors.New("es client is not initialized")
	}

	res, err := fn(normalizeContext(ctx), EsClient)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return readResponseError(res)
	}

	return nil
}

func readResponseError(res *esapi.Response) error {
	body, _ := io.ReadAll(res.Body)
	if len(body) == 0 {
		return fmt.Errorf("es request failed: %s", res.Status())
	}
	return fmt.Errorf("es request failed: %s: %s", res.Status(), strings.TrimSpace(string(body)))
}
