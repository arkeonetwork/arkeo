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
  metadataURI: string;
  metadataNonce: number;
  status: ProviderStatus;
  minContractDuration: number;
  maxContractDuration: number;
  subscriptionRate: number;
  payAsYouGoRate: number;
  bond: string;
  lastUpdate: number;
}

export interface Contract {
  providerPubKey: string;
  chain: string;
  client: string;
  type: ContractType;
  height: number;
  duration: number;
  rate: number;
  deposit: string;
  paid: string;
  nonce: number;
  closedHeight: number;
}

export interface ContractExpiration {
  providerPubKey: string;
  chain: string;
  client: string;
}

export interface ContractExpirationSet {
  height: number;
  contracts: ContractExpiration[];
}

const baseProvider: object = {
  pubKey: "",
  chain: "",
  metadataURI: "",
  metadataNonce: 0,
  status: 0,
  minContractDuration: 0,
  maxContractDuration: 0,
  subscriptionRate: 0,
  payAsYouGoRate: 0,
  bond: "",
  lastUpdate: 0,
};

export const Provider = {
  encode(message: Provider, writer: Writer = Writer.create()): Writer {
    if (message.pubKey !== "") {
      writer.uint32(10).string(message.pubKey);
    }
    if (message.chain !== "") {
      writer.uint32(18).string(message.chain);
    }
    if (message.metadataURI !== "") {
      writer.uint32(26).string(message.metadataURI);
    }
    if (message.metadataNonce !== 0) {
      writer.uint32(32).uint64(message.metadataNonce);
    }
    if (message.status !== 0) {
      writer.uint32(40).int32(message.status);
    }
    if (message.minContractDuration !== 0) {
      writer.uint32(48).int64(message.minContractDuration);
    }
    if (message.maxContractDuration !== 0) {
      writer.uint32(56).int64(message.maxContractDuration);
    }
    if (message.subscriptionRate !== 0) {
      writer.uint32(64).int64(message.subscriptionRate);
    }
    if (message.payAsYouGoRate !== 0) {
      writer.uint32(72).int64(message.payAsYouGoRate);
    }
    if (message.bond !== "") {
      writer.uint32(82).string(message.bond);
    }
    if (message.lastUpdate !== 0) {
      writer.uint32(88).int64(message.lastUpdate);
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
          message.metadataURI = reader.string();
          break;
        case 4:
          message.metadataNonce = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.status = reader.int32() as any;
          break;
        case 6:
          message.minContractDuration = longToNumber(reader.int64() as Long);
          break;
        case 7:
          message.maxContractDuration = longToNumber(reader.int64() as Long);
          break;
        case 8:
          message.subscriptionRate = longToNumber(reader.int64() as Long);
          break;
        case 9:
          message.payAsYouGoRate = longToNumber(reader.int64() as Long);
          break;
        case 10:
          message.bond = reader.string();
          break;
        case 11:
          message.lastUpdate = longToNumber(reader.int64() as Long);
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
    if (object.metadataURI !== undefined && object.metadataURI !== null) {
      message.metadataURI = String(object.metadataURI);
    } else {
      message.metadataURI = "";
    }
    if (object.metadataNonce !== undefined && object.metadataNonce !== null) {
      message.metadataNonce = Number(object.metadataNonce);
    } else {
      message.metadataNonce = 0;
    }
    if (object.status !== undefined && object.status !== null) {
      message.status = providerStatusFromJSON(object.status);
    } else {
      message.status = 0;
    }
    if (
      object.minContractDuration !== undefined &&
      object.minContractDuration !== null
    ) {
      message.minContractDuration = Number(object.minContractDuration);
    } else {
      message.minContractDuration = 0;
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
      message.bond = String(object.bond);
    } else {
      message.bond = "";
    }
    if (object.lastUpdate !== undefined && object.lastUpdate !== null) {
      message.lastUpdate = Number(object.lastUpdate);
    } else {
      message.lastUpdate = 0;
    }
    return message;
  },

  toJSON(message: Provider): unknown {
    const obj: any = {};
    message.pubKey !== undefined && (obj.pubKey = message.pubKey);
    message.chain !== undefined && (obj.chain = message.chain);
    message.metadataURI !== undefined &&
      (obj.metadataURI = message.metadataURI);
    message.metadataNonce !== undefined &&
      (obj.metadataNonce = message.metadataNonce);
    message.status !== undefined &&
      (obj.status = providerStatusToJSON(message.status));
    message.minContractDuration !== undefined &&
      (obj.minContractDuration = message.minContractDuration);
    message.maxContractDuration !== undefined &&
      (obj.maxContractDuration = message.maxContractDuration);
    message.subscriptionRate !== undefined &&
      (obj.subscriptionRate = message.subscriptionRate);
    message.payAsYouGoRate !== undefined &&
      (obj.payAsYouGoRate = message.payAsYouGoRate);
    message.bond !== undefined && (obj.bond = message.bond);
    message.lastUpdate !== undefined && (obj.lastUpdate = message.lastUpdate);
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
    if (object.metadataURI !== undefined && object.metadataURI !== null) {
      message.metadataURI = object.metadataURI;
    } else {
      message.metadataURI = "";
    }
    if (object.metadataNonce !== undefined && object.metadataNonce !== null) {
      message.metadataNonce = object.metadataNonce;
    } else {
      message.metadataNonce = 0;
    }
    if (object.status !== undefined && object.status !== null) {
      message.status = object.status;
    } else {
      message.status = 0;
    }
    if (
      object.minContractDuration !== undefined &&
      object.minContractDuration !== null
    ) {
      message.minContractDuration = object.minContractDuration;
    } else {
      message.minContractDuration = 0;
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
      message.bond = "";
    }
    if (object.lastUpdate !== undefined && object.lastUpdate !== null) {
      message.lastUpdate = object.lastUpdate;
    } else {
      message.lastUpdate = 0;
    }
    return message;
  },
};

const baseContract: object = {
  providerPubKey: "",
  chain: "",
  client: "",
  type: 0,
  height: 0,
  duration: 0,
  rate: 0,
  deposit: "",
  paid: "",
  nonce: 0,
  closedHeight: 0,
};

export const Contract = {
  encode(message: Contract, writer: Writer = Writer.create()): Writer {
    if (message.providerPubKey !== "") {
      writer.uint32(10).string(message.providerPubKey);
    }
    if (message.chain !== "") {
      writer.uint32(18).string(message.chain);
    }
    if (message.client !== "") {
      writer.uint32(26).string(message.client);
    }
    if (message.type !== 0) {
      writer.uint32(32).int32(message.type);
    }
    if (message.height !== 0) {
      writer.uint32(40).int64(message.height);
    }
    if (message.duration !== 0) {
      writer.uint32(48).int64(message.duration);
    }
    if (message.rate !== 0) {
      writer.uint32(56).int64(message.rate);
    }
    if (message.deposit !== "") {
      writer.uint32(66).string(message.deposit);
    }
    if (message.paid !== "") {
      writer.uint32(74).string(message.paid);
    }
    if (message.nonce !== 0) {
      writer.uint32(80).int64(message.nonce);
    }
    if (message.closedHeight !== 0) {
      writer.uint32(88).int64(message.closedHeight);
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
          message.client = reader.string();
          break;
        case 4:
          message.type = reader.int32() as any;
          break;
        case 5:
          message.height = longToNumber(reader.int64() as Long);
          break;
        case 6:
          message.duration = longToNumber(reader.int64() as Long);
          break;
        case 7:
          message.rate = longToNumber(reader.int64() as Long);
          break;
        case 8:
          message.deposit = reader.string();
          break;
        case 9:
          message.paid = reader.string();
          break;
        case 10:
          message.nonce = longToNumber(reader.int64() as Long);
          break;
        case 11:
          message.closedHeight = longToNumber(reader.int64() as Long);
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
    if (object.client !== undefined && object.client !== null) {
      message.client = String(object.client);
    } else {
      message.client = "";
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
    if (object.deposit !== undefined && object.deposit !== null) {
      message.deposit = String(object.deposit);
    } else {
      message.deposit = "";
    }
    if (object.paid !== undefined && object.paid !== null) {
      message.paid = String(object.paid);
    } else {
      message.paid = "";
    }
    if (object.nonce !== undefined && object.nonce !== null) {
      message.nonce = Number(object.nonce);
    } else {
      message.nonce = 0;
    }
    if (object.closedHeight !== undefined && object.closedHeight !== null) {
      message.closedHeight = Number(object.closedHeight);
    } else {
      message.closedHeight = 0;
    }
    return message;
  },

  toJSON(message: Contract): unknown {
    const obj: any = {};
    message.providerPubKey !== undefined &&
      (obj.providerPubKey = message.providerPubKey);
    message.chain !== undefined && (obj.chain = message.chain);
    message.client !== undefined && (obj.client = message.client);
    message.type !== undefined && (obj.type = contractTypeToJSON(message.type));
    message.height !== undefined && (obj.height = message.height);
    message.duration !== undefined && (obj.duration = message.duration);
    message.rate !== undefined && (obj.rate = message.rate);
    message.deposit !== undefined && (obj.deposit = message.deposit);
    message.paid !== undefined && (obj.paid = message.paid);
    message.nonce !== undefined && (obj.nonce = message.nonce);
    message.closedHeight !== undefined &&
      (obj.closedHeight = message.closedHeight);
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
    if (object.client !== undefined && object.client !== null) {
      message.client = object.client;
    } else {
      message.client = "";
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
    if (object.deposit !== undefined && object.deposit !== null) {
      message.deposit = object.deposit;
    } else {
      message.deposit = "";
    }
    if (object.paid !== undefined && object.paid !== null) {
      message.paid = object.paid;
    } else {
      message.paid = "";
    }
    if (object.nonce !== undefined && object.nonce !== null) {
      message.nonce = object.nonce;
    } else {
      message.nonce = 0;
    }
    if (object.closedHeight !== undefined && object.closedHeight !== null) {
      message.closedHeight = object.closedHeight;
    } else {
      message.closedHeight = 0;
    }
    return message;
  },
};

const baseContractExpiration: object = {
  providerPubKey: "",
  chain: "",
  client: "",
};

export const ContractExpiration = {
  encode(
    message: ContractExpiration,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.providerPubKey !== "") {
      writer.uint32(10).string(message.providerPubKey);
    }
    if (message.chain !== "") {
      writer.uint32(18).string(message.chain);
    }
    if (message.client !== "") {
      writer.uint32(26).string(message.client);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ContractExpiration {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseContractExpiration } as ContractExpiration;
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
          message.client = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ContractExpiration {
    const message = { ...baseContractExpiration } as ContractExpiration;
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
    if (object.client !== undefined && object.client !== null) {
      message.client = String(object.client);
    } else {
      message.client = "";
    }
    return message;
  },

  toJSON(message: ContractExpiration): unknown {
    const obj: any = {};
    message.providerPubKey !== undefined &&
      (obj.providerPubKey = message.providerPubKey);
    message.chain !== undefined && (obj.chain = message.chain);
    message.client !== undefined && (obj.client = message.client);
    return obj;
  },

  fromPartial(object: DeepPartial<ContractExpiration>): ContractExpiration {
    const message = { ...baseContractExpiration } as ContractExpiration;
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
    if (object.client !== undefined && object.client !== null) {
      message.client = object.client;
    } else {
      message.client = "";
    }
    return message;
  },
};

const baseContractExpirationSet: object = { height: 0 };

export const ContractExpirationSet = {
  encode(
    message: ContractExpirationSet,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.height !== 0) {
      writer.uint32(8).int64(message.height);
    }
    for (const v of message.contracts) {
      ContractExpiration.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ContractExpirationSet {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseContractExpirationSet } as ContractExpirationSet;
    message.contracts = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.height = longToNumber(reader.int64() as Long);
          break;
        case 2:
          message.contracts.push(
            ContractExpiration.decode(reader, reader.uint32())
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ContractExpirationSet {
    const message = { ...baseContractExpirationSet } as ContractExpirationSet;
    message.contracts = [];
    if (object.height !== undefined && object.height !== null) {
      message.height = Number(object.height);
    } else {
      message.height = 0;
    }
    if (object.contracts !== undefined && object.contracts !== null) {
      for (const e of object.contracts) {
        message.contracts.push(ContractExpiration.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: ContractExpirationSet): unknown {
    const obj: any = {};
    message.height !== undefined && (obj.height = message.height);
    if (message.contracts) {
      obj.contracts = message.contracts.map((e) =>
        e ? ContractExpiration.toJSON(e) : undefined
      );
    } else {
      obj.contracts = [];
    }
    return obj;
  },

  fromPartial(
    object: DeepPartial<ContractExpirationSet>
  ): ContractExpirationSet {
    const message = { ...baseContractExpirationSet } as ContractExpirationSet;
    message.contracts = [];
    if (object.height !== undefined && object.height !== null) {
      message.height = object.height;
    } else {
      message.height = 0;
    }
    if (object.contracts !== undefined && object.contracts !== null) {
      for (const e of object.contracts) {
        message.contracts.push(ContractExpiration.fromPartial(e));
      }
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
