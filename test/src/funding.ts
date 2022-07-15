import BigNumber from "bignumber.js";
import {
  ADDRESS_ALICE,
  ADDRESS_BOB,
  ADDRESS_CHARLIE,
  alice,
  bob,
  charlie,
} from "./helpers/accounts";
import { lcdClient, restartChain } from "./helpers/utils";
import { MAX_FUNDERS } from "./helpers/constants";
import { constants } from "@kyve/sdk";

const getDefaultPool = async () =>
  (await lcdClient.kyve.registry.v1beta1.pool({ id: "0" })).pool ??
  (() => {
    throw new Error("Pool doesn't exist");
  })();

export const funding = () => {
  // disable timeout
  afterAll(async () => {
    await restartChain();
    await alice.init();
    await bob.init();
    await charlie.init();
  });
  jest.setTimeout(24 * 60 * 60 * 1000);

  test("Test if default pool exist", async () => {
    const data = await lcdClient.kyve.registry.v1beta1.pools({ paused: false });
    expect(data.pools).not.toBeUndefined();
    expect(data.pools).toHaveLength(1);
  });

  test("fund more than available balance", async () => {
    // get default pool
    let pool = await getDefaultPool();
    // get funders list
    let funders_list = await lcdClient.kyve.registry.v1beta1.fundersList({
      pool_id: "0",
    });

    // get balance before funding
    const preBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    const amount = new BigNumber(preBalance).plus(1);
    // check if funding amount is zero
    expect(pool.funders).toHaveLength(0);
    expect(funders_list.funders).toHaveLength(0);
    expect(pool.total_funds).toEqual("0");
    expect(pool.lowest_funder).toBe("");

    const request = alice.client.kyve.v1beta1.base
      .fundPool({
        amount: amount.toString(),
        id: pool.id,
      })
      .then((tx) => tx.execute());
    await expect(request).rejects.toThrow(/insufficient funds/);
    // refetch pool
    pool = await getDefaultPool();

    // refetch funders list
    funders_list = await lcdClient.kyve.registry.v1beta1.fundersList({
      pool_id: "0",
    });
    // get balance before funding
    const postBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();
    // check if funds went though
    expect(pool.funders).toHaveLength(0);
    expect(pool.lowest_funder).toBe("");

    expect(funders_list.funders).toHaveLength(0);
    expect(pool.total_funds).toEqual("0");

    // check if balance has not changed
    expect(preBalance).toEqual(postBalance);
  });

  test("fund 80 KYVE with alice", async () => {
    // define amount
    const amount = new BigNumber(80).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );

    // get default pool
    let pool = await getDefaultPool();

    // get funders list
    let fundersListResponse = await lcdClient.kyve.registry.v1beta1.fundersList(
      { pool_id: "0" }
    );

    // get balance before funding
    const preBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // check if funding amount is zero
    expect(pool.funders).toHaveLength(0);
    expect(fundersListResponse.funders).toHaveLength(0);
    expect(pool.total_funds).toEqual("0");
    expect(pool.lowest_funder).toBe("");

    // fund 80 KYVE
    const tx = await alice.client.kyve.v1beta1.base.fundPool({
      amount: amount.toString(),
      id: pool.id,
    });
    const receipt = await tx.execute();

    // 0 means transaction was successful
    expect(receipt.code).toEqual(0);

    // refetch pool
    pool = await getDefaultPool();

    // refetch funders list
    fundersListResponse = await lcdClient.kyve.registry.v1beta1.fundersList({
      pool_id: "0",
    });

    // get balance before funding
    const postBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // check if funds went though
    expect(pool.funders).toHaveLength(1);
    expect(pool.funders[0]).toEqual(ADDRESS_ALICE);
    expect(pool.lowest_funder).toBe(ADDRESS_ALICE);

    expect(fundersListResponse.funders).toHaveLength(1);
    expect(fundersListResponse.funders[0]).toEqual({
      account: ADDRESS_ALICE,
      pool_id: "0",
      amount: amount.toString(),
    });

    expect(pool.total_funds).toEqual(amount.toString());

    // check if balance was decreased correct
    expect(
      new BigNumber(preBalance)
        .minus(postBalance)
        .minus(tx.fee.amount[0].amount)
    ).toEqual(amount);
  });

  test("fund additional 20 KYVE with alice", async () => {
    // define amount
    const amount = new BigNumber(20).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );
    const total = new BigNumber(100).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );

    // get default pool
    let pool = await getDefaultPool();

    // get balance before funding
    const preBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // fund 100 KYVE
    const tx = await alice.client.kyve.v1beta1.base.fundPool({
      id: pool.id,
      amount: amount.toString(),
    });
    const receipt = await tx.execute();
    // 0 means transaction was successful
    expect(receipt.code).toEqual(0);

    // refetch pool
    pool = await getDefaultPool();

    // refetch funders list
    const fundersListResponse =
      await lcdClient.kyve.registry.v1beta1.fundersList({ pool_id: "0" });

    // get balance before funding
    const postBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // check if funds went though
    expect(pool.funders).toHaveLength(1);
    expect(pool.funders[0]).toEqual(ADDRESS_ALICE);
    expect(pool.lowest_funder).toBe(ADDRESS_ALICE);

    expect(fundersListResponse.funders).toHaveLength(1);
    expect(fundersListResponse.funders[0]).toEqual({
      account: ADDRESS_ALICE,
      pool_id: "0",
      amount: total.toString(),
    });

    expect(pool.total_funds).toEqual(total.toString());

    // check if balance was decreased correct
    expect(
      new BigNumber(preBalance)
        .minus(postBalance)
        .minus(tx.fee.amount[0].amount)
    ).toEqual(amount);
  });

  test("defund 80 KYVE with alice", async () => {
    // define amount
    const amount = new BigNumber(80).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );
    const remaining = new BigNumber(20).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );

    // get default pool
    let pool = await getDefaultPool();

    // get balance before funding
    const preBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // defund 80 KYVE
    const tx = await alice.client.kyve.v1beta1.base.defundPool({
      id: pool.id,
      amount: amount.toString(),
    });
    const receipt = await tx.execute();
    // 0 means transaction was successful
    expect(receipt.code).toEqual(0);

    // refetch pool
    pool = await getDefaultPool();

    // refetch funders list
    let fundersListResponse = await lcdClient.kyve.registry.v1beta1.fundersList(
      { pool_id: "0" }
    );

    // get balance before funding
    const postBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // check if funds went though
    expect(pool.funders).toHaveLength(1);
    expect(pool.funders[0]).toEqual(ADDRESS_ALICE);
    expect(pool.lowest_funder).toBe(ADDRESS_ALICE);

    expect(fundersListResponse.funders).toHaveLength(1);
    expect(fundersListResponse.funders[0]).toEqual({
      account: ADDRESS_ALICE,
      pool_id: "0",
      amount: remaining.toString(),
    });

    expect(pool.total_funds).toEqual(remaining.toString());

    // check if balance was decreased correct
    expect(new BigNumber(postBalance).minus(preBalance)).toEqual(
      amount.minus(tx.fee.amount[0].amount)
    );
  });

  test("defund more than funding balance with alice", async () => {
    // define amount
    const amount = new BigNumber(50).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );
    const remaining = new BigNumber(20).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );
    // get default pool
    let pool = await getDefaultPool();

    // get balance before funding
    const preBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // defund 50 KYVE
    const request = alice.client.kyve.v1beta1.base.defundPool({
      id: pool.id,
      amount: amount.toString(),
    });
    //check that transaction  unsuccessful
    await expect(request).rejects.toThrow(/maximum defunding amount of/);

    // refetch pool
    pool = await getDefaultPool();

    // refetch funders list
    let fundersListResponse = await lcdClient.kyve.registry.v1beta1.fundersList(
      { pool_id: "0" }
    );

    // get balance before funding
    const postBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // check if funds went though
    expect(pool.funders).toHaveLength(1);
    expect(pool.funders[0]).toEqual(ADDRESS_ALICE);
    expect(pool.lowest_funder).toBe(ADDRESS_ALICE);

    expect(fundersListResponse.funders).toHaveLength(1);
    expect(fundersListResponse.funders[0]).toEqual({
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
    const amount = new BigNumber(20).multipliedBy(
      10 ** constants.KYVE_DECIMALS
    );

    // get default pool
    let pool = await getDefaultPool();

    // get funders list

    // get balance before funding
    const preBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // defund 20 KYVE
    const tx = await alice.client.kyve.v1beta1.base.defundPool({
      id: pool.id,
      amount: amount.toString(),
    });
    const receipt = await tx.execute();
    // 0 means transaction was successful
    expect(receipt.code).toEqual(0);

    // refetch pool
    pool = await getDefaultPool();

    // refetch funders list
    let fundersListResponse = await lcdClient.kyve.registry.v1beta1.fundersList(
      { pool_id: "0" }
    );

    // get balance before funding
    const postBalance = await alice.client.kyve.v1beta1.base.getKyveBalance();

    // check if funds went though
    expect(pool.funders).toHaveLength(0);
    expect(fundersListResponse.funders).toHaveLength(0);
    expect(pool.total_funds).toEqual("0");
    expect(pool.lowest_funder).toBe("");

    // check if balance was decreased correct
    expect(new BigNumber(postBalance).minus(preBalance)).toEqual(
      amount.minus(tx.fee.amount[0].amount)
    );
  });

  test("fund with multiple funders", async () => {
    // define amounts
    const aliceAmount = new BigNumber(200).multipliedBy(10 ** 9);
    const bobAmount = new BigNumber(100).multipliedBy(10 ** 9);
    const charlieAmount = new BigNumber(300).multipliedBy(10 ** 9);

    // get default pool
    let pool = await getDefaultPool();

    // get funders list
    let fundersListResponse = await lcdClient.kyve.registry.v1beta1.fundersList(
      { pool_id: "0" }
    );

    // check if funding amount is zero
    expect(pool.funders).toHaveLength(0);
    expect(fundersListResponse.funders).toHaveLength(0);
    expect(pool.total_funds).toEqual("0");

    // fund alice
    await alice.client.kyve.v1beta1.base
      .fundPool({
        id: pool.id,
        amount: aliceAmount.toString(),
      })
      .then((tx) => tx.execute());

    // fund bob
    await bob.client.kyve.v1beta1.base
      .fundPool({
        id: pool.id,
        amount: bobAmount.toString(),
      })
      .then((tx) => tx.execute());

    // fund charlie
    await charlie.client.kyve.v1beta1.base
      .fundPool({
        id: pool.id,
        amount: charlieAmount.toString(),
      })
      .then((tx) => tx.execute());

    // refetch pool
    pool = await getDefaultPool();

    // refetch funders list
    fundersListResponse = await lcdClient.kyve.registry.v1beta1.fundersList({
      pool_id: "0",
    });

    // check if funds went though
    expect(pool.funders).toHaveLength(Math.min(3, MAX_FUNDERS));
    expect(pool.funders).toContain(ADDRESS_ALICE);
    expect(pool.funders).toContain(ADDRESS_BOB);
    expect(pool.funders).toContain(ADDRESS_CHARLIE);
    expect(pool.lowest_funder).toBe(ADDRESS_BOB);
    expect(pool.total_funds).toEqual(
      new BigNumber(aliceAmount).plus(bobAmount).plus(charlieAmount).toString()
    );

    expect(fundersListResponse.funders).toHaveLength(Math.min(3, MAX_FUNDERS));
  });
};
