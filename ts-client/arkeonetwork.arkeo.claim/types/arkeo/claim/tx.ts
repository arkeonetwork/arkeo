/* eslint-disable */
import _m0 from "protobufjs/minimal";

export const protobufPackage = "arkeonetwork.arkeo.claim";

export interface MsgClaimEth {
  creator: string;
  ethAdress: string;
  signature: string;
}

export interface MsgClaimEthResponse {
}

function createBaseMsgClaimEth(): MsgClaimEth {
  return { creator: "", ethAdress: "", signature: "" };
}

export const MsgClaimEth = {
  encode(message: MsgClaimEth, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.creator !== "") {
      writer.uint32(10).string(message.creator);
    }
    if (message.ethAdress !== "") {
      writer.uint32(18).string(message.ethAdress);
    }
    if (message.signature !== "") {
      writer.uint32(26).string(message.signature);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgClaimEth {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgClaimEth();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.creator = reader.string();
          break;
        case 2:
          message.ethAdress = reader.string();
          break;
        case 3:
          message.signature = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): MsgClaimEth {
    return {
      creator: isSet(object.creator) ? String(object.creator) : "",
      ethAdress: isSet(object.ethAdress) ? String(object.ethAdress) : "",
      signature: isSet(object.signature) ? String(object.signature) : "",
    };
  },

  toJSON(message: MsgClaimEth): unknown {
    const obj: any = {};
    message.creator !== undefined && (obj.creator = message.creator);
    message.ethAdress !== undefined && (obj.ethAdress = message.ethAdress);
    message.signature !== undefined && (obj.signature = message.signature);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgClaimEth>, I>>(object: I): MsgClaimEth {
    const message = createBaseMsgClaimEth();
    message.creator = object.creator ?? "";
    message.ethAdress = object.ethAdress ?? "";
    message.signature = object.signature ?? "";
    return message;
  },
};

function createBaseMsgClaimEthResponse(): MsgClaimEthResponse {
  return {};
}

export const MsgClaimEthResponse = {
  encode(_: MsgClaimEthResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): MsgClaimEthResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseMsgClaimEthResponse();
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

  fromJSON(_: any): MsgClaimEthResponse {
    return {};
  },

  toJSON(_: MsgClaimEthResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<MsgClaimEthResponse>, I>>(_: I): MsgClaimEthResponse {
    const message = createBaseMsgClaimEthResponse();
    return message;
  },
};

/** Msg defines the Msg service. */
export interface Msg {
  ClaimEth(request: MsgClaimEth): Promise<MsgClaimEthResponse>;
}

export class MsgClientImpl implements Msg {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.ClaimEth = this.ClaimEth.bind(this);
  }
  ClaimEth(request: MsgClaimEth): Promise<MsgClaimEthResponse> {
    const data = MsgClaimEth.encode(request).finish();
    const promise = this.rpc.request("arkeonetwork.arkeo.claim.Msg", "ClaimEth", data);
    return promise.then((data) => MsgClaimEthResponse.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
