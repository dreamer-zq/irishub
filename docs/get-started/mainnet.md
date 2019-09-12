---
order: 2
---

# Join The Mainnet

:::tip
We assume that you already have `iris` installed, or you need to [install iris](../software/How-to-install-irishub.md) first.
:::

## Run a Full Node

```bash
# initialize node configurations
iris init --moniker=<your_custom_name> --chain-id=irishub

# download mainnet public config.toml and genesis.json
curl -o ~/.iris/config/config.toml https://raw.githubusercontent.com/irisnet/betanet/master/config/config.toml
curl -o ~/.iris/config/genesis.json https://raw.githubusercontent.com/irisnet/betanet/master/config/genesis.json

# start the node (you can also use "nohup" to run in the background)
iris start
```

:::tip
You may see some connection errors, it does not matter, the P2P network is trying to find available connections

[Advanced Configurations](#TODO)

[Community Peers](https://github.com/irisnet/betanet/blob/master/config/community-peers.md)
:::

:::tip
It will take a long time to catch up the latest block, you can also download the [mainnet data snapshot](#TODO) to reduce the time spent on synchronization
:::

## Upgrade to Validator Node

### Create a Wallet

You can [create a new wallet](../cli-client/keys/add.md#create-a-new-key) or [import an existing one](../cli-client/keys/add.md#recover-an-existing-key), then get some IRIS from the exchanges or anywhere else into the wallet you just created, .e.g.

```bash
# create a new wallet
iriscli keys add <key_name>
```

:::warning
**Important**

write the seed phrase in a safe place! It is the only way to recover your account if you ever forget your password.
:::

### Confirm your node has caught-up

```bash
# if you have not installed jq
# apt-get update && apt-get install -y jq

# if the output is false, means your node has caught-up
iriscli status | jq .sync_info.catching_up
```

### Create Validator

Only if your node has caught-up, you can run the following command to upgrade your node to be a validator.

```bash
iriscli stake create-validator \
    --pubkey=$(iris tendermint show-validator) \
    --moniker=<your_validator_name> \
    --amount=<amount_to_be_delegated, e.g. 10000iris> \
    --commission-rate=0.1 \
    --gas=100000 \
    --fee=0.6iris \
    --chain-id=irishub \
    --from=<key_name> \
    --commit
```

:::warning
**Important**

Backup the `config` directory located in your iris home (default ~/.iris/) carefully! It is the only way to recover your validator.
:::

If there are no errors, then your node is now a validator or candidate (depending on whether your delegation amount is in the top 100)

More details:

- [Concepts](#TODO)
- [Validator FAQ](#TODO)