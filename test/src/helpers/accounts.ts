import { localSdk } from "./utils";
import { KyveClient } from "@kyve/sdk";

class Account {
  private cl: KyveClient | null;
  readonly mnemonic: string;
  constructor(mnemonic: string) {
    this.mnemonic = mnemonic;
    this.cl = null;
  }
  async init() {
    this.cl = await localSdk.fromMnemonic(this.mnemonic);
  }
  get client() {
    if (!this.cl) throw new Error("client call without initializing");
    return this.cl;
  }
}

export const ADDRESS_ALICE = "kyve1jq304cthpx0lwhpqzrdjrcza559ukyy3zsl2vd";
export const MNEMONIC_ALICE =
  "worry grief loyal smoke pencil arrow trap focus high pioneer tomato hedgehog essence purchase dove pond knee custom phone gentle sunset addict mother fabric";
export const alice = new Account(MNEMONIC_ALICE);

export const ADDRESS_BOB = "kyve1hvg7zsnrj6h29q9ss577mhrxa04rn94h7zjugq";
export const MNEMONIC_BOB =
  "crash sick toilet stumble join cash erode glory door weird diagram away lizard solid segment apple urge joy annual able tank define candy demise";
export const bob = new Account(MNEMONIC_BOB);

export const ADDRESS_CHARLIE = "kyve1ay22rr3kz659fupu0tcswlagq4ql6rwm4nuv0s";
export const MNEMONIC_CHARLIE =
  "shoot inject fragile width trend satisfy army enact volcano crowd message strike true divorce search rich office shoulder sport relax rhythm symbol gadget size";
export const charlie = new Account(MNEMONIC_CHARLIE);
