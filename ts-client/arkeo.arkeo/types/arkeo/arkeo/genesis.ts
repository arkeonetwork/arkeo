/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Contract, ContractExpirationSet, Provider, UserContractSet } from "./keeper";
import { Params } from "./params";

export const protobufPackage = "arkeo.arkeo";

/** GenesisState defines the arkeo module's genesis state. */
export interface GenesisState {
  params: Params | undefined;
  providers: Provider[];
  contracts: Contract[];
  nextContractId: number;
  contractExpirationSets: ContractExpirationSet[];
  userContractSets: UserContractSet[];
  /** this line is used by starport scaffolding # genesis/proto/state */
  version: number;
}

function createBaseGenesisState(): GenesisState {
  return {
    params: undefined,
    providers: [],
    contracts: [],
    nextContractId: 0,
    contractExpirationSets: [],
    userContractSets: [],
    version: 0,
  };
}

export const GenesisState = {
  encode(message: GenesisState, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.providers) {
      Provider.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.contracts) {
      Contract.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    if (message.nextContractId !== 0) {
      writer.uint32(32).uint64(message.nextContractId);
    }
    for (const v of message.contractExpirationSets) {
      ContractExpirationSet.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.userContractSets) {
      UserContractSet.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    if (message.version !== 0) {
      writer.uint32(56).int64(message.version);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGenesisState();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        case 2:
          message.providers.push(Provider.decode(reader, reader.uint32()));
          break;
        case 3:
          message.contracts.push(Contract.decode(reader, reader.uint32()));
          break;
        case 4:
          message.nextContractId = longToNumber(reader.uint64() as Long);
          break;
        case 5:
          message.contractExpirationSets.push(ContractExpirationSet.decode(reader, reader.uint32()));
          break;
        case 6:
          message.userContractSets.push(UserContractSet.decode(reader, reader.uint32()));
          break;
        case 7:
          message.version = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GenesisState {
    return {
      params: isSet(object.params) ? Params.fromJSON(object.params) : undefined,
      providers: Array.isArray(object?.providers) ? object.providers.map((e: any) => Provider.fromJSON(e)) : [],
      contracts: Array.isArray(object?.contracts) ? object.contracts.map((e: any) => Contract.fromJSON(e)) : [],
      nextContractId: isSet(object.nextContractId) ? Number(object.nextContractId) : 0,
      contractExpirationSets: Array.isArray(object?.contractExpirationSets)
        ? object.contractExpirationSets.map((e: any) => ContractExpirationSet.fromJSON(e))
        : [],
      userContractSets: Array.isArray(object?.userContractSets)
        ? object.userContractSets.map((e: any) => UserContractSet.fromJSON(e))
        : [],
      version: isSet(object.version) ? Number(object.version) : 0,
    };
  },

  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    if (message.providers) {
      obj.providers = message.providers.map((e) => e ? Provider.toJSON(e) : undefined);
    } else {
      obj.providers = [];
    }
    if (message.contracts) {
      obj.contracts = message.contracts.map((e) => e ? Contract.toJSON(e) : undefined);
    } else {
      obj.contracts = [];
    }
    message.nextContractId !== undefined && (obj.nextContractId = Math.round(message.nextContractId));
    if (message.contractExpirationSets) {
      obj.contractExpirationSets = message.contractExpirationSets.map((e) =>
        e ? ContractExpirationSet.toJSON(e) : undefined
      );
    } else {
      obj.contractExpirationSets = [];
    }
    if (message.userContractSets) {
      obj.userContractSets = message.userContractSets.map((e) => e ? UserContractSet.toJSON(e) : undefined);
    } else {
      obj.userContractSets = [];
    }
    message.version !== undefined && (obj.version = Math.round(message.version));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GenesisState>, I>>(object: I): GenesisState {
    const message = createBaseGenesisState();
    message.params = (object.params !== undefined && object.params !== null)
      ? Params.fromPartial(object.params)
      : undefined;
    message.providers = object.providers?.map((e) => Provider.fromPartial(e)) || [];
    message.contracts = object.contracts?.map((e) => Contract.fromPartial(e)) || [];
    message.nextContractId = object.nextContractId ?? 0;
    message.contractExpirationSets = object.contractExpirationSets?.map((e) => ContractExpirationSet.fromPartial(e))
      || [];
    message.userContractSets = object.userContractSets?.map((e) => UserContractSet.fromPartial(e)) || [];
    message.version = object.version ?? 0;
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
