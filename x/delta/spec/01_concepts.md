 <!-- order: 1 -->

# Concepts

## Delta

Delta is an EVM chain in Crypto.org ecosystem. It allows users to deploy and interact with smart contracts. Powered by Ethermint, it is built using cosmos SDK which allows to leverage to full capability of the cosmos ecosystem. It is also connected to the ethereum network with the gravity bridge integration (WIP).

Bridging asset from cosmos or ethereum are automatically converted to a CRC20 asset when moved to Delta which make it extremely easy to integrate with existing web3 tools.

The Delta module glues IBC, gravity bridge and ethermint together through hooks and token mapping.

## Gas Token

Delta uses CRO as its gas token.

The CRO assets on the Delta Chain need to be transferred from the Crypto.org Chain through an IBC channel. The tokens arrived at the Delta Chain as IBC tokens first, then are automatically converted to the gas token. 

### Decimal Places Conversion

The CRO on the Crypto.org Chain has 8 decimals, while the Delta gas token has 18 decimals (to keep compatibility with Ethereum). So a conversion is done automatically:

- When transferring CRO to Delta chain, the decimal places of the amount are expanded to 18.
- When transferring CRO from Delta chain, the amount is truncated to 8 decimals, the remaining part is left in Delta, so the total value is preserved.

## Native Token

Native token is a token managed by cosmos native bank module, there are several kinds of native tokens in Delta:

- gas token. used to pay the gas fee.
- staking token. used to do PoA consensus.
- IBC token. tokens come from IBC channels.
- gravity token. tokens come from the gravity bridge.

## CRC20 Token

CRC20 token is Delta's equivalence of ERC20 token on Ethereum with some extensions, they can be mapped with native tokens and support transfer to/from native tokens, and potentially transfer to/from Ethereum and another cosmos chain through gravity bridge and IBC.

## Auto-deployed Contract

A contract whose byte code is embedded in Delta module and deployed by it, and some operations in it are only authorized to Delta module. These contracts can be trusted by Delta module directly. Currently, Delta module support auto-deploy a minimal CRC20 contract to support automatic wrapping native token in EVM.

+++ https://github.com/delta-chain/delta/blob/v0.6.0-testnet/contracts/src/ModuleCRC20.sol#L5-L52

## Token Mapping

To support transfer tokens between native tokens and EVM tokens, the Delta module maintains two mappings between native denom to contract address, one for auto-deployed contracts, one for external contracts.

When auto-deployment is enabled, incoming IBC and gravity native tokens are wrapped to an auto-deployed CRC20 contract automatically.

One can also register an external contract mapping for the denom, either through the governance process or an authorized transaction.
