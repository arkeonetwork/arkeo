/* eslint-disable */
import { Reader, Writer } from "protobufjs/minimal";

export const protobufPackage = "mercury.mercury";

export interface MsgBondProvider {
  creator: string;
  pubKey: string;
  chain: string;
  bond: string;
}

export interface MsgBondProviderResponse {}

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

/** Msg defines the Msg service. */
export interface Msg {
  /** this line is used by starport scaffolding # proto/tx/rpc */
  BondProvider(request: MsgBondProvider): Promise<MsgBondProviderResponse>;
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
}

interface Rpc {
  request(
    service: string,
    method: string,
    data: Uint8Array
  ): Promise<Uint8Array>;
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
