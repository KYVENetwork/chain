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
} from "./helpers/utils";
import { BASE_URL } from "./helpers/constants";

export const staking = () => {
  // disable timeout
  jest.setTimeout(24 * 60 * 60 * 1000);

  test("Test if default pool exist", async () => {
    const { data } = await axios.get(`${BASE_URL}/kyve/registry/v1beta1/pools`);

    expect(data.pools).not.toBeUndefined();
    expect(data.pools).toHaveLength(1);
  });

  test("stake more than available balance", async () => {
    // get default pool
    let pool = await getPoolById(0);

    // get stakers list
    let stakers_list = await getStakersList(0);

    // get balance before staking
    const preBalance = await getBalanceByAddress(ADDRESS_ALICE);
    const amount = preBalance.plus(1);

    // check if staking amount is zero
    expect(pool.stakers).toHaveLength(0);
    expect(stakers_list).toHaveLength(0);
    expect(pool.total_stake).toEqual("0");
    expect(pool.lowest_staker).toBe("");

    // stake 100 KYVE
    const { transactionBroadcast } = await alice.stake(pool.id, amount);
    const receipt = await transactionBroadcast;

    // 0 means transaction was successful
    expect(receipt.code).not.toEqual(0);

    // refetch pool
    pool = await getPoolById(0);

    // refetch stakers list
    stakers_list = await getStakersList(0);

    // get balance before staking
    const postBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // check if stake went though
    expect(pool.stakers).toHaveLength(0);
    expect(pool.lowest_staker).toBe("");

    expect(stakers_list).toHaveLength(0);
    expect(pool.total_stake).toEqual("0");

    // check if balance has not changed
    expect(preBalance).toEqual(postBalance);
  });

  test("stake 80 KYVE with alice", async () => {
    // define amount
    const amount = new BigNumber(80).multipliedBy(10 ** 9);

    // get default pool
    let pool = await getPoolById(0);

    // get stakers list
    let stakers_list = await getStakersList(0);

    // get balance before staking
    const preBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // check if staking amount is zero
    expect(pool.stakers).toHaveLength(0);
    expect(stakers_list).toHaveLength(0);
    expect(pool.total_stake).toEqual("0");
    expect(pool.lowest_staker).toBe("");

    // stake 80 KYVE
    const { transactionBroadcast } = await alice.stake(pool.id, amount);
    const receipt = await transactionBroadcast;

    // 0 means transaction was successful
    expect(receipt.code).toEqual(0);

    // refetch pool
    pool = await getPoolById(0);

    // refetch stakers list
    stakers_list = await getStakersList(0);

    // get balance before staking
    const postBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // check if stake went though
    expect(pool.stakers).toHaveLength(1);
    expect(pool.stakers[0]).toEqual(ADDRESS_ALICE);
    expect(pool.lowest_staker).toBe(ADDRESS_ALICE);

    expect(stakers_list).toHaveLength(1);
    expect(stakers_list[0].amount).toEqual(amount.toString());

    expect(pool.total_stake).toEqual(amount.toString());

    // check if balance was decreased correct
    expect(preBalance.minus(postBalance)).toEqual(amount);
  });

  test("stake additional 20 KYVE with alice", async () => {
    // define amount
    const amount = new BigNumber(20).multipliedBy(10 ** 9);
    const total = new BigNumber(100).multipliedBy(10 ** 9);

    // get default pool
    let pool = await getPoolById(0);

    // get stakers list
    let stakers_list = await getStakersList(0);

    // get balance before staking
    const preBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // stake 100 KYVE
    const { transactionBroadcast } = await alice.stake(pool.id, amount);
    const receipt = await transactionBroadcast;

    // 0 means transaction was successful
    expect(receipt.code).toEqual(0);

    // refetch pool
    pool = await getPoolById(0);

    // refetch stakers list
    stakers_list = await getStakersList(0);

    // get balance before staking
    const postBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // check if stake went though
    expect(pool.stakers).toHaveLength(1);
    expect(pool.stakers[0]).toEqual(ADDRESS_ALICE);
    expect(pool.lowest_staker).toBe(ADDRESS_ALICE);

    expect(stakers_list).toHaveLength(1);
    expect(stakers_list[0].amount).toEqual(total.toString());

    expect(pool.total_stake).toEqual(total.toString());

    // check if balance was decreased correct
    expect(preBalance.minus(postBalance)).toEqual(amount);
  });

  test("unstake 80 KYVE with alice", async () => {
    // define amount
    const amount = new BigNumber(80).multipliedBy(10 ** 9);
    const remaining = new BigNumber(20).multipliedBy(10 ** 9);

    // get default pool
    let pool = await getPoolById(0);

    // get stakers list
    let stakers_list: any[];

    // get balance before staking
    const preBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // unstake 80 KYVE
    const { transactionBroadcast } = await alice.unstake(pool.id, amount);
    const receipt = await transactionBroadcast;

    // 0 means transaction was successful
    expect(receipt.code).toEqual(0);

    await sleep(10000);

    // refetch pool
    pool = await getPoolById(0);

    // refetch stakers list
    stakers_list = await getStakersList(0);

    // get balance before staking
    const postBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // check if stake went though
    expect(pool.stakers).toHaveLength(1);
    expect(stakers_list).toHaveLength(1);

    expect(pool.stakers[0]).toEqual(ADDRESS_ALICE);
    expect(pool.lowest_staker).toBe(ADDRESS_ALICE);

    expect(stakers_list).toHaveLength(1);
    expect(stakers_list[0].amount).toEqual(remaining.toString());

    expect(pool.total_stake).toEqual(remaining.toString());

    // check if balance was decreased correct
    expect(postBalance.minus(preBalance)).toEqual(amount);
  });

  test("unstake more than staking balance with alice", async () => {
    // define amount
    const amount = new BigNumber(50).multipliedBy(10 ** 9);
    const remaining = new BigNumber(20).multipliedBy(10 ** 9);

    // get default pool
    let pool = await getPoolById(0);

    // get balance before staking
    const preBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // unstake 50 KYVE
    const { transactionBroadcast } = await alice.unstake(pool.id, amount);
    const receipt = await transactionBroadcast;

    // 0 means transaction was successful
    expect(receipt.code).not.toEqual(0);

    // refetch pool
    pool = await getPoolById(0);

    // refetch stakers list
    let stakers_list = await getStakersList(0);

    // get balance before staking
    const postBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // check if stake went though
    expect(pool.stakers).toHaveLength(1);
    expect(pool.stakers[0]).toEqual(ADDRESS_ALICE);
    expect(pool.lowest_staker).toBe(ADDRESS_ALICE);

    expect(stakers_list).toHaveLength(1);
    expect(stakers_list[0].amount).toEqual(remaining.toString());

    expect(pool.total_stake).toEqual(remaining.toString());

    // check if balance was decreased correct
    expect(postBalance).toEqual(preBalance);
  });

  test("unstake all with alice", async () => {
    // define amount
    const amount = new BigNumber(20).multipliedBy(10 ** 9);

    // get default pool
    let pool = await getPoolById(0);

    // get balance before staking
    const preBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // unstake 20 KYVE
    const { transactionBroadcast } = await alice.unstake(pool.id, amount);
    const receipt = await transactionBroadcast;

    // 0 means transaction was successful
    expect(receipt.code).toEqual(0);

    await sleep(10000);

    // refetch pool
    pool = await getPoolById(0);

    // refetch stakers list
    let stakers_list = await getStakersList(0);

    // get balance before staking
    const postBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // check if stake went though
    expect(pool.stakers).toHaveLength(0);
    expect(stakers_list).toHaveLength(0);
    expect(pool.total_stake).toEqual("0");
    expect(pool.lowest_staker).toBe("");

    // check if balance was decreased correct
    expect(postBalance.minus(preBalance)).toEqual(amount);
  });

  test("stake with multiple stakers", async () => {
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

  test("unstake with multiple stakers", async () => {
    // define amounts
    const aliceAmount = new BigNumber(200).multipliedBy(10 ** 9);
    const bobAmount = new BigNumber(100).multipliedBy(10 ** 9);
    const charlieAmount = new BigNumber(300).multipliedBy(10 ** 9);

    // get default pool
    let pool = await getPoolById(0);

    // get stakers list
    let stakers_list = await getStakersList(0);

    // check if staking amount is zero
    expect(pool.stakers).toHaveLength(3);
    expect(stakers_list).toHaveLength(3);
    expect(pool.total_stake).toEqual("600000000000");

    // stake alice
    const { transactionBroadcast: aliceTx } = await alice.unstake(
      pool.id,
      aliceAmount
    );
    await aliceTx;

    // stake alice
    const { transactionBroadcast: bobTx } = await bob.unstake(
      pool.id,
      bobAmount
    );
    await bobTx;

    // stake alice
    const { transactionBroadcast: charlieTx } = await charlie.unstake(
      pool.id,
      charlieAmount
    );
    await charlieTx;

    await sleep(10000);

    // refetch pool
    pool = await getPoolById(0);

    // refetch stakers list
    stakers_list = await getStakersList(0);

    expect(stakers_list).toHaveLength(0);

    // check if stake went though
    expect(pool.stakers).toHaveLength(0);
    expect(pool.lowest_staker).toEqual("");
    expect(pool.total_stake).toEqual("0");
  });
};
