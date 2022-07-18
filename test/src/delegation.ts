import axios from "axios";
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
const getDefaultPool = async () =>
  (await lcdClient.kyve.registry.v1beta1.pool({ id: "0" })).pool ??
  (() => {
    throw new Error("Pool doesn't exist");
  })();

export const delegation = () => {
  afterAll(async () => {
    await restartChain();
    await alice.init();
    await bob.init();
    await charlie.init();
  });
  // disable timeout
  jest.setTimeout(24 * 60 * 60 * 1000);

  test("Test if default pool exist", async () => {
    const data = await lcdClient.kyve.registry.v1beta1.pools({ paused: false });
    expect(data.pools).not.toBeUndefined();
    expect(data.pools).toHaveLength(1);
  });

  test("Create three stakers", async () => {
    // define amounts
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
    expect(pool.stakers).toHaveLength(0);
    expect(stakersListResponse.stakers).toHaveLength(0);
    expect(pool.total_stake).toEqual("0");

    await alice.client.kyve.v1beta1.base
      .stakePool({
        id: pool.id,
        amount: aliceAmount.toString(),
      })
      .then((tx) => tx.execute());
    await bob.client.kyve.v1beta1.base
      .stakePool({
        id: pool.id,
        amount: bobAmount.toString(),
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
      aliceAmount.plus(bobAmount).plus(charlieAmount).toString()
    );
  });

  test("No Self delegation", async () => {
    // get default pool
    let pool = await getDefaultPool();

    // get stakers list
    let stakersListResponse = await lcdClient.kyve.registry.v1beta1.stakersList(
      { pool_id: "0", status: 1 }
    );

    const amount = new BigNumber(50).multipliedBy(10 ** 9);

    const preBalanceAlice =
      await alice.client.kyve.v1beta1.base.getKyveBalance();
    const preBalanceBob = await bob.client.kyve.v1beta1.base.getKyveBalance();
    const preBalanceCharlie =
      await charlie.client.kyve.v1beta1.base.getKyveBalance();
    const selfDelegationErr = /self delegation not allowed/;
    const aliceDelegateReq = alice.client.kyve.v1beta1.base
      .delegatePool({
        id: pool.id,
        amount: amount.toString(),
        staker: ADDRESS_ALICE,
      })
      .then((tx) => tx.execute());
    await expect(aliceDelegateReq).rejects.toThrow(selfDelegationErr);
    const bobDelegateReq = bob.client.kyve.v1beta1.base
      .delegatePool({
        id: pool.id,
        amount: amount.toString(),
        staker: ADDRESS_BOB,
      })
      .then((tx) => tx.execute());
    await expect(bobDelegateReq).rejects.toThrow(selfDelegationErr);
    const charlieDelegateReq = charlie.client.kyve.v1beta1.base
      .delegatePool({
        id: pool.id,
        amount: amount.toString(),
        staker: ADDRESS_CHARLIE,
      })
      .then((tx) => tx.execute());
    await expect(charlieDelegateReq).rejects.toThrow(selfDelegationErr);
    const postBalanceAlice =
      await alice.client.kyve.v1beta1.base.getKyveBalance();
    const postBalanceBob = await bob.client.kyve.v1beta1.base.getKyveBalance();
    const postBalanceCharlie =
      await charlie.client.kyve.v1beta1.base.getKyveBalance();

    expect(preBalanceAlice).toEqual(postBalanceAlice);
    expect(preBalanceBob).toEqual(postBalanceBob);
    expect(preBalanceCharlie).toEqual(postBalanceCharlie);

    // refetch pool
    pool = await getDefaultPool();

    expect(pool.total_delegation).toEqual("0");
  });

  test("Delegate into Alice", async () => {
    // Delegation amounts to delegate into alice
    const bobDelegation = new BigNumber(100).multipliedBy(10 ** 9);
    const charlieDelegation = new BigNumber(300).multipliedBy(10 ** 9);

    // PreBalances
    const preBalanceAlice =
      await alice.client.kyve.v1beta1.base.getKyveBalance();
    const preBalanceBob = await bob.client.kyve.v1beta1.base.getKyveBalance();
    const preBalanceCharlie =
      await charlie.client.kyve.v1beta1.base.getKyveBalance();

    const txBob = await bob.client.kyve.v1beta1.base.delegatePool({
      id: "0",
      amount: bobDelegation.toString(),
      staker: ADDRESS_ALICE,
    });
    const receiptBob = await txBob.execute();
    expect(receiptBob.code).toEqual(0);

    const txCharlie = await charlie.client.kyve.v1beta1.base.delegatePool({
      id: "0",
      amount: charlieDelegation.toString(),
      staker: ADDRESS_ALICE,
    });
    const receiptCharlie = await txCharlie.execute();
    expect(receiptCharlie.code).toEqual(0);

    // Check post balances
    const postBalanceAlice =
      await alice.client.kyve.v1beta1.base.getKyveBalance();
    const postBalanceBob = await bob.client.kyve.v1beta1.base.getKyveBalance();
    const postBalanceCharlie =
      await charlie.client.kyve.v1beta1.base.getKyveBalance();

    expect(preBalanceAlice).toEqual(postBalanceAlice);
    expect(
      new BigNumber(preBalanceBob)
        .minus(bobDelegation)
        .minus(txBob.fee.amount[0].amount)
        .toString()
    ).toEqual(postBalanceBob);
    expect(
      new BigNumber(preBalanceCharlie)
        .minus(charlieDelegation)
        .minus(txCharlie.fee.amount[0].amount)
        .toString()
    ).toEqual(postBalanceCharlie);

    // Check pool delegation entry
    const delegationPoolData =
      await lcdClient.kyve.registry.v1beta1.delegatorsByPoolAndStaker({
        staker: ADDRESS_ALICE,
        pool_id: "0",
      });
    expect(delegationPoolData.delegation_pool_data?.current_rewards).toEqual(
      "0"
    );
    expect(delegationPoolData.pool?.total_delegation).toEqual(
      bobDelegation.plus(charlieDelegation).toString()
    );
    expect(delegationPoolData.delegation_pool_data?.delegator_count).toEqual(
      "2"
    );
  });

  test("Undelegate everything from Alice", async () => {
    // Delegation amounts to delegate into alice
    const bobDelegation = new BigNumber(100).multipliedBy(10 ** 9);
    const charlieDelegation = new BigNumber(300).multipliedBy(10 ** 9);

    // PreBalances
    const preBalanceAlice =
      await alice.client.kyve.v1beta1.base.getKyveBalance();
    const preBalanceBob = await bob.client.kyve.v1beta1.base.getKyveBalance();
    const preBalanceCharlie =
      await charlie.client.kyve.v1beta1.base.getKyveBalance();

    // undelegate 100 $KYVE from Bob into Alice
    const txBob = await bob.client.kyve.v1beta1.base.undelegatePool({
      id: "0",
      staker: ADDRESS_ALICE,
      amount: bobDelegation.toString(),
    });
    const receiptBob = await txBob.execute();
    expect(receiptBob.code).toEqual(0);

    // undelegate 300 $KYVE from Charlie into Alice
    const txCharlie = await charlie.client.kyve.v1beta1.base.undelegatePool({
      id: "0",
      staker: ADDRESS_ALICE,
      amount: charlieDelegation.toString(),
    });
    const receiptCharlie = await txCharlie.execute();
    expect(receiptCharlie.code).toEqual(0);

    await sleep(5 * 1000);

    // Check post balances
    const postBalanceAlice =
      await alice.client.kyve.v1beta1.base.getKyveBalance();
    const postBalanceBob = await bob.client.kyve.v1beta1.base.getKyveBalance();
    const postBalanceCharlie =
      await charlie.client.kyve.v1beta1.base.getKyveBalance();
    const unbondingsAlice =
      await lcdClient.kyve.registry.v1beta1.accountDelegationUnbondings({
        address: ADDRESS_ALICE,
      });
    const unbondingsBob =
      await lcdClient.kyve.registry.v1beta1.accountDelegationUnbondings({
        address: ADDRESS_BOB,
      });
    const unbondingsCharlie =
      await lcdClient.kyve.registry.v1beta1.accountDelegationUnbondings({
        address: ADDRESS_CHARLIE,
      });
    const unbondingsTotalBob = unbondingsBob.unbondings.reduce(
      (acc, cur) => acc.plus(cur.amount.toString()),
      new BigNumber(0)
    );
    const unbondingsTotalCharlie = unbondingsCharlie.unbondings.reduce(
      (acc, cur) => acc.plus(cur.amount.toString()),
      new BigNumber(0)
    );
    expect(postBalanceAlice).toEqual(preBalanceAlice);
    expect(postBalanceCharlie).toEqual(
      new BigNumber(preBalanceCharlie)
        .minus(txCharlie.fee.amount[0].amount)
        .toString()
    );
    expect(postBalanceBob).toEqual(
      new BigNumber(preBalanceBob).minus(txBob.fee.amount[0].amount).toString()
    );
    expect(unbondingsTotalBob.toString()).toEqual(bobDelegation.toString());
    expect(unbondingsTotalCharlie.toString()).toEqual(
      charlieDelegation.toString()
    );

    // Check pool delegation entry
    const delegationPoolData =
      await lcdClient.kyve.registry.v1beta1.delegatorsByPoolAndStaker({
        staker: ADDRESS_ALICE,
        pool_id: "0",
      });
    expect(delegationPoolData.delegation_pool_data).toEqual({
      current_rewards: "0",
      delegator_count: "0",
      id: "0",
      latest_index_k: "0",
      latest_index_was_undelegation: false,
      staker: "",
      total_delegation: "0",
    });
  });
};
