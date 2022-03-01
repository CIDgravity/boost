// Code generated by github.com/whyrusleeping/cbor-gen. DO NOT EDIT.

package types

import (
	"fmt"
	"io"
	"math"
	"sort"

	abi "github.com/filecoin-project/go-state-types/abi"
	cid "github.com/ipfs/go-cid"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

var _ = xerrors.Errorf
var _ = cid.Undef
var _ = math.E
var _ = sort.Sort

func (t *StorageAsk) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{165}); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Price (big.Int) (struct)
	if len("Price") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Price\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Price"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Price")); err != nil {
		return err
	}

	if err := t.Price.MarshalCBOR(w); err != nil {
		return err
	}

	// t.VerifiedPrice (big.Int) (struct)
	if len("VerifiedPrice") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"VerifiedPrice\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("VerifiedPrice"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("VerifiedPrice")); err != nil {
		return err
	}

	if err := t.VerifiedPrice.MarshalCBOR(w); err != nil {
		return err
	}

	// t.MinPieceSize (abi.PaddedPieceSize) (uint64)
	if len("MinPieceSize") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"MinPieceSize\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("MinPieceSize"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("MinPieceSize")); err != nil {
		return err
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.MinPieceSize)); err != nil {
		return err
	}

	// t.MaxPieceSize (abi.PaddedPieceSize) (uint64)
	if len("MaxPieceSize") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"MaxPieceSize\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("MaxPieceSize"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("MaxPieceSize")); err != nil {
		return err
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.MaxPieceSize)); err != nil {
		return err
	}

	// t.Miner (address.Address) (struct)
	if len("Miner") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Miner\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Miner"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Miner")); err != nil {
		return err
	}

	if err := t.Miner.MarshalCBOR(w); err != nil {
		return err
	}
	return nil
}

func (t *StorageAsk) UnmarshalCBOR(r io.Reader) error {
	*t = StorageAsk{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("StorageAsk: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadStringBuf(br, scratch)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.Price (big.Int) (struct)
		case "Price":

			{

				if err := t.Price.UnmarshalCBOR(br); err != nil {
					return xerrors.Errorf("unmarshaling t.Price: %w", err)
				}

			}
			// t.VerifiedPrice (big.Int) (struct)
		case "VerifiedPrice":

			{

				if err := t.VerifiedPrice.UnmarshalCBOR(br); err != nil {
					return xerrors.Errorf("unmarshaling t.VerifiedPrice: %w", err)
				}

			}
			// t.MinPieceSize (abi.PaddedPieceSize) (uint64)
		case "MinPieceSize":

			{

				maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
				if err != nil {
					return err
				}
				if maj != cbg.MajUnsignedInt {
					return fmt.Errorf("wrong type for uint64 field")
				}
				t.MinPieceSize = abi.PaddedPieceSize(extra)

			}
			// t.MaxPieceSize (abi.PaddedPieceSize) (uint64)
		case "MaxPieceSize":

			{

				maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
				if err != nil {
					return err
				}
				if maj != cbg.MajUnsignedInt {
					return fmt.Errorf("wrong type for uint64 field")
				}
				t.MaxPieceSize = abi.PaddedPieceSize(extra)

			}
			// t.Miner (address.Address) (struct)
		case "Miner":

			{

				if err := t.Miner.UnmarshalCBOR(br); err != nil {
					return xerrors.Errorf("unmarshaling t.Miner: %w", err)
				}

			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
func (t *DealParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{164}); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.DealUUID (uuid.UUID) (array)
	if len("DealUUID") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"DealUUID\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("DealUUID"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("DealUUID")); err != nil {
		return err
	}

	if len(t.DealUUID) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.DealUUID was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.DealUUID))); err != nil {
		return err
	}

	if _, err := w.Write(t.DealUUID[:]); err != nil {
		return err
	}

	// t.ClientDealProposal (market.ClientDealProposal) (struct)
	if len("ClientDealProposal") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"ClientDealProposal\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("ClientDealProposal"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("ClientDealProposal")); err != nil {
		return err
	}

	if err := t.ClientDealProposal.MarshalCBOR(w); err != nil {
		return err
	}

	// t.DealDataRoot (cid.Cid) (struct)
	if len("DealDataRoot") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"DealDataRoot\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("DealDataRoot"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("DealDataRoot")); err != nil {
		return err
	}

	if err := cbg.WriteCidBuf(scratch, w, t.DealDataRoot); err != nil {
		return xerrors.Errorf("failed to write cid field t.DealDataRoot: %w", err)
	}

	// t.Transfer (types.Transfer) (struct)
	if len("Transfer") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Transfer\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Transfer"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Transfer")); err != nil {
		return err
	}

	if err := t.Transfer.MarshalCBOR(w); err != nil {
		return err
	}

	return nil
}

func (t *DealParams) UnmarshalCBOR(r io.Reader) error {
	*t = DealParams{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("DealParams: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadStringBuf(br, scratch)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.DealUUID (uuid.UUID) (array)
		case "DealUUID":

			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}

			if extra > cbg.ByteArrayMaxLen {
				return fmt.Errorf("t.DealUUID: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra != 16 {
				return fmt.Errorf("expected array to have 16 elements")
			}

			t.DealUUID = [16]uint8{}

			if _, err := io.ReadFull(br, t.DealUUID[:]); err != nil {
				return err
			}
			// t.ClientDealProposal (market.ClientDealProposal) (struct)
		case "ClientDealProposal":

			{

				if err := t.ClientDealProposal.UnmarshalCBOR(br); err != nil {
					return xerrors.Errorf("unmarshaling t.ClientDealProposal: %w", err)
				}

			}
			// t.DealDataRoot (cid.Cid) (struct)
		case "DealDataRoot":

			{

				c, err := cbg.ReadCid(br)
				if err != nil {
					return xerrors.Errorf("failed to read cid field t.DealDataRoot: %w", err)
				}

				t.DealDataRoot = c

			}
			// t.Transfer (types.Transfer) (struct)
		case "Transfer":

			{

				if err := t.Transfer.UnmarshalCBOR(br); err != nil {
					return xerrors.Errorf("unmarshaling t.Transfer: %w", err)
				}

			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
func (t *Transfer) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{164}); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Type (string) (string)
	if len("Type") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Type\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Type"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Type")); err != nil {
		return err
	}

	if len(t.Type) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Type was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len(t.Type))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Type)); err != nil {
		return err
	}

	// t.ClientID (string) (string)
	if len("ClientID") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"ClientID\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("ClientID"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("ClientID")); err != nil {
		return err
	}

	if len(t.ClientID) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.ClientID was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len(t.ClientID))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.ClientID)); err != nil {
		return err
	}

	// t.Params ([]uint8) (slice)
	if len("Params") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Params\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Params"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Params")); err != nil {
		return err
	}

	if len(t.Params) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.Params was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.Params))); err != nil {
		return err
	}

	if _, err := w.Write(t.Params[:]); err != nil {
		return err
	}

	// t.Size (uint64) (uint64)
	if len("Size") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Size\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Size"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Size")); err != nil {
		return err
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajUnsignedInt, uint64(t.Size)); err != nil {
		return err
	}

	return nil
}

func (t *Transfer) UnmarshalCBOR(r io.Reader) error {
	*t = Transfer{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("Transfer: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadStringBuf(br, scratch)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.Type (string) (string)
		case "Type":

			{
				sval, err := cbg.ReadStringBuf(br, scratch)
				if err != nil {
					return err
				}

				t.Type = string(sval)
			}
			// t.ClientID (string) (string)
		case "ClientID":

			{
				sval, err := cbg.ReadStringBuf(br, scratch)
				if err != nil {
					return err
				}

				t.ClientID = string(sval)
			}
			// t.Params ([]uint8) (slice)
		case "Params":

			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}

			if extra > cbg.ByteArrayMaxLen {
				return fmt.Errorf("t.Params: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra > 0 {
				t.Params = make([]uint8, extra)
			}

			if _, err := io.ReadFull(br, t.Params[:]); err != nil {
				return err
			}
			// t.Size (uint64) (uint64)
		case "Size":

			{

				maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
				if err != nil {
					return err
				}
				if maj != cbg.MajUnsignedInt {
					return fmt.Errorf("wrong type for uint64 field")
				}
				t.Size = uint64(extra)

			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
func (t *DealResponse) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{162}); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.Accepted (bool) (bool)
	if len("Accepted") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Accepted\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Accepted"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Accepted")); err != nil {
		return err
	}

	if err := cbg.WriteBool(w, t.Accepted); err != nil {
		return err
	}

	// t.Message (string) (string)
	if len("Message") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Message\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Message"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Message")); err != nil {
		return err
	}

	if len(t.Message) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Message was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len(t.Message))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Message)); err != nil {
		return err
	}
	return nil
}

func (t *DealResponse) UnmarshalCBOR(r io.Reader) error {
	*t = DealResponse{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("DealResponse: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadStringBuf(br, scratch)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.Accepted (bool) (bool)
		case "Accepted":

			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}
			if maj != cbg.MajOther {
				return fmt.Errorf("booleans must be major type 7")
			}
			switch extra {
			case 20:
				t.Accepted = false
			case 21:
				t.Accepted = true
			default:
				return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
			}
			// t.Message (string) (string)
		case "Message":

			{
				sval, err := cbg.ReadStringBuf(br, scratch)
				if err != nil {
					return err
				}

				t.Message = string(sval)
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
func (t *DealStatusRequest) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{161}); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.DealUUID (uuid.UUID) (array)
	if len("DealUUID") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"DealUUID\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("DealUUID"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("DealUUID")); err != nil {
		return err
	}

	if len(t.DealUUID) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.DealUUID was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.DealUUID))); err != nil {
		return err
	}

	if _, err := w.Write(t.DealUUID[:]); err != nil {
		return err
	}
	return nil
}

func (t *DealStatusRequest) UnmarshalCBOR(r io.Reader) error {
	*t = DealStatusRequest{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("DealStatusRequest: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadStringBuf(br, scratch)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.DealUUID (uuid.UUID) (array)
		case "DealUUID":

			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}

			if extra > cbg.ByteArrayMaxLen {
				return fmt.Errorf("t.DealUUID: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra != 16 {
				return fmt.Errorf("expected array to have 16 elements")
			}

			t.DealUUID = [16]uint8{}

			if _, err := io.ReadFull(br, t.DealUUID[:]); err != nil {
				return err
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
func (t *DealStatusResponse) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{163}); err != nil {
		return err
	}

	scratch := make([]byte, 9)

	// t.DealUUID (uuid.UUID) (array)
	if len("DealUUID") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"DealUUID\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("DealUUID"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("DealUUID")); err != nil {
		return err
	}

	if len(t.DealUUID) > cbg.ByteArrayMaxLen {
		return xerrors.Errorf("Byte array in field t.DealUUID was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajByteString, uint64(len(t.DealUUID))); err != nil {
		return err
	}

	if _, err := w.Write(t.DealUUID[:]); err != nil {
		return err
	}

	// t.Error (string) (string)
	if len("Error") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"Error\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("Error"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("Error")); err != nil {
		return err
	}

	if len(t.Error) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.Error was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len(t.Error))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.Error)); err != nil {
		return err
	}

	// t.DealStatus (string) (string)
	if len("DealStatus") > cbg.MaxLength {
		return xerrors.Errorf("Value in field \"DealStatus\" was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len("DealStatus"))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string("DealStatus")); err != nil {
		return err
	}

	if len(t.DealStatus) > cbg.MaxLength {
		return xerrors.Errorf("Value in field t.DealStatus was too long")
	}

	if err := cbg.WriteMajorTypeHeaderBuf(scratch, w, cbg.MajTextString, uint64(len(t.DealStatus))); err != nil {
		return err
	}
	if _, err := io.WriteString(w, string(t.DealStatus)); err != nil {
		return err
	}
	return nil
}

func (t *DealStatusResponse) UnmarshalCBOR(r io.Reader) error {
	*t = DealStatusResponse{}

	br := cbg.GetPeeker(r)
	scratch := make([]byte, 8)

	maj, extra, err := cbg.CborReadHeaderBuf(br, scratch)
	if err != nil {
		return err
	}
	if maj != cbg.MajMap {
		return fmt.Errorf("cbor input should be of type map")
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("DealStatusResponse: map struct too large (%d)", extra)
	}

	var name string
	n := extra

	for i := uint64(0); i < n; i++ {

		{
			sval, err := cbg.ReadStringBuf(br, scratch)
			if err != nil {
				return err
			}

			name = string(sval)
		}

		switch name {
		// t.DealUUID (uuid.UUID) (array)
		case "DealUUID":

			maj, extra, err = cbg.CborReadHeaderBuf(br, scratch)
			if err != nil {
				return err
			}

			if extra > cbg.ByteArrayMaxLen {
				return fmt.Errorf("t.DealUUID: byte array too large (%d)", extra)
			}
			if maj != cbg.MajByteString {
				return fmt.Errorf("expected byte array")
			}

			if extra != 16 {
				return fmt.Errorf("expected array to have 16 elements")
			}

			t.DealUUID = [16]uint8{}

			if _, err := io.ReadFull(br, t.DealUUID[:]); err != nil {
				return err
			}
			// t.Error (string) (string)
		case "Error":

			{
				sval, err := cbg.ReadStringBuf(br, scratch)
				if err != nil {
					return err
				}

				t.Error = string(sval)
			}
			// t.DealStatus (string) (string)
		case "DealStatus":

			{
				sval, err := cbg.ReadStringBuf(br, scratch)
				if err != nil {
					return err
				}

				t.DealStatus = string(sval)
			}

		default:
			// Field doesn't exist on this type, so ignore it
			cbg.ScanForLinks(r, func(cid.Cid) {})
		}
	}

	return nil
}
