/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "mercury.mercury";

export enum ProviderStatus {
  Offline = 0,
  Online = 1,
  UNRECOGNIZED = -1,
}

export function providerStatusFromJSON(object: any): ProviderStatus {
  switch (object) {
    case 0:
    case "Offline":
      return ProviderStatus.Offline;
    case 1:
    case "Online":
      return ProviderStatus.Online;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ProviderStatus.UNRECOGNIZED;
  }
}

export function providerStatusToJSON(object: ProviderStatus): string {
  switch (object) {
    case ProviderStatus.Offline:
      return "Offline";
    case ProviderStatus.Online:
      return "Online";
    default:
      return "UNKNOWN";
  }
}

export enum ContractType {
  Subscription = 0,
  PayAsYouGo = 1,
  UNRECOGNIZED = -1,
}

export function contractTypeFromJSON(object: any): ContractType {
  switch (object) {
    case 0:
    case "Subscription":
      return ContractType.Subscription;
    case 1:
    case "PayAsYouGo":
      return ContractType.PayAsYouGo;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ContractType.UNRECOGNIZED;
  }
}

export function contractTypeToJSON(object: ContractType): string {
  switch (object) {
    case ContractType.Subscription:
      return "Subscription";
    case ContractType.PayAsYouGo:
      return "PayAsYouGo";
    default:
      return "UNKNOWN";
  }
}

export interface Provider {
  pubKey: string;
  chain: string;
  metadataUri: string;
  metadataNonce: number;
  type: ProviderStatus;
  maxContractDuration: number;
  subscriptionRate: number;
  payAsYouGoRate: number;
  bond: number;
}

export interface Contract {
  providerPubKey: string;
  chain: string;
  clientAddress: Uint8Array;
  type: ContractType;
  height: number;
  duration: number;
  rate: number;
}

const baseProvider: object = {
  pubKey: "",
  chain: "",
  metadataUri: "",
  metadataNonce: 0,
  type: 0,
  maxContractDuration: 0,
  subscriptionRate: 0,
  payAsYouGoRate: 0,
  bond: 0,
};

export const Provider = {
  encode(message: Provider, writer: Writer = Writer.create()): Writer {
    if (message.pubKey !== "") {
      writer.uint32(10).string(message.pubKey);
    }
    if (message.chain !== "") {
      writer.uint32(18).string(message.chain);
    }
    if (message.metadataUri !== "") {
      writer.uint32(26).string(message.metadataUri);
    }
    if (message.metadataNonce !== 0) {
      writer.uint32(32).uint64(message.metadataNonce);
    }
    if (message.type !== 0) {
      writer.uint32(40).int32(message.type);
    }
    if (message.maxContractDuration !== 0) {
      writer.uint32(48).uint64(message.maxContractDuration);
    }
    if (message.subscriptionRate !== 0) {
      writer.uint32(56).uint64(message.subscriptionRate);
    }
    if (message.payAsYouGoRate !== 0) {
      writer.uint32(64).uint64(message.payAsYouGoRate);
    }
    if (message.bond !== 0) {
      writer.uint32(72).uint64(message.bond);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): Provider {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseProvider } as Provider;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pubKey = reader.string();
          break;
        case 2:
          message.chain = reader.string();
          break;
        case 3:
          message.metadataUri = reader.string();
          break;
        case 4:
          message.metadataNonce = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.type = reader.int32() as any;
          break;
        case 6:
          message.maxContractDuration = longToNumber(reader.uint64() as Long);
          break;
        case 7:
          message.subscriptionRate = longToNumber(reader.uint64() as Long);
          break;
        case 8:
          message.payAsYouGoRate = longToNumber(reader.uint64() as Long);
          break;
        case 9:
          message.bond = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Provider {
    const message = { ...baseProvider } as Provider;
    if (object.pubKey !== undefined && object.pubKey !== null) {
      message.pubKey = String(object.pubKey);
    } else {
      message.pubKey = "";
    }
    if (object.chain !== undefined && object.chain !== null) {
      message.chain = String(object.chain);
    } else {
      message.chain = "";
    }
    if (object.metadataUri !== undefined && object.metadataUri !== null) {
      message.metadataUri = String(object.metadataUri);
    } else {
      message.metadataUri = "";
    }
    if (object.metadataNonce !== undefined && object.metadataNonce !== null) {
      message.metadataNonce = Number(object.metadataNonce);
    } else {
      message.metadataNonce = 0;
    }
    if (object.type !== undefined && object.type !== null) {
      message.type = providerStatusFromJSON(object.type);
    } else {
      message.type = 0;
    }
    if (
      object.maxContractDuration !== undefined &&
      object.maxContractDuration !== null
    ) {
      message.maxContractDuration = Number(object.maxContractDuration);
    } else {
      message.maxContractDuration = 0;
    }
    if (
      object.subscriptionRate !== undefined &&
      object.subscriptionRate !== null
    ) {
      message.subscriptionRate = Number(object.subscriptionRate);
    } else {
      message.subscriptionRate = 0;
    }
    if (object.payAsYouGoRate !== undefined && object.payAsYouGoRate !== null) {
      message.payAsYouGoRate = Number(object.payAsYouGoRate);
    } else {
      message.payAsYouGoRate = 0;
    }
    if (object.bond !== undefined && object.bond !== null) {
      message.bond = Number(object.bond);
    } else {
      message.bond = 0;
    }
    return message;
  },

  toJSON(message: Provider): unknown {
    const obj: any = {};
    message.pubKey !== undefined && (obj.pubKey = message.pubKey);
    message.chain !== undefined && (obj.chain = message.chain);
    message.metadataUri !== undefined &&
      (obj.metadataUri = message.metadataUri);
    message.metadataNonce !== undefined &&
      (obj.metadataNonce = message.metadataNonce);
    message.type !== undefined &&
      (obj.type = providerStatusToJSON(message.type));
    message.maxContractDuration !== undefined &&
      (obj.maxContractDuration = message.maxContractDuration);
    message.subscriptionRate !== undefined &&
      (obj.subscriptionRate = message.subscriptionRate);
    message.payAsYouGoRate !== undefined &&
      (obj.payAsYouGoRate = message.payAsYouGoRate);
    message.bond !== undefined && (obj.bond = message.bond);
    return obj;
  },

  fromPartial(object: DeepPartial<Provider>): Provider {
    const message = { ...baseProvider } as Provider;
    if (object.pubKey !== undefined && object.pubKey !== null) {
      message.pubKey = object.pubKey;
    } else {
      message.pubKey = "";
    }
    if (object.chain !== undefined && object.chain !== null) {
      message.chain = object.chain;
    } else {
      message.chain = "";
    }
    if (object.metadataUri !== undefined && object.metadataUri !== null) {
      message.metadataUri = object.metadataUri;
    } else {
      message.metadataUri = "";
    }
    if (object.metadataNonce !== undefined && object.metadataNonce !== null) {
      message.metadataNonce = object.metadataNonce;
    } else {
      message.metadataNonce = 0;
    }
    if (object.type !== undefined && object.type !== null) {
      message.type = object.type;
    } else {
      message.type = 0;
    }
    if (
      object.maxContractDuration !== undefined &&
      object.maxContractDuration !== null
    ) {
      message.maxContractDuration = object.maxContractDuration;
    } else {
      message.maxContractDuration = 0;
    }
    if (
      object.subscriptionRate !== undefined &&
      object.subscriptionRate !== null
    ) {
      message.subscriptionRate = object.subscriptionRate;
    } else {
      message.subscriptionRate = 0;
    }
    if (object.payAsYouGoRate !== undefined && object.payAsYouGoRate !== null) {
      message.payAsYouGoRate = object.payAsYouGoRate;
    } else {
      message.payAsYouGoRate = 0;
    }
    if (object.bond !== undefined && object.bond !== null) {
      message.bond = object.bond;
    } else {
      message.bond = 0;
    }
    return message;
  },
};

const baseContract: object = {
  providerPubKey: "",
  chain: "",
  type: 0,
  height: 0,
  duration: 0,
  rate: 0,
};

export const Contract = {
  encode(message: Contract, writer: Writer = Writer.create()): Writer {
    if (message.providerPubKey !== "") {
      writer.uint32(10).string(message.providerPubKey);
    }
    if (message.chain !== "") {
      writer.uint32(18).string(message.chain);
    }
    if (message.clientAddress.length !== 0) {
      writer.uint32(26).bytes(message.clientAddress);
    }
    if (message.type !== 0) {
      writer.uint32(32).int32(message.type);
    }
    if (message.height !== 0) {
      writer.uint32(40).uint64(message.height);
    }
    if (message.duration !== 0) {
      writer.uint32(48).uint64(message.duration);
    }
    if (message.rate !== 0) {
      writer.uint32(56).uint64(message.rate);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): Contract {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseContract } as Contract;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.providerPubKey = reader.string();
          break;
        case 2:
          message.chain = reader.string();
          break;
        case 3:
          message.clientAddress = reader.bytes();
          break;
        case 4:
          message.type = reader.int32() as any;
          break;
        case 5:
          message.height = longToNumber(reader.uint64() as Long);
          break;
        case 6:
          message.duration = longToNumber(reader.uint64() as Long);
          break;
        case 7:
          message.rate = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Contract {
    const message = { ...baseContract } as Contract;
    if (object.providerPubKey !== undefined && object.providerPubKey !== null) {
      message.providerPubKey = String(object.providerPubKey);
    } else {
      message.providerPubKey = "";
    }
    if (object.chain !== undefined && object.chain !== null) {
      message.chain = String(object.chain);
    } else {
      message.chain = "";
    }
    if (object.clientAddress !== undefined && object.clientAddress !== null) {
      message.clientAddress = bytesFromBase64(object.clientAddress);
    }
    if (object.type !== undefined && object.type !== null) {
      message.type = contractTypeFromJSON(object.type);
    } else {
      message.type = 0;
    }
    if (object.height !== undefined && object.height !== null) {
      message.height = Number(object.height);
    } else {
      message.height = 0;
    }
    if (object.duration !== undefined && object.duration !== null) {
      message.duration = Number(object.duration);
    } else {
      message.duration = 0;
    }
    if (object.rate !== undefined && object.rate !== null) {
      message.rate = Number(object.rate);
    } else {
      message.rate = 0;
    }
    return message;
  },

  toJSON(message: Contract): unknown {
    const obj: any = {};
    message.providerPubKey !== undefined &&
      (obj.providerPubKey = message.providerPubKey);
    message.chain !== undefined && (obj.chain = message.chain);
    message.clientAddress !== undefined &&
      (obj.clientAddress = base64FromBytes(
        message.clientAddress !== undefined
          ? message.clientAddress
          : new Uint8Array()
      ));
    message.type !== undefined && (obj.type = contractTypeToJSON(message.type));
    message.height !== undefined && (obj.height = message.height);
    message.duration !== undefined && (obj.duration = message.duration);
    message.rate !== undefined && (obj.rate = message.rate);
    return obj;
  },

  fromPartial(object: DeepPartial<Contract>): Contract {
    const message = { ...baseContract } as Contract;
    if (object.providerPubKey !== undefined && object.providerPubKey !== null) {
      message.providerPubKey = object.providerPubKey;
    } else {
      message.providerPubKey = "";
    }
    if (object.chain !== undefined && object.chain !== null) {
      message.chain = object.chain;
    } else {
      message.chain = "";
    }
    if (object.clientAddress !== undefined && object.clientAddress !== null) {
      message.clientAddress = object.clientAddress;
    } else {
      message.clientAddress = new Uint8Array();
    }
    if (object.type !== undefined && object.type !== null) {
      message.type = object.type;
    } else {
      message.type = 0;
    }
    if (object.height !== undefined && object.height !== null) {
      message.height = object.height;
    } else {
      message.height = 0;
    }
    if (object.duration !== undefined && object.duration !== null) {
      message.duration = object.duration;
    } else {
      message.duration = 0;
    }
    if (object.rate !== undefined && object.rate !== null) {
      message.rate = object.rate;
    } else {
      message.rate = 0;
    }
    return message;
  },
};

declare var self: any | undefined;
declare var window: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") return globalThis;
  if (typeof self !== "undefined") return self;
  if (typeof window !== "undefined") return window;
  if (typeof global !== "undefined") return global;
  throw "Unable to locate global object";
})();

const atob: (b64: string) => string =
  globalThis.atob ||
  ((b64) => globalThis.Buffer.from(b64, "base64").toString("binary"));
function bytesFromBase64(b64: string): Uint8Array {
  const bin = atob(b64);
  const arr = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; ++i) {
    arr[i] = bin.charCodeAt(i);
  }
  return arr;
}

const btoa: (bin: string) => string =
  globalThis.btoa ||
  ((bin) => globalThis.Buffer.from(bin, "binary").toString("base64"));
function base64FromBytes(arr: Uint8Array): string {
  const bin: string[] = [];
  for (let i = 0; i < arr.byteLength; ++i) {
    bin.push(String.fromCharCode(arr[i]));
  }
  return btoa(bin.join(""));
}

type Builtin = Date | Function | Uint8Array | string | number | undefined;
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (util.Long !== Long) {
  util.Long = Long as any;
  configure();
}
