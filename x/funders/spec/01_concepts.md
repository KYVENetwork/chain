<!--
order: 1
-->

# Concepts

This module contains the logic of maintaining and tracking the funds
of funders. Funders are special actors who provide liquidity to a pool
and basically pay for the rewards the validators earn for their work.
Funders would usually be stakeholders of the data that is being archived and 
therefore have a strong interest  in further archiving the data. Once a valid 
bundle is produced and the reward is paid out the pool module takes care of 
correctly deducting the funds equally from each funder in order to guarantee 
a steady pool economy.

## Funding Slots

Currently, the KYVE protocol allows at maximum 50 funders per pool to limit 
gas consumption. If the slots are full and a funder wants to join anyway he 
has to fund more than the current lowest funder. By doing so the funds of the
lowest funder will be automatically returned to the lowest funder's wallet
and the new funder can join.

## Multiple Coin Funding

To ease the process of funding KYVE pools we introduced the concept of multiple
coin fundings. This implies that funders can not only provide the native $KYVE
coin as funds, but also other coins which are provided over IBC. This affects
the community pool, the protocol validators and the delegators who not only receive 
$KYVE as rewards but also the other coins with which the pool was funded with.

## Price per Bundle

Funders can choose for themselves which total amount they want to contribute and
how much funds they want to distribute per validated and archived bundle of data.
This gives funders huge flexibility and promotes competition between pools in order
to attract validators who typically choose the pool with the highest provided funds.
