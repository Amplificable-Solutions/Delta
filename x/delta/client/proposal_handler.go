package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/delta-chain/delta/v2/x/delta/client/cli"
)

// ProposalHandler is the token mapping change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewSubmitTokenMappingChangeProposalTxCmd)
