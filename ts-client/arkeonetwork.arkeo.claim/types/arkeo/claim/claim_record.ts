/* eslint-disable */
import _m0 from "protobufjs/minimal";
import { Coin } from "../../cosmos/base/v1beta1/coin";

export const protobufPackage = "arkeonetwork.arkeo.claim";

/** actions for arkeo chain */
export enum Action {
  ActionVote = 0,
  ActionDelegateStake = 1,
  UNRECOGNIZED = -1,
}

export function actionFromJSON(object: any): Action {
  switch (object) {
    case 0:
    case "ActionVote":
      return Action.ActionVote;
    case 1:
    case "ActionDelegateStake":
      return Action.ActionDelegateStake;
    case -1:
    case "UNRECOGNIZED":
    default:
      return Action.UNRECOGNIZED;
  }
}

export function actionToJSON(object: Action): string {
  switch (object) {
    case Action.ActionVote:
      return "ActionVote";
    case Action.ActionDelegateStake:
      return "ActionDelegateStake";
    case Action.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/** actions for chains other than arkeo, limited currently to claiming */
export enum ForeignChainAction {
  ForeignChainActionClaim = 0,
  UNRECOGNIZED = -1,
}

export function foreignChainActionFromJSON(object: any): ForeignChainAction {
  switch (object) {
    case 0:
    case "ForeignChainActionClaim":
      return ForeignChainAction.ForeignChainActionClaim;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ForeignChainAction.UNRECOGNIZED;
  }
}

export function foreignChainActionToJSON(object: ForeignChainAction): string {
  switch (object) {
    case ForeignChainAction.ForeignChainActionClaim:
      return "ForeignChainActionClaim";
    case ForeignChainAction.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export enum Chain {
  ARKEO = 0,
  ETHEREUM = 1,
  THORCHAIN = 2,
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
    case 2:
    case "THORCHAIN":
      return Chain.THORCHAIN;
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
    case Chain.THORCHAIN:
      return "THORCHAIN";
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
  /** total initial claimable amount for the user */
  initialClaimableAmount: Coin[];
  /**
   * true if action is completed
   * index of bool in array refers to action enum #
   */
  actionCompleted: boolean[];
}

function createBaseClaimRecord(): ClaimRecord {
  return { chain: 0, address: "", initialClaimableAmount: [], actionCompleted: [] };
}

export const ClaimRecord = {
  encode(message: ClaimRecord, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.chain !== 0) {
      writer.uint32(8).int32(message.chain);
    }
    if (message.address !== "") {
      writer.uint32(18).string(message.address);
    }
    for (const v of message.initialClaimableAmount) {
      Coin.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    writer.uint32(34).fork();
    for (const v of message.actionCompleted) {
      writer.bool(v);
    }
    writer.ldelim();
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
          message.initialClaimableAmount.push(Coin.decode(reader, reader.uint32()));
          break;
        case 4:
          if ((tag & 7) === 2) {
            const end2 = reader.uint32() + reader.pos;
            while (reader.pos < end2) {
              message.actionCompleted.push(reader.bool());
            }
          } else {
            message.actionCompleted.push(reader.bool());
          }
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
      initialClaimableAmount: Array.isArray(object?.initialClaimableAmount)
        ? object.initialClaimableAmount.map((e: any) => Coin.fromJSON(e))
        : [],
      actionCompleted: Array.isArray(object?.actionCompleted) ? object.actionCompleted.map((e: any) => Boolean(e)) : [],
    };
  },

  toJSON(message: ClaimRecord): unknown {
    const obj: any = {};
    message.chain !== undefined && (obj.chain = chainToJSON(message.chain));
    message.address !== undefined && (obj.address = message.address);
    if (message.initialClaimableAmount) {
      obj.initialClaimableAmount = message.initialClaimableAmount.map((e) => e ? Coin.toJSON(e) : undefined);
    } else {
      obj.initialClaimableAmount = [];
    }
    if (message.actionCompleted) {
      obj.actionCompleted = message.actionCompleted.map((e) => e);
    } else {
      obj.actionCompleted = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ClaimRecord>, I>>(object: I): ClaimRecord {
    const message = createBaseClaimRecord();
    message.chain = object.chain ?? 0;
    message.address = object.address ?? "";
    message.initialClaimableAmount = object.initialClaimableAmount?.map((e) => Coin.fromPartial(e)) || [];
    message.actionCompleted = object.actionCompleted?.map((e) => e) || [];
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
