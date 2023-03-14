/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import {
  ContractType,
  contractTypeFromJSON,
  contractTypeToJSON,
  ProviderStatus,
  providerStatusFromJSON,
  providerStatusToJSON,
} from "./keeper";

export const protobufPackage = "arkeo.arkeo";

export interface MsgBondProvider {
  creator: string;
  pubKey: string;
  service: string;
  bond: string;
}

export interface MsgBondProviderResponse {
}

export interface MsgModProvider {
  creator: string;
  pubKey: string;
  service: string;
  metadataURI: string;
  metadataNonce: number;
  status: ProviderStatus;
  minContractDuration: number;
  maxContractDuration: number;
  subscriptionRate: number;
  payAsYouGoRate: number;
}

export interface MsgModProviderResponse {
}

export interface MsgOpenContract {
  creator: string;
  pubKey: string;
  service: string;
  client: string;
  delegate: string;
  contractType: ContractType;
  duration: number;
  rate: number;
  deposit: string;
}

export interface MsgOpenContractResponse {
}

export interface MsgCloseContract {
  creator: string;
  pubKey: string;
  service: string;
  client: string;
  delegate: string;
}

export interface MsgCloseContractResponse {
}

export interface MsgClaimContractIncome {
  creator: string;
  pubKey: string;
  service: string;
  spender: string;
  signature: Uint8Array;
  nonce: number;
  height: number;
}

export interface MsgClaimContractIncomeResponse {
}

function createBaseMsgBondProvider(): MsgBondProvider {
  return { creator: "", pubKey: "", service: "", bond: "" };
}

export const MsgBondProvider = {
  encode(message: MsgBondProvider, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.pubKey !== "") {
      writer.uint32(18).string(message.pubKey);
    }
    if (message.service !== "") {
      writer.uint32(26).string(message.service);
    }
    if (message.bond !== "") {
      writer.uint32(34).string(message.bond);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgBondProvider {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgBondProvider();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.pubKey = reader.string();
          break;
        case 3:
          message.service = reader.string();
          break;
        case 4:
          message.bond = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgBondProvider {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      pubKey: isSet(object.pubKey) ? String(object.pubKey) : "",
      service: isSet(object.service) ? String(object.service) : "",
      bond: isSet(object.bond) ? String(object.bond) : "",
    };
  },

  toJSON(message: MsgBondProvider): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.pubKey !== undefined && (obj.pubKey = message.pubKey);
    message.service !== undefined && (obj.service = message.service);
    message.bond !== undefined && (obj.bond = message.bond);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgBondProvider>, I>>(object: I): MsgBondProvider {
    const message = createBaseMsgBondProvider();
    message.creator = object.creator ?? "";
    message.pubKey = object.pubKey ?? "";
    message.service = object.service ?? "";
    message.bond = object.bond ?? "";
    return message;
  },
};

function createBaseMsgBondProviderResponse(): MsgBondProviderResponse {
  return {};
}

export const MsgBondProviderResponse = {
  encode(_: MsgBondProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgBondProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgBondProviderResponse();
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

  fromJSON(_: any): MsgBondProviderResponse {
    return {};
  },

  toJSON(_: MsgBondProviderResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgBondProviderResponse>, I>>(_: I): MsgBondProviderResponse {
    const message = createBaseMsgBondProviderResponse();
    return message;
  },
};

function createBaseMsgModProvider(): MsgModProvider {
  return {
    creator: "",
    pubKey: "",
    service: "",
    metadataURI: "",
    metadataNonce: 0,
    status: 0,
    minContractDuration: 0,
    maxContractDuration: 0,
    subscriptionRate: 0,
    payAsYouGoRate: 0,
  };
}

export const MsgModProvider = {
  encode(message: MsgModProvider, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.pubKey !== "") {
      writer.uint32(18).string(message.pubKey);
    }
    if (message.service !== "") {
      writer.uint32(26).string(message.service);
    }
    if (message.metadataURI !== "") {
      writer.uint32(34).string(message.metadataURI);
    }
    if (message.metadataNonce !== 0) {
      writer.uint32(40).uint64(message.metadataNonce);
    }
    if (message.status !== 0) {
      writer.uint32(48).int32(message.status);
    }
    if (message.minContractDuration !== 0) {
      writer.uint32(56).int64(message.minContractDuration);
    }
    if (message.maxContractDuration !== 0) {
      writer.uint32(64).int64(message.maxContractDuration);
    }
    if (message.subscriptionRate !== 0) {
      writer.uint32(72).int64(message.subscriptionRate);
    }
    if (message.payAsYouGoRate !== 0) {
      writer.uint32(80).int64(message.payAsYouGoRate);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgModProvider {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgModProvider();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.pubKey = reader.string();
          break;
        case 3:
          message.service = reader.string();
          break;
        case 4:
          message.metadataURI = reader.string();
          break;
        case 5:
          message.metadataNonce = longToNumber(reader.uint64() as Long);
          break;
        case 6:
          message.status = reader.int32() as any;
          break;
        case 7:
          message.minContractDuration = longToNumber(reader.int64() as Long);
          break;
        case 8:
          message.maxContractDuration = longToNumber(reader.int64() as Long);
          break;
        case 9:
          message.subscriptionRate = longToNumber(reader.int64() as Long);
          break;
        case 10:
          message.payAsYouGoRate = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgModProvider {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      pubKey: isSet(object.pubKey) ? String(object.pubKey) : "",
      service: isSet(object.service) ? String(object.service) : "",
      metadataURI: isSet(object.metadataURI) ? String(object.metadataURI) : "",
      metadataNonce: isSet(object.metadataNonce) ? Number(object.metadataNonce) : 0,
      status: isSet(object.status) ? providerStatusFromJSON(object.status) : 0,
      minContractDuration: isSet(object.minContractDuration) ? Number(object.minContractDuration) : 0,
      maxContractDuration: isSet(object.maxContractDuration) ? Number(object.maxContractDuration) : 0,
      subscriptionRate: isSet(object.subscriptionRate) ? Number(object.subscriptionRate) : 0,
      payAsYouGoRate: isSet(object.payAsYouGoRate) ? Number(object.payAsYouGoRate) : 0,
    };
  },

  toJSON(message: MsgModProvider): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.pubKey !== undefined && (obj.pubKey = message.pubKey);
    message.service !== undefined && (obj.service = message.service);
    message.metadataURI !== undefined && (obj.metadataURI = message.metadataURI);
    message.metadataNonce !== undefined && (obj.metadataNonce = Math.round(message.metadataNonce));
    message.status !== undefined && (obj.status = providerStatusToJSON(message.status));
    message.minContractDuration !== undefined && (obj.minContractDuration = Math.round(message.minContractDuration));
    message.maxContractDuration !== undefined && (obj.maxContractDuration = Math.round(message.maxContractDuration));
    message.subscriptionRate !== undefined && (obj.subscriptionRate = Math.round(message.subscriptionRate));
    message.payAsYouGoRate !== undefined && (obj.payAsYouGoRate = Math.round(message.payAsYouGoRate));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgModProvider>, I>>(object: I): MsgModProvider {
    const message = createBaseMsgModProvider();
    message.creator = object.creator ?? "";
    message.pubKey = object.pubKey ?? "";
    message.service = object.service ?? "";
    message.metadataURI = object.metadataURI ?? "";
    message.metadataNonce = object.metadataNonce ?? 0;
    message.status = object.status ?? 0;
    message.minContractDuration = object.minContractDuration ?? 0;
    message.maxContractDuration = object.maxContractDuration ?? 0;
    message.subscriptionRate = object.subscriptionRate ?? 0;
    message.payAsYouGoRate = object.payAsYouGoRate ?? 0;
    return message;
  },
};

function createBaseMsgModProviderResponse(): MsgModProviderResponse {
  return {};
}

export const MsgModProviderResponse = {
  encode(_: MsgModProviderResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgModProviderResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgModProviderResponse();
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

  fromJSON(_: any): MsgModProviderResponse {
    return {};
  },

  toJSON(_: MsgModProviderResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgModProviderResponse>, I>>(_: I): MsgModProviderResponse {
    const message = createBaseMsgModProviderResponse();
    return message;
  },
};

function createBaseMsgOpenContract(): MsgOpenContract {
  return { creator: "", pubKey: "", service: "", client: "", delegate: "", contractType: 0, duration: 0, rate: 0, deposit: "" };
}

export const MsgOpenContract = {
  encode(message: MsgOpenContract, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.pubKey !== "") {
      writer.uint32(18).string(message.pubKey);
    }
    if (message.service !== "") {
      writer.uint32(26).string(message.service);
    }
    if (message.client !== "") {
      writer.uint32(34).string(message.client);
    }
    if (message.delegate !== "") {
      writer.uint32(42).string(message.delegate);
    }
    if (message.contractType !== 0) {
      writer.uint32(48).int32(message.contractType);
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
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgOpenContract {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgOpenContract();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.pubKey = reader.string();
          break;
        case 3:
          message.service = reader.string();
          break;
        case 4:
          message.client = reader.string();
          break;
        case 5:
          message.delegate = reader.string();
          break;
        case 6:
          message.contractType = reader.int32() as any;
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
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgOpenContract {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      pubKey: isSet(object.pubKey) ? String(object.pubKey) : "",
      service: isSet(object.service) ? String(object.service) : "",
      client: isSet(object.client) ? String(object.client) : "",
      delegate: isSet(object.delegate) ? String(object.delegate) : "",
      contractType: isSet(object.contractType) ? contractTypeFromJSON(object.contractType) : 0,
      duration: isSet(object.duration) ? Number(object.duration) : 0,
      rate: isSet(object.rate) ? Number(object.rate) : 0,
      deposit: isSet(object.deposit) ? String(object.deposit) : "",
    };
  },

  toJSON(message: MsgOpenContract): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.pubKey !== undefined && (obj.pubKey = message.pubKey);
    message.service !== undefined && (obj.service = message.service);
    message.client !== undefined && (obj.client = message.client);
    message.delegate !== undefined && (obj.delegate = message.delegate);
    message.contractType !== undefined && (obj.contractType = contractTypeToJSON(message.contractType));
    message.duration !== undefined && (obj.duration = Math.round(message.duration));
    message.rate !== undefined && (obj.rate = Math.round(message.rate));
    message.deposit !== undefined && (obj.deposit = message.deposit);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgOpenContract>, I>>(object: I): MsgOpenContract {
    const message = createBaseMsgOpenContract();
    message.creator = object.creator ?? "";
    message.pubKey = object.pubKey ?? "";
    message.service = object.service ?? "";
    message.client = object.client ?? "";
    message.delegate = object.delegate ?? "";
    message.contractType = object.contractType ?? 0;
    message.duration = object.duration ?? 0;
    message.rate = object.rate ?? 0;
    message.deposit = object.deposit ?? "";
    return message;
  },
};

function createBaseMsgOpenContractResponse(): MsgOpenContractResponse {
  return {};
}

export const MsgOpenContractResponse = {
  encode(_: MsgOpenContractResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgOpenContractResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgOpenContractResponse();
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

  fromJSON(_: any): MsgOpenContractResponse {
    return {};
  },

  toJSON(_: MsgOpenContractResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgOpenContractResponse>, I>>(_: I): MsgOpenContractResponse {
    const message = createBaseMsgOpenContractResponse();
    return message;
  },
};

function createBaseMsgCloseContract(): MsgCloseContract {
  return { creator: "", pubKey: "", service: "", client: "", delegate: "" };
}

export const MsgCloseContract = {
  encode(message: MsgCloseContract, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.pubKey !== "") {
      writer.uint32(18).string(message.pubKey);
    }
    if (message.service !== "") {
      writer.uint32(26).string(message.service);
    }
    if (message.client !== "") {
      writer.uint32(34).string(message.client);
    }
    if (message.delegate !== "") {
      writer.uint32(42).string(message.delegate);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCloseContract {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCloseContract();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.pubKey = reader.string();
          break;
        case 3:
          message.service = reader.string();
          break;
        case 4:
          message.client = reader.string();
          break;
        case 5:
          message.delegate = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgCloseContract {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      pubKey: isSet(object.pubKey) ? String(object.pubKey) : "",
      service: isSet(object.service) ? String(object.service) : "",
      client: isSet(object.client) ? String(object.client) : "",
      delegate: isSet(object.delegate) ? String(object.delegate) : "",
    };
  },

  toJSON(message: MsgCloseContract): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.pubKey !== undefined && (obj.pubKey = message.pubKey);
    message.service !== undefined && (obj.service = message.service);
    message.client !== undefined && (obj.client = message.client);
    message.delegate !== undefined && (obj.delegate = message.delegate);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCloseContract>, I>>(object: I): MsgCloseContract {
    const message = createBaseMsgCloseContract();
    message.creator = object.creator ?? "";
    message.pubKey = object.pubKey ?? "";
    message.service = object.service ?? "";
    message.client = object.client ?? "";
    message.delegate = object.delegate ?? "";
    return message;
  },
};

function createBaseMsgCloseContractResponse(): MsgCloseContractResponse {
  return {};
}

export const MsgCloseContractResponse = {
  encode(_: MsgCloseContractResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgCloseContractResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgCloseContractResponse();
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

  fromJSON(_: any): MsgCloseContractResponse {
    return {};
  },

  toJSON(_: MsgCloseContractResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgCloseContractResponse>, I>>(_: I): MsgCloseContractResponse {
    const message = createBaseMsgCloseContractResponse();
    return message;
  },
};

function createBaseMsgClaimContractIncome(): MsgClaimContractIncome {
  return { creator: "", pubKey: "", service: "", spender: "", signature: new Uint8Array(), nonce: 0, height: 0 };
}

export const MsgClaimContractIncome = {
  encode(message: MsgClaimContractIncome, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.pubKey !== "") {
      writer.uint32(18).string(message.pubKey);
    }
    if (message.service !== "") {
      writer.uint32(26).string(message.service);
    }
    if (message.spender !== "") {
      writer.uint32(34).string(message.spender);
    }
    if (message.signature.length !== 0) {
      writer.uint32(42).bytes(message.signature);
    }
    if (message.nonce !== 0) {
      writer.uint32(48).int64(message.nonce);
    }
    if (message.height !== 0) {
      writer.uint32(56).int64(message.height);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgClaimContractIncome {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgClaimContractIncome();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.pubKey = reader.string();
          break;
        case 3:
          message.service = reader.string();
          break;
        case 4:
          message.spender = reader.string();
          break;
        case 5:
          message.signature = reader.bytes();
          break;
        case 6:
          message.nonce = longToNumber(reader.int64() as Long);
          break;
        case 7:
          message.height = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgClaimContractIncome {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      pubKey: isSet(object.pubKey) ? String(object.pubKey) : "",
      service: isSet(object.service) ? String(object.service) : "",
      spender: isSet(object.spender) ? String(object.spender) : "",
      signature: isSet(object.signature) ? bytesFromBase64(object.signature) : new Uint8Array(),
      nonce: isSet(object.nonce) ? Number(object.nonce) : 0,
      height: isSet(object.height) ? Number(object.height) : 0,
    };
  },

  toJSON(message: MsgClaimContractIncome): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.pubKey !== undefined && (obj.pubKey = message.pubKey);
    message.service !== undefined && (obj.service = message.service);
    message.spender !== undefined && (obj.spender = message.spender);
    message.signature !== undefined
      && (obj.signature = base64FromBytes(message.signature !== undefined ? message.signature : new Uint8Array()));
    message.nonce !== undefined && (obj.nonce = Math.round(message.nonce));
    message.height !== undefined && (obj.height = Math.round(message.height));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgClaimContractIncome>, I>>(object: I): MsgClaimContractIncome {
    const message = createBaseMsgClaimContractIncome();
    message.creator = object.creator ?? "";
    message.pubKey = object.pubKey ?? "";
    message.service = object.service ?? "";
    message.spender = object.spender ?? "";
    message.signature = object.signature ?? new Uint8Array();
    message.nonce = object.nonce ?? 0;
    message.height = object.height ?? 0;
    return message;
  },
};

function createBaseMsgClaimContractIncomeResponse(): MsgClaimContractIncomeResponse {
  return {};
}

export const MsgClaimContractIncomeResponse = {
  encode(_: MsgClaimContractIncomeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgClaimContractIncomeResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgClaimContractIncomeResponse();
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

  fromJSON(_: any): MsgClaimContractIncomeResponse {
    return {};
  },

  toJSON(_: MsgClaimContractIncomeResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgClaimContractIncomeResponse>, I>>(_: I): MsgClaimContractIncomeResponse {
    const message = createBaseMsgClaimContractIncomeResponse();
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  BondProvider(request: MsgBondProvider): Promise<MsgBondProviderResponse>;
  ModProvider(request: MsgModProvider): Promise<MsgModProviderResponse>;
  OpenContract(request: MsgOpenContract): Promise<MsgOpenContractResponse>;
  CloseContract(request: MsgCloseContract): Promise<MsgCloseContractResponse>;
  /** this line is used by starport scaffolding # proto/tx/rpc */
  ClaimContractIncome(request: MsgClaimContractIncome): Promise<MsgClaimContractIncomeResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.BondProvider = this.BondProvider.bind(this);
    this.ModProvider = this.ModProvider.bind(this);
    this.OpenContract = this.OpenContract.bind(this);
    this.CloseContract = this.CloseContract.bind(this);
    this.ClaimContractIncome = this.ClaimContractIncome.bind(this);
  }
  BondProvider(request: MsgBondProvider): Promise<MsgBondProviderResponse> {
    const data = MsgBondProvider.encode(request).finish();
    const promise = this.rpc.request("arkeo.arkeo.Msg", "BondProvider", data);
    return promise.then((data) => MsgBondProviderResponse.decode(new _m0.Reader(data)));
  }

  ModProvider(request: MsgModProvider): Promise<MsgModProviderResponse> {
    const data = MsgModProvider.encode(request).finish();
    const promise = this.rpc.request("arkeo.arkeo.Msg", "ModProvider", data);
    return promise.then((data) => MsgModProviderResponse.decode(new _m0.Reader(data)));
  }

  OpenContract(request: MsgOpenContract): Promise<MsgOpenContractResponse> {
    const data = MsgOpenContract.encode(request).finish();
    const promise = this.rpc.request("arkeo.arkeo.Msg", "OpenContract", data);
    return promise.then((data) => MsgOpenContractResponse.decode(new _m0.Reader(data)));
  }

  CloseContract(request: MsgCloseContract): Promise<MsgCloseContractResponse> {
    const data = MsgCloseContract.encode(request).finish();
    const promise = this.rpc.request("arkeo.arkeo.Msg", "CloseContract", data);
    return promise.then((data) => MsgCloseContractResponse.decode(new _m0.Reader(data)));
  }

  ClaimContractIncome(request: MsgClaimContractIncome): Promise<MsgClaimContractIncomeResponse> {
    const data = MsgClaimContractIncome.encode(request).finish();
    const promise = this.rpc.request("arkeo.arkeo.Msg", "ClaimContractIncome", data);
    return promise.then((data) => MsgClaimContractIncomeResponse.decode(new _m0.Reader(data)));
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
