import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgClaimArkeo } from "./types/arkeo/claim/tx";
import { MsgClaimEth } from "./types/arkeo/claim/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/arkeonetwork.arkeo.claim.MsgClaimArkeo", MsgClaimArkeo],
    ["/arkeonetwork.arkeo.claim.MsgClaimEth", MsgClaimEth],
    
];

export { msgTypes }