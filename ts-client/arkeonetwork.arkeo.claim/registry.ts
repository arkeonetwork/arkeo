import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgClaimEth } from "./types/arkeo/claim/tx";
import { MsgClaimArkeo } from "./types/arkeo/claim/tx";
import { MsgTransferClaim } from "./types/arkeo/claim/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/arkeonetwork.arkeo.claim.MsgClaimEth", MsgClaimEth],
    ["/arkeonetwork.arkeo.claim.MsgClaimArkeo", MsgClaimArkeo],
    ["/arkeonetwork.arkeo.claim.MsgTransferClaim", MsgTransferClaim],
    
];

export { msgTypes }