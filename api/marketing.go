//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../oapi.cfg.yaml ../openapi.yaml

package api

import (
	"context"
	"errors"

	"github.com/MukeshGKastala/marketing/db"
)

type server struct {
	store db.Store
}

func NewServer(store db.Store) *server {
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
	promotion, err := s.store.CreatePromotionTx(ctx, db.CreatePromotionTxParams{
		CreatePromotionParams: db.CreatePromotionParams{
			PromotionCode: req.Body.PromotionCode,
			StartDate:     req.Body.StartDate,
			EndDate:       req.Body.EndDate,
		},
	})
	if err != nil {
		if errors.Is(err, db.ErrDuplicateEntry) {
			return CreatePromotion400JSONResponse{
				BadRequestJSONResponse{Message: err.Error()},
			}, nil
		}
		return CreatePromotion500JSONResponse{
			InternalServerErrorJSONResponse{Message: err.Error()},
		}, nil
	}

	return CreatePromotion201JSONResponse{
		Id:            int(promotion.ID),
		PromotionCode: promotion.PromotionCode,
		StartDate:     promotion.StartDate,
		EndDate:       promotion.EndDate,
		CreatedAt:     promotion.CreatedAt,
		UpdatedAt:     promotion.UpdatedAt,
	}, nil
}
