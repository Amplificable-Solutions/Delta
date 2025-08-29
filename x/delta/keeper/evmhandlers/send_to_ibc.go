package evmhandler

import (
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"

	deltakeeper "github.com/delta-chain/delta/v2/x/delta/keeper"
	"github.com/delta-chain/delta/v2/x/delta/types"
)

var _ types.EvmLogHandler = SendToIbcHandler{}

const SendToIbcEventName = "__DeltaSendToIbc"

// SendToIbcEvent represent the signature of
// `event __DeltaSendToIbc(address sender, string recipient, uint256 amount)`
var SendToIbcEvent abi.Event

func init() {
	addressType, _ := abi.NewType("address", "", nil)
	uint256Type, _ := abi.NewType("uint256", "", nil)
	stringType, _ := abi.NewType("string", "", nil)

	SendToIbcEvent = abi.NewEvent(
		SendToIbcEventName,
		SendToIbcEventName,
		false,
		abi.Arguments{abi.Argument{
			Name:    "sender",
			Type:    addressType,
			Indexed: false,
		}, abi.Argument{
			Name:    "recipient",
			Type:    stringType,
			Indexed: false,
		}, abi.Argument{
			Name:    "amount",
			Type:    uint256Type,
			Indexed: false,
		}},
	)
}

// SendToIbcHandler handles `__DeltaSendToIbc` log
type SendToIbcHandler struct {
	bankKeeper  types.BankKeeper
	deltaKeeper deltakeeper.Keeper
}

func NewSendToIbcHandler(bankKeeper types.BankKeeper, deltaKeeper deltakeeper.Keeper) *SendToIbcHandler {
	return &SendToIbcHandler{
		bankKeeper:  bankKeeper,
		deltaKeeper: deltaKeeper,
	}
}

func (h SendToIbcHandler) EventID() common.Hash {
	return SendToIbcEvent.ID
}

func (h SendToIbcHandler) Handle(
	ctx sdk.Context,
	contract common.Address,
	topics []common.Hash,
	data []byte,
	_ func(contractAddress common.Address, logSig common.Hash, logData []byte),
) error {
	unpacked, err := SendToIbcEvent.Inputs.Unpack(data)
	if err != nil {
		// log and ignore
		h.deltaKeeper.Logger(ctx).Error("log signature matches but failed to decode", "error", err)
		return nil
	}
	sender := unpacked[0].(common.Address)
	recipient := unpacked[1].(string)
	amount := unpacked[2].(*big.Int)
	return h.handle(ctx, contract, sender, recipient, amount, nil)
}

func (h SendToIbcHandler) handle(
	ctx sdk.Context,
	contract common.Address,
	senderAddress common.Address,
	recipient string,
	amountInt *big.Int,
	id *big.Int,
) error {
	denom, found := h.deltaKeeper.GetDenomByContract(ctx, contract)
	if !found {
		return fmt.Errorf("contract %s is not connected to native token", contract)
	}

	if !types.IsValidIBCDenom(denom) && !types.IsValidDeltaDenom(denom) {
		return fmt.Errorf("the native token associated with the contract %s is neither an ibc voucher or a delta token", contract)
	}

	contractAddr := sdk.AccAddress(contract.Bytes())
	sender := sdk.AccAddress(senderAddress.Bytes())
	amount := sdkmath.NewIntFromBigInt(amountInt)
	coins := sdk.NewCoins(sdk.NewCoin(denom, amount))

	var err error
	if types.IsSourceCoin(denom) {
		// it is a source token, we need to mint coins
		if err = h.bankKeeper.MintCoins(ctx, types.ModuleName, coins); err != nil {
			return err
		}
		// send the coin to the user
		if err = h.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, coins); err != nil {
			return err
		}
	} else {
		// First, transfer IBC coin to user so that he will be the refunded address if transfer fails
		if err = h.bankKeeper.SendCoins(ctx, contractAddr, sender, coins); err != nil {
			return err
		}
	}

	channelId := ""
	if id != nil {
		channelId = "channel-" + id.String()
	}
	// Initiate IBC transfer from sender account
	if err = h.deltaKeeper.IbcTransferCoins(ctx, sender.String(), recipient, coins, channelId); err != nil {
		return err
	}
	return nil
}
