/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";

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

export interface Provider {
  pubKey: string;
  service: number;
  metadataUri: string;
  metadataNonce: number;
  status: ProviderStatus;
  minContractDuration: number;
  maxContractDuration: number;
  subscriptionRate: number;
  payAsYouGoRate: number;
  bond: string;
  lastUpdate: number;
  settlementDuration: number;
}

export interface Contract {
  provider: string;
  service: number;
  client: string;
  delegate: string;
  type: ContractType;
  height: number;
  duration: number;
  rate: number;
  deposit: string;
  paid: string;
  nonce: number;
  settlementHeight: number;
  id: number;
  settlementDuration: number;
}

export interface ContractSet {
  contractIds: number[];
}

export interface ContractExpirationSet {
  height: number;
  contractSet: ContractSet | undefined;
}

export interface UserContractSet {
  user: string;
  contractSet: ContractSet | undefined;
}

function createBaseProvider(): Provider {
  return {
    pubKey: "",
    service: 0,
    metadataUri: "",
    metadataNonce: 0,
    status: 0,
    minContractDuration: 0,
    maxContractDuration: 0,
    subscriptionRate: 0,
    payAsYouGoRate: 0,
    bond: "",
    lastUpdate: 0,
    settlementDuration: 0,
  };
}

export const Provider = {
  encode(message: Provider, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pubKey !== "") {
      writer.uint32(10).string(message.pubKey);
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
          message.pubKey = reader.string();
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
      pubKey: isSet(object.pubKey) ? String(object.pubKey) : "",
      service: isSet(object.service) ? Number(object.service) : 0,
      metadataUri: isSet(object.metadataUri) ? String(object.metadataUri) : "",
      metadataNonce: isSet(object.metadataNonce) ? Number(object.metadataNonce) : 0,
      status: isSet(object.status) ? providerStatusFromJSON(object.status) : 0,
      minContractDuration: isSet(object.minContractDuration) ? Number(object.minContractDuration) : 0,
      maxContractDuration: isSet(object.maxContractDuration) ? Number(object.maxContractDuration) : 0,
      subscriptionRate: isSet(object.subscriptionRate) ? Number(object.subscriptionRate) : 0,
      payAsYouGoRate: isSet(object.payAsYouGoRate) ? Number(object.payAsYouGoRate) : 0,
      bond: isSet(object.bond) ? String(object.bond) : "",
      lastUpdate: isSet(object.lastUpdate) ? Number(object.lastUpdate) : 0,
      settlementDuration: isSet(object.settlementDuration) ? Number(object.settlementDuration) : 0,
    };
  },

  toJSON(message: Provider): unknown {
    const obj: any = {};
    message.pubKey !== undefined && (obj.pubKey = message.pubKey);
    message.service !== undefined && (obj.service = Math.round(message.service));
    message.metadataUri !== undefined && (obj.metadataUri = message.metadataUri);
    message.metadataNonce !== undefined && (obj.metadataNonce = Math.round(message.metadataNonce));
    message.status !== undefined && (obj.status = providerStatusToJSON(message.status));
    message.minContractDuration !== undefined && (obj.minContractDuration = Math.round(message.minContractDuration));
    message.maxContractDuration !== undefined && (obj.maxContractDuration = Math.round(message.maxContractDuration));
    message.subscriptionRate !== undefined && (obj.subscriptionRate = Math.round(message.subscriptionRate));
    message.payAsYouGoRate !== undefined && (obj.payAsYouGoRate = Math.round(message.payAsYouGoRate));
    message.bond !== undefined && (obj.bond = message.bond);
    message.lastUpdate !== undefined && (obj.lastUpdate = Math.round(message.lastUpdate));
    message.settlementDuration !== undefined && (obj.settlementDuration = Math.round(message.settlementDuration));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Provider>, I>>(object: I): Provider {
    const message = createBaseProvider();
    message.pubKey = object.pubKey ?? "";
    message.service = object.service ?? 0;
    message.metadataUri = object.metadataUri ?? "";
    message.metadataNonce = object.metadataNonce ?? 0;
    message.status = object.status ?? 0;
    message.minContractDuration = object.minContractDuration ?? 0;
    message.maxContractDuration = object.maxContractDuration ?? 0;
    message.subscriptionRate = object.subscriptionRate ?? 0;
    message.payAsYouGoRate = object.payAsYouGoRate ?? 0;
    message.bond = object.bond ?? "";
    message.lastUpdate = object.lastUpdate ?? 0;
    message.settlementDuration = object.settlementDuration ?? 0;
    return message;
  },
};

function createBaseContract(): Contract {
  return {
    provider: "",
    service: 0,
    client: "",
    delegate: "",
    type: 0,
    height: 0,
    duration: 0,
    rate: 0,
    deposit: "",
    paid: "",
    nonce: 0,
    settlementHeight: 0,
    id: 0,
    settlementDuration: 0,
  };
}

export const Contract = {
  encode(message: Contract, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.provider !== "") {
      writer.uint32(10).string(message.provider);
    }
    if (message.service !== 0) {
      writer.uint32(16).int32(message.service);
    }
    if (message.client !== "") {
      writer.uint32(26).string(message.client);
    }
    if (message.delegate !== "") {
      writer.uint32(34).string(message.delegate);
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
    if (message.rate !== 0) {
      writer.uint32(64).int64(message.rate);
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
          message.provider = reader.string();
          break;
        case 2:
          message.service = reader.int32();
          break;
        case 3:
          message.client = reader.string();
          break;
        case 4:
          message.delegate = reader.string();
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
          message.rate = longToNumber(reader.int64() as Long);
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
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Contract {
    return {
      provider: isSet(object.provider) ? String(object.provider) : "",
      service: isSet(object.service) ? Number(object.service) : 0,
      client: isSet(object.client) ? String(object.client) : "",
      delegate: isSet(object.delegate) ? String(object.delegate) : "",
      type: isSet(object.type) ? contractTypeFromJSON(object.type) : 0,
      height: isSet(object.height) ? Number(object.height) : 0,
      duration: isSet(object.duration) ? Number(object.duration) : 0,
      rate: isSet(object.rate) ? Number(object.rate) : 0,
      deposit: isSet(object.deposit) ? String(object.deposit) : "",
      paid: isSet(object.paid) ? String(object.paid) : "",
      nonce: isSet(object.nonce) ? Number(object.nonce) : 0,
      settlementHeight: isSet(object.settlementHeight) ? Number(object.settlementHeight) : 0,
      id: isSet(object.id) ? Number(object.id) : 0,
      settlementDuration: isSet(object.settlementDuration) ? Number(object.settlementDuration) : 0,
    };
  },

  toJSON(message: Contract): unknown {
    const obj: any = {};
    message.provider !== undefined && (obj.provider = message.provider);
    message.service !== undefined && (obj.service = Math.round(message.service));
    message.client !== undefined && (obj.client = message.client);
    message.delegate !== undefined && (obj.delegate = message.delegate);
    message.type !== undefined && (obj.type = contractTypeToJSON(message.type));
    message.height !== undefined && (obj.height = Math.round(message.height));
    message.duration !== undefined && (obj.duration = Math.round(message.duration));
    message.rate !== undefined && (obj.rate = Math.round(message.rate));
    message.deposit !== undefined && (obj.deposit = message.deposit);
    message.paid !== undefined && (obj.paid = message.paid);
    message.nonce !== undefined && (obj.nonce = Math.round(message.nonce));
    message.settlementHeight !== undefined && (obj.settlementHeight = Math.round(message.settlementHeight));
    message.id !== undefined && (obj.id = Math.round(message.id));
    message.settlementDuration !== undefined && (obj.settlementDuration = Math.round(message.settlementDuration));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Contract>, I>>(object: I): Contract {
    const message = createBaseContract();
    message.provider = object.provider ?? "";
    message.service = object.service ?? 0;
    message.client = object.client ?? "";
    message.delegate = object.delegate ?? "";
    message.type = object.type ?? 0;
    message.height = object.height ?? 0;
    message.duration = object.duration ?? 0;
    message.rate = object.rate ?? 0;
    message.deposit = object.deposit ?? "";
    message.paid = object.paid ?? "";
    message.nonce = object.nonce ?? 0;
    message.settlementHeight = object.settlementHeight ?? 0;
    message.id = object.id ?? 0;
    message.settlementDuration = object.settlementDuration ?? 0;
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
  return { user: "", contractSet: undefined };
}

export const UserContractSet = {
  encode(message: UserContractSet, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.user !== "") {
      writer.uint32(10).string(message.user);
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
          message.user = reader.string();
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
      user: isSet(object.user) ? String(object.user) : "",
      contractSet: isSet(object.contractSet) ? ContractSet.fromJSON(object.contractSet) : undefined,
    };
  },

  toJSON(message: UserContractSet): unknown {
    const obj: any = {};
    message.user !== undefined && (obj.user = message.user);
    message.contractSet !== undefined
      && (obj.contractSet = message.contractSet ? ContractSet.toJSON(message.contractSet) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UserContractSet>, I>>(object: I): UserContractSet {
    const message = createBaseUserContractSet();
    message.user = object.user ?? "";
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
