import { ChildProcessWithoutNullStreams } from "child_process";
import { funding } from "./funding";
import { sleep, startChain } from "./helpers/utils";
import { staking } from "./staking";
import { delegation } from "./delegation";

describe("chain", () => {
  // disable timeout
  jest.setTimeout(24 * 60 * 60 * 1000);

  // define chain process
  let chain: ChildProcessWithoutNullStreams;

  beforeAll(async () => {
    // start local chain
    chain = await startChain();

    await sleep(5000);
  });

  // funding
  describe("Funding", funding);
  // staking
  describe("Staking", staking);
  // delegation
  describe("Delegation", delegation);

  afterAll(() => {
    // stop local chain
    chain.kill();
  });
});
