/* eslint-disable */
import * as Long from "long";
import { util, configure, Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "mercury.mercury";

export interface ProtoInt64 {
  value: number;
}

export interface ProtoUint64 {
  value: number;
}

export interface ProtoAccAddresses {
  value: Uint8Array[];
}

export interface ProtoStrings {
  value: string[];
}

export interface ProtoBools {
  value: boolean[];
}

const baseProtoInt64: object = { value: 0 };

export const ProtoInt64 = {
  encode(message: ProtoInt64, writer: Writer = Writer.create()): Writer {
    if (message.value !== 0) {
      writer.uint32(8).int64(message.value);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ProtoInt64 {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseProtoInt64 } as ProtoInt64;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.value = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ProtoInt64 {
    const message = { ...baseProtoInt64 } as ProtoInt64;
    if (object.value !== undefined && object.value !== null) {
      message.value = Number(object.value);
    } else {
      message.value = 0;
    }
    return message;
  },

  toJSON(message: ProtoInt64): unknown {
    const obj: any = {};
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  fromPartial(object: DeepPartial<ProtoInt64>): ProtoInt64 {
    const message = { ...baseProtoInt64 } as ProtoInt64;
    if (object.value !== undefined && object.value !== null) {
      message.value = object.value;
    } else {
      message.value = 0;
    }
    return message;
  },
};

const baseProtoUint64: object = { value: 0 };

export const ProtoUint64 = {
  encode(message: ProtoUint64, writer: Writer = Writer.create()): Writer {
    if (message.value !== 0) {
      writer.uint32(8).uint64(message.value);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ProtoUint64 {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseProtoUint64 } as ProtoUint64;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.value = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ProtoUint64 {
    const message = { ...baseProtoUint64 } as ProtoUint64;
    if (object.value !== undefined && object.value !== null) {
      message.value = Number(object.value);
    } else {
      message.value = 0;
    }
    return message;
  },

  toJSON(message: ProtoUint64): unknown {
    const obj: any = {};
    message.value !== undefined && (obj.value = message.value);
    return obj;
  },

  fromPartial(object: DeepPartial<ProtoUint64>): ProtoUint64 {
    const message = { ...baseProtoUint64 } as ProtoUint64;
    if (object.value !== undefined && object.value !== null) {
      message.value = object.value;
    } else {
      message.value = 0;
    }
    return message;
  },
};

const baseProtoAccAddresses: object = {};

export const ProtoAccAddresses = {
  encode(message: ProtoAccAddresses, writer: Writer = Writer.create()): Writer {
    for (const v of message.value) {
      writer.uint32(10).bytes(v!);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ProtoAccAddresses {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseProtoAccAddresses } as ProtoAccAddresses;
    message.value = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.value.push(reader.bytes());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ProtoAccAddresses {
    const message = { ...baseProtoAccAddresses } as ProtoAccAddresses;
    message.value = [];
    if (object.value !== undefined && object.value !== null) {
      for (const e of object.value) {
        message.value.push(bytesFromBase64(e));
      }
    }
    return message;
  },

  toJSON(message: ProtoAccAddresses): unknown {
    const obj: any = {};
    if (message.value) {
      obj.value = message.value.map((e) =>
        base64FromBytes(e !== undefined ? e : new Uint8Array())
      );
    } else {
      obj.value = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<ProtoAccAddresses>): ProtoAccAddresses {
    const message = { ...baseProtoAccAddresses } as ProtoAccAddresses;
    message.value = [];
    if (object.value !== undefined && object.value !== null) {
      for (const e of object.value) {
        message.value.push(e);
      }
    }
    return message;
  },
};

const baseProtoStrings: object = { value: "" };

export const ProtoStrings = {
  encode(message: ProtoStrings, writer: Writer = Writer.create()): Writer {
    for (const v of message.value) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ProtoStrings {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseProtoStrings } as ProtoStrings;
    message.value = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.value.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ProtoStrings {
    const message = { ...baseProtoStrings } as ProtoStrings;
    message.value = [];
    if (object.value !== undefined && object.value !== null) {
      for (const e of object.value) {
        message.value.push(String(e));
      }
    }
    return message;
  },

  toJSON(message: ProtoStrings): unknown {
    const obj: any = {};
    if (message.value) {
      obj.value = message.value.map((e) => e);
    } else {
      obj.value = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<ProtoStrings>): ProtoStrings {
    const message = { ...baseProtoStrings } as ProtoStrings;
    message.value = [];
    if (object.value !== undefined && object.value !== null) {
      for (const e of object.value) {
        message.value.push(e);
      }
    }
    return message;
  },
};

const baseProtoBools: object = { value: false };

export const ProtoBools = {
  encode(message: ProtoBools, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).fork();
    for (const v of message.value) {
      writer.bool(v);
    }
    writer.ldelim();
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ProtoBools {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseProtoBools } as ProtoBools;
    message.value = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.value.push(reader.bool());
            }
          } else {
            message.value.push(reader.bool());
          }
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ProtoBools {
    const message = { ...baseProtoBools } as ProtoBools;
    message.value = [];
    if (object.value !== undefined && object.value !== null) {
      for (const e of object.value) {
        message.value.push(Boolean(e));
      }
    }
    return message;
  },

  toJSON(message: ProtoBools): unknown {
    const obj: any = {};
    if (message.value) {
      obj.value = message.value.map((e) => e);
    } else {
      obj.value = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<ProtoBools>): ProtoBools {
    const message = { ...baseProtoBools } as ProtoBools;
    message.value = [];
    if (object.value !== undefined && object.value !== null) {
      for (const e of object.value) {
        message.value.push(e);
      }
    }
    return message;
  },
};

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
