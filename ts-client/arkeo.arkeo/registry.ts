import { GeneratedType } from "@cosmjs/proto-signing";
import { Params } from "./types/arkeo/arkeo/params";
import { QueryFetchContractRequest } from "./types/arkeo/arkeo/query";
import { MsgBondProvider } from "./types/arkeo/arkeo/tx";
import { MsgClaimContractIncome } from "./types/arkeo/arkeo/tx";
import { EventModProvider } from "./types/arkeo/arkeo/events";
import { QueryAllContractResponse } from "./types/arkeo/arkeo/query";
import { ProtoStrings } from "./types/arkeo/arkeo/misc";
import { GenesisState } from "./types/arkeo/arkeo/genesis";
import { Contract } from "./types/arkeo/arkeo/keeper";
import { ContractExpirationSet } from "./types/arkeo/arkeo/keeper";
import { MsgOpenContractResponse } from "./types/arkeo/arkeo/tx";
import { EventOpenContract } from "./types/arkeo/arkeo/events";
import { QueryParamsResponse } from "./types/arkeo/arkeo/query";
import { MsgCloseContractResponse } from "./types/arkeo/arkeo/tx";
import { ProtoInt64 } from "./types/arkeo/arkeo/misc";
import { ProtoUint64 } from "./types/arkeo/arkeo/misc";
import { QueryAllProviderRequest } from "./types/arkeo/arkeo/query";
import { MsgClaimContractIncomeResponse } from "./types/arkeo/arkeo/tx";
import { MsgSetVersion } from "./types/arkeo/arkeo/tx";
import { UserContractSet } from "./types/arkeo/arkeo/keeper";
import { QueryActiveContractRequest } from "./types/arkeo/arkeo/query";
import { MsgModProviderResponse } from "./types/arkeo/arkeo/tx";
import { MsgSetVersionResponse } from "./types/arkeo/arkeo/tx";
import { MsgModProvider } from "./types/arkeo/arkeo/tx";
import { EventCloseContract } from "./types/arkeo/arkeo/events";
import { QueryAllContractRequest } from "./types/arkeo/arkeo/query";
import { ProtoBools } from "./types/arkeo/arkeo/misc";
import { EventSettleContract } from "./types/arkeo/arkeo/events";
import { QueryActiveContractResponse } from "./types/arkeo/arkeo/query";
import { EventBondProvider } from "./types/arkeo/arkeo/events";
import { QueryFetchContractResponse } from "./types/arkeo/arkeo/query";
import { MsgBondProviderResponse } from "./types/arkeo/arkeo/tx";
import { EventValidatorPayout } from "./types/arkeo/arkeo/events";
import { QueryAllProviderResponse } from "./types/arkeo/arkeo/query";
import { MsgOpenContract } from "./types/arkeo/arkeo/tx";
import { QueryParamsRequest } from "./types/arkeo/arkeo/query";
import { QueryFetchProviderRequest } from "./types/arkeo/arkeo/query";
import { ProtoAccAddresses } from "./types/arkeo/arkeo/misc";
import { MsgCloseContract } from "./types/arkeo/arkeo/tx";
import { Provider } from "./types/arkeo/arkeo/keeper";
import { ContractSet } from "./types/arkeo/arkeo/keeper";
import { QueryFetchProviderResponse } from "./types/arkeo/arkeo/query";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/arkeo.arkeo.Params", Params],
    ["/arkeo.arkeo.QueryFetchContractRequest", QueryFetchContractRequest],
    ["/arkeo.arkeo.MsgBondProvider", MsgBondProvider],
    ["/arkeo.arkeo.MsgClaimContractIncome", MsgClaimContractIncome],
    ["/arkeo.arkeo.EventModProvider", EventModProvider],
    ["/arkeo.arkeo.QueryAllContractResponse", QueryAllContractResponse],
    ["/arkeo.arkeo.ProtoStrings", ProtoStrings],
    ["/arkeo.arkeo.GenesisState", GenesisState],
    ["/arkeo.arkeo.Contract", Contract],
    ["/arkeo.arkeo.ContractExpirationSet", ContractExpirationSet],
    ["/arkeo.arkeo.MsgOpenContractResponse", MsgOpenContractResponse],
    ["/arkeo.arkeo.EventOpenContract", EventOpenContract],
    ["/arkeo.arkeo.QueryParamsResponse", QueryParamsResponse],
    ["/arkeo.arkeo.MsgCloseContractResponse", MsgCloseContractResponse],
    ["/arkeo.arkeo.ProtoInt64", ProtoInt64],
    ["/arkeo.arkeo.ProtoUint64", ProtoUint64],
    ["/arkeo.arkeo.QueryAllProviderRequest", QueryAllProviderRequest],
    ["/arkeo.arkeo.MsgClaimContractIncomeResponse", MsgClaimContractIncomeResponse],
    ["/arkeo.arkeo.MsgSetVersion", MsgSetVersion],
    ["/arkeo.arkeo.UserContractSet", UserContractSet],
    ["/arkeo.arkeo.QueryActiveContractRequest", QueryActiveContractRequest],
    ["/arkeo.arkeo.MsgModProviderResponse", MsgModProviderResponse],
    ["/arkeo.arkeo.MsgSetVersionResponse", MsgSetVersionResponse],
    ["/arkeo.arkeo.MsgModProvider", MsgModProvider],
    ["/arkeo.arkeo.EventCloseContract", EventCloseContract],
    ["/arkeo.arkeo.QueryAllContractRequest", QueryAllContractRequest],
    ["/arkeo.arkeo.ProtoBools", ProtoBools],
    ["/arkeo.arkeo.EventSettleContract", EventSettleContract],
    ["/arkeo.arkeo.QueryActiveContractResponse", QueryActiveContractResponse],
    ["/arkeo.arkeo.EventBondProvider", EventBondProvider],
    ["/arkeo.arkeo.QueryFetchContractResponse", QueryFetchContractResponse],
    ["/arkeo.arkeo.MsgBondProviderResponse", MsgBondProviderResponse],
    ["/arkeo.arkeo.EventValidatorPayout", EventValidatorPayout],
    ["/arkeo.arkeo.QueryAllProviderResponse", QueryAllProviderResponse],
    ["/arkeo.arkeo.MsgOpenContract", MsgOpenContract],
    ["/arkeo.arkeo.QueryParamsRequest", QueryParamsRequest],
    ["/arkeo.arkeo.QueryFetchProviderRequest", QueryFetchProviderRequest],
    ["/arkeo.arkeo.ProtoAccAddresses", ProtoAccAddresses],
    ["/arkeo.arkeo.MsgCloseContract", MsgCloseContract],
    ["/arkeo.arkeo.Provider", Provider],
    ["/arkeo.arkeo.ContractSet", ContractSet],
    ["/arkeo.arkeo.QueryFetchProviderResponse", QueryFetchProviderResponse],
    
];

export { msgTypes }