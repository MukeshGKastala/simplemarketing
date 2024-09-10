//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../oapi.cfg.yaml ../openapi.yaml

package api

import (
	"context"

	sqlc "github.com/MukeshGKastala/marketing/db"
)

type Store interface {
	ListPromotions(ctx context.Context) ([]sqlc.Promotion, error)
}

type server struct {
	store Store
}

func NewServer(store Store) *server {
	return &server{store: store}
}

var _ StrictServerInterface = (*server)(nil)

func (s *server) GetPromotions(ctx context.Context, _ GetPromotionsRequestObject) (GetPromotionsResponseObject, error) {
	promotions, err := s.store.ListPromotions(ctx)
	if err != nil {
		errMsg := err.Error()
		return GetPromotions500JSONResponse{
			InternalServerErrorJSONResponse{Message: &errMsg},
		}, nil
	}

	var ret GetPromotions200JSONResponse
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
