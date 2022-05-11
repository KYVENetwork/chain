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
  getFundersList,
  getPoolById,
} from "./helpers/utils";
import { BASE_URL, MAX_FUNDERS } from "./helpers/constants";

export const funding = () => {
  // disable timeout
  jest.setTimeout(24 * 60 * 60 * 1000);

  test("Test if default pool exist", async () => {
    const { data } = await axios.get(`${BASE_URL}/kyve/registry/v1beta1/pools`);

    expect(data.pools).not.toBeUndefined();
    expect(data.pools).toHaveLength(1);
  });

  test("fund more than available balance", async () => {
    // get default pool
    let pool = await getPoolById(0);

    // get funders list
    let funders_list = await getFundersList(0);

    // get balance before funding
    const preBalance = await getBalanceByAddress(ADDRESS_ALICE);
    const amount = preBalance.plus(1);

    // check if funding amount is zero
    expect(pool.funders).toHaveLength(0);
    expect(funders_list).toHaveLength(0);
    expect(pool.total_funds).toEqual("0");
    expect(pool.lowest_funder).toBe("");

    // fund 100 KYVE
    const { transactionBroadcast } = await alice.fund(pool.id, amount);
    const receipt = await transactionBroadcast;

    // 0 means transaction was successful
    expect(receipt.code).not.toEqual(0);

    // refetch pool
    pool = await getPoolById(0);

    // refetch funders list
    funders_list = await getFundersList(0);

    // get balance before funding
    const postBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // check if funds went though
    expect(pool.funders).toHaveLength(0);
    expect(pool.lowest_funder).toBe("");

    expect(funders_list).toHaveLength(0);
    expect(pool.total_funds).toEqual("0");

    // check if balance has not changed
    expect(preBalance).toEqual(postBalance);
  });

  test("fund 80 KYVE with alice", async () => {
    // define amount
    const amount = new BigNumber(80).multipliedBy(10 ** 9);

    // get default pool
    let pool = await getPoolById(0);

    // get funders list
    let funders_list = await getFundersList(0);

    // get balance before funding
    const preBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // check if funding amount is zero
    expect(pool.funders).toHaveLength(0);
    expect(funders_list).toHaveLength(0);
    expect(pool.total_funds).toEqual("0");
    expect(pool.lowest_funder).toBe("");

    // fund 80 KYVE
    const { transactionBroadcast } = await alice.fund(pool.id, amount);
    const receipt = await transactionBroadcast;

    // 0 means transaction was successful
    expect(receipt.code).toEqual(0);

    // refetch pool
    pool = await getPoolById(0);

    // refetch funders list
    funders_list = await getFundersList(0);

    // get balance before funding
    const postBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // check if funds went though
    expect(pool.funders).toHaveLength(1);
    expect(pool.funders[0]).toEqual(ADDRESS_ALICE);
    expect(pool.lowest_funder).toBe(ADDRESS_ALICE);

    expect(funders_list).toHaveLength(1);
    expect(funders_list[0]).toEqual({
      account: ADDRESS_ALICE,
      pool_id: "0",
      amount: amount.toString(),
    });

    expect(pool.total_funds).toEqual(amount.toString());

    // check if balance was decreased correct
    expect(preBalance.minus(postBalance)).toEqual(amount);
  });

  test("fund additional 20 KYVE with alice", async () => {
    // define amount
    const amount = new BigNumber(20).multipliedBy(10 ** 9);
    const total = new BigNumber(100).multipliedBy(10 ** 9);

    // get default pool
    let pool = await getPoolById(0);

    // get funders list
    let funders_list = await getFundersList(0);

    // get balance before funding
    const preBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // fund 100 KYVE
    const { transactionBroadcast } = await alice.fund(pool.id, amount);
    const receipt = await transactionBroadcast;

    // 0 means transaction was successful
    expect(receipt.code).toEqual(0);

    // refetch pool
    pool = await getPoolById(0);

    // refetch funders list
    funders_list = await getFundersList(0);

    // get balance before funding
    const postBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // check if funds went though
    expect(pool.funders).toHaveLength(1);
    expect(pool.funders[0]).toEqual(ADDRESS_ALICE);
    expect(pool.lowest_funder).toBe(ADDRESS_ALICE);

    expect(funders_list).toHaveLength(1);
    expect(funders_list[0]).toEqual({
      account: ADDRESS_ALICE,
      pool_id: "0",
      amount: total.toString(),
    });

    expect(pool.total_funds).toEqual(total.toString());

    // check if balance was decreased correct
    expect(preBalance.minus(postBalance)).toEqual(amount);
  });

  test("defund 80 KYVE with alice", async () => {
    // define amount
    const amount = new BigNumber(80).multipliedBy(10 ** 9);
    const remaining = new BigNumber(20).multipliedBy(10 ** 9);

    // get default pool
    let pool = await getPoolById(0);

    // get funders list
    let funders_list = await getFundersList(0);

    // get balance before funding
    const preBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // defund 80 KYVE
    const { transactionBroadcast } = await alice.defund(pool.id, amount);
    const receipt = await transactionBroadcast;

    // 0 means transaction was successful
    expect(receipt.code).toEqual(0);

    // refetch pool
    pool = await getPoolById(0);

    // refetch funders list
    funders_list = await getFundersList(0);

    // get balance before funding
    const postBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // check if funds went though
    expect(pool.funders).toHaveLength(1);
    expect(pool.funders[0]).toEqual(ADDRESS_ALICE);
    expect(pool.lowest_funder).toBe(ADDRESS_ALICE);

    expect(funders_list).toHaveLength(1);
    expect(funders_list[0]).toEqual({
      account: ADDRESS_ALICE,
      pool_id: "0",
      amount: remaining.toString(),
    });

    expect(pool.total_funds).toEqual(remaining.toString());

    // check if balance was decreased correct
    expect(postBalance.minus(preBalance)).toEqual(amount);
  });

  test("defund more than funding balance with alice", async () => {
    // define amount
    const amount = new BigNumber(50).multipliedBy(10 ** 9);
    const remaining = new BigNumber(20).multipliedBy(10 ** 9);

    // get default pool
    let pool = await getPoolById(0);

    // get funders list
    let funders_list = await getFundersList(0);

    // get balance before funding
    const preBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // defund 50 KYVE
    const { transactionBroadcast } = await alice.defund(pool.id, amount);
    const receipt = await transactionBroadcast;

    // 0 means transaction was successful
    expect(receipt.code).not.toEqual(0);

    // refetch pool
    pool = await getPoolById(0);

    // refetch funders list
    funders_list = await getFundersList(0);

    // get balance before funding
    const postBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // check if funds went though
    expect(pool.funders).toHaveLength(1);
    expect(pool.funders[0]).toEqual(ADDRESS_ALICE);
    expect(pool.lowest_funder).toBe(ADDRESS_ALICE);

    expect(funders_list).toHaveLength(1);
    expect(funders_list[0]).toEqual({
      account: ADDRESS_ALICE,
      pool_id: "0",
      amount: remaining.toString(),
    });

    expect(pool.total_funds).toEqual(remaining.toString());

    // check if balance was decreased correct
    expect(postBalance).toEqual(preBalance);
  });

  test("defund all with alice", async () => {
    // define amount
    const amount = new BigNumber(20).multipliedBy(10 ** 9);

    // get default pool
    let pool = await getPoolById(0);

    // get funders list
    let funders_list = await getFundersList(0);

    // get balance before funding
    const preBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // defund 20 KYVE
    const { transactionBroadcast } = await alice.defund(pool.id, amount);
    const receipt = await transactionBroadcast;

    // 0 means transaction was successful
    expect(receipt.code).toEqual(0);

    // refetch pool
    pool = await getPoolById(0);

    // refetch funders list
    funders_list = await getFundersList(0);

    // get balance before funding
    const postBalance = await getBalanceByAddress(ADDRESS_ALICE);

    // check if funds went though
    expect(pool.funders).toHaveLength(0);
    expect(funders_list).toHaveLength(0);
    expect(pool.total_funds).toEqual("0");
    expect(pool.lowest_funder).toBe("");

    // check if balance was decreased correct
    expect(postBalance.minus(preBalance)).toEqual(amount);
  });

  test.skip("fund with multiple funders", async () => {
    // define amounts
    const aliceAmount = new BigNumber(200).multipliedBy(10 ** 9);
    const bobAmount = new BigNumber(100).multipliedBy(10 ** 9);
    const charlieAmount = new BigNumber(300).multipliedBy(10 ** 9);

    // get default pool
    let pool = await getPoolById(0);

    // get funders list
    let funders_list = await getFundersList(0);

    // check if funding amount is zero
    expect(pool.funders).toHaveLength(0);
    expect(funders_list).toHaveLength(0);
    expect(pool.total_funds).toEqual("0");

    // fund alice
    const { transactionBroadcast: aliceTx } = await alice.fund(
      pool.id,
      aliceAmount
    );
    await aliceTx;

    // fund alice
    const { transactionBroadcast: bobTx } = await bob.fund(pool.id, bobAmount);
    await bobTx;

    // fund alice
    const { transactionBroadcast: charlieTx } = await charlie.fund(
      pool.id,
      charlieAmount
    );
    await charlieTx;

    // refetch pool
    pool = await getPoolById(0);

    // refetch funders list
    funders_list = await getFundersList(0);

    // check if funds went though
    expect(pool.funders).toHaveLength(Math.min(3, MAX_FUNDERS));
    if (MAX_FUNDERS == 2) {
      expect(pool.funders).toContain(ADDRESS_ALICE);
      expect(pool.funders).toContain(ADDRESS_CHARLIE);
      expect(pool.lowest_funder).toBe(ADDRESS_ALICE);
      expect(pool.total_funds).toEqual(
        aliceAmount.plus(charlieAmount).toString()
      );
    } else if (MAX_FUNDERS >= 3) {
      expect(pool.funders).toContain(ADDRESS_ALICE);
      expect(pool.funders).toContain(ADDRESS_BOB);
      expect(pool.funders).toContain(ADDRESS_CHARLIE);
      expect(pool.lowest_funder).toBe(ADDRESS_BOB);
      expect(pool.total_funds).toEqual(
        aliceAmount.plus(bobAmount).plus(charlieAmount).toString()
      );
    }

    expect(funders_list).toHaveLength(Math.min(3, MAX_FUNDERS));
  });
};
