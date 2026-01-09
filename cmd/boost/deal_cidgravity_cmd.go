package main

import (
	"strings"
	"golang.org/x/xerrors"
	"net/http"
	"io"
	"encoding/json"
	"bytes"

	bcli "github.com/filecoin-project/boost/cli"
	clinode "github.com/filecoin-project/boost/cli/node"
	"github.com/filecoin-project/boost/cmd"
	"github.com/filecoin-project/boost/storagemarket/types"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/boost/storagemarket/types/legacytypes/network"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	chain_types "github.com/filecoin-project/lotus/chain/types"
	lcli "github.com/filecoin-project/lotus/cli"
	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"github.com/multiformats/go-multiaddr"
	"github.com/urfave/cli/v2"
)

const cidGravityEncryptLabelUrl = "https://service.cidgravity.com/private/v1/get-erc20-encoded-label"

var dealCIDgravityFlags = []cli.Flag{
	&cli.StringFlag{
		Name:  "label",
		Usage: "label to be specified in the proposal (default: root CID) will be ignored if erc20-deal is set to true",
		Value: "",
	},
	&cli.StringFlag{
		Name:     "provider",
		Usage:    "storage provider on-chain address",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "commp",
		Usage:    "commp of the CAR file",
		Required: true,
	},
	&cli.Uint64Flag{
		Name:     "piece-size",
		Usage:    "size of the CAR file as a padded piece",
		Required: true,
	},
	&cli.Uint64Flag{
		Name:     "car-size",
		Usage:    "size of the CAR file",
		Required: true,
	},
	&cli.StringFlag{
		Name:     "payload-cid",
		Usage:    "root CID of the CAR file",
		Required: true,
	},
	&cli.IntFlag{
		Name:        "start-epoch",
		Usage:       "start epoch by when the deal should be proved by provider on-chain",
		DefaultText: "current chain head + 2 days",
	},
	&cli.IntFlag{
		Name:  "duration",
		Usage: "duration of the deal in epochs",
		Value: 518400, // default is 2880 * 180 == 180 days
	},
	&cli.IntFlag{
		Name:  "provider-collateral",
		Usage: "deal collateral that storage miner must put in escrow; if empty, the min collateral for the given piece size will be used",
	},
	&cli.Int64Flag{
		Name:  "storage-price",
		Usage: "storage price in attoFIL per epoch per GiB (can be in USDFC/TiB/30d is flag --erc20-deal is true)",
		Value: 1,
	},
	&cli.BoolFlag{
		Name:  "verified",
		Usage: "whether the deal funds should come from verified client data-cap",
		Value: true,
	},
	&cli.BoolFlag{
		Name:  "erc20-deal",
		Usage: "whether the deal should be an ERC20 deal and price sent in USDFC (to use this, cidgravity-token is required)",
		Value: false,
	},
	&cli.StringFlag{
		Name:     "cidgravity-token",
		Usage:    "your CIDgravity token to send API requests (required to use erc20-deal flag)",
		Required: false,
	},
}

type CheckStatus string

const (
	Available     CheckStatus = "available"
	Unavailable   CheckStatus = "unavailable"
	InternalError CheckStatus = "internal_error"
	Unknown       CheckStatus = "unknown"
)

type dealCidGravityResponse struct {
	Status  CheckStatus `json:"status"`
	Reason  string      `json:"reason"`
	Message string      `json:"message"`

	Multiaddresses []multiaddr.Multiaddr `json:"multiaddresses"`
	PeerId         peer.ID               `json:"peerId"`

	GetAskPricePerGib         string `json:"getAskPricePerGib"`
	GetAskVerifiedPricePerGib string `json:"getAskVerifiedPricePerGib"`
	GetAskMinPieceSize        string `json:"getAskMinPieceSize"`
	GetAskMaxPieceSize        string `json:"getAskMaxPieceSize"`
	GetAskSectorSize          string `json:"getAskSectorSize"`

	DealProtocolsSupported protocol.ID `json:"dealProtocolsSupported"`
}

type cidGravityEncryptLabelPayload struct {
	Currency	string `json:"currency"`
	PieceCID	string `json:"pieceCID"`
	Price		string `json:"price"`
}

type cidGravityEncryptLabelResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`

	Result struct {
		EncodedLabel string `json:"encodedLabel"`
	} `json:"result"`
}

var dealCidGravityCmd = &cli.Command{
	Name:   "cidgravity-deal",
	Usage:  "Make an Keep Alive deal with Boost using CIDgravity proposal format (offline deal)",
	Flags:  dealCIDgravityFlags,
	Before: before,
	Action: func(cctx *cli.Context) error {
		return dealCidGravityCmdAction(cctx)
	},
}

func dealCidGravityCmdAction(cctx *cli.Context) error {
	ctx := bcli.ReqContext(cctx)

	// check required param for ERC20
	if cctx.Bool("erc20-deal") && cctx.String("cidgravity-token") == "" {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:  InternalError,
			Reason:  "ERR_INVALID_PARAMS",
			Message: "cidgravity-token must be provided to send ERC20 deals",
		})
	}

	n, err := clinode.Setup(cctx.String(cmd.FlagRepo.Name))
	if err != nil {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:  InternalError,
			Reason:  "ERR_API_GATEWAY",
			Message: err.Error(),
		})
	}

	api, closer, err := lcli.GetGatewayAPIV1(cctx)
	if err != nil {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:  InternalError,
			Reason:  "ERR_API_GATEWAY",
			Message: err.Error(),
		})
	}
	defer closer()

	walletAddr, err := n.GetProvidedOrDefaultWallet(ctx, cctx.String("wallet"))
	if err != nil {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:  InternalError,
			Reason:  "ERR_INVALID_WALLET",
			Message: err.Error(),
		})
	}

	log.Debugw("selected wallet", "wallet", walletAddr)

	maddr, err := address.NewFromString(cctx.String("provider"))
	if err != nil {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:  InternalError,
			Reason:  "ERR_INVALID_PROVIDER",
			Message: err.Error(),
		})
	}

	addrInfo, sectorSize, reason, err := cmd.GetAddrInfoCid(ctx, api, maddr)
	if err != nil {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:         Unavailable,
			Reason:         reason,
			Message:        err.Error(),
			Multiaddresses: []multiaddr.Multiaddr{},
		})
	}

	log.Debugw("found storage provider", "id", addrInfo.ID, "multiaddrs", addrInfo.Addrs, "addr", maddr)

	if err := n.Host.Connect(ctx, *addrInfo); err != nil {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:         Unavailable,
			Reason:         "ERR_CONNECT_MINER_PEER_ID",
			Message:        err.Error(),
			Multiaddresses: addrInfo.Addrs,
			PeerId:         addrInfo.ID,
		})
	}

	x, err := n.Host.Peerstore().FirstSupportedProtocol(addrInfo.ID, DealProtocolv120)
	if err != nil {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:                 Unavailable,
			Reason:                 "ERR_DEAL_PROTOCOL_UNSUPPORTED",
			Message:                err.Error(),
			Multiaddresses:         addrInfo.Addrs,
			PeerId:                 addrInfo.ID,
			DealProtocolsSupported: x,
		})
	}

	if len(x) == 0 {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:                 Unavailable,
			Reason:                 "ERR_NO_MATCHING_DEAL_PROTOCOL_SUPPORTED",
			Message:                "boost client cannot make a deal with storage provider because it does not support protocol version 1.2",
			Multiaddresses:         addrInfo.Addrs,
			PeerId:                 addrInfo.ID,
			DealProtocolsSupported: x,
		})
	}

	dealUuid := uuid.New()

	commP := cctx.String("commp")
	pieceCid, err := cid.Parse(commP)
	if err != nil {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:                 InternalError,
			Reason:                 "ERR_INVALID_PARAM_COMMP",
			Message:                "unable to parse commP : " + err.Error(),
			Multiaddresses:         addrInfo.Addrs,
			PeerId:                 addrInfo.ID,
			DealProtocolsSupported: x,
		})
	}

	pieceSize := cctx.Uint64("piece-size")
	if pieceSize == 0 {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:                 InternalError,
			Reason:                 "ERR_INVALID_PARAM_PIECE_SIZE",
			Message:                "must provide piece-size parameter for CAR url",
			Multiaddresses:         addrInfo.Addrs,
			PeerId:                 addrInfo.ID,
			DealProtocolsSupported: x,
		})
	}

	payloadCidStr := cctx.String("payload-cid")
	rootCid, err := cid.Parse(payloadCidStr)
	if err != nil {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:                 InternalError,
			Reason:                 "ERR_INVALID_PARAM_PAYLOAD_CID",
			Message:                "unable to parse payload cid : " + err.Error(),
			Multiaddresses:         addrInfo.Addrs,
			PeerId:                 addrInfo.ID,
			DealProtocolsSupported: x,
		})
	}

	carFileSize := cctx.Uint64("car-size")
	if carFileSize == 0 {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:                 InternalError,
			Reason:                 "ERR_INVALID_PARAM_CAR_SIZE",
			Message:                "size of car file cannot be 0",
			Multiaddresses:         addrInfo.Addrs,
			PeerId:                 addrInfo.ID,
			DealProtocolsSupported: x,
		})
	}

	transfer := types.Transfer{
		Size: carFileSize,
	}

	var providerCollateral abi.TokenAmount
	if cctx.IsSet("provider-collateral") {
		providerCollateral = abi.NewTokenAmount(cctx.Int64("provider-collateral"))
	} else {
		bounds, err := api.StateDealProviderCollateralBounds(ctx, abi.PaddedPieceSize(pieceSize), cctx.Bool("verified"), chain_types.EmptyTSK)
		if err != nil {
			return cmd.PrintJson(dealCidGravityResponse{
				Status:                 InternalError,
				Reason:                 "ERR_API_GATEWAY_CANNOT_RETRIEVE_MINIMUM_COLLATERAL",
				Message:                "node error getting collateral bounds : " + err.Error(),
				Multiaddresses:         addrInfo.Addrs,
				PeerId:                 addrInfo.ID,
				DealProtocolsSupported: x,
			})
		}

		providerCollateral = big.Div(big.Mul(bounds.Min, big.NewInt(6)), big.NewInt(5)) // add 20%
	}

	var startEpoch abi.ChainEpoch
	if cctx.IsSet("start-epoch") {
		startEpoch = abi.ChainEpoch(cctx.Int("start-epoch"))
	} else {
		tipset, err := api.ChainHead(ctx)
		if err != nil {
			return cmd.PrintJson(dealCidGravityResponse{
				Status:                 InternalError,
				Reason:                 "ERR_API_GATEWAY_CANNOT_RETRIEVE_CURRENT_CHAIN_HEAD",
				Message:                "unable to get chain head : " + err.Error(),
				Multiaddresses:         addrInfo.Addrs,
				PeerId:                 addrInfo.ID,
				DealProtocolsSupported: x,
			})
		}

		head := tipset.Height()

		log.Debugw("current block height", "number", head)

		startEpoch = head + abi.ChainEpoch(5760) // head + 2 days
	}

	// If flag erc20-deal is set to true
	// We need to encrypt the label with an API call to CIDgravity services
	// and use this encrypted label in the proposal
	// otherwise we take the label provided in the command line.
	label := cctx.String("label")

	if cctx.Bool("erc20-deal") {
		log.Debugw("about to send an API call to CIDgravity to get encrypted label")

		storagePrice := abi.NewTokenAmount(cctx.Int64("storage-price")).String()
		encryptedLabel, err := getEncryptedLabel(pieceCid.String(), storagePrice, cctx.String("cidgravity-token"))
		if err != nil {
			return cmd.PrintJson(dealCidGravityResponse{
				Status:                 Unavailable,
				Reason:                 "ERR_ENCRYPTED_LABEL",
				Message:                "failed to get encrypted label: " + err.Error(),
				Multiaddresses:         addrInfo.Addrs,
				PeerId:                 addrInfo.ID,
				DealProtocolsSupported: x,
			})
		}

		label = *encryptedLabel
	}

	// Send a get ask request for the final check
	log.Debugw("about to send a get ask request to", "address", maddr.String())

	streamGetAsk, err := n.Host.NewStream(ctx, addrInfo.ID, AskProtocolID)

	if err != nil {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:                 Unavailable,
			Reason:                 "ERR_GET_ASK",
			Message:                "failed to open stream to address : " + err.Error(),
			Multiaddresses:         addrInfo.Addrs,
			PeerId:                 addrInfo.ID,
			DealProtocolsSupported: x,
		})
	}

	defer streamGetAsk.Close()
	var resp network.AskResponse

	askRequest := network.AskRequest{
		Miner: maddr,
	}

	if err := doRpc(ctx, streamGetAsk, &askRequest, &resp); err != nil {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:                 Unavailable,
			Reason:                 "ERR_GET_ASK",
			Message:                "error while send get ask request rpc: " + err.Error(),
			Multiaddresses:         addrInfo.Addrs,
			PeerId:                 addrInfo.ID,
			DealProtocolsSupported: x,
		})
	}

	ask := resp.Ask.Ask

	log.Debugw("found get ask information",
		"Price per GiB", chain_types.FIL(ask.Price),
		"Verified Price per GiB", chain_types.FIL(ask.VerifiedPrice),
		"Min Piece size", chain_types.SizeStr(chain_types.NewInt(uint64(ask.MaxPieceSize))),
		"Max Piece size", chain_types.SizeStr(chain_types.NewInt(uint64(ask.MinPieceSize))),
		"Sector size", sectorSize.ShortString())

	// To work with CIDgravity, every price must be 0
	if !ask.Price.Equals(abi.NewTokenAmount(0)) || !ask.VerifiedPrice.Equals(abi.NewTokenAmount(0)) {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:                    Unavailable,
			Reason:                    "ERR_GET_ASK_PRICES_NOT_SET_TO_ZERO",
			Message:                   "get-ask price has to be set to 0",
			Multiaddresses:            addrInfo.Addrs,
			PeerId:                    addrInfo.ID,
			DealProtocolsSupported:    x,
			GetAskPricePerGib:         chain_types.FIL(ask.Price).String(),
			GetAskVerifiedPricePerGib: chain_types.FIL(ask.VerifiedPrice).String(),
			GetAskMinPieceSize:        chain_types.SizeStr(chain_types.NewInt(uint64(ask.MinPieceSize))),
			GetAskMaxPieceSize:        chain_types.SizeStr(chain_types.NewInt(uint64(ask.MaxPieceSize))),
			GetAskSectorSize:          sectorSize.ShortString(),
		})
	}

	// To work with CIDgravity, the min size must be set >= 256B and max size must be set to {sectorSize}
	if ask.MinPieceSize > 256 || ask.MaxPieceSize != abi.PaddedPieceSize(*sectorSize) {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:                    Unavailable,
			Reason:                    "ERR_GET_ASK_SIZES_NOT_PROPERLY_SET",
			Message:                   "get-ask accepting size must be min<=256B and max=" + sectorSize.ShortString(),
			Multiaddresses:            addrInfo.Addrs,
			PeerId:                    addrInfo.ID,
			DealProtocolsSupported:    x,
			GetAskPricePerGib:         chain_types.FIL(ask.Price).String(),
			GetAskVerifiedPricePerGib: chain_types.FIL(ask.VerifiedPrice).String(),
			GetAskMinPieceSize:        chain_types.SizeStr(chain_types.NewInt(uint64(ask.MinPieceSize))),
			GetAskMaxPieceSize:        chain_types.SizeStr(chain_types.NewInt(uint64(ask.MaxPieceSize))),
			GetAskSectorSize:          sectorSize.ShortString(),
		})
	}

	// Build the proposal
	dealProposal, err := dealProposal(ctx, label, n, walletAddr, rootCid, abi.PaddedPieceSize(pieceSize), pieceCid, maddr, startEpoch, cctx.Int("duration"), cctx.Bool("verified"), providerCollateral, abi.NewTokenAmount(cctx.Int64("storage-price")))
	if err != nil {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:                    InternalError,
			Reason:                    "ERR_SEND_PROPOSAL",
			Message:                   "unable to create a deal proposal : " + err.Error(),
			Multiaddresses:            addrInfo.Addrs,
			PeerId:                    addrInfo.ID,
			DealProtocolsSupported:    x,
			GetAskPricePerGib:         chain_types.FIL(ask.Price).String(),
			GetAskVerifiedPricePerGib: chain_types.FIL(ask.VerifiedPrice).String(),
			GetAskMinPieceSize:        chain_types.SizeStr(chain_types.NewInt(uint64(ask.MinPieceSize))),
			GetAskMaxPieceSize:        chain_types.SizeStr(chain_types.NewInt(uint64(ask.MaxPieceSize))),
			GetAskSectorSize:          sectorSize.ShortString(),
		})
	}

	// Generate the final proposal
	dealParams := types.DealParams{
		DealUUID:           dealUuid,
		ClientDealProposal: *dealProposal,
		DealDataRoot:       rootCid,
		IsOffline:          true,
		Transfer:           transfer,
		RemoveUnsealedCopy: true,
		SkipIPNIAnnounce:   false,
	}

	log.Debugw("about to submit deal proposal", "uuid", dealUuid.String())

	streamSendProposal, err := n.Host.NewStream(ctx, addrInfo.ID, DealProtocolv120)

	if err != nil {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:                    Unavailable,
			Reason:                    "ERR_SEND_PROPOSAL",
			Message:                   "failed to open stream to peer : " + err.Error(),
			Multiaddresses:            addrInfo.Addrs,
			PeerId:                    addrInfo.ID,
			DealProtocolsSupported:    x,
			GetAskPricePerGib:         chain_types.FIL(ask.Price).String(),
			GetAskVerifiedPricePerGib: chain_types.FIL(ask.VerifiedPrice).String(),
			GetAskMinPieceSize:        chain_types.SizeStr(chain_types.NewInt(uint64(ask.MinPieceSize))),
			GetAskMaxPieceSize:        chain_types.SizeStr(chain_types.NewInt(uint64(ask.MaxPieceSize))),
			GetAskSectorSize:          sectorSize.ShortString(),
		})
	}

	defer streamSendProposal.Close()

	var respDealResponse types.DealResponse
	if err := doRpc(ctx, streamSendProposal, &dealParams, &respDealResponse); err != nil {
		return cmd.PrintJson(dealCidGravityResponse{
			Status:                    Unavailable,
			Reason:                    "ERR_SEND_PROPOSAL",
			Message:                   "error while send proposal rpc : " + err.Error(),
			Multiaddresses:            addrInfo.Addrs,
			PeerId:                    addrInfo.ID,
			DealProtocolsSupported:    x,
			GetAskPricePerGib:         chain_types.FIL(ask.Price).String(),
			GetAskVerifiedPricePerGib: chain_types.FIL(ask.VerifiedPrice).String(),
			GetAskMinPieceSize:        chain_types.SizeStr(chain_types.NewInt(uint64(ask.MinPieceSize))),
			GetAskMaxPieceSize:        chain_types.SizeStr(chain_types.NewInt(uint64(ask.MaxPieceSize))),
			GetAskSectorSize:          sectorSize.ShortString(),
		})
	}

	// For CIDgravity miner-status-check the proposal will be rejected every time
	if !respDealResponse.Accepted {

		// So if the message contains the return code from miner-status-check service, we shouldn't consider as error
		if strings.Contains(respDealResponse.Message, "CIDgravity miner status check successful") {
			return cmd.PrintJson(dealCidGravityResponse{
				Status:                    Available,
				Reason:                    "DIAGNOSIS_SUCCESS",
				Message:                   "deal proposal successfully sent and received",
				Multiaddresses:            addrInfo.Addrs,
				PeerId:                    addrInfo.ID,
				DealProtocolsSupported:    x,
				GetAskPricePerGib:         chain_types.FIL(ask.Price).String(),
				GetAskVerifiedPricePerGib: chain_types.FIL(ask.VerifiedPrice).String(),
				GetAskMinPieceSize:        chain_types.SizeStr(chain_types.NewInt(uint64(ask.MinPieceSize))),
				GetAskMaxPieceSize:        chain_types.SizeStr(chain_types.NewInt(uint64(ask.MaxPieceSize))),
				GetAskSectorSize:          sectorSize.ShortString(),
			})
		}

		// Otherwise, return ERR_CIDGRAVITY_CONNECTOR_MISCONFIGURED
		return cmd.PrintJson(dealCidGravityResponse{
			Status:                    Unavailable,
			Reason:                    "ERR_CIDGRAVITY_CONNECTOR_MISCONFIGURED",
			Message:                   respDealResponse.Message,
			Multiaddresses:            addrInfo.Addrs,
			PeerId:                    addrInfo.ID,
			DealProtocolsSupported:    x,
			GetAskPricePerGib:         chain_types.FIL(ask.Price).String(),
			GetAskVerifiedPricePerGib: chain_types.FIL(ask.VerifiedPrice).String(),
			GetAskMinPieceSize:        chain_types.SizeStr(chain_types.NewInt(uint64(ask.MinPieceSize))),
			GetAskMaxPieceSize:        chain_types.SizeStr(chain_types.NewInt(uint64(ask.MaxPieceSize))),
			GetAskSectorSize:          sectorSize.ShortString(),
		})
	}

	// In our case, this block will never be reached (because all proposals will be rejected)
	return cmd.PrintJson(dealCidGravityResponse{
		Status:                    Unknown,
		Reason:                    "SENT",
		Message:                   "this case should normally never happend : unknown",
		Multiaddresses:            addrInfo.Addrs,
		PeerId:                    addrInfo.ID,
		DealProtocolsSupported:    x,
		GetAskPricePerGib:         chain_types.FIL(ask.Price).String(),
		GetAskVerifiedPricePerGib: chain_types.FIL(ask.VerifiedPrice).String(),
		GetAskMinPieceSize:        chain_types.SizeStr(chain_types.NewInt(uint64(ask.MinPieceSize))),
		GetAskMaxPieceSize:        chain_types.SizeStr(chain_types.NewInt(uint64(ask.MaxPieceSize))),
		GetAskSectorSize:          sectorSize.ShortString(),
	})
}

func getEncryptedLabel(pieceCID, storagePrice, cidgravityToken string) (*string, error) {
	data := cidGravityEncryptLabelPayload{
		Currency: 	"usdfc",
		PieceCID: 	pieceCID,
		Price: 		storagePrice,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, xerrors.Errorf("error encoding payload: %w", err)
	}

	// creating http client
	client := &http.Client{}

	// creating request
	req, err := http.NewRequest("POST", cidGravityEncryptLabelUrl, bytes.NewBuffer(payload))
	if err != nil {
		return nil, xerrors.Errorf("error creating request: %w", err)
	}

	// set authentication header
	req.Header.Set("X-API-KEY", cidgravityToken)

	// execute the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, xerrors.Errorf("error making request: %w", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Errorf("error closing response body: %w", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, xerrors.Errorf("invalid response from CIDgravity to get encrypted label")
	}

	body := new(bytes.Buffer)
	_, err = body.ReadFrom(resp.Body)
	if err != nil {
		return nil, xerrors.Errorf("error reading response body: %w", err)
	}

	response := cidGravityEncryptLabelResponse{}
	err = json.Unmarshal(body.Bytes(), &response)
	if err != nil {
		return nil, xerrors.Errorf("error parsing response body: %w", err)
	}

	return &response.Result.EncodedLabel, nil
}