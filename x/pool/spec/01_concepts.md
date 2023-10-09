<!--
order: 1
-->

# Concepts

This module contains the building block of validating and archiving
data with KYVE, the storage pools. KYVE allows multiple pools to exist
at once, all validating different kinds of data sources. A pool for
example could be responsible for validating Bitcoin data, another
pool for Ethereum data. Each staker then can join multiple pools at once
to validate more data and in return earn more rewards in $KYVE.

## Storage Pool

A storage pool is responsible for validating and archiving
a single type of data. As of now each pool can have up to 50 validators, where
the requirement of validating data in a pool is that those validators have a cumulative stake
greater or equal to the specified minimum stake.

## Keeping Pools Funded

Furthermore, funders are special actors who provide liquidity to a pool and basically pay
for the rewards the validators earn for their work. Funders would usually be
stakeholders of the data that is being archived and therefore have a strong interest
in further archiving the data. Once a valid bundle is produced and the reward is paid
out the pool module takes care of correctly deducting the funds equally from each funder
in order to guarantee a steady pool economy.

## Inflation Splitting

In order to support funders inflation splitting was introduced where a part of the block inflation
goes to the protocol and is paid out with the funds from the funders. This relieves the burden of the
funders to keep a pool alive and allows a pool to even run without any funds.
