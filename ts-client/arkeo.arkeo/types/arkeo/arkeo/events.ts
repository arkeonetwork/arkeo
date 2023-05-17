/* eslint-disable */
import Long from "long";
import _m0 from "protobufjs/minimal";
import { Coin } from "../../cosmos/base/v1beta1/coin";
import {
  ContractAuthorization,
  contractAuthorizationFromJSON,
  contractAuthorizationToJSON,
  ContractType,
  contractTypeFromJSON,
  contractTypeToJSON,
  ProviderStatus,
  providerStatusFromJSON,
  providerStatusToJSON,
} from "./keeper";

export const protobufPackage = "arkeo.arkeo";

export interface EventBondProvider {
  provider: Uint8Array;
  service: string;
  bondRel: string;
  bondAbs: string;
}

export interface EventModProvider {
  creator: Uint8Array;
  provider: Uint8Array;
  service: string;
  metadataUri: string;
  metadataNonce: number;
  status: ProviderStatus;
  minContractDuration: number;
  maxContractDuration: number;
  subscriptionRate: Coin[];
  payAsYouGoRate: Coin[];
  bond: string;
  settlementDuration: number;
}

export interface EventOpenContract {
  provider: Uint8Array;
  contractId: number;
  service: string;
  client: Uint8Array;
  delegate: Uint8Array;
  type: ContractType;
  height: number;
  duration: number;
  rate: Coin | undefined;
  openCost: number;
  deposit: string;
  settlementDuration: number;
  authorization: ContractAuthorization;
  queriesPerMinute: number;
}

export interface EventSettleContract {
  provider: Uint8Array;
  contractId: number;
  service: string;
  client: Uint8Array;
  delegate: Uint8Array;
  type: ContractType;
  nonce: number;
  height: number;
  paid: string;
  reserve: string;
}

export interface EventCloseContract {
  contractId: number;
  provider: Uint8Array;
  service: string;
  client: Uint8Array;
  delegate: Uint8Array;
}

export interface EventValidatorPayout {
  validator: Uint8Array;
  reward: string;
}

function createBaseEventBondProvider(): EventBondProvider {
  return { provider: new Uint8Array(), service: "", bondRel: "", bondAbs: "" };
}

export const EventBondProvider = {
  encode(message: EventBondProvider, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.provider.length !== 0) {
      writer.uint32(10).bytes(message.provider);
    }
    if (message.service !== "") {
      writer.uint32(18).string(message.service);
    }
    if (message.bondRel !== "") {
      writer.uint32(26).string(message.bondRel);
    }
    if (message.bondAbs !== "") {
      writer.uint32(34).string(message.bondAbs);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EventBondProvider {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventBondProvider();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.provider = reader.bytes();
          break;
        case 2:
          message.service = reader.string();
          break;
        case 3:
          message.bondRel = reader.string();
          break;
        case 4:
          message.bondAbs = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): EventBondProvider {
    return {
      provider: isSet(object.provider) ? bytesFromBase64(object.provider) : new Uint8Array(),
      service: isSet(object.service) ? String(object.service) : "",
      bondRel: isSet(object.bondRel) ? String(object.bondRel) : "",
      bondAbs: isSet(object.bondAbs) ? String(object.bondAbs) : "",
    };
  },

  toJSON(message: EventBondProvider): unknown {
    const obj: any = {};
    message.provider !== undefined
      && (obj.provider = base64FromBytes(message.provider !== undefined ? message.provider : new Uint8Array()));
    message.service !== undefined && (obj.service = message.service);
    message.bondRel !== undefined && (obj.bondRel = message.bondRel);
    message.bondAbs !== undefined && (obj.bondAbs = message.bondAbs);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<EventBondProvider>, I>>(object: I): EventBondProvider {
    const message = createBaseEventBondProvider();
    message.provider = object.provider ?? new Uint8Array();
    message.service = object.service ?? "";
    message.bondRel = object.bondRel ?? "";
    message.bondAbs = object.bondAbs ?? "";
    return message;
  },
};

function createBaseEventModProvider(): EventModProvider {
  return {
    creator: new Uint8Array(),
    provider: new Uint8Array(),
    service: "",
    metadataUri: "",
    metadataNonce: 0,
    status: 0,
    minContractDuration: 0,
    maxContractDuration: 0,
    subscriptionRate: [],
    payAsYouGoRate: [],
    bond: "",
    settlementDuration: 0,
  };
}

export const EventModProvider = {
  encode(message: EventModProvider, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator.length !== 0) {
      writer.uint32(10).bytes(message.creator);
    }
    if (message.provider.length !== 0) {
      writer.uint32(18).bytes(message.provider);
    }
    if (message.service !== "") {
      writer.uint32(26).string(message.service);
    }
    if (message.metadataUri !== "") {
      writer.uint32(34).string(message.metadataUri);
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
    for (const v of message.subscriptionRate) {
      Coin.encode(v!, writer.uint32(74).fork()).ldelim();
    }
    for (const v of message.payAsYouGoRate) {
      Coin.encode(v!, writer.uint32(82).fork()).ldelim();
    }
    if (message.bond !== "") {
      writer.uint32(90).string(message.bond);
    }
    if (message.settlementDuration !== 0) {
      writer.uint32(96).int64(message.settlementDuration);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EventModProvider {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventModProvider();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.bytes();
          break;
        case 2:
          message.provider = reader.bytes();
          break;
        case 3:
          message.service = reader.string();
          break;
        case 4:
          message.metadataUri = reader.string();
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
          message.subscriptionRate.push(Coin.decode(reader, reader.uint32()));
          break;
        case 10:
          message.payAsYouGoRate.push(Coin.decode(reader, reader.uint32()));
          break;
        case 11:
          message.bond = reader.string();
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

  fromJSON(object: any): EventModProvider {
    return {
      creator: isSet(object.creator) ? bytesFromBase64(object.creator) : new Uint8Array(),
      provider: isSet(object.provider) ? bytesFromBase64(object.provider) : new Uint8Array(),
      service: isSet(object.service) ? String(object.service) : "",
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
      settlementDuration: isSet(object.settlementDuration) ? Number(object.settlementDuration) : 0,
    };
  },

  toJSON(message: EventModProvider): unknown {
    const obj: any = {};
    message.creator !== undefined
      && (obj.creator = base64FromBytes(message.creator !== undefined ? message.creator : new Uint8Array()));
    message.provider !== undefined
      && (obj.provider = base64FromBytes(message.provider !== undefined ? message.provider : new Uint8Array()));
    message.service !== undefined && (obj.service = message.service);
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
    message.settlementDuration !== undefined && (obj.settlementDuration = Math.round(message.settlementDuration));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<EventModProvider>, I>>(object: I): EventModProvider {
    const message = createBaseEventModProvider();
    message.creator = object.creator ?? new Uint8Array();
    message.provider = object.provider ?? new Uint8Array();
    message.service = object.service ?? "";
    message.metadataUri = object.metadataUri ?? "";
    message.metadataNonce = object.metadataNonce ?? 0;
    message.status = object.status ?? 0;
    message.minContractDuration = object.minContractDuration ?? 0;
    message.maxContractDuration = object.maxContractDuration ?? 0;
    message.subscriptionRate = object.subscriptionRate?.map((e) => Coin.fromPartial(e)) || [];
    message.payAsYouGoRate = object.payAsYouGoRate?.map((e) => Coin.fromPartial(e)) || [];
    message.bond = object.bond ?? "";
    message.settlementDuration = object.settlementDuration ?? 0;
    return message;
  },
};

function createBaseEventOpenContract(): EventOpenContract {
  return {
    provider: new Uint8Array(),
    contractId: 0,
    service: "",
    client: new Uint8Array(),
    delegate: new Uint8Array(),
    type: 0,
    height: 0,
    duration: 0,
    rate: undefined,
    openCost: 0,
    deposit: "",
    settlementDuration: 0,
    authorization: 0,
    queriesPerMinute: 0,
  };
}

export const EventOpenContract = {
  encode(message: EventOpenContract, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.provider.length !== 0) {
      writer.uint32(10).bytes(message.provider);
    }
    if (message.contractId !== 0) {
      writer.uint32(16).uint64(message.contractId);
    }
    if (message.service !== "") {
      writer.uint32(26).string(message.service);
    }
    if (message.client.length !== 0) {
      writer.uint32(34).bytes(message.client);
    }
    if (message.delegate.length !== 0) {
      writer.uint32(42).bytes(message.delegate);
    }
    if (message.type !== 0) {
      writer.uint32(48).int32(message.type);
    }
    if (message.height !== 0) {
      writer.uint32(56).int64(message.height);
    }
    if (message.duration !== 0) {
      writer.uint32(64).int64(message.duration);
    }
    if (message.rate !== undefined) {
      Coin.encode(message.rate, writer.uint32(74).fork()).ldelim();
    }
    if (message.openCost !== 0) {
      writer.uint32(80).int64(message.openCost);
    }
    if (message.deposit !== "") {
      writer.uint32(90).string(message.deposit);
    }
    if (message.settlementDuration !== 0) {
      writer.uint32(96).int64(message.settlementDuration);
    }
    if (message.authorization !== 0) {
      writer.uint32(104).int32(message.authorization);
    }
    if (message.queriesPerMinute !== 0) {
      writer.uint32(112).int64(message.queriesPerMinute);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EventOpenContract {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventOpenContract();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.provider = reader.bytes();
          break;
        case 2:
          message.contractId = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.service = reader.string();
          break;
        case 4:
          message.client = reader.bytes();
          break;
        case 5:
          message.delegate = reader.bytes();
          break;
        case 6:
          message.type = reader.int32() as any;
          break;
        case 7:
          message.height = longToNumber(reader.int64() as Long);
          break;
        case 8:
          message.duration = longToNumber(reader.int64() as Long);
          break;
        case 9:
          message.rate = Coin.decode(reader, reader.uint32());
          break;
        case 10:
          message.openCost = longToNumber(reader.int64() as Long);
          break;
        case 11:
          message.deposit = reader.string();
          break;
        case 12:
          message.settlementDuration = longToNumber(reader.int64() as Long);
          break;
        case 13:
          message.authorization = reader.int32() as any;
          break;
        case 14:
          message.queriesPerMinute = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): EventOpenContract {
    return {
      provider: isSet(object.provider) ? bytesFromBase64(object.provider) : new Uint8Array(),
      contractId: isSet(object.contractId) ? Number(object.contractId) : 0,
      service: isSet(object.service) ? String(object.service) : "",
      client: isSet(object.client) ? bytesFromBase64(object.client) : new Uint8Array(),
      delegate: isSet(object.delegate) ? bytesFromBase64(object.delegate) : new Uint8Array(),
      type: isSet(object.type) ? contractTypeFromJSON(object.type) : 0,
      height: isSet(object.height) ? Number(object.height) : 0,
      duration: isSet(object.duration) ? Number(object.duration) : 0,
      rate: isSet(object.rate) ? Coin.fromJSON(object.rate) : undefined,
      openCost: isSet(object.openCost) ? Number(object.openCost) : 0,
      deposit: isSet(object.deposit) ? String(object.deposit) : "",
      settlementDuration: isSet(object.settlementDuration) ? Number(object.settlementDuration) : 0,
      authorization: isSet(object.authorization) ? contractAuthorizationFromJSON(object.authorization) : 0,
      queriesPerMinute: isSet(object.queriesPerMinute) ? Number(object.queriesPerMinute) : 0,
    };
  },

  toJSON(message: EventOpenContract): unknown {
    const obj: any = {};
    message.provider !== undefined
      && (obj.provider = base64FromBytes(message.provider !== undefined ? message.provider : new Uint8Array()));
    message.contractId !== undefined && (obj.contractId = Math.round(message.contractId));
    message.service !== undefined && (obj.service = message.service);
    message.client !== undefined
      && (obj.client = base64FromBytes(message.client !== undefined ? message.client : new Uint8Array()));
    message.delegate !== undefined
      && (obj.delegate = base64FromBytes(message.delegate !== undefined ? message.delegate : new Uint8Array()));
    message.type !== undefined && (obj.type = contractTypeToJSON(message.type));
    message.height !== undefined && (obj.height = Math.round(message.height));
    message.duration !== undefined && (obj.duration = Math.round(message.duration));
    message.rate !== undefined && (obj.rate = message.rate ? Coin.toJSON(message.rate) : undefined);
    message.openCost !== undefined && (obj.openCost = Math.round(message.openCost));
    message.deposit !== undefined && (obj.deposit = message.deposit);
    message.settlementDuration !== undefined && (obj.settlementDuration = Math.round(message.settlementDuration));
    message.authorization !== undefined && (obj.authorization = contractAuthorizationToJSON(message.authorization));
    message.queriesPerMinute !== undefined && (obj.queriesPerMinute = Math.round(message.queriesPerMinute));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<EventOpenContract>, I>>(object: I): EventOpenContract {
    const message = createBaseEventOpenContract();
    message.provider = object.provider ?? new Uint8Array();
    message.contractId = object.contractId ?? 0;
    message.service = object.service ?? "";
    message.client = object.client ?? new Uint8Array();
    message.delegate = object.delegate ?? new Uint8Array();
    message.type = object.type ?? 0;
    message.height = object.height ?? 0;
    message.duration = object.duration ?? 0;
    message.rate = (object.rate !== undefined && object.rate !== null) ? Coin.fromPartial(object.rate) : undefined;
    message.openCost = object.openCost ?? 0;
    message.deposit = object.deposit ?? "";
    message.settlementDuration = object.settlementDuration ?? 0;
    message.authorization = object.authorization ?? 0;
    message.queriesPerMinute = object.queriesPerMinute ?? 0;
    return message;
  },
};

function createBaseEventSettleContract(): EventSettleContract {
  return {
    provider: new Uint8Array(),
    contractId: 0,
    service: "",
    client: new Uint8Array(),
    delegate: new Uint8Array(),
    type: 0,
    nonce: 0,
    height: 0,
    paid: "",
    reserve: "",
  };
}

export const EventSettleContract = {
  encode(message: EventSettleContract, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.provider.length !== 0) {
      writer.uint32(10).bytes(message.provider);
    }
    if (message.contractId !== 0) {
      writer.uint32(16).uint64(message.contractId);
    }
    if (message.service !== "") {
      writer.uint32(26).string(message.service);
    }
    if (message.client.length !== 0) {
      writer.uint32(34).bytes(message.client);
    }
    if (message.delegate.length !== 0) {
      writer.uint32(42).bytes(message.delegate);
    }
    if (message.type !== 0) {
      writer.uint32(48).int32(message.type);
    }
    if (message.nonce !== 0) {
      writer.uint32(56).int64(message.nonce);
    }
    if (message.height !== 0) {
      writer.uint32(64).int64(message.height);
    }
    if (message.paid !== "") {
      writer.uint32(74).string(message.paid);
    }
    if (message.reserve !== "") {
      writer.uint32(82).string(message.reserve);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EventSettleContract {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventSettleContract();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.provider = reader.bytes();
          break;
        case 2:
          message.contractId = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.service = reader.string();
          break;
        case 4:
          message.client = reader.bytes();
          break;
        case 5:
          message.delegate = reader.bytes();
          break;
        case 6:
          message.type = reader.int32() as any;
          break;
        case 7:
          message.nonce = longToNumber(reader.int64() as Long);
          break;
        case 8:
          message.height = longToNumber(reader.int64() as Long);
          break;
        case 9:
          message.paid = reader.string();
          break;
        case 10:
          message.reserve = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): EventSettleContract {
    return {
      provider: isSet(object.provider) ? bytesFromBase64(object.provider) : new Uint8Array(),
      contractId: isSet(object.contractId) ? Number(object.contractId) : 0,
      service: isSet(object.service) ? String(object.service) : "",
      client: isSet(object.client) ? bytesFromBase64(object.client) : new Uint8Array(),
      delegate: isSet(object.delegate) ? bytesFromBase64(object.delegate) : new Uint8Array(),
      type: isSet(object.type) ? contractTypeFromJSON(object.type) : 0,
      nonce: isSet(object.nonce) ? Number(object.nonce) : 0,
      height: isSet(object.height) ? Number(object.height) : 0,
      paid: isSet(object.paid) ? String(object.paid) : "",
      reserve: isSet(object.reserve) ? String(object.reserve) : "",
    };
  },

  toJSON(message: EventSettleContract): unknown {
    const obj: any = {};
    message.provider !== undefined
      && (obj.provider = base64FromBytes(message.provider !== undefined ? message.provider : new Uint8Array()));
    message.contractId !== undefined && (obj.contractId = Math.round(message.contractId));
    message.service !== undefined && (obj.service = message.service);
    message.client !== undefined
      && (obj.client = base64FromBytes(message.client !== undefined ? message.client : new Uint8Array()));
    message.delegate !== undefined
      && (obj.delegate = base64FromBytes(message.delegate !== undefined ? message.delegate : new Uint8Array()));
    message.type !== undefined && (obj.type = contractTypeToJSON(message.type));
    message.nonce !== undefined && (obj.nonce = Math.round(message.nonce));
    message.height !== undefined && (obj.height = Math.round(message.height));
    message.paid !== undefined && (obj.paid = message.paid);
    message.reserve !== undefined && (obj.reserve = message.reserve);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<EventSettleContract>, I>>(object: I): EventSettleContract {
    const message = createBaseEventSettleContract();
    message.provider = object.provider ?? new Uint8Array();
    message.contractId = object.contractId ?? 0;
    message.service = object.service ?? "";
    message.client = object.client ?? new Uint8Array();
    message.delegate = object.delegate ?? new Uint8Array();
    message.type = object.type ?? 0;
    message.nonce = object.nonce ?? 0;
    message.height = object.height ?? 0;
    message.paid = object.paid ?? "";
    message.reserve = object.reserve ?? "";
    return message;
  },
};

function createBaseEventCloseContract(): EventCloseContract {
  return {
    contractId: 0,
    provider: new Uint8Array(),
    service: "",
    client: new Uint8Array(),
    delegate: new Uint8Array(),
  };
}

export const EventCloseContract = {
  encode(message: EventCloseContract, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.contractId !== 0) {
      writer.uint32(8).uint64(message.contractId);
    }
    if (message.provider.length !== 0) {
      writer.uint32(18).bytes(message.provider);
    }
    if (message.service !== "") {
      writer.uint32(26).string(message.service);
    }
    if (message.client.length !== 0) {
      writer.uint32(34).bytes(message.client);
    }
    if (message.delegate.length !== 0) {
      writer.uint32(42).bytes(message.delegate);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EventCloseContract {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventCloseContract();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.contractId = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.provider = reader.bytes();
          break;
        case 3:
          message.service = reader.string();
          break;
        case 4:
          message.client = reader.bytes();
          break;
        case 5:
          message.delegate = reader.bytes();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): EventCloseContract {
    return {
      contractId: isSet(object.contractId) ? Number(object.contractId) : 0,
      provider: isSet(object.provider) ? bytesFromBase64(object.provider) : new Uint8Array(),
      service: isSet(object.service) ? String(object.service) : "",
      client: isSet(object.client) ? bytesFromBase64(object.client) : new Uint8Array(),
      delegate: isSet(object.delegate) ? bytesFromBase64(object.delegate) : new Uint8Array(),
    };
  },

  toJSON(message: EventCloseContract): unknown {
    const obj: any = {};
    message.contractId !== undefined && (obj.contractId = Math.round(message.contractId));
    message.provider !== undefined
      && (obj.provider = base64FromBytes(message.provider !== undefined ? message.provider : new Uint8Array()));
    message.service !== undefined && (obj.service = message.service);
    message.client !== undefined
      && (obj.client = base64FromBytes(message.client !== undefined ? message.client : new Uint8Array()));
    message.delegate !== undefined
      && (obj.delegate = base64FromBytes(message.delegate !== undefined ? message.delegate : new Uint8Array()));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<EventCloseContract>, I>>(object: I): EventCloseContract {
    const message = createBaseEventCloseContract();
    message.contractId = object.contractId ?? 0;
    message.provider = object.provider ?? new Uint8Array();
    message.service = object.service ?? "";
    message.client = object.client ?? new Uint8Array();
    message.delegate = object.delegate ?? new Uint8Array();
    return message;
  },
};

function createBaseEventValidatorPayout(): EventValidatorPayout {
  return { validator: new Uint8Array(), reward: "" };
}

export const EventValidatorPayout = {
  encode(message: EventValidatorPayout, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.validator.length !== 0) {
      writer.uint32(10).bytes(message.validator);
    }
    if (message.reward !== "") {
      writer.uint32(18).string(message.reward);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EventValidatorPayout {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEventValidatorPayout();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.validator = reader.bytes();
          break;
        case 2:
          message.reward = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): EventValidatorPayout {
    return {
      validator: isSet(object.validator) ? bytesFromBase64(object.validator) : new Uint8Array(),
      reward: isSet(object.reward) ? String(object.reward) : "",
    };
  },

  toJSON(message: EventValidatorPayout): unknown {
    const obj: any = {};
    message.validator !== undefined
      && (obj.validator = base64FromBytes(message.validator !== undefined ? message.validator : new Uint8Array()));
    message.reward !== undefined && (obj.reward = message.reward);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<EventValidatorPayout>, I>>(object: I): EventValidatorPayout {
    const message = createBaseEventValidatorPayout();
    message.validator = object.validator ?? new Uint8Array();
    message.reward = object.reward ?? "";
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
