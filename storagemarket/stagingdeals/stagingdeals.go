package stagingdeals

import (
	"context"

	"github.com/filecoin-project/boost/db"
	"github.com/filecoin-project/boost/storagemarket/types"
	"github.com/filecoin-project/boost/storagemarket/types/dealcheckpoints"
	abi "github.com/filecoin-project/go-state-types/abi"
)

type CheckpointState struct {
	Deals          int64               `json:"Deals"`
	CumulativeSize abi.PaddedPieceSize `json:"CumulativeSize"`
}

type Status struct {
	Accepted    map[string]CheckpointState `json:"Accepted"` // string will represent a TransferType value
	Transferred CheckpointState            `json:"Transferred"`
	Published   CheckpointState            `json:"Published"`
}

// GetStatus compute the staging deals status
// Accepted: get number of deals and cumulative size for each transfer types (online only)
// Transferred and Published : get only total number of deals and cumulative size (for both online and offline)
func GetStatus(ctx context.Context, dealsDB *db.DealsDB) (*Status, error) {
	acceptedCheckpointState := make(map[string]CheckpointState)

	for _, transferType := range [3]string{"http", "libp2p", "graphsync"} {
		dealList, err := dealsDB.ByCheckpointAndTransferTypeOnlineDeals(ctx, dealcheckpoints.Accepted.String(), transferType)

		if err != nil {
			return nil, err
		}

		acceptedCheckpointState[transferType] = CheckpointState{
			Deals:          int64(len(dealList)),
			CumulativeSize: GetCumulativeSizeFromDealList(dealList),
		}
	}

	// For "Transferred" checkpoint
	dealListTransferred, err := dealsDB.ByCheckpoint(ctx, dealcheckpoints.Transferred.String())

	if err != nil {
		return nil, err
	}

	// For "Published" checkpoint
	dealListPublished, err := dealsDB.ByCheckpoint(ctx, dealcheckpoints.Published.String())

	if err != nil {
		return nil, err
	}

	return &Status{
		Accepted: acceptedCheckpointState,
		Transferred: CheckpointState{
			Deals:          int64(len(dealListTransferred)),
			CumulativeSize: GetCumulativeSizeFromDealList(dealListTransferred),
		},
		Published: CheckpointState{
			Deals:          int64(len(dealListPublished)),
			CumulativeSize: GetCumulativeSizeFromDealList(dealListPublished),
		},
	}, nil
}

// GetCumulativeSizeFromDealList Sum the PieceSize for each deal in a provided list
func GetCumulativeSizeFromDealList(deals []*types.ProviderDealState) abi.PaddedPieceSize {
	var cumulativePieceSize abi.PaddedPieceSize

	for _, deal := range deals {
		if deal != nil {
			cumulativePieceSize += deal.ClientDealProposal.Proposal.PieceSize
		}
	}

	return cumulativePieceSize
}
