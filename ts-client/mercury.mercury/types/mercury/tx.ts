/* eslint-disable */
import { Reader, Writer } from "protobufjs/minimal";

export const protobufPackage = "mercury.mercury";

export interface MsgRegisterProvider {
  creator: string;
  chain: string;
  pubkey: string;
}

export interface MsgRegisterProviderResponse {}

const baseMsgRegisterProvider: object = { creator: "", chain: "", pubkey: "" };

export const MsgRegisterProvider = {
  encode(
    message: MsgRegisterProvider,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.chain !== "") {
      writer.uint32(18).string(message.chain);
    }
    if (message.pubkey !== "") {
      writer.uint32(26).string(message.pubkey);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): MsgRegisterProvider {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseMsgRegisterProvider } as MsgRegisterProvider;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.chain = reader.string();
          break;
        case 3:
          message.pubkey = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgRegisterProvider {
    const message = { ...baseMsgRegisterProvider } as MsgRegisterProvider;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = String(object.creator);
    } else {
      message.creator = "";
    }
    if (object.chain !== undefined && object.chain !== null) {
      message.chain = String(object.chain);
    } else {
      message.chain = "";
    }
    if (object.pubkey !== undefined && object.pubkey !== null) {
      message.pubkey = String(object.pubkey);
    } else {
      message.pubkey = "";
    }
    return message;
  },

  toJSON(message: MsgRegisterProvider): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.chain !== undefined && (obj.chain = message.chain);
    message.pubkey !== undefined && (obj.pubkey = message.pubkey);
    return obj;
  },

  fromPartial(object: DeepPartial<MsgRegisterProvider>): MsgRegisterProvider {
    const message = { ...baseMsgRegisterProvider } as MsgRegisterProvider;
    if (object.creator !== undefined && object.creator !== null) {
      message.creator = object.creator;
    } else {
      message.creator = "";
    }
    if (object.chain !== undefined && object.chain !== null) {
      message.chain = object.chain;
    } else {
      message.chain = "";
    }
    if (object.pubkey !== undefined && object.pubkey !== null) {
      message.pubkey = object.pubkey;
    } else {
      message.pubkey = "";
    }
    return message;
  },
};

const baseMsgRegisterProviderResponse: object = {};

export const MsgRegisterProviderResponse = {
  encode(
    _: MsgRegisterProviderResponse,
    writer: Writer = Writer.create()
  ): Writer {
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): MsgRegisterProviderResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseMsgRegisterProviderResponse,
    } as MsgRegisterProviderResponse;
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

  fromJSON(_: any): MsgRegisterProviderResponse {
    const message = {
      ...baseMsgRegisterProviderResponse,
    } as MsgRegisterProviderResponse;
    return message;
  },

  toJSON(_: MsgRegisterProviderResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<MsgRegisterProviderResponse>
  ): MsgRegisterProviderResponse {
    const message = {
      ...baseMsgRegisterProviderResponse,
    } as MsgRegisterProviderResponse;
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  /** this line is used by starport scaffolding # proto/tx/rpc */
  RegisterProvider(
    request: MsgRegisterProvider
  ): Promise<MsgRegisterProviderResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
  }
  RegisterProvider(
    request: MsgRegisterProvider
  ): Promise<MsgRegisterProviderResponse> {
    const data = MsgRegisterProvider.encode(request).finish();
    const promise = this.rpc.request(
      "mercury.mercury.Msg",
      "RegisterProvider",
      data
    );
    return promise.then((data) =>
      MsgRegisterProviderResponse.decode(new Reader(data))
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
