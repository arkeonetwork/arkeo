import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgAddClaim } from "./types/arkeo/claim/tx";
import { MsgClaimEth } from "./types/arkeo/claim/tx";
import { MsgTransferClaim } from "./types/arkeo/claim/tx";
import { MsgClaimArkeo } from "./types/arkeo/claim/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/arkeonetwork.arkeo.claim.MsgAddClaim", MsgAddClaim],
    ["/arkeonetwork.arkeo.claim.MsgClaimEth", MsgClaimEth],
    ["/arkeonetwork.arkeo.claim.MsgTransferClaim", MsgTransferClaim],
    ["/arkeonetwork.arkeo.claim.MsgClaimArkeo", MsgClaimArkeo],
    
];

export { msgTypes }