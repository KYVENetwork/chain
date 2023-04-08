# `x/oracle` Module

## IBC Middleware

Inspired by [Osmosis' IBC Hooks](https://github.com/osmosis-labs/osmosis/tree/main/x/ibc-hooks), we have implemented our oracle as an [ICS-20](https://ibc.cosmos.network/main/apps/transfer/overview.html) [middleware](https://ibc.cosmos.network/main/ibc/middleware/overview.html). This was the cleanest and most backwards-compatible solution for implementing paid interchain queries.

A valid oracle packet is a standard [`MsgTransfer`](https://ibc.cosmos.network/main/apps/transfer/messages.html#msgtransfer) message sent from any chain (**"Querier Chain"**) to the KYVE chain (**"Host Chain"**). The receiver of the transferred tokens must be the `x/oracle` module address (`kyve1jgp27m8fykex4e4jtt0l7ze8q528ux2lxl4pxd`).

### ICS-20 Memo Specification

We utilise the memo field of an ICS-20 packet to trigger a query. The memo field must follow the following format. Please note that you can only query the latest bundle summary of a specific KYVE data pool for now. We will enable more queries in the future.

```json lines
{
  "query": {
    "latestSummary": {
      "poolId": 0 // any data pool ID.
    }
  }
}
```
