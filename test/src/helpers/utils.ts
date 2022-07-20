import {
  spawn,
  ChildProcessWithoutNullStreams,
  execFileSync,
} from "child_process";
import { NETWORK } from "./constants";
import KyveSDK from "@kyve/sdk";

export const sleep = (milliseconds: number) => {
  return new Promise((resolve) => setTimeout(resolve, milliseconds));
};
declare global {
  var chain: ChildProcessWithoutNullStreams;
}
type StartChainType = {
  (): Promise<ChildProcessWithoutNullStreams>;
  isIgniteMod: boolean;
};

export const startChain = <StartChainType>(() => {
  return new Promise<ChildProcessWithoutNullStreams>((resolve, reject) => {
    console.log("Starting chain on localhost ...");
    console.log("This may take up to a minute");
    if (startChain.isIgniteMod) {
      global.chain = spawn("ignite", ["chain", "serve", "--reset-once"]);
    } else {
      execFileSync(process.env.COSMOS_BINARY as string, [
        "tendermint",
        "unsafe-reset-all",
        "--home",
        process.env.COSMOS_DATA as string,
      ]);
      global.chain = spawn(process.env.COSMOS_BINARY as string, [
        "start",
        "--home",
        process.env.COSMOS_DATA as string,
      ]);
    }

    function processData(data: Buffer) {
      if (startChain.isIgniteMod && data.toString().includes(`Token faucet`)) {
        console.log(`Ignite chain started`);
        setTimeout(() => resolve(chain), 2000);
      } else if (data.toString().includes(`Starting RPC HTTP server`)) {
        console.log(`Binary chain started`);
        setTimeout(() => resolve(chain), 2000);
      }
    }
    function processError() {
      reject(new Error("Process error"));
    }
    if (startChain.isIgniteMod) {
      chain.stdout.on("data", processData);
      chain.stderr.on("data", processError);
    } else {
      //cosmos binary data to stderr
      chain.stderr.on("data", processData);
      chain.stdout.on("data", processError);
    }
  });
});
export const restartChain = () => {
  return new Promise<ChildProcessWithoutNullStreams>((resolve) => {
    global.chain.on("close", function test() {
      setTimeout(async () => {
        await startChain();
        global.chain.removeAllListeners("close");
        resolve(global.chain);
      }, 2000);
    });
    global.chain.kill();
  });
};

export const localSdk = new KyveSDK(NETWORK);
export const lcdClient = localSdk.createLCDClient();
