import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgClaimArkeo } from "./types/arkeo/claim/tx";
import { MsgAddClaim } from "./types/arkeo/claim/tx";
import { MsgTransferClaim } from "./types/arkeo/claim/tx";
import { MsgClaimEth } from "./types/arkeo/claim/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/arkeo.claim.MsgClaimArkeo", MsgClaimArkeo],
    ["/arkeo.claim.MsgAddClaim", MsgAddClaim],
    ["/arkeo.claim.MsgTransferClaim", MsgTransferClaim],
    ["/arkeo.claim.MsgClaimEth", MsgClaimEth],
    
];

export { msgTypes }