/* eslint-disable */
import { Params } from "../mercury/params";
import { Provider, Contract } from "../mercury/keeper";
import { Writer, Reader } from "protobufjs/minimal";

export const protobufPackage = "mercury.mercury";

/** GenesisState defines the mercury module's genesis state. */
export interface GenesisState {
  params: Params | undefined;
  providers: Provider[];
  /** this line is used by starport scaffolding # genesis/proto/state */
  contracts: Contract[];
}

const baseGenesisState: object = {};

export const GenesisState = {
  encode(message: GenesisState, writer: Writer = Writer.create()): Writer {
    if (message.params !== undefined) {
      Params.encode(message.params, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.providers) {
      Provider.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    for (const v of message.contracts) {
      Contract.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): GenesisState {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseGenesisState } as GenesisState;
    message.providers = [];
    message.contracts = [];
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
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GenesisState {
    const message = { ...baseGenesisState } as GenesisState;
    message.providers = [];
    message.contracts = [];
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromJSON(object.params);
    } else {
      message.params = undefined;
    }
    if (object.providers !== undefined && object.providers !== null) {
      for (const e of object.providers) {
        message.providers.push(Provider.fromJSON(e));
      }
    }
    if (object.contracts !== undefined && object.contracts !== null) {
      for (const e of object.contracts) {
        message.contracts.push(Contract.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: GenesisState): unknown {
    const obj: any = {};
    message.params !== undefined &&
      (obj.params = message.params ? Params.toJSON(message.params) : undefined);
    if (message.providers) {
      obj.providers = message.providers.map((e) =>
        e ? Provider.toJSON(e) : undefined
      );
    } else {
      obj.providers = [];
    }
    if (message.contracts) {
      obj.contracts = message.contracts.map((e) =>
        e ? Contract.toJSON(e) : undefined
      );
    } else {
      obj.contracts = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<GenesisState>): GenesisState {
    const message = { ...baseGenesisState } as GenesisState;
    message.providers = [];
    message.contracts = [];
    if (object.params !== undefined && object.params !== null) {
      message.params = Params.fromPartial(object.params);
    } else {
      message.params = undefined;
    }
    if (object.providers !== undefined && object.providers !== null) {
      for (const e of object.providers) {
        message.providers.push(Provider.fromPartial(e));
      }
    }
    if (object.contracts !== undefined && object.contracts !== null) {
      for (const e of object.contracts) {
        message.contracts.push(Contract.fromPartial(e));
      }
    }
    return message;
  },
};

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
