package mock

import (
	store "github.com/mo-mohamed/txparser/storage"
)

type BlockchainMock struct {
	ParseBlockFunc         func(block int) ([]store.Transaction, error)
	LatestNetworkBlockFunc func() int
}

func (b *BlockchainMock) ParseBlock(block int) ([]store.Transaction, error) {
	return b.ParseBlockFunc(block)
}

func (b *BlockchainMock) LatestNetworkBlock() int {
	return b.LatestNetworkBlockFunc()
}
