# Delta Chain

This guide provides step-by-step instructions to initialize, configure, and start the `delta-mainnet`.

---

## 📋 Table of Contents
1. [System Requirements](#-system-requirements)  
2. [Network Initialization](#-network-initialization)  
3. [Genesis Configuration](#-genesis-configuration)  
4. [Node Configuration](#-node-configuration)  
5. [Start the Node](#️-start-the-node)  
6. [Notes](#-notes)  

---

## 🖥️ System Requirements

**Software**
- [Go](https://go.dev/) (v1.23 recommended)

**Hardware (recommended for validators)**
- 4+ vCPUs  
- 16GB RAM  
- 100GB+ SSD storage  
- Stable internet connection  

---

## 🚀 Network Initialization

```bash
# Initialize the chain
deltad init delta-mainnet   --chain-id deltamainnet_1313-1   --default-denom adlta

#Modify Genesis with custom parameter
Replace the generated ~/.delta/config/genesis.json with the provided genesis.json and change genesis_time as per your genesis.

# Add a validator key
deltad keys add validator1
```

---

## 📜 Genesis Configuration

```bash
# Add a genesis account (1 Billion DLTA with 18 decimals)
deltad genesis add-genesis-account $(deltad keys show validator1 -a) 1000000000000000000000000000adlta

# Generate stake transaction (100,000 DLTA stake with 18 decimals)
deltad genesis gentx validator1 100000000000000000000000adlta   --chain-id deltamainnet_1313-1

# Collect and validate genesis transactions
deltad genesis collect-gentxs
deltad genesis validate-genesis
```

---

## ⚙️ Node Configuration

Edit the following configuration files located in `~/.delta/config/`.

### `app.toml`

```toml
# Minimum gas price
minimum-gas-prices = "0.025adlta"

# JSON-RPC & WebSocket endpoints
address = "0.0.0.0:8545"
ws-address = "0.0.0.0:8546"

# Enable JSON-RPC APIs
api = "eth,net,web3,txpool,debug"

# REST API configuration
enable = true
swagger = true
address = "tcp://0.0.0.0:1317"
```

### `config.toml`

```toml
# RPC endpoint
laddr = "tcp://0.0.0.0:26657"

# Block commit timeout
timeout_commit = "40s" # 🟢 Optional
```

---

## ▶️ Start the Node

```bash
deltad start
```

---

## ✅ Notes
- **Token Denomination**: `adlta` (1 DLTA = 10^18 adlta)  
- **Chain ID**: `deltamainnet_1313-1`  
- **Validator Stake**: 100,000 DLTA  
- **Genesis Allocation Total Supply**: 1 Billion DLTA  

---
