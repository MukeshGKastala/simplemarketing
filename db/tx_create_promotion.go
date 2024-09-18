package db

import (
	"context"
	"fmt"
)

type CreatePromotionTxParams struct {
	CreatePromotionParams
}

type CreatePromotionTxResult struct {
	Promotion
}

func (store *SQLStore) CreatePromotionTx(ctx context.Context, arg CreatePromotionTxParams) (CreatePromotionTxResult, error) {
	var result CreatePromotionTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		taken, err := q.IsPromotionCodeTaken(ctx, IsPromotionCodeTakenParams{
			PromotionCode: arg.PromotionCode,
			StartDate:     arg.StartDate,
			EndDate:       arg.EndDate,
		})
		if err != nil {
			return err
		}
		if taken {
			return fmt.Errorf("promotion code is in use: %w", ErrDuplicateEntry)
		}

		res, err := q.CreatePromotion(ctx, CreatePromotionParams{
			PromotionCode: arg.PromotionCode,
			StartDate:     arg.StartDate,
			EndDate:       arg.EndDate,
		})
		if err != nil {
			return err
		}

		id, err := res.LastInsertId()
		if err != nil {
			return err
		}

		p, err := q.GetPromotion(ctx, int32(id))
		if err != nil {
			return err
		}

		result.Promotion = p
		return nil
	})

	return result, err
}
