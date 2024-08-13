import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgTransferClaim } from "./types/arkeo/claim/tx";
import { QueryClaimRecordRequest } from "./types/arkeo/claim/query";
import { GenesisState } from "./types/arkeo/claim/genesis";
import { MsgClaimEth } from "./types/arkeo/claim/tx";
import { MsgClaimArkeo } from "./types/arkeo/claim/tx";
import { MsgTransferClaimResponse } from "./types/arkeo/claim/tx";
import { ClaimRecord } from "./types/arkeo/claim/claim_record";
import { QueryParamsRequest } from "./types/arkeo/claim/query";
import { QueryParamsResponse } from "./types/arkeo/claim/query";
import { QueryClaimRecordResponse } from "./types/arkeo/claim/query";
import { MsgClaimArkeoResponse } from "./types/arkeo/claim/tx";
import { MsgAddClaimResponse } from "./types/arkeo/claim/tx";
import { Params } from "./types/arkeo/claim/params";
import { MsgClaimEthResponse } from "./types/arkeo/claim/tx";
import { MsgAddClaim } from "./types/arkeo/claim/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/arkeo.claim.MsgTransferClaim", MsgTransferClaim],
    ["/arkeo.claim.QueryClaimRecordRequest", QueryClaimRecordRequest],
    ["/arkeo.claim.GenesisState", GenesisState],
    ["/arkeo.claim.MsgClaimEth", MsgClaimEth],
    ["/arkeo.claim.MsgClaimArkeo", MsgClaimArkeo],
    ["/arkeo.claim.MsgTransferClaimResponse", MsgTransferClaimResponse],
    ["/arkeo.claim.ClaimRecord", ClaimRecord],
    ["/arkeo.claim.QueryParamsRequest", QueryParamsRequest],
    ["/arkeo.claim.QueryParamsResponse", QueryParamsResponse],
    ["/arkeo.claim.QueryClaimRecordResponse", QueryClaimRecordResponse],
    ["/arkeo.claim.MsgClaimArkeoResponse", MsgClaimArkeoResponse],
    ["/arkeo.claim.MsgAddClaimResponse", MsgAddClaimResponse],
    ["/arkeo.claim.Params", Params],
    ["/arkeo.claim.MsgClaimEthResponse", MsgClaimEthResponse],
    ["/arkeo.claim.MsgAddClaim", MsgAddClaim],
    
];

export { msgTypes }