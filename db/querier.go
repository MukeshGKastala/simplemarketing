// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"context"
)

type Querier interface {
	ListPromotions(ctx context.Context) ([]Promotion, error)
}

var _ Querier = (*Queries)(nil)
