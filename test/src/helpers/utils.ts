import axios from "axios";
import BigNumber from "bignumber.js";
import { spawn, ChildProcessWithoutNullStreams } from "child_process";
import { BASE_URL } from "./constants";

export const sleep = (milliseconds: number) => {
  return new Promise((resolve) => setTimeout(resolve, milliseconds));
};

export const startChain = (): Promise<ChildProcessWithoutNullStreams> => {
  return new Promise<ChildProcessWithoutNullStreams>((resolve, reject) => {
    console.log("Starting chain on localhost ...");
    console.log("This may take up to a minute");

    const chain = spawn("ignite", ["chain", "serve", "--reset-once"]);

    chain.stdout.on("data", (data: Buffer) => {
      if (data.toString().includes(`Blockchain API: ${BASE_URL}`)) {
        console.log(`Chain started on ${BASE_URL}`);
        setTimeout(() => resolve(chain), 5000);
      }
    });

    chain.stderr.on("data", () => {
      reject();
    });

    chain.on("close", () => {
      reject();
    });
  });
};

export const getBalanceByAddress = async (
  address: string
): Promise<BigNumber> => {
  const { data } = await axios.get(
    `${BASE_URL}/cosmos/bank/v1beta1/balances/${address}/by_denom?denom=tkyve`
  );
  return new BigNumber(data.balance.amount);
};

export const getPoolById = async (id: string | number): Promise<any> => {
  const { data } = await axios.get(
    `${BASE_URL}/kyve/registry/v1beta1/pool/${id}`
  );
  return data.pool;
};

export const getFundersList = async (id: string | number): Promise<any[]> => {
  const { data } = await axios.get(
    `${BASE_URL}/kyve/registry/v1beta1/funders_list/${id}`
  );
  return data.funders;
};

export const getStakersList = async (id: string | number): Promise<any[]> => {
  const { data } = await axios.get(
    `${BASE_URL}/kyve/registry/v1beta1/stakers_list/${id}`
  );
  return data.stakers;
};

export const getDelegationPoolData = async (
  staker: string,
  poolId: string | number
): Promise<any> => {
  try {
    const { data } = await axios.get(
      `${BASE_URL}/kyve/registry/v1beta1/delegators_by_pool_and_staker/${poolId}/${staker}`
    );
    return data.delegation_pool_data;
  } catch (e) {
    return undefined;
  }
};
