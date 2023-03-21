/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { PageRequest, PageResponse } from "../../cosmos/base/query/v1beta1/pagination";
import { Contract, Provider } from "./keeper";
import { Params } from "./params";

export const protobufPackage = "arkeo.arkeo";

/** QueryParamsRequest is request type for the Query/Params RPC method. */
export interface QueryParamsRequest {
}

/** QueryParamsResponse is response type for the Query/Params RPC method. */
export interface QueryParamsResponse {
  /** params holds all the parameters of this module. */
  params: Params | undefined;
}

export interface QueryFetchProviderRequest {
  pubkey: string;
  service: string;
}

export interface QueryFetchProviderResponse {
  provider: Provider | undefined;
}

export interface QueryAllProviderRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllProviderResponse {
  provider: Provider[];
  pagination: PageResponse | undefined;
}

export interface QueryFetchContractRequest {
  contractId: number;
}

export interface QueryFetchContractResponse {
  contract: Contract | undefined;
}

export interface QueryAllContractRequest {
  pagination: PageRequest | undefined;
}

export interface QueryAllContractResponse {
  contract: Contract[];
  pagination: PageResponse | undefined;
}

/** this line is used by starport scaffolding # 3 */
export interface QueryActiveContractRequest {
  spender: string;
  provider: string;
  service: string;
}

export interface QueryActiveContractResponse {
  contract: Contract | undefined;
}

function createBaseQueryParamsRequest(): QueryParamsRequest {
  return {};
}

export const QueryParamsRequest = {
  encode(_: QueryParamsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): QueryParamsRequest {
    return {};
  },

  toJSON(_: QueryParamsRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryParamsRequest>, I>>(_: I): QueryParamsRequest {
    const message = createBaseQueryParamsRequest();
    return message;
  },
};

function createBaseQueryParamsResponse(): QueryParamsResponse {
  return { params: undefined };
}

export const QueryParamsResponse = {
  encode(message: QueryParamsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryParamsResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryParamsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.params = Params.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryParamsResponse {
    return { params: isSet(object.params) ? Params.fromJSON(object.params) : undefined };
  },

  toJSON(message: QueryParamsResponse): unknown {
    const obj: any = {};
    message.params !== undefined && (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryParamsResponse>, I>>(object: I): QueryParamsResponse {
    const message = createBaseQueryParamsResponse();
    message.params = (object.params !== undefined && object.params !== null)
      ? Params.fromPartial(object.params)
      : undefined;
    return message;
  },
};

function createBaseQueryFetchProviderRequest(): QueryFetchProviderRequest {
  return { pubkey: "", service: "" };
}

export const QueryFetchProviderRequest = {
  encode(message: QueryFetchProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pubkey !== "") {
      writer.uint32(10).string(message.pubkey);
    }
    if (message.service !== "") {
      writer.uint32(18).string(message.service);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFetchProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFetchProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pubkey = reader.string();
          break;
        case 2:
          message.service = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryFetchProviderRequest {
    return {
      pubkey: isSet(object.pubkey) ? String(object.pubkey) : "",
      service: isSet(object.service) ? String(object.service) : "",
    };
  },

  toJSON(message: QueryFetchProviderRequest): unknown {
    const obj: any = {};
    message.pubkey !== undefined && (obj.pubkey = message.pubkey);
    message.service !== undefined && (obj.service = message.service);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryFetchProviderRequest>, I>>(object: I): QueryFetchProviderRequest {
    const message = createBaseQueryFetchProviderRequest();
    message.pubkey = object.pubkey ?? "";
    message.service = object.service ?? "";
    return message;
  },
};

function createBaseQueryFetchProviderResponse(): QueryFetchProviderResponse {
  return { provider: undefined };
}

export const QueryFetchProviderResponse = {
  encode(message: QueryFetchProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.provider !== undefined) {
      Provider.encode(message.provider, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFetchProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFetchProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.provider = Provider.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryFetchProviderResponse {
    return { provider: isSet(object.provider) ? Provider.fromJSON(object.provider) : undefined };
  },

  toJSON(message: QueryFetchProviderResponse): unknown {
    const obj: any = {};
    message.provider !== undefined && (obj.provider = message.provider ? Provider.toJSON(message.provider) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryFetchProviderResponse>, I>>(object: I): QueryFetchProviderResponse {
    const message = createBaseQueryFetchProviderResponse();
    message.provider = (object.provider !== undefined && object.provider !== null)
      ? Provider.fromPartial(object.provider)
      : undefined;
    return message;
  },
};

function createBaseQueryAllProviderRequest(): QueryAllProviderRequest {
  return { pagination: undefined };
}

export const QueryAllProviderRequest = {
  encode(message: QueryAllProviderRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllProviderRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllProviderRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllProviderRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryAllProviderRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllProviderRequest>, I>>(object: I): QueryAllProviderRequest {
    const message = createBaseQueryAllProviderRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllProviderResponse(): QueryAllProviderResponse {
  return { provider: [], pagination: undefined };
}

export const QueryAllProviderResponse = {
  encode(message: QueryAllProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.provider) {
      Provider.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllProviderResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.provider.push(Provider.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllProviderResponse {
    return {
      provider: Array.isArray(object?.provider) ? object.provider.map((e: any) => Provider.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryAllProviderResponse): unknown {
    const obj: any = {};
    if (message.provider) {
      obj.provider = message.provider.map((e) => e ? Provider.toJSON(e) : undefined);
    } else {
      obj.provider = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllProviderResponse>, I>>(object: I): QueryAllProviderResponse {
    const message = createBaseQueryAllProviderResponse();
    message.provider = object.provider?.map((e) => Provider.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryFetchContractRequest(): QueryFetchContractRequest {
  return { contractId: 0 };
}

export const QueryFetchContractRequest = {
  encode(message: QueryFetchContractRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.contractId !== 0) {
      writer.uint32(8).uint64(message.contractId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFetchContractRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFetchContractRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.contractId = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryFetchContractRequest {
    return { contractId: isSet(object.contractId) ? Number(object.contractId) : 0 };
  },

  toJSON(message: QueryFetchContractRequest): unknown {
    const obj: any = {};
    message.contractId !== undefined && (obj.contractId = Math.round(message.contractId));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryFetchContractRequest>, I>>(object: I): QueryFetchContractRequest {
    const message = createBaseQueryFetchContractRequest();
    message.contractId = object.contractId ?? 0;
    return message;
  },
};

function createBaseQueryFetchContractResponse(): QueryFetchContractResponse {
  return { contract: undefined };
}

export const QueryFetchContractResponse = {
  encode(message: QueryFetchContractResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.contract !== undefined) {
      Contract.encode(message.contract, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryFetchContractResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryFetchContractResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.contract = Contract.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryFetchContractResponse {
    return { contract: isSet(object.contract) ? Contract.fromJSON(object.contract) : undefined };
  },

  toJSON(message: QueryFetchContractResponse): unknown {
    const obj: any = {};
    message.contract !== undefined && (obj.contract = message.contract ? Contract.toJSON(message.contract) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryFetchContractResponse>, I>>(object: I): QueryFetchContractResponse {
    const message = createBaseQueryFetchContractResponse();
    message.contract = (object.contract !== undefined && object.contract !== null)
      ? Contract.fromPartial(object.contract)
      : undefined;
    return message;
  },
};

function createBaseQueryAllContractRequest(): QueryAllContractRequest {
  return { pagination: undefined };
}

export const QueryAllContractRequest = {
  encode(message: QueryAllContractRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.pagination !== undefined) {
      PageRequest.encode(message.pagination, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllContractRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllContractRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.pagination = PageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllContractRequest {
    return { pagination: isSet(object.pagination) ? PageRequest.fromJSON(object.pagination) : undefined };
  },

  toJSON(message: QueryAllContractRequest): unknown {
    const obj: any = {};
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageRequest.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllContractRequest>, I>>(object: I): QueryAllContractRequest {
    const message = createBaseQueryAllContractRequest();
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageRequest.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryAllContractResponse(): QueryAllContractResponse {
  return { contract: [], pagination: undefined };
}

export const QueryAllContractResponse = {
  encode(message: QueryAllContractResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.contract) {
      Contract.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.pagination !== undefined) {
      PageResponse.encode(message.pagination, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryAllContractResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryAllContractResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.contract.push(Contract.decode(reader, reader.uint32()));
          break;
        case 2:
          message.pagination = PageResponse.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryAllContractResponse {
    return {
      contract: Array.isArray(object?.contract) ? object.contract.map((e: any) => Contract.fromJSON(e)) : [],
      pagination: isSet(object.pagination) ? PageResponse.fromJSON(object.pagination) : undefined,
    };
  },

  toJSON(message: QueryAllContractResponse): unknown {
    const obj: any = {};
    if (message.contract) {
      obj.contract = message.contract.map((e) => e ? Contract.toJSON(e) : undefined);
    } else {
      obj.contract = [];
    }
    message.pagination !== undefined
      && (obj.pagination = message.pagination ? PageResponse.toJSON(message.pagination) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryAllContractResponse>, I>>(object: I): QueryAllContractResponse {
    const message = createBaseQueryAllContractResponse();
    message.contract = object.contract?.map((e) => Contract.fromPartial(e)) || [];
    message.pagination = (object.pagination !== undefined && object.pagination !== null)
      ? PageResponse.fromPartial(object.pagination)
      : undefined;
    return message;
  },
};

function createBaseQueryActiveContractRequest(): QueryActiveContractRequest {
  return { spender: "", provider: "", service: "" };
}

export const QueryActiveContractRequest = {
  encode(message: QueryActiveContractRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.spender !== "") {
      writer.uint32(10).string(message.spender);
    }
    if (message.provider !== "") {
      writer.uint32(18).string(message.provider);
    }
    if (message.service !== "") {
      writer.uint32(26).string(message.service);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryActiveContractRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryActiveContractRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.spender = reader.string();
          break;
        case 2:
          message.provider = reader.string();
          break;
        case 3:
          message.service = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryActiveContractRequest {
    return {
      spender: isSet(object.spender) ? String(object.spender) : "",
      provider: isSet(object.provider) ? String(object.provider) : "",
      service: isSet(object.service) ? String(object.service) : "",
    };
  },

  toJSON(message: QueryActiveContractRequest): unknown {
    const obj: any = {};
    message.spender !== undefined && (obj.spender = message.spender);
    message.provider !== undefined && (obj.provider = message.provider);
    message.service !== undefined && (obj.service = message.service);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryActiveContractRequest>, I>>(object: I): QueryActiveContractRequest {
    const message = createBaseQueryActiveContractRequest();
    message.spender = object.spender ?? "";
    message.provider = object.provider ?? "";
    message.service = object.service ?? "";
    return message;
  },
};

function createBaseQueryActiveContractResponse(): QueryActiveContractResponse {
  return { contract: undefined };
}

export const QueryActiveContractResponse = {
  encode(message: QueryActiveContractResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.contract !== undefined) {
      Contract.encode(message.contract, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): QueryActiveContractResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseQueryActiveContractResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.contract = Contract.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): QueryActiveContractResponse {
    return { contract: isSet(object.contract) ? Contract.fromJSON(object.contract) : undefined };
  },

  toJSON(message: QueryActiveContractResponse): unknown {
    const obj: any = {};
    message.contract !== undefined && (obj.contract = message.contract ? Contract.toJSON(message.contract) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<QueryActiveContractResponse>, I>>(object: I): QueryActiveContractResponse {
    const message = createBaseQueryActiveContractResponse();
    message.contract = (object.contract !== undefined && object.contract !== null)
      ? Contract.fromPartial(object.contract)
      : undefined;
    return message;
  },
};

/** Query defines the gRPC querier service. */
export interface Query {
  /** Parameters queries the parameters of the module. */
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse>;
  FetchProvider(request: QueryFetchProviderRequest): Promise<QueryFetchProviderResponse>;
  ProviderAll(request: QueryAllProviderRequest): Promise<QueryAllProviderResponse>;
  FetchContract(request: QueryFetchContractRequest): Promise<QueryFetchContractResponse>;
  ContractAll(request: QueryAllContractRequest): Promise<QueryAllContractResponse>;
  /** Queries an active contract by spender, provider and service. */
  ActiveContract(request: QueryActiveContractRequest): Promise<QueryActiveContractResponse>;
}

export class QueryClientImpl implements Query {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Params = this.Params.bind(this);
    this.FetchProvider = this.FetchProvider.bind(this);
    this.ProviderAll = this.ProviderAll.bind(this);
    this.FetchContract = this.FetchContract.bind(this);
    this.ContractAll = this.ContractAll.bind(this);
    this.ActiveContract = this.ActiveContract.bind(this);
  }
  Params(request: QueryParamsRequest): Promise<QueryParamsResponse> {
    const data = QueryParamsRequest.encode(request).finish();
    const promise = this.rpc.request("arkeo.arkeo.Query", "Params", data);
    return promise.then((data) => QueryParamsResponse.decode(new _m0.Reader(data)));
  }

  FetchProvider(request: QueryFetchProviderRequest): Promise<QueryFetchProviderResponse> {
    const data = QueryFetchProviderRequest.encode(request).finish();
    const promise = this.rpc.request("arkeo.arkeo.Query", "FetchProvider", data);
    return promise.then((data) => QueryFetchProviderResponse.decode(new _m0.Reader(data)));
  }

  ProviderAll(request: QueryAllProviderRequest): Promise<QueryAllProviderResponse> {
    const data = QueryAllProviderRequest.encode(request).finish();
    const promise = this.rpc.request("arkeo.arkeo.Query", "ProviderAll", data);
    return promise.then((data) => QueryAllProviderResponse.decode(new _m0.Reader(data)));
  }

  FetchContract(request: QueryFetchContractRequest): Promise<QueryFetchContractResponse> {
    const data = QueryFetchContractRequest.encode(request).finish();
    const promise = this.rpc.request("arkeo.arkeo.Query", "FetchContract", data);
    return promise.then((data) => QueryFetchContractResponse.decode(new _m0.Reader(data)));
  }

  ContractAll(request: QueryAllContractRequest): Promise<QueryAllContractResponse> {
    const data = QueryAllContractRequest.encode(request).finish();
    const promise = this.rpc.request("arkeo.arkeo.Query", "ContractAll", data);
    return promise.then((data) => QueryAllContractResponse.decode(new _m0.Reader(data)));
  }

  ActiveContract(request: QueryActiveContractRequest): Promise<QueryActiveContractResponse> {
    const data = QueryActiveContractRequest.encode(request).finish();
    const promise = this.rpc.request("arkeo.arkeo.Query", "ActiveContract", data);
    return promise.then((data) => QueryActiveContractResponse.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

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
