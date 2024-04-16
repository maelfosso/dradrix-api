package storage

import "context"

type Querier interface {
	GetUserByUsername(ctx context.Context, arg GetUserByUsernameParams) error
}

type QuerierTx interface {
}

var _ Querier = (*Queries)(nil)
