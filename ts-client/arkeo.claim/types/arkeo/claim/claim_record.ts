/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Coin } from "../../cosmos/base/v1beta1/coin";

export const protobufPackage = "arkeo.claim";

/** actions for arkeo chain */
export enum Action {
  ACTION_CLAIM = 0,
  ACTION_VOTE = 1,
  ACTION_DELEGATE = 2,
  UNRECOGNIZED = -1,
}

export function actionFromJSON(object: any): Action {
  switch (object) {
    case 0:
    case "ACTION_CLAIM":
      return Action.ACTION_CLAIM;
    case 1:
    case "ACTION_VOTE":
      return Action.ACTION_VOTE;
    case 2:
    case "ACTION_DELEGATE":
      return Action.ACTION_DELEGATE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return Action.UNRECOGNIZED;
  }
}

export function actionToJSON(object: Action): string {
  switch (object) {
    case Action.ACTION_CLAIM:
      return "ACTION_CLAIM";
    case Action.ACTION_VOTE:
      return "ACTION_VOTE";
    case Action.ACTION_DELEGATE:
      return "ACTION_DELEGATE";
    case Action.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum Chain {
  ARKEO = 0,
  ETHEREUM = 1,
  UNRECOGNIZED = -1,
}

export function chainFromJSON(object: any): Chain {
  switch (object) {
    case 0:
    case "ARKEO":
      return Chain.ARKEO;
    case 1:
    case "ETHEREUM":
      return Chain.ETHEREUM;
    case -1:
    case "UNRECOGNIZED":
    default:
      return Chain.UNRECOGNIZED;
  }
}

export function chainToJSON(object: Chain): string {
  switch (object) {
    case Chain.ARKEO:
      return "ARKEO";
    case Chain.ETHEREUM:
      return "ETHEREUM";
    case Chain.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/** A Claim Records is the metadata of claim data per address */
export interface ClaimRecord {
  chain: Chain;
  /** arkeo address of claim user */
  address: string;
  /**
   * claimable amount per action (claim, vote, delegate - changed to 0 after
   * action completed)
   */
  amountClaim: Coin | undefined;
  amountVote: Coin | undefined;
  amountDelegate: Coin | undefined;
  isTransferable: boolean;
}

function createBaseClaimRecord(): ClaimRecord {
  return {
    chain: 0,
    address: "",
    amountClaim: undefined,
    amountVote: undefined,
    amountDelegate: undefined,
    isTransferable: false,
  };
}

export const ClaimRecord = {
  encode(message: ClaimRecord, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chain !== 0) {
      writer.uint32(8).int32(message.chain);
    }
    if (message.address !== "") {
      writer.uint32(18).string(message.address);
    }
    if (message.amountClaim !== undefined) {
      Coin.encode(message.amountClaim, writer.uint32(26).fork()).ldelim();
    }
    if (message.amountVote !== undefined) {
      Coin.encode(message.amountVote, writer.uint32(34).fork()).ldelim();
    }
    if (message.amountDelegate !== undefined) {
      Coin.encode(message.amountDelegate, writer.uint32(42).fork()).ldelim();
    }
    if (message.isTransferable === true) {
      writer.uint32(48).bool(message.isTransferable);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ClaimRecord {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseClaimRecord();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.chain = reader.int32() as any;
          break;
        case 2:
          message.address = reader.string();
          break;
        case 3:
          message.amountClaim = Coin.decode(reader, reader.uint32());
          break;
        case 4:
          message.amountVote = Coin.decode(reader, reader.uint32());
          break;
        case 5:
          message.amountDelegate = Coin.decode(reader, reader.uint32());
          break;
        case 6:
          message.isTransferable = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ClaimRecord {
    return {
      chain: isSet(object.chain) ? chainFromJSON(object.chain) : 0,
      address: isSet(object.address) ? String(object.address) : "",
      amountClaim: isSet(object.amountClaim) ? Coin.fromJSON(object.amountClaim) : undefined,
      amountVote: isSet(object.amountVote) ? Coin.fromJSON(object.amountVote) : undefined,
      amountDelegate: isSet(object.amountDelegate) ? Coin.fromJSON(object.amountDelegate) : undefined,
      isTransferable: isSet(object.isTransferable) ? Boolean(object.isTransferable) : false,
    };
  },

  toJSON(message: ClaimRecord): unknown {
    const obj: any = {};
    message.chain !== undefined && (obj.chain = chainToJSON(message.chain));
    message.address !== undefined && (obj.address = message.address);
    message.amountClaim !== undefined
      && (obj.amountClaim = message.amountClaim ? Coin.toJSON(message.amountClaim) : undefined);
    message.amountVote !== undefined
      && (obj.amountVote = message.amountVote ? Coin.toJSON(message.amountVote) : undefined);
    message.amountDelegate !== undefined
      && (obj.amountDelegate = message.amountDelegate ? Coin.toJSON(message.amountDelegate) : undefined);
    message.isTransferable !== undefined && (obj.isTransferable = message.isTransferable);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ClaimRecord>, I>>(object: I): ClaimRecord {
    const message = createBaseClaimRecord();
    message.chain = object.chain ?? 0;
    message.address = object.address ?? "";
    message.amountClaim = (object.amountClaim !== undefined && object.amountClaim !== null)
      ? Coin.fromPartial(object.amountClaim)
      : undefined;
    message.amountVote = (object.amountVote !== undefined && object.amountVote !== null)
      ? Coin.fromPartial(object.amountVote)
      : undefined;
    message.amountDelegate = (object.amountDelegate !== undefined && object.amountDelegate !== null)
      ? Coin.fromPartial(object.amountDelegate)
      : undefined;
    message.isTransferable = object.isTransferable ?? false;
    return message;
  },
};

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
