import axios from "axios";
import { BASE_URL, UNSTAKING_TIME } from "./helpers/constants";

export const delegation = () => {
  // disable timeout
  jest.setTimeout(24 * 60 * 60 * 1000);

  test("Test if default pool exist", async () => {
    const { data } = await axios.get(`${BASE_URL}/kyve/registry/v1beta1/pools`);

    expect(data.pools).not.toBeUndefined();
    expect(data.pools).toHaveLength(1);
  });

  // test("Claim uploader role without stake", async () => {
  //
  //     // get default pool
  //     let pool = await getPoolById(0);
  //
  //     await alice.
  //
  //     // stake alice
  //     const { transactionBroadcast: aliceTx } = await alice.stake(
  //         pool.id,
  //         aliceAmount
  //     );
  //     await aliceTx;
  //
  //     // stake alice
  //     const { transactionBroadcast: bobTx } = await bob.stake(pool.id, bobAmount);
  //     await bobTx;
  //
  //     // stake alice
  //     const { transactionBroadcast: charlieTx } = await charlie.stake(
  //         pool.id,
  //         charlieAmount
  //     );
  //     await charlieTx;
  //
  //     // refetch pool
  //     pool = await getPoolById(0);
  //
  //     // refetch stakers list
  //     stakers_list = await getStakersList(0);
  //
  //     // check if stake went though
  //     expect(pool.stakers).toHaveLength(3);
  //     expect(pool.stakers).toContain(ADDRESS_ALICE);
  //     expect(pool.stakers).toContain(ADDRESS_BOB);
  //     expect(pool.stakers).toContain(ADDRESS_CHARLIE);
  //     expect(pool.lowest_staker).toBe(ADDRESS_BOB);
  //
  //     expect(stakers_list).toHaveLength(3);
  //
  //     expect(pool.total_stake).toEqual(
  //         aliceAmount.plus(bobAmount).plus(charlieAmount).toString()
  //     );
  // });
};
