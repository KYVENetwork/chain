import { funding } from "./funding";
import { startChain } from "./helpers/utils";
import { staking } from "./staking";
import { alice, bob, charlie } from "./helpers/accounts";
import { delegation } from "./delegation";
import * as dotenv from "dotenv";
dotenv.config();
describe("chain", () => {
  // disable timeout
  jest.setTimeout(24 * 60 * 60 * 1000);

  beforeAll(async () => {
    startChain.isIgniteMod = process.env.IGNITE_MODE === "true";
    if (!startChain.isIgniteMod && !process.env.COSMOS_BINARY?.length) {
      console.error("COSMOS_BINARY isn't set");
      process.exit(1);
    }
    if (!startChain.isIgniteMod && !process.env.COSMOS_DATA?.length) {
      console.error("COSMOS_DATA isn't set");
      process.exit(1);
    }
    await startChain();
    await alice.init();
    await bob.init();
    await charlie.init();
  });

  // funding
  describe("Funding", funding);
  // staking
  describe("Staking", staking);
  // // delegation
  describe("Delegation", delegation);
  afterAll(() => {
    // stop local chain
    global.chain.kill();
  });
});
