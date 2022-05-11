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
import {
  getBalanceByAddress,
  getStakersList,
  getPoolById,
  sleep,
  getDelegationPoolData,
} from "./helpers/utils";
import { BASE_URL } from "./helpers/constants";

export const delegation = () => {
  // disable timeout
  jest.setTimeout(24 * 60 * 60 * 1000);

  test("Test if default pool exist", async () => {
    const { data } = await axios.get(`${BASE_URL}/kyve/registry/v1beta1/pools`);

    expect(data.pools).not.toBeUndefined();
    expect(data.pools).toHaveLength(1);
  });

  test("Create three stakers", async () => {
    // define amounts
    const aliceAmount = new BigNumber(200).multipliedBy(10 ** 9);
    const bobAmount = new BigNumber(100).multipliedBy(10 ** 9);
    const charlieAmount = new BigNumber(300).multipliedBy(10 ** 9);

    // get default pool
    let pool = await getPoolById(0);

    // get stakers list
    let stakers_list = await getStakersList(0);

    // check if staking amount is zero
    expect(pool.stakers).toHaveLength(0);
    expect(stakers_list).toHaveLength(0);
    expect(pool.total_stake).toEqual("0");

    // stake alice
    const { transactionBroadcast: aliceTx } = await alice.stake(
      pool.id,
      aliceAmount
    );
    await aliceTx;

    // stake alice
    const { transactionBroadcast: bobTx } = await bob.stake(pool.id, bobAmount);
    await bobTx;

    // stake alice
    const { transactionBroadcast: charlieTx } = await charlie.stake(
      pool.id,
      charlieAmount
    );
    await charlieTx;

    // refetch pool
    pool = await getPoolById(0);

    // refetch stakers list
    stakers_list = await getStakersList(0);

    // check if stake went though
    expect(pool.stakers).toHaveLength(3);
    expect(pool.stakers).toContain(ADDRESS_ALICE);
    expect(pool.stakers).toContain(ADDRESS_BOB);
    expect(pool.stakers).toContain(ADDRESS_CHARLIE);
    expect(pool.lowest_staker).toBe(ADDRESS_BOB);

    expect(stakers_list).toHaveLength(3);

    expect(pool.total_stake).toEqual(
      aliceAmount.plus(bobAmount).plus(charlieAmount).toString()
    );
  });

  test("No Self delegation", async () => {
    // get default pool
    let pool = await getPoolById(0);

    // get stakers list
    let stakers_list = await getStakersList(0);

    const amount = new BigNumber(50).multipliedBy(10 ** 9);

    const preBalanceAlice = await getBalanceByAddress(ADDRESS_ALICE);
    const preBalanceBob = await getBalanceByAddress(ADDRESS_BOB);
    const preBalanceCharlie = await getBalanceByAddress(ADDRESS_CHARLIE);

    const receiptAlice = await (
      await alice.delegate(0, ADDRESS_ALICE, amount)
    ).transactionBroadcast;
    const receiptBob = await (
      await bob.delegate(0, ADDRESS_BOB, amount)
    ).transactionBroadcast;
    const receiptCharlie = await (
      await charlie.delegate(0, ADDRESS_CHARLIE, amount)
    ).transactionBroadcast;

    expect(receiptAlice.code).not.toEqual(0);
    expect(receiptBob.code).not.toEqual(0);
    expect(receiptCharlie.code).not.toEqual(0);

    const postBalanceAlice = await getBalanceByAddress(ADDRESS_ALICE);
    const postBalanceBob = await getBalanceByAddress(ADDRESS_BOB);
    const postBalanceCharlie = await getBalanceByAddress(ADDRESS_CHARLIE);

    expect(preBalanceAlice).toEqual(postBalanceAlice);
    expect(preBalanceBob).toEqual(postBalanceBob);
    expect(preBalanceCharlie).toEqual(postBalanceCharlie);

    // refetch pool
    pool = await getPoolById(0);

    expect(pool.total_delegation).toEqual("0");
  });

  test("Delegate into Alice", async () => {
    // Delegation amounts to delegate into alice
    const bobDelegation = new BigNumber(100).multipliedBy(10 ** 9);
    const charlieDelegation = new BigNumber(300).multipliedBy(10 ** 9);

    // PreBalances
    const preBalanceAlice = await getBalanceByAddress(ADDRESS_ALICE);
    const preBalanceBob = await getBalanceByAddress(ADDRESS_BOB);
    const preBalanceCharlie = await getBalanceByAddress(ADDRESS_CHARLIE);

    // Delegate 100 $KYVE from Bob into Alice
    const receiptBob = await (
      await bob.delegate(0, ADDRESS_ALICE, bobDelegation)
    ).transactionBroadcast;
    expect(receiptBob.code).toEqual(0);

    // Delegate 300 $KYVE from Charlie into Alice
    const receiptCharlie = await (
      await charlie.delegate(0, ADDRESS_ALICE, charlieDelegation)
    ).transactionBroadcast;
    expect(receiptCharlie.code).toEqual(0);

    // Check post balances
    const postBalanceAlice = await getBalanceByAddress(ADDRESS_ALICE);
    const postBalanceBob = await getBalanceByAddress(ADDRESS_BOB);
    const postBalanceCharlie = await getBalanceByAddress(ADDRESS_CHARLIE);

    expect(preBalanceAlice).toEqual(postBalanceAlice);
    expect(preBalanceBob.minus(bobDelegation)).toEqual(postBalanceBob);
    expect(preBalanceCharlie.minus(charlieDelegation)).toEqual(
      postBalanceCharlie
    );

    // Check pool delegation entry
    const delegationPoolData = await getDelegationPoolData(ADDRESS_ALICE, 0);
    expect(delegationPoolData.current_rewards).toEqual("0");
    expect(delegationPoolData.total_delegation).toEqual(
      bobDelegation.plus(charlieDelegation).toString()
    );
    expect(delegationPoolData.delegator_count).toEqual("2");
  });

  test("Undelegate everything from Alice", async () => {
    // Delegation amounts to delegate into alice
    const bobDelegation = new BigNumber(100).multipliedBy(10 ** 9);
    const charlieDelegation = new BigNumber(300).multipliedBy(10 ** 9);

    // PreBalances
    const preBalanceAlice = await getBalanceByAddress(ADDRESS_ALICE);
    const preBalanceBob = await getBalanceByAddress(ADDRESS_BOB);
    const preBalanceCharlie = await getBalanceByAddress(ADDRESS_CHARLIE);

    // Delegate 100 $KYVE from Bob into Alice
    const receiptBob = await (
      await bob.undelegate(0, ADDRESS_ALICE, bobDelegation)
    ).transactionBroadcast;
    expect(receiptBob.code).toEqual(0);

    // Delegate 300 $KYVE from Charlie into Alice
    const receiptCharlie = await (
      await charlie.undelegate(0, ADDRESS_ALICE, charlieDelegation)
    ).transactionBroadcast;
    expect(receiptCharlie.code).toEqual(0);

    await sleep(5 * 1000);

    // Check post balances
    const postBalanceAlice = await getBalanceByAddress(ADDRESS_ALICE);
    const postBalanceBob = await getBalanceByAddress(ADDRESS_BOB);
    const postBalanceCharlie = await getBalanceByAddress(ADDRESS_CHARLIE);

    expect(preBalanceAlice).toEqual(postBalanceAlice);
    expect(preBalanceBob.plus(bobDelegation)).toEqual(postBalanceBob);
    expect(preBalanceCharlie.plus(charlieDelegation)).toEqual(
      postBalanceCharlie
    );

    // Check pool delegation entry
    const delegationPoolData = await getDelegationPoolData(ADDRESS_ALICE, 0);
    expect(delegationPoolData).toEqual({
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
