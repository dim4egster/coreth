// (c) 2021-2022, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package handlers

import (
	"context"

	"github.com/dim4egster/coreth/core/state/snapshot"
	"github.com/dim4egster/coreth/core/types"
	"github.com/dim4egster/coreth/plugin/evm/message"
	"github.com/dim4egster/coreth/sync/handlers/stats"
	"github.com/dim4egster/coreth/trie"
	"github.com/dim4egster/qmallgo/codec"
	"github.com/dim4egster/qmallgo/ids"
	"github.com/ethereum/go-ethereum/common"
)

var _ message.RequestHandler = &syncHandler{}

type BlockProvider interface {
	GetBlock(common.Hash, uint64) *types.Block
}

type SnapshotProvider interface {
	Snapshots() *snapshot.Tree
}

type SyncDataProvider interface {
	BlockProvider
	SnapshotProvider
}

type syncHandler struct {
	stateTrieLeafsRequestHandler  *LeafsRequestHandler
	atomicTrieLeafsRequestHandler *LeafsRequestHandler
	blockRequestHandler           *BlockRequestHandler
	codeRequestHandler            *CodeRequestHandler
}

// NewSyncHandler constructs the handler for serving state sync.
func NewSyncHandler(
	provider SyncDataProvider,
	evmTrieDB *trie.Database,
	atomicTrieDB *trie.Database,
	networkCodec codec.Manager,
	stats stats.HandlerStats,
) message.RequestHandler {
	return &syncHandler{
		stateTrieLeafsRequestHandler:  NewLeafsRequestHandler(evmTrieDB, provider, networkCodec, stats),
		atomicTrieLeafsRequestHandler: NewLeafsRequestHandler(atomicTrieDB, nil, networkCodec, stats),
		blockRequestHandler:           NewBlockRequestHandler(provider, networkCodec, stats),
		codeRequestHandler:            NewCodeRequestHandler(evmTrieDB.DiskDB(), networkCodec, stats),
	}
}

func (s *syncHandler) HandleStateTrieLeafsRequest(ctx context.Context, nodeID ids.NodeID, requestID uint32, leafsRequest message.LeafsRequest) ([]byte, error) {
	return s.stateTrieLeafsRequestHandler.OnLeafsRequest(ctx, nodeID, requestID, leafsRequest)
}

func (s *syncHandler) HandleAtomicTrieLeafsRequest(ctx context.Context, nodeID ids.NodeID, requestID uint32, leafsRequest message.LeafsRequest) ([]byte, error) {
	return s.atomicTrieLeafsRequestHandler.OnLeafsRequest(ctx, nodeID, requestID, leafsRequest)
}

func (s *syncHandler) HandleBlockRequest(ctx context.Context, nodeID ids.NodeID, requestID uint32, blockRequest message.BlockRequest) ([]byte, error) {
	return s.blockRequestHandler.OnBlockRequest(ctx, nodeID, requestID, blockRequest)
}

func (s *syncHandler) HandleCodeRequest(ctx context.Context, nodeID ids.NodeID, requestID uint32, codeRequest message.CodeRequest) ([]byte, error) {
	return s.codeRequestHandler.OnCodeRequest(ctx, nodeID, requestID, codeRequest)
}
