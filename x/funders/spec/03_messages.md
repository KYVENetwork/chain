<!--
order: 3
-->

# Messages

## MsgCreateFunder

MsgCreateFunder creates the funder with all relevant information about the funder, this includes a moniker (human 
readable name), an identity (keybase.io identity with avatar), a website, a contact and a description.
Note that this transaction fails if a funder with the same address already exists or if the moniker is empty.

## MsgUpdateFunder

MsgUpdateFunder can update all properties of the funder except the funder's address.

## MsgFundPool

MsgFundPool commits funds from the funder to a specific pool. The funder can fund multiple coins at the same time, but
they have to be whitelisted by the KYVE protocol. For each coin the funder funds, the 
`amount_per_bundle` has to be specified, too. This parameter specifies how much of each coin gets distributed
per finalized bundle to the protocol validators.

## MsgDefundPool

MsgDefundPool can withdraw remaining funds on the protocol if the funder decides to get his funds back or to allocate
them to a different pool. It takes a list of coins, so multiple coins can also be withdrawn in a single transaction.
If a funder wants to defund a coin which was removed from the whitelist since he funded he has to defund the entire
amount of that coin, else the transaction will fail.

## MsgUpdateParams

MsgUpdateParams is a gov transaction and can be only called by the governance authority. To submit this transaction
someone has to create a MsgUpdateParams governance proposal.

This will update the parameters of the funders module, containing the coin whitelist and the `MinFundingMultiple`
parameter.
