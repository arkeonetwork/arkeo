/* eslint-disable */
import {
  ProviderStatus,
  ContractType,
  providerStatusFromJSON,
  providerStatusToJSON,
  contractTypeFromJSON,
  contractTypeToJSON,
} from "../mercury/keeper";
import { Reader, util, configure, Writer } from "protobufjs/minimal";
import * as Long from "long";

export const protobufPackage = "mercury.mercury";

export interface MsgBondProvider {
  creator: string;
  pubKey: string;
  chain: string;
  bond: string;
}

export interface MsgBondProviderResponse {}

export interface MsgModProvider {
  creator: string;
  pubKey: string;
  chain: string;
  metadataURI: string;
  metadataNonce: number;
  status: ProviderStatus;
  minContractDuration: number;
  maxContractDuration: number;
  subscriptionRate: number;
  payAsYouGoRate: number;
}

export interface MsgModProviderResponse {}

export interface MsgOpenContract {
  creator: string;
  pubKey: string;
  chain: string;
  client: string;
  delegate: string;
  cType: ContractType;
  duration: number;
  rate: number;
  deposit: string;
}

export interface MsgOpenContractResponse {}

export interface MsgCloseContract {
  creator: string;
  pubKey: string;
  chain: string;
  client: string;
  delegate: string;
}

export interface MsgCloseContractResponse {}

export interface MsgClaimContractIncome {
  creator: string;
  pubKey: string;
  chain: string;
  spender: string;
  signature: Uint8Array;
  nonce: number;
  height: number;
}

export interface MsgClaimContractIncomeResponse {}

const baseMsgBondProvider: object = {
  creator: "",
  pubKey: "",
  chain: "",
  bond: "",
};

export const MsgBondProvider = {
  encode(message: MsgBondProvider, writer: Writer = Writer.create()): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.pubKey !== "") {
      writer.uint32(18).string(message.pubKey);
    }
    if (message.chain !== "") {
      writer.uint32(26).string(message.chain);
    }
    if (message.bond !== "") {
      writer.uint32(34).string(message.bond);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgBondProvider {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgBondProvider } as MsgBondProvider;
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
          message.chain = reader.string();
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
    const message = { ...baseMsgBondProvider } as MsgBondProvider;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
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
    if (object.bond !== undefined && object.bond !== null) {
      message.bond = String(object.bond);
    } else {
      message.bond = "";
    }
    return message;
  },

  toJSON(message: MsgBondProvider): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.pubKey !== undefined && (obj.pubKey = message.pubKey);
    message.chain !== undefined && (obj.chain = message.chain);
    message.bond !== undefined && (obj.bond = message.bond);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgBondProvider>): MsgBondProvider {
    const message = { ...baseMsgBondProvider } as MsgBondProvider;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
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
    if (object.bond !== undefined && object.bond !== null) {
      message.bond = object.bond;
    } else {
      message.bond = "";
    }
    return message;
  },
};

const baseMsgBondProviderResponse: object = {};

export const MsgBondProviderResponse = {
  encode(_: MsgBondProviderResponse, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgBondProviderResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgBondProviderResponse,
    } as MsgBondProviderResponse;
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
    const message = {
      ...baseMsgBondProviderResponse,
    } as MsgBondProviderResponse;
    return message;
  },

  toJSON(_: MsgBondProviderResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgBondProviderResponse>
  ): MsgBondProviderResponse {
    const message = {
      ...baseMsgBondProviderResponse,
    } as MsgBondProviderResponse;
    return message;
  },
};

const baseMsgModProvider: object = {
  creator: "",
  pubKey: "",
  chain: "",
  metadataURI: "",
  metadataNonce: 0,
  status: 0,
  minContractDuration: 0,
  maxContractDuration: 0,
  subscriptionRate: 0,
  payAsYouGoRate: 0,
};

export const MsgModProvider = {
  encode(message: MsgModProvider, writer: Writer = Writer.create()): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.pubKey !== "") {
      writer.uint32(18).string(message.pubKey);
    }
    if (message.chain !== "") {
      writer.uint32(26).string(message.chain);
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

  decode(input: Reader | Uint8Array, length?: number): MsgModProvider {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgModProvider } as MsgModProvider;
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
          message.chain = reader.string();
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
    const message = { ...baseMsgModProvider } as MsgModProvider;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
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
    return message;
  },

  toJSON(message: MsgModProvider): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
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
    return obj;
  },

  fromPartial(object: DeepPartial<MsgModProvider>): MsgModProvider {
    const message = { ...baseMsgModProvider } as MsgModProvider;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
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
    return message;
  },
};

const baseMsgModProviderResponse: object = {};

export const MsgModProviderResponse = {
  encode(_: MsgModProviderResponse, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgModProviderResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgModProviderResponse } as MsgModProviderResponse;
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
    const message = { ...baseMsgModProviderResponse } as MsgModProviderResponse;
    return message;
  },

  toJSON(_: MsgModProviderResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<MsgModProviderResponse>): MsgModProviderResponse {
    const message = { ...baseMsgModProviderResponse } as MsgModProviderResponse;
    return message;
  },
};

const baseMsgOpenContract: object = {
  creator: "",
  pubKey: "",
  chain: "",
  client: "",
  delegate: "",
  cType: 0,
  duration: 0,
  rate: 0,
  deposit: "",
};

export const MsgOpenContract = {
  encode(message: MsgOpenContract, writer: Writer = Writer.create()): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.pubKey !== "") {
      writer.uint32(18).string(message.pubKey);
    }
    if (message.chain !== "") {
      writer.uint32(26).string(message.chain);
    }
    if (message.client !== "") {
      writer.uint32(34).string(message.client);
    }
    if (message.delegate !== "") {
      writer.uint32(42).string(message.delegate);
    }
    if (message.cType !== 0) {
      writer.uint32(48).int32(message.cType);
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

  decode(input: Reader | Uint8Array, length?: number): MsgOpenContract {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgOpenContract } as MsgOpenContract;
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
          message.chain = reader.string();
          break;
        case 4:
          message.client = reader.string();
          break;
        case 5:
          message.delegate = reader.string();
          break;
        case 6:
          message.cType = reader.int32() as any;
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
    const message = { ...baseMsgOpenContract } as MsgOpenContract;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
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
    if (object.client !== undefined && object.client !== null) {
      message.client = String(object.client);
    } else {
      message.client = "";
    }
    if (object.delegate !== undefined && object.delegate !== null) {
      message.delegate = String(object.delegate);
    } else {
      message.delegate = "";
    }
    if (object.cType !== undefined && object.cType !== null) {
      message.cType = contractTypeFromJSON(object.cType);
    } else {
      message.cType = 0;
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
    return message;
  },

  toJSON(message: MsgOpenContract): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.pubKey !== undefined && (obj.pubKey = message.pubKey);
    message.chain !== undefined && (obj.chain = message.chain);
    message.client !== undefined && (obj.client = message.client);
    message.delegate !== undefined && (obj.delegate = message.delegate);
    message.cType !== undefined &&
      (obj.cType = contractTypeToJSON(message.cType));
    message.duration !== undefined && (obj.duration = message.duration);
    message.rate !== undefined && (obj.rate = message.rate);
    message.deposit !== undefined && (obj.deposit = message.deposit);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgOpenContract>): MsgOpenContract {
    const message = { ...baseMsgOpenContract } as MsgOpenContract;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
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
    if (object.client !== undefined && object.client !== null) {
      message.client = object.client;
    } else {
      message.client = "";
    }
    if (object.delegate !== undefined && object.delegate !== null) {
      message.delegate = object.delegate;
    } else {
      message.delegate = "";
    }
    if (object.cType !== undefined && object.cType !== null) {
      message.cType = object.cType;
    } else {
      message.cType = 0;
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
    return message;
  },
};

const baseMsgOpenContractResponse: object = {};

export const MsgOpenContractResponse = {
  encode(_: MsgOpenContractResponse, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgOpenContractResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgOpenContractResponse,
    } as MsgOpenContractResponse;
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
    const message = {
      ...baseMsgOpenContractResponse,
    } as MsgOpenContractResponse;
    return message;
  },

  toJSON(_: MsgOpenContractResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgOpenContractResponse>
  ): MsgOpenContractResponse {
    const message = {
      ...baseMsgOpenContractResponse,
    } as MsgOpenContractResponse;
    return message;
  },
};

const baseMsgCloseContract: object = {
  creator: "",
  pubKey: "",
  chain: "",
  client: "",
  delegate: "",
};

export const MsgCloseContract = {
  encode(message: MsgCloseContract, writer: Writer = Writer.create()): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.pubKey !== "") {
      writer.uint32(18).string(message.pubKey);
    }
    if (message.chain !== "") {
      writer.uint32(26).string(message.chain);
    }
    if (message.client !== "") {
      writer.uint32(34).string(message.client);
    }
    if (message.delegate !== "") {
      writer.uint32(42).string(message.delegate);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgCloseContract {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgCloseContract } as MsgCloseContract;
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
          message.chain = reader.string();
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
    const message = { ...baseMsgCloseContract } as MsgCloseContract;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
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
    if (object.client !== undefined && object.client !== null) {
      message.client = String(object.client);
    } else {
      message.client = "";
    }
    if (object.delegate !== undefined && object.delegate !== null) {
      message.delegate = String(object.delegate);
    } else {
      message.delegate = "";
    }
    return message;
  },

  toJSON(message: MsgCloseContract): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.pubKey !== undefined && (obj.pubKey = message.pubKey);
    message.chain !== undefined && (obj.chain = message.chain);
    message.client !== undefined && (obj.client = message.client);
    message.delegate !== undefined && (obj.delegate = message.delegate);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgCloseContract>): MsgCloseContract {
    const message = { ...baseMsgCloseContract } as MsgCloseContract;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
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
    if (object.client !== undefined && object.client !== null) {
      message.client = object.client;
    } else {
      message.client = "";
    }
    if (object.delegate !== undefined && object.delegate !== null) {
      message.delegate = object.delegate;
    } else {
      message.delegate = "";
    }
    return message;
  },
};

const baseMsgCloseContractResponse: object = {};

export const MsgCloseContractResponse = {
  encode(
    _: MsgCloseContractResponse,
    writer: Writer = Writer.create()
  ): Writer {
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): MsgCloseContractResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgCloseContractResponse,
    } as MsgCloseContractResponse;
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
    const message = {
      ...baseMsgCloseContractResponse,
    } as MsgCloseContractResponse;
    return message;
  },

  toJSON(_: MsgCloseContractResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgCloseContractResponse>
  ): MsgCloseContractResponse {
    const message = {
      ...baseMsgCloseContractResponse,
    } as MsgCloseContractResponse;
    return message;
  },
};

const baseMsgClaimContractIncome: object = {
  creator: "",
  pubKey: "",
  chain: "",
  spender: "",
  nonce: 0,
  height: 0,
};

export const MsgClaimContractIncome = {
  encode(
    message: MsgClaimContractIncome,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.pubKey !== "") {
      writer.uint32(18).string(message.pubKey);
    }
    if (message.chain !== "") {
      writer.uint32(26).string(message.chain);
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

  decode(input: Reader | Uint8Array, length?: number): MsgClaimContractIncome {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgClaimContractIncome } as MsgClaimContractIncome;
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
          message.chain = reader.string();
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
    const message = { ...baseMsgClaimContractIncome } as MsgClaimContractIncome;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
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
    if (object.spender !== undefined && object.spender !== null) {
      message.spender = String(object.spender);
    } else {
      message.spender = "";
    }
    if (object.signature !== undefined && object.signature !== null) {
      message.signature = bytesFromBase64(object.signature);
    }
    if (object.nonce !== undefined && object.nonce !== null) {
      message.nonce = Number(object.nonce);
    } else {
      message.nonce = 0;
    }
    if (object.height !== undefined && object.height !== null) {
      message.height = Number(object.height);
    } else {
      message.height = 0;
    }
    return message;
  },

  toJSON(message: MsgClaimContractIncome): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.pubKey !== undefined && (obj.pubKey = message.pubKey);
    message.chain !== undefined && (obj.chain = message.chain);
    message.spender !== undefined && (obj.spender = message.spender);
    message.signature !== undefined &&
      (obj.signature = base64FromBytes(
        message.signature !== undefined ? message.signature : new Uint8Array()
      ));
    message.nonce !== undefined && (obj.nonce = message.nonce);
    message.height !== undefined && (obj.height = message.height);
    return obj;
  },

  fromPartial(
    object: DeepPartial<MsgClaimContractIncome>
  ): MsgClaimContractIncome {
    const message = { ...baseMsgClaimContractIncome } as MsgClaimContractIncome;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
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
    if (object.spender !== undefined && object.spender !== null) {
      message.spender = object.spender;
    } else {
      message.spender = "";
    }
    if (object.signature !== undefined && object.signature !== null) {
      message.signature = object.signature;
    } else {
      message.signature = new Uint8Array();
    }
    if (object.nonce !== undefined && object.nonce !== null) {
      message.nonce = object.nonce;
    } else {
      message.nonce = 0;
    }
    if (object.height !== undefined && object.height !== null) {
      message.height = object.height;
    } else {
      message.height = 0;
    }
    return message;
  },
};

const baseMsgClaimContractIncomeResponse: object = {};

export const MsgClaimContractIncomeResponse = {
  encode(
    _: MsgClaimContractIncomeResponse,
    writer: Writer = Writer.create()
  ): Writer {
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): MsgClaimContractIncomeResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgClaimContractIncomeResponse,
    } as MsgClaimContractIncomeResponse;
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
    const message = {
      ...baseMsgClaimContractIncomeResponse,
    } as MsgClaimContractIncomeResponse;
    return message;
  },

  toJSON(_: MsgClaimContractIncomeResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgClaimContractIncomeResponse>
  ): MsgClaimContractIncomeResponse {
    const message = {
      ...baseMsgClaimContractIncomeResponse,
    } as MsgClaimContractIncomeResponse;
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
  ClaimContractIncome(
    request: MsgClaimContractIncome
  ): Promise<MsgClaimContractIncomeResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
  }
  BondProvider(request: MsgBondProvider): Promise<MsgBondProviderResponse> {
    const data = MsgBondProvider.encode(request).finish();
    const promise = this.rpc.request(
      "mercury.mercury.Msg",
      "BondProvider",
      data
    );
    return promise.then((data) =>
      MsgBondProviderResponse.decode(new Reader(data))
    );
  }

  ModProvider(request: MsgModProvider): Promise<MsgModProviderResponse> {
    const data = MsgModProvider.encode(request).finish();
    const promise = this.rpc.request(
      "mercury.mercury.Msg",
      "ModProvider",
      data
    );
    return promise.then((data) =>
      MsgModProviderResponse.decode(new Reader(data))
    );
  }

  OpenContract(request: MsgOpenContract): Promise<MsgOpenContractResponse> {
    const data = MsgOpenContract.encode(request).finish();
    const promise = this.rpc.request(
      "mercury.mercury.Msg",
      "OpenContract",
      data
    );
    return promise.then((data) =>
      MsgOpenContractResponse.decode(new Reader(data))
    );
  }

  CloseContract(request: MsgCloseContract): Promise<MsgCloseContractResponse> {
    const data = MsgCloseContract.encode(request).finish();
    const promise = this.rpc.request(
      "mercury.mercury.Msg",
      "CloseContract",
      data
    );
    return promise.then((data) =>
      MsgCloseContractResponse.decode(new Reader(data))
    );
  }

  ClaimContractIncome(
    request: MsgClaimContractIncome
  ): Promise<MsgClaimContractIncomeResponse> {
    const data = MsgClaimContractIncome.encode(request).finish();
    const promise = this.rpc.request(
      "mercury.mercury.Msg",
      "ClaimContractIncome",
      data
    );
    return promise.then((data) =>
      MsgClaimContractIncomeResponse.decode(new Reader(data))
    );
  }
}

interface Rpc {
  request(
    service: string,
    method: string,
    data: Uint8Array
  ): Promise<Uint8Array>;
}

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
