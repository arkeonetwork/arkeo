import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgTransferClaim } from "./types/arkeo/claim/tx";
import { MsgAddClaim } from "./types/arkeo/claim/tx";
import { MsgClaimArkeo } from "./types/arkeo/claim/tx";
import { MsgClaimEth } from "./types/arkeo/claim/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/arkeo.claim.MsgTransferClaim", MsgTransferClaim],
    ["/arkeo.claim.MsgAddClaim", MsgAddClaim],
    ["/arkeo.claim.MsgClaimArkeo", MsgClaimArkeo],
    ["/arkeo.claim.MsgClaimEth", MsgClaimEth],
    
];

export { msgTypes }