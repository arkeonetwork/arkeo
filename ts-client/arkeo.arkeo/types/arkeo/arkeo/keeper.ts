/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Coin } from "../../cosmos/base/v1beta1/coin";

export const protobufPackage = "arkeo.arkeo";

export enum ProviderStatus {
  OFFLINE = 0,
  ONLINE = 1,
  UNRECOGNIZED = -1,
}

export function providerStatusFromJSON(object: any): ProviderStatus {
  switch (object) {
    case 0:
    case "OFFLINE":
      return ProviderStatus.OFFLINE;
    case 1:
    case "ONLINE":
      return ProviderStatus.ONLINE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ProviderStatus.UNRECOGNIZED;
  }
}

export function providerStatusToJSON(object: ProviderStatus): string {
  switch (object) {
    case ProviderStatus.OFFLINE:
      return "OFFLINE";
    case ProviderStatus.ONLINE:
      return "ONLINE";
    case ProviderStatus.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum ContractType {
  SUBSCRIPTION = 0,
  PAY_AS_YOU_GO = 1,
  UNRECOGNIZED = -1,
}

export function contractTypeFromJSON(object: any): ContractType {
  switch (object) {
    case 0:
    case "SUBSCRIPTION":
      return ContractType.SUBSCRIPTION;
    case 1:
    case "PAY_AS_YOU_GO":
      return ContractType.PAY_AS_YOU_GO;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ContractType.UNRECOGNIZED;
  }
}

export function contractTypeToJSON(object: ContractType): string {
  switch (object) {
    case ContractType.SUBSCRIPTION:
      return "SUBSCRIPTION";
    case ContractType.PAY_AS_YOU_GO:
      return "PAY_AS_YOU_GO";
    case ContractType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum ContractAuthorization {
  STRICT = 0,
  OPEN = 1,
  UNRECOGNIZED = -1,
}

export function contractAuthorizationFromJSON(object: any): ContractAuthorization {
  switch (object) {
    case 0:
    case "STRICT":
      return ContractAuthorization.STRICT;
    case 1:
    case "OPEN":
      return ContractAuthorization.OPEN;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ContractAuthorization.UNRECOGNIZED;
  }
}

export function contractAuthorizationToJSON(object: ContractAuthorization): string {
  switch (object) {
    case ContractAuthorization.STRICT:
      return "STRICT";
    case ContractAuthorization.OPEN:
      return "OPEN";
    case ContractAuthorization.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface Provider {
  pubKey: Uint8Array;
  service: number;
  metadataUri: string;
  metadataNonce: number;
  status: ProviderStatus;
  minContractDuration: number;
  maxContractDuration: number;
  subscriptionRate: Coin[];
  payAsYouGoRate: Coin[];
  bond: string;
  lastUpdate: number;
  settlementDuration: number;
}

export interface Contract {
  provider: Uint8Array;
  service: number;
  client: Uint8Array;
  delegate: Uint8Array;
  type: ContractType;
  height: number;
  duration: number;
  rate: Coin | undefined;
  deposit: string;
  paid: string;
  nonce: number;
  settlementHeight: number;
  id: number;
  settlementDuration: number;
  authorization: ContractAuthorization;
  queriesPerMinute: number;
}

export interface ContractSet {
  contractIds: number[];
}

export interface ContractExpirationSet {
  height: number;
  contractSet: ContractSet | undefined;
}

export interface UserContractSet {
  user: Uint8Array;
  contractSet: ContractSet | undefined;
}

function createBaseProvider(): Provider {
  return {
    pubKey: new Uint8Array(),
    service: 0,
    metadataUri: "",
    metadataNonce: 0,
    status: 0,
    minContractDuration: 0,
    maxContractDuration: 0,
    subscriptionRate: [],
    payAsYouGoRate: [],
    bond: "",
    lastUpdate: 0,
    settlementDuration: 0,
  };
}

export const Provider = {
  encode(message: Provider, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pubKey.length !== 0) {
      writer.uint32(10).bytes(message.pubKey);
    }
    if (message.service !== 0) {
      writer.uint32(16).int32(message.service);
    }
    if (message.metadataUri !== "") {
      writer.uint32(26).string(message.metadataUri);
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
    for (const v of message.subscriptionRate) {
      Coin.encode(v!, writer.uint32(66).fork()).ldelim();
    }
    for (const v of message.payAsYouGoRate) {
      Coin.encode(v!, writer.uint32(74).fork()).ldelim();
    }
    if (message.bond !== "") {
      writer.uint32(82).string(message.bond);
    }
    if (message.lastUpdate !== 0) {
      writer.uint32(88).int64(message.lastUpdate);
    }
    if (message.settlementDuration !== 0) {
      writer.uint32(96).int64(message.settlementDuration);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Provider {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseProvider();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pubKey = reader.bytes();
          break;
        case 2:
          message.service = reader.int32();
          break;
        case 3:
          message.metadataUri = reader.string();
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
          message.subscriptionRate.push(Coin.decode(reader, reader.uint32()));
          break;
        case 9:
          message.payAsYouGoRate.push(Coin.decode(reader, reader.uint32()));
          break;
        case 10:
          message.bond = reader.string();
          break;
        case 11:
          message.lastUpdate = longToNumber(reader.int64() as Long);
          break;
        case 12:
          message.settlementDuration = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Provider {
    return {
      pubKey: isSet(object.pubKey) ? bytesFromBase64(object.pubKey) : new Uint8Array(),
      service: isSet(object.service) ? Number(object.service) : 0,
      metadataUri: isSet(object.metadataUri) ? String(object.metadataUri) : "",
      metadataNonce: isSet(object.metadataNonce) ? Number(object.metadataNonce) : 0,
      status: isSet(object.status) ? providerStatusFromJSON(object.status) : 0,
      minContractDuration: isSet(object.minContractDuration) ? Number(object.minContractDuration) : 0,
      maxContractDuration: isSet(object.maxContractDuration) ? Number(object.maxContractDuration) : 0,
      subscriptionRate: Array.isArray(object?.subscriptionRate)
        ? object.subscriptionRate.map((e: any) => Coin.fromJSON(e))
        : [],
      payAsYouGoRate: Array.isArray(object?.payAsYouGoRate)
        ? object.payAsYouGoRate.map((e: any) => Coin.fromJSON(e))
        : [],
      bond: isSet(object.bond) ? String(object.bond) : "",
      lastUpdate: isSet(object.lastUpdate) ? Number(object.lastUpdate) : 0,
      settlementDuration: isSet(object.settlementDuration) ? Number(object.settlementDuration) : 0,
    };
  },

  toJSON(message: Provider): unknown {
    const obj: any = {};
    message.pubKey !== undefined
      && (obj.pubKey = base64FromBytes(message.pubKey !== undefined ? message.pubKey : new Uint8Array()));
    message.service !== undefined && (obj.service = Math.round(message.service));
    message.metadataUri !== undefined && (obj.metadataUri = message.metadataUri);
    message.metadataNonce !== undefined && (obj.metadataNonce = Math.round(message.metadataNonce));
    message.status !== undefined && (obj.status = providerStatusToJSON(message.status));
    message.minContractDuration !== undefined && (obj.minContractDuration = Math.round(message.minContractDuration));
    message.maxContractDuration !== undefined && (obj.maxContractDuration = Math.round(message.maxContractDuration));
    if (message.subscriptionRate) {
      obj.subscriptionRate = message.subscriptionRate.map((e) => e ? Coin.toJSON(e) : undefined);
    } else {
      obj.subscriptionRate = [];
    }
    if (message.payAsYouGoRate) {
      obj.payAsYouGoRate = message.payAsYouGoRate.map((e) => e ? Coin.toJSON(e) : undefined);
    } else {
      obj.payAsYouGoRate = [];
    }
    message.bond !== undefined && (obj.bond = message.bond);
    message.lastUpdate !== undefined && (obj.lastUpdate = Math.round(message.lastUpdate));
    message.settlementDuration !== undefined && (obj.settlementDuration = Math.round(message.settlementDuration));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Provider>, I>>(object: I): Provider {
    const message = createBaseProvider();
    message.pubKey = object.pubKey ?? new Uint8Array();
    message.service = object.service ?? 0;
    message.metadataUri = object.metadataUri ?? "";
    message.metadataNonce = object.metadataNonce ?? 0;
    message.status = object.status ?? 0;
    message.minContractDuration = object.minContractDuration ?? 0;
    message.maxContractDuration = object.maxContractDuration ?? 0;
    message.subscriptionRate = object.subscriptionRate?.map((e) => Coin.fromPartial(e)) || [];
    message.payAsYouGoRate = object.payAsYouGoRate?.map((e) => Coin.fromPartial(e)) || [];
    message.bond = object.bond ?? "";
    message.lastUpdate = object.lastUpdate ?? 0;
    message.settlementDuration = object.settlementDuration ?? 0;
    return message;
  },
};

function createBaseContract(): Contract {
  return {
    provider: new Uint8Array(),
    service: 0,
    client: new Uint8Array(),
    delegate: new Uint8Array(),
    type: 0,
    height: 0,
    duration: 0,
    rate: undefined,
    deposit: "",
    paid: "",
    nonce: 0,
    settlementHeight: 0,
    id: 0,
    settlementDuration: 0,
    authorization: 0,
    queriesPerMinute: 0,
  };
}

export const Contract = {
  encode(message: Contract, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.provider.length !== 0) {
      writer.uint32(10).bytes(message.provider);
    }
    if (message.service !== 0) {
      writer.uint32(16).int32(message.service);
    }
    if (message.client.length !== 0) {
      writer.uint32(26).bytes(message.client);
    }
    if (message.delegate.length !== 0) {
      writer.uint32(34).bytes(message.delegate);
    }
    if (message.type !== 0) {
      writer.uint32(40).int32(message.type);
    }
    if (message.height !== 0) {
      writer.uint32(48).int64(message.height);
    }
    if (message.duration !== 0) {
      writer.uint32(56).int64(message.duration);
    }
    if (message.rate !== undefined) {
      Coin.encode(message.rate, writer.uint32(66).fork()).ldelim();
    }
    if (message.deposit !== "") {
      writer.uint32(74).string(message.deposit);
    }
    if (message.paid !== "") {
      writer.uint32(82).string(message.paid);
    }
    if (message.nonce !== 0) {
      writer.uint32(88).int64(message.nonce);
    }
    if (message.settlementHeight !== 0) {
      writer.uint32(96).int64(message.settlementHeight);
    }
    if (message.id !== 0) {
      writer.uint32(104).uint64(message.id);
    }
    if (message.settlementDuration !== 0) {
      writer.uint32(112).int64(message.settlementDuration);
    }
    if (message.authorization !== 0) {
      writer.uint32(120).int32(message.authorization);
    }
    if (message.queriesPerMinute !== 0) {
      writer.uint32(128).int64(message.queriesPerMinute);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Contract {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseContract();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.provider = reader.bytes();
          break;
        case 2:
          message.service = reader.int32();
          break;
        case 3:
          message.client = reader.bytes();
          break;
        case 4:
          message.delegate = reader.bytes();
          break;
        case 5:
          message.type = reader.int32() as any;
          break;
        case 6:
          message.height = longToNumber(reader.int64() as Long);
          break;
        case 7:
          message.duration = longToNumber(reader.int64() as Long);
          break;
        case 8:
          message.rate = Coin.decode(reader, reader.uint32());
          break;
        case 9:
          message.deposit = reader.string();
          break;
        case 10:
          message.paid = reader.string();
          break;
        case 11:
          message.nonce = longToNumber(reader.int64() as Long);
          break;
        case 12:
          message.settlementHeight = longToNumber(reader.int64() as Long);
          break;
        case 13:
          message.id = longToNumber(reader.uint64() as Long);
          break;
        case 14:
          message.settlementDuration = longToNumber(reader.int64() as Long);
          break;
        case 15:
          message.authorization = reader.int32() as any;
          break;
        case 16:
          message.queriesPerMinute = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Contract {
    return {
      provider: isSet(object.provider) ? bytesFromBase64(object.provider) : new Uint8Array(),
      service: isSet(object.service) ? Number(object.service) : 0,
      client: isSet(object.client) ? bytesFromBase64(object.client) : new Uint8Array(),
      delegate: isSet(object.delegate) ? bytesFromBase64(object.delegate) : new Uint8Array(),
      type: isSet(object.type) ? contractTypeFromJSON(object.type) : 0,
      height: isSet(object.height) ? Number(object.height) : 0,
      duration: isSet(object.duration) ? Number(object.duration) : 0,
      rate: isSet(object.rate) ? Coin.fromJSON(object.rate) : undefined,
      deposit: isSet(object.deposit) ? String(object.deposit) : "",
      paid: isSet(object.paid) ? String(object.paid) : "",
      nonce: isSet(object.nonce) ? Number(object.nonce) : 0,
      settlementHeight: isSet(object.settlementHeight) ? Number(object.settlementHeight) : 0,
      id: isSet(object.id) ? Number(object.id) : 0,
      settlementDuration: isSet(object.settlementDuration) ? Number(object.settlementDuration) : 0,
      authorization: isSet(object.authorization) ? contractAuthorizationFromJSON(object.authorization) : 0,
      queriesPerMinute: isSet(object.queriesPerMinute) ? Number(object.queriesPerMinute) : 0,
    };
  },

  toJSON(message: Contract): unknown {
    const obj: any = {};
    message.provider !== undefined
      && (obj.provider = base64FromBytes(message.provider !== undefined ? message.provider : new Uint8Array()));
    message.service !== undefined && (obj.service = Math.round(message.service));
    message.client !== undefined
      && (obj.client = base64FromBytes(message.client !== undefined ? message.client : new Uint8Array()));
    message.delegate !== undefined
      && (obj.delegate = base64FromBytes(message.delegate !== undefined ? message.delegate : new Uint8Array()));
    message.type !== undefined && (obj.type = contractTypeToJSON(message.type));
    message.height !== undefined && (obj.height = Math.round(message.height));
    message.duration !== undefined && (obj.duration = Math.round(message.duration));
    message.rate !== undefined && (obj.rate = message.rate ? Coin.toJSON(message.rate) : undefined);
    message.deposit !== undefined && (obj.deposit = message.deposit);
    message.paid !== undefined && (obj.paid = message.paid);
    message.nonce !== undefined && (obj.nonce = Math.round(message.nonce));
    message.settlementHeight !== undefined && (obj.settlementHeight = Math.round(message.settlementHeight));
    message.id !== undefined && (obj.id = Math.round(message.id));
    message.settlementDuration !== undefined && (obj.settlementDuration = Math.round(message.settlementDuration));
    message.authorization !== undefined && (obj.authorization = contractAuthorizationToJSON(message.authorization));
    message.queriesPerMinute !== undefined && (obj.queriesPerMinute = Math.round(message.queriesPerMinute));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Contract>, I>>(object: I): Contract {
    const message = createBaseContract();
    message.provider = object.provider ?? new Uint8Array();
    message.service = object.service ?? 0;
    message.client = object.client ?? new Uint8Array();
    message.delegate = object.delegate ?? new Uint8Array();
    message.type = object.type ?? 0;
    message.height = object.height ?? 0;
    message.duration = object.duration ?? 0;
    message.rate = (object.rate !== undefined && object.rate !== null) ? Coin.fromPartial(object.rate) : undefined;
    message.deposit = object.deposit ?? "";
    message.paid = object.paid ?? "";
    message.nonce = object.nonce ?? 0;
    message.settlementHeight = object.settlementHeight ?? 0;
    message.id = object.id ?? 0;
    message.settlementDuration = object.settlementDuration ?? 0;
    message.authorization = object.authorization ?? 0;
    message.queriesPerMinute = object.queriesPerMinute ?? 0;
    return message;
  },
};

function createBaseContractSet(): ContractSet {
  return { contractIds: [] };
}

export const ContractSet = {
  encode(message: ContractSet, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    writer.uint32(10).fork();
    for (const v of message.contractIds) {
      writer.uint64(v);
    }
    writer.ldelim();
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ContractSet {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseContractSet();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.contractIds.push(longToNumber(reader.uint64() as Long));
            }
          } else {
            message.contractIds.push(longToNumber(reader.uint64() as Long));
          }
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ContractSet {
    return { contractIds: Array.isArray(object?.contractIds) ? object.contractIds.map((e: any) => Number(e)) : [] };
  },

  toJSON(message: ContractSet): unknown {
    const obj: any = {};
    if (message.contractIds) {
      obj.contractIds = message.contractIds.map((e) => Math.round(e));
    } else {
      obj.contractIds = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ContractSet>, I>>(object: I): ContractSet {
    const message = createBaseContractSet();
    message.contractIds = object.contractIds?.map((e) => e) || [];
    return message;
  },
};

function createBaseContractExpirationSet(): ContractExpirationSet {
  return { height: 0, contractSet: undefined };
}

export const ContractExpirationSet = {
  encode(message: ContractExpirationSet, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.height !== 0) {
      writer.uint32(8).int64(message.height);
    }
    if (message.contractSet !== undefined) {
      ContractSet.encode(message.contractSet, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ContractExpirationSet {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseContractExpirationSet();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.height = longToNumber(reader.int64() as Long);
          break;
        case 2:
          message.contractSet = ContractSet.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ContractExpirationSet {
    return {
      height: isSet(object.height) ? Number(object.height) : 0,
      contractSet: isSet(object.contractSet) ? ContractSet.fromJSON(object.contractSet) : undefined,
    };
  },

  toJSON(message: ContractExpirationSet): unknown {
    const obj: any = {};
    message.height !== undefined && (obj.height = Math.round(message.height));
    message.contractSet !== undefined
      && (obj.contractSet = message.contractSet ? ContractSet.toJSON(message.contractSet) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ContractExpirationSet>, I>>(object: I): ContractExpirationSet {
    const message = createBaseContractExpirationSet();
    message.height = object.height ?? 0;
    message.contractSet = (object.contractSet !== undefined && object.contractSet !== null)
      ? ContractSet.fromPartial(object.contractSet)
      : undefined;
    return message;
  },
};

function createBaseUserContractSet(): UserContractSet {
  return { user: new Uint8Array(), contractSet: undefined };
}

export const UserContractSet = {
  encode(message: UserContractSet, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.user.length !== 0) {
      writer.uint32(10).bytes(message.user);
    }
    if (message.contractSet !== undefined) {
      ContractSet.encode(message.contractSet, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UserContractSet {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUserContractSet();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.user = reader.bytes();
          break;
        case 2:
          message.contractSet = ContractSet.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UserContractSet {
    return {
      user: isSet(object.user) ? bytesFromBase64(object.user) : new Uint8Array(),
      contractSet: isSet(object.contractSet) ? ContractSet.fromJSON(object.contractSet) : undefined,
    };
  },

  toJSON(message: UserContractSet): unknown {
    const obj: any = {};
    message.user !== undefined
      && (obj.user = base64FromBytes(message.user !== undefined ? message.user : new Uint8Array()));
    message.contractSet !== undefined
      && (obj.contractSet = message.contractSet ? ContractSet.toJSON(message.contractSet) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserContractSet>, I>>(object: I): UserContractSet {
    const message = createBaseUserContractSet();
    message.user = object.user ?? new Uint8Array();
    message.contractSet = (object.contractSet !== undefined && object.contractSet !== null)
      ? ContractSet.fromPartial(object.contractSet)
      : undefined;
    return message;
  },
};

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

function bytesFromBase64(b64: string): Uint8Array {
  if (globalThis.Buffer) {
    return Uint8Array.from(globalThis.Buffer.from(b64, "base64"));
  } else {
    const bin = globalThis.atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; ++i) {
      arr[i] = bin.charCodeAt(i);
    }
    return arr;
  }
}

function base64FromBytes(arr: Uint8Array): string {
  if (globalThis.Buffer) {
    return globalThis.Buffer.from(arr).toString("base64");
  } else {
    const bin: string[] = [];
    arr.forEach((byte) => {
      bin.push(String.fromCharCode(byte));
    });
    return globalThis.btoa(bin.join(""));
  }
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
