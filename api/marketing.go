//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../oapi.cfg.yaml ../openapi.yaml

package api

import (
	"context"

	"github.com/MukeshGKastala/marketing/db"
)

type server struct {
	store db.Querier
}

func NewServer(store db.Querier) *server {
	return &server{store: store}
}

var _ StrictServerInterface = (*server)(nil)

func (s *server) ListPromotions(ctx context.Context, _ ListPromotionsRequestObject) (ListPromotionsResponseObject, error) {
	promotions, err := s.store.ListPromotions(ctx)
	if err != nil {
		return ListPromotions500JSONResponse{
			InternalServerErrorJSONResponse{Message: err.Error()},
		}, nil
	}

	var ret ListPromotions200JSONResponse
	for _, p := range promotions {
		ret = append(ret, Promotion{
			Id:            int(p.ID),
			PromotionCode: p.PromotionCode,
			StartDate:     p.StartDate,
			EndDate:       p.EndDate,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
		})
	}

	return ret, nil
}

func (s *server) CreatePromotion(ctx context.Context, req CreatePromotionRequestObject) (CreatePromotionResponseObject, error) {
	taken, err := s.store.IsPromotionCodeTaken(ctx, db.IsPromotionCodeTakenParams{
		PromotionCode: req.Body.PromotionCode,
		StartDate:     req.Body.StartDate,
		EndDate:       req.Body.EndDate,
	})
	if err != nil {
		return CreatePromotion500JSONResponse{
			InternalServerErrorJSONResponse{Message: err.Error()},
		}, nil
	}
	if taken {
		return CreatePromotion400JSONResponse{
			BadRequestJSONResponse{Message: "promotion_code is taken"},
		}, nil
	}

	res, err := s.store.CreatePromotion(ctx, db.CreatePromotionParams{
		PromotionCode: req.Body.PromotionCode,
		StartDate:     req.Body.StartDate,
		EndDate:       req.Body.EndDate,
	})
	if err != nil {
		return CreatePromotion500JSONResponse{
			InternalServerErrorJSONResponse{Message: err.Error()},
		}, nil
	}

	id, err := res.LastInsertId()
	if err != nil {
		return CreatePromotion500JSONResponse{
			InternalServerErrorJSONResponse{Message: err.Error()},
		}, nil
	}

	p, err := s.store.GetPromotion(ctx, int32(id))
	if err != nil {
		return CreatePromotion500JSONResponse{
			InternalServerErrorJSONResponse{Message: err.Error()},
		}, nil
	}

	return CreatePromotion201JSONResponse{
		Id:            int(id),
		PromotionCode: p.PromotionCode,
		StartDate:     p.StartDate,
		EndDate:       p.EndDate,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}, nil
}
