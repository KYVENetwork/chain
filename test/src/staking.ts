import BigNumber from "bignumber.js";
import {
  ADDRESS_ALICE,
  ADDRESS_BOB,
  ADDRESS_CHARLIE,
  alice,
  bob,
  charlie,
} from "./helpers/accounts";
import { lcdClient, restartChain, sleep } from "./helpers/utils";
import { constants } from "@kyve/sdk";
import { funding } from "./funding";

const getDefaultPool = async () =>
  (await lcdClient.kyve.registry.v1beta1.pool({ id: "0" })).pool ??
  (() => {
    throw new Error("Pool doesn't exist");
  })();

export const staking = () => {
  afterAll(async () => {
    await restartChain();
    await alice.init();
    await bob.init();
    await charlie.init();
  });
  // disable timeout
  jest.setTimeout(24 * 60 * 60 * 1000);

  test("Test if default pool exist", async () => {
    const poolsResponse = await lcdClient.kyve.registry.v1beta1.pools({
      paused: false,
    });

    expect(poolsResponse.pools).not.toBeUndefined();
    expect(poolsResponse.pools).toHaveLength(1);
  });

  test("stake more than available balance", async () => {
    // get default pool
    let pool = await getDefaultPool();

    // get stakers list
    let stakersListResponse = await lcdClient.kyve.registry.v1beta1.stakersList(
      { pool_id: "0", status: 1 }
    );

    // get balance before staking
    const preBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();
    const amount = new BigNumber(preBalance).plus(1);

    // check if staking amount is zero
    expect(pool.stakers).toHaveLength(0);
    expect(stakersListResponse.stakers).toHaveLength(0);
    expect(pool.total_stake).toEqual("0");
    expect(pool.lowest_staker).toBe("");

    const request = alice.client.kyve.v1beta1.base
      .stakePool({
        id: "0",
        amount: amount.toString(),
      })
      .then((tx) => tx.execute());
    //insufficient funds
    await expect(request).rejects.toThrow(/insufficient funds/);
    // refetch pool
    pool = await getDefaultPool();
    // refetch stakers list
    stakersListResponse = await lcdClient.kyve.registry.v1beta1.stakersList({
      pool_id: "0",
      status: 1,
    });

    // get balance before staking
    const postBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // check if stake went though
    expect(pool.stakers).toHaveLength(0);
    expect(pool.lowest_staker).toBe("");

    expect(stakersListResponse.stakers).toHaveLength(0);
    expect(pool.total_stake).toEqual("0");

    // check if balance has not changed
    expect(preBalance).toEqual(postBalance);
  });

  test("stake 80 KYVE with alice", async () => {
    // define amount
    const amount = new BigNumber(80).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );

    // get default pool
    let pool = await getDefaultPool();

    // get stakers list
    let stakersListResponse = await lcdClient.kyve.registry.v1beta1.stakersList(
      { pool_id: "0", status: 1 }
    );
    // get balance before staking
    const preBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // check if staking amount is zero
    expect(pool.stakers).toHaveLength(0);
    expect(stakersListResponse.stakers).toHaveLength(0);
    expect(pool.total_stake).toEqual("0");
    expect(pool.lowest_staker).toBe("");
    // stake 80 KYVE
    const tx = await alice.client.kyve.v1beta1.base.stakePool({
      id: pool.id,
      amount: amount.toString(),
    });
    const receipt = await tx.execute();
    // 0 means transaction was successful
    expect(receipt.code).toEqual(0);

    // refetch pool
    pool = await getDefaultPool();

    // refetch stakers list
    stakersListResponse = await lcdClient.kyve.registry.v1beta1.stakersList({
      pool_id: "0",
      status: 1,
    });
    // get balance before staking
    const postBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // check if stake went though
    expect(pool.stakers).toHaveLength(1);
    expect(pool.stakers[0]).toEqual(ADDRESS_ALICE);
    expect(pool.lowest_staker).toBe(ADDRESS_ALICE);

    expect(stakersListResponse.stakers).toHaveLength(1);
    expect(stakersListResponse.stakers[0].amount).toEqual(amount.toString());
    expect(pool.total_stake).toEqual(amount.toString());

    // check if balance was decreased correct
    expect(
      new BigNumber(preBalance)
        .minus(postBalance)
        .minus(tx.fee.amount[0].amount)
    ).toEqual(amount);
  });

  test("stake additional 20 KYVE with alice", async () => {
    // define amount
    const amount = new BigNumber(20).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );
    const total = new BigNumber(100).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );

    // get default pool
    let pool = await getDefaultPool();

    // get stakers list

    // get balance before staking
    const preBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // stake 100 KYVE
    const tx = await alice.client.kyve.v1beta1.base.stakePool({
      id: pool.id,
      amount: amount.toString(),
    });
    const receipt = await tx.execute();
    // 0 means transaction was successful
    expect(receipt.code).toEqual(0);

    // refetch pool
    pool = await getDefaultPool();

    // refetch stakers list
    let stakersListResponse = await lcdClient.kyve.registry.v1beta1.stakersList(
      { pool_id: "0", status: 1 }
    );

    // get balance before staking
    const postBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // check if stake went though
    expect(pool.stakers).toHaveLength(1);
    expect(pool.stakers[0]).toEqual(ADDRESS_ALICE);
    expect(pool.lowest_staker).toBe(ADDRESS_ALICE);

    expect(stakersListResponse.stakers).toHaveLength(1);
    expect(stakersListResponse.stakers[0].amount).toEqual(total.toString());

    expect(pool.total_stake).toEqual(total.toString());

    // check if balance was decreased correct
    expect(
      new BigNumber(preBalance)
        .minus(postBalance)
        .minus(tx.fee.amount[0].amount)
    ).toEqual(amount);
  });

  test("unstake 80 KYVE with alice", async () => {
    // define amount
    const amount = new BigNumber(80).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );
    const remaining = new BigNumber(20).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );

    // get default pool
    let pool = await getDefaultPool();

    // get balance before staking
    const preBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // unstake 80 KYVE
    const tx = await alice.client.kyve.v1beta1.base.unstakePool({
      id: pool.id,
      amount: amount.toString(),
    });
    const receipt = await tx.execute();
    // 0 means transaction was successful
    expect(receipt.code).toEqual(0);

    await sleep(10000);

    // refetch pool
    pool = await getDefaultPool();

    // refetch stakers list
    let stakersListResponse = await lcdClient.kyve.registry.v1beta1.stakersList(
      { pool_id: "0", status: 1 }
    );
    //check unbonding
    const unbondings =
      await lcdClient.kyve.registry.v1beta1.accountStakingUnbonding({
        address: ADDRESS_ALICE,
      });
    // get balance before staking
    const postBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();
    const unbondingsTotal = unbondings.unbondings.reduce(
      (acc, cur) => acc.plus(cur.amount.toString()),
      new BigNumber(0)
    );
    // check if stake went though
    expect(pool.stakers).toHaveLength(1);
    expect(stakersListResponse.stakers).toHaveLength(1);

    expect(pool.stakers[0]).toEqual(ADDRESS_ALICE);
    expect(pool.lowest_staker).toBe(ADDRESS_ALICE);
    expect(unbondingsTotal.toString()).toEqual(amount.toString());
    expect(stakersListResponse.stakers).toHaveLength(1);
    expect(stakersListResponse.stakers[0].unbonding_amount).toEqual(
      amount.toString()
    );
    expect(pool.total_stake).toEqual(amount.plus(remaining).toString());
    expect(postBalance.toString()).toEqual(
      new BigNumber(preBalance).minus(tx.fee.amount[0].amount).toString()
    );
  });

  test("unstake more than staking balance with alice", async () => {
    // define amount
    const amount = new BigNumber(50).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );
    const remaining = new BigNumber(100).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );

    // get default pool
    let pool = await getDefaultPool();

    // get balance before staking
    const preBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // unstake 50 KYVE
    const request = alice.client.kyve.v1beta1.base
      .unstakePool({
        id: pool.id,
        amount: amount.toString(),
      })
      .then((tx) => tx.execute());
    await expect(request).rejects.toThrow(/maximum unstaking amount/);

    // refetch pool
    pool = await getDefaultPool();

    // refetch stakers list
    let stakersListResponse = await lcdClient.kyve.registry.v1beta1.stakersList(
      { pool_id: "0", status: 1 }
    );

    // get balance before staking
    const postBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // check if stake went though
    expect(pool.stakers).toHaveLength(1);
    expect(pool.stakers[0]).toEqual(ADDRESS_ALICE);
    expect(pool.lowest_staker).toBe(ADDRESS_ALICE);

    expect(stakersListResponse.stakers).toHaveLength(1);
    expect(stakersListResponse.stakers[0].amount).toEqual(remaining.toString());

    expect(pool.total_stake).toEqual(remaining.toString());

    // check if balance was decreased correct
    expect(postBalance).toEqual(preBalance);
  });

  test("unstake all with alice", async () => {
    // define amount
    const unstake = new BigNumber(20).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );
    const totalStaked = new BigNumber(100).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );

    // get default pool
    let pool = await getDefaultPool();

    // get balance before staking
    const preBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // unstake 20 KYVE
    const tx = await alice.client.kyve.v1beta1.base.unstakePool({
      id: pool.id,
      amount: unstake.toString(),
    });
    const receipt = await tx.execute();
    // 0 means transaction was successful
    expect(receipt.code).toEqual(0);

    await sleep(10000);

    // refetch pool
    pool = await getDefaultPool();

    let stakersListResponse = await lcdClient.kyve.registry.v1beta1.stakersList(
      { pool_id: "0", status: 1 }
    );

    const postBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();
    //check unbonding
    const unbondings =
      await lcdClient.kyve.registry.v1beta1.accountStakingUnbonding({
        address: ADDRESS_ALICE,
      });
    const unbondingsTotal = unbondings.unbondings.reduce(
      (acc, cur) => acc.plus(cur.amount.toString()),
      new BigNumber(0)
    );

    // check if stake went though
    expect(unbondingsTotal.toString()).toEqual(totalStaked.toString());
    expect(pool.stakers).toHaveLength(1);
    expect(stakersListResponse.stakers).toHaveLength(1);
    expect(stakersListResponse.stakers[0].unbonding_amount).toEqual(
      totalStaked.toString()
    );
    expect(pool.total_stake).toEqual(totalStaked.toString());
    expect(pool.lowest_staker).toBe(ADDRESS_ALICE);
    // check if balance was  correct
    expect(postBalance.toString()).toEqual(
      new BigNumber(preBalance).minus(tx.fee.amount[0].amount).toString()
    );
  });

  test("stake with multiple stakers", async () => {
    // define amounts
    const totalStaked = new BigNumber(100).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );
    const aliceAmount = new BigNumber(200).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );
    const bobAmount = new BigNumber(100).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );
    const charlieAmount = new BigNumber(300).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );

    // get default pool
    let pool = await getDefaultPool();

    // get stakers list
    let stakersListResponse = await lcdClient.kyve.registry.v1beta1.stakersList(
      { pool_id: "0", status: 1 }
    );

    // check if staking amount is zero
    expect(pool.stakers).toHaveLength(1);
    expect(stakersListResponse.stakers).toHaveLength(1);
    expect(pool.total_stake).toEqual(totalStaked.toString());
    await bob.client.kyve.v1beta1.base
      .stakePool({
        id: pool.id,
        amount: bobAmount.toString(),
      })
      .then((tx) => tx.execute());
    await alice.client.kyve.v1beta1.base
      .stakePool({
        id: pool.id,
        amount: aliceAmount.toString(),
      })
      .then((tx) => tx.execute());
    await charlie.client.kyve.v1beta1.base
      .stakePool({
        id: pool.id,
        amount: charlieAmount.toString(),
      })
      .then((tx) => tx.execute());
    // refetch pool
    pool = await getDefaultPool();

    // refetch stakers list
    stakersListResponse = await lcdClient.kyve.registry.v1beta1.stakersList({
      pool_id: "0",
      status: 1,
    });

    // check if stake went though
    expect(pool.stakers).toHaveLength(3);
    expect(pool.stakers).toContain(ADDRESS_ALICE);
    expect(pool.stakers).toContain(ADDRESS_BOB);
    expect(pool.stakers).toContain(ADDRESS_CHARLIE);
    expect(pool.lowest_staker).toBe(ADDRESS_BOB);

    expect(stakersListResponse.stakers).toHaveLength(3);

    expect(pool.total_stake).toEqual(
      totalStaked
        .plus(aliceAmount)
        .plus(bobAmount)
        .plus(charlieAmount)
        .toString()
    );
  });

  test("unstake with multiple stakers", async () => {
    // define amounts
    const totalStaked = new BigNumber(700).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );
    const totalUnbonding = new BigNumber(100).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );
    const aliceAmount = new BigNumber(200).multipliedBy(10 ** 9);
    const bobAmount = new BigNumber(100).multipliedBy(10 ** 9);
    const charlieAmount = new BigNumber(300).multipliedBy(10 ** 9);

    // get default pool
    let pool = await getDefaultPool();

    // get stakers list
    let stakersListResponse = await lcdClient.kyve.registry.v1beta1.stakersList(
      { pool_id: "0", status: 1 }
    );

    // check if staking amount is zero
    expect(pool.stakers).toHaveLength(3);
    expect(stakersListResponse.stakers).toHaveLength(3);
    expect(pool.total_stake).toEqual(totalStaked.toString());

    await alice.client.kyve.v1beta1.base
      .unstakePool({
        id: pool.id,
        amount: aliceAmount.toString(),
      })
      .then((tx) => tx.execute());
    await bob.client.kyve.v1beta1.base
      .unstakePool({
        id: pool.id,
        amount: bobAmount.toString(),
      })
      .then((tx) => tx.execute());
    await charlie.client.kyve.v1beta1.base
      .unstakePool({
        id: pool.id,
        amount: charlieAmount.toString(),
      })
      .then((tx) => tx.execute());

    await sleep(10000);

    // refetch pool
    pool = await getDefaultPool();

    // refetch stakers list
    stakersListResponse = await lcdClient.kyve.registry.v1beta1.stakersList({
      pool_id: "0",
      status: 1,
    });
    const unbondingsAlice =
      await lcdClient.kyve.registry.v1beta1.accountStakingUnbonding({
        address: ADDRESS_ALICE,
      });
    const unbondingsBob =
      await lcdClient.kyve.registry.v1beta1.accountStakingUnbonding({
        address: ADDRESS_BOB,
      });
    const unbondingsCharlie =
      await lcdClient.kyve.registry.v1beta1.accountStakingUnbonding({
        address: ADDRESS_CHARLIE,
      });
    const unbondingsTotalAlice = unbondingsAlice.unbondings.reduce(
      (acc, cur) => acc.plus(cur.amount.toString()),
      new BigNumber(0)
    );
    const unbondingsTotalBob = unbondingsBob.unbondings.reduce(
      (acc, cur) => acc.plus(cur.amount.toString()),
      new BigNumber(0)
    );
    const unbondingsTotalCharlie = unbondingsCharlie.unbondings.reduce(
      (acc, cur) => acc.plus(cur.amount.toString()),
      new BigNumber(0)
    );
    expect(stakersListResponse.stakers).toHaveLength(3);

    // check if stake went though
    expect(pool.stakers).toHaveLength(3);
    expect(pool.total_stake).toEqual(totalStaked.toString());
    expect(unbondingsTotalAlice.toString()).toEqual(
      totalUnbonding.plus(aliceAmount).toString()
    );
    expect(unbondingsTotalBob.toString()).toEqual(bobAmount.toString());
    expect(unbondingsTotalCharlie.toString()).toEqual(charlieAmount.toString());
  });
};
