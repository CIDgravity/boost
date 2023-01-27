package stagingdeals

import (
	"context"

	"github.com/filecoin-project/boost/db"
	abi "github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/boost/storagemarket/types"
)

type CheckpointStatus struct {
	Checkpoint string
	Count CheckpointCounters
	CumulativePieceSize  CheckpointPieceSize
}

type CheckpointCounters struct {
	Total int64
	InError int64
}

type CheckpointPieceSize struct {
	Total abi.PaddedPieceSize
	InError abi.PaddedPieceSize 
}

type Status struct {
	CheckpointsStatus []*CheckpointStatus
}

func GetStatus(ctx context.Context, dealsDB *db.DealsDB) (*Status, error) {
	checkpoints := [5]string{"Accepted", "Transferred", "Published", "PublishConfirmed", "AddedPiece"}

	var checkpointsStatus []*CheckpointStatus

	for _, checkpoint := range checkpoints {

		// Gets deals without err for current checkpoint
		dealsWithoutErr, err := dealsDB.ByErrorNullAndCheckpoint(ctx, checkpoint)

		if err != nil {
			return nil, err
		}

		// Gets deals with err for current checkpoint
		dealsWithErr, err := dealsDB.ByErrorNotNullAndCheckpoint(ctx, checkpoint)

		if err != nil {
			return nil, err
		}

		checkpointsStatus = append(checkpointsStatus, 
			&CheckpointStatus{ 
				Checkpoint: checkpoint,  
				Count: CheckpointCounters {
					Total: int64(len(dealsWithoutErr)) + int64(len(dealsWithErr)),
					InError: int64(len(dealsWithErr)),
				},
				CumulativePieceSize: CheckpointPieceSize {
					Total: GetCumulativePieceSize(dealsWithErr) + GetCumulativePieceSize(dealsWithoutErr),
					InError: GetCumulativePieceSize(dealsWithErr),
				},
			},
		)
	}

	return &Status{
		CheckpointsStatus: checkpointsStatus,
	}, nil
}

func GetCumulativePieceSize(deals []*types.ProviderDealState) abi.PaddedPieceSize {
	var cumulativePieceSize abi.PaddedPieceSize

	for _, deal := range deals {
		if (deal != nil) {
			cumulativePieceSize += deal.ClientDealProposal.Proposal.PieceSize
		}
	}

	return cumulativePieceSize
}