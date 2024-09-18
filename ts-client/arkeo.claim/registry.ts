import { GeneratedType } from "@cosmjs/proto-signing";
import { QueryParamsResponse } from "./types/arkeo/claim/query";
import { MsgAddClaimResponse } from "./types/arkeo/claim/tx";
import { MsgClaimThorchainResponse } from "./types/arkeo/claim/tx";
import { MsgClaimArkeoResponse } from "./types/arkeo/claim/tx";
import { MsgTransferClaim } from "./types/arkeo/claim/tx";
import { QueryClaimRecordRequest } from "./types/arkeo/claim/query";
import { ClaimRecord } from "./types/arkeo/claim/claim_record";
import { Params } from "./types/arkeo/claim/params";
import { GenesisState } from "./types/arkeo/claim/genesis";
import { MsgClaimEth } from "./types/arkeo/claim/tx";
import { MsgClaimEthResponse } from "./types/arkeo/claim/tx";
import { MsgTransferClaimResponse } from "./types/arkeo/claim/tx";
import { MsgAddClaim } from "./types/arkeo/claim/tx";
import { QueryParamsRequest } from "./types/arkeo/claim/query";
import { QueryClaimRecordResponse } from "./types/arkeo/claim/query";
import { MsgClaimThorchain } from "./types/arkeo/claim/tx";
import { MsgClaimArkeo } from "./types/arkeo/claim/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/arkeo.claim.QueryParamsResponse", QueryParamsResponse],
    ["/arkeo.claim.MsgAddClaimResponse", MsgAddClaimResponse],
    ["/arkeo.claim.MsgClaimThorchainResponse", MsgClaimThorchainResponse],
    ["/arkeo.claim.MsgClaimArkeoResponse", MsgClaimArkeoResponse],
    ["/arkeo.claim.MsgTransferClaim", MsgTransferClaim],
    ["/arkeo.claim.QueryClaimRecordRequest", QueryClaimRecordRequest],
    ["/arkeo.claim.ClaimRecord", ClaimRecord],
    ["/arkeo.claim.Params", Params],
    ["/arkeo.claim.GenesisState", GenesisState],
    ["/arkeo.claim.MsgClaimEth", MsgClaimEth],
    ["/arkeo.claim.MsgClaimEthResponse", MsgClaimEthResponse],
    ["/arkeo.claim.MsgTransferClaimResponse", MsgTransferClaimResponse],
    ["/arkeo.claim.MsgAddClaim", MsgAddClaim],
    ["/arkeo.claim.QueryParamsRequest", QueryParamsRequest],
    ["/arkeo.claim.QueryClaimRecordResponse", QueryClaimRecordResponse],
    ["/arkeo.claim.MsgClaimThorchain", MsgClaimThorchain],
    ["/arkeo.claim.MsgClaimArkeo", MsgClaimArkeo],
    
];

export { msgTypes }