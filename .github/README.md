# The KYVE Network

###### v1.0.0-rc0

The KYVE consensus layer is the backbone of the KYVE ecosystem. The layer is a
sovereign Delegated Proof of Stake network built using the
[Cosmos SDK](https://github.com/cosmos/cosmos-sdk) and
[Tendermint Core (BFT Consensus)](https://github.com/tendermint/tendermint).

## Building from Source

```shell
make build
```

You can find the `kyved` binary in the `./build` directory.

If you need binaries for alternative architectures than your host:

```shell
make release
```

The different binaries can be found in the `./release` directory.
