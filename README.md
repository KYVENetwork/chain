# The KYVE Network

###### v1.4.0

The KYVE consensus layer is the backbone of the KYVE ecosystem. This layer is a
sovereign Delegated Proof of Stake network built using the
[Cosmos SDK](https://github.com/cosmos/cosmos-sdk) and
[CometBFT (BFT Consensus)](https://github.com/cometbft/cometbft).

## Building from source

```shell
make build
```

You can find the `kyved` binary in the `./build` directory.

If you need binaries for alternative architectures than your host:

```shell
make release ENV=mainnet
```

The different binaries can be found in the `./release` directory.

### Building for the Kaon testnet

If you want to build the binary for the Kaon testnet, you will need to specify
its build environment. This is important because mainnet and testnet use
different denoms for the native token.

```shell
make build ENV=kaon
```

You can verify the build information using the following command:

```shell
kyved info
```

### Building docker image

#### Root
To create a regular `kyve-network/kyve:${VERSION}` docker image with `kyved` binary only execute:
```bash
make docker-image
```
To create the corresponding debug image containing a `sh` shell execute `make docker-image-debug`.

#### Nonroot
To create a nonroot docker image `kyve-network/kyve:${VERSION}-nonroot` running with user `nonroot:65532`:
```bash
make docker-image-nonroot
```
To create the corresponding debug image `kyve-network/kyve:${VERSION}-debug-nonroot` execute `make docker-image-debug-nonroot`.