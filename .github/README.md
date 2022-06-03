# The KYVE Chain

###### v0.4.0

The chain nodes are the backbone of KYVE. The chain layer is a completely sovereign
[Proof of Stake](https://en.wikipedia.org/wiki/Proof_of_stake) blockchain build with
[Cosmos SDK](https://github.com/cosmos/cosmos-sdk) using the [Ignite CLI](https://ignt.com/cli). This blockchain is run
by independent nodes we call _Chain Nodes_ since they're running on the chain level. The native currency of the KYVE
chain is [$KYVE](https://docs.kyve.network/basics/kyve.html), it secures the chain and allows chain nodes to stake and
other users to delegate into them.

---

## Building from source

To build from source, the [Ignite CLI](https://ignt.com/cli) is required.

```sh
ignite chain build --release --release.prefix kyve
```

The output can be found in `./release`.

If you need to build for different architectures, use the `-t` flag, e.g. `-t linux:amd64,linux:arm64`.

## Running a chain node

Full documentation for setting up a chain node are provided [here](https://docs.kyve.network/intro/chain-node.html).
