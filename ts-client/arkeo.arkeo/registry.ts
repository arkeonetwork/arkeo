import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgClaimContractIncomeResponse } from "./types/arkeo/arkeo/tx";
import { EventValidatorPayout } from "./types/arkeo/arkeo/events";
import { QueryParamsRequest } from "./types/arkeo/arkeo/query";
import { ProtoInt64 } from "./types/arkeo/arkeo/misc";
import { Provider } from "./types/arkeo/arkeo/keeper";
import { MsgOpenContract } from "./types/arkeo/arkeo/tx";
import { MsgOpenContractResponse } from "./types/arkeo/arkeo/tx";
import { QueryAllProviderResponse } from "./types/arkeo/arkeo/query";
import { QueryActiveContractRequest } from "./types/arkeo/arkeo/query";
import { MsgModProvider } from "./types/arkeo/arkeo/tx";
import { MsgSetVersionResponse } from "./types/arkeo/arkeo/tx";
import { MsgModProviderResponse } from "./types/arkeo/arkeo/tx";
import { EventCloseContract } from "./types/arkeo/arkeo/events";
import { ProtoBools } from "./types/arkeo/arkeo/misc";
import { MsgBondProvider } from "./types/arkeo/arkeo/tx";
import { MsgCloseContract } from "./types/arkeo/arkeo/tx";
import { EventModProvider } from "./types/arkeo/arkeo/events";
import { ProtoUint64 } from "./types/arkeo/arkeo/misc";
import { GenesisState } from "./types/arkeo/arkeo/genesis";
import { ContractSet } from "./types/arkeo/arkeo/keeper";
import { Params } from "./types/arkeo/arkeo/params";
import { MsgSetVersion } from "./types/arkeo/arkeo/tx";
import { QueryAllContractResponse } from "./types/arkeo/arkeo/query";
import { ContractExpirationSet } from "./types/arkeo/arkeo/keeper";
import { QueryFetchProviderRequest } from "./types/arkeo/arkeo/query";
import { QueryActiveContractResponse } from "./types/arkeo/arkeo/query";
import { ProtoStrings } from "./types/arkeo/arkeo/misc";
import { MsgCloseContractResponse } from "./types/arkeo/arkeo/tx";
import { QueryFetchContractResponse } from "./types/arkeo/arkeo/query";
import { MsgClaimContractIncome } from "./types/arkeo/arkeo/tx";
import { ProtoAccAddresses } from "./types/arkeo/arkeo/misc";
import { Contract } from "./types/arkeo/arkeo/keeper";
import { QueryAllProviderRequest } from "./types/arkeo/arkeo/query";
import { UserContractSet } from "./types/arkeo/arkeo/keeper";
import { MsgBondProviderResponse } from "./types/arkeo/arkeo/tx";
import { EventBondProvider } from "./types/arkeo/arkeo/events";
import { EventOpenContract } from "./types/arkeo/arkeo/events";
import { QueryFetchProviderResponse } from "./types/arkeo/arkeo/query";
import { EventSettleContract } from "./types/arkeo/arkeo/events";
import { QueryParamsResponse } from "./types/arkeo/arkeo/query";
import { QueryFetchContractRequest } from "./types/arkeo/arkeo/query";
import { QueryAllContractRequest } from "./types/arkeo/arkeo/query";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/arkeo.arkeo.MsgClaimContractIncomeResponse", MsgClaimContractIncomeResponse],
    ["/arkeo.arkeo.EventValidatorPayout", EventValidatorPayout],
    ["/arkeo.arkeo.QueryParamsRequest", QueryParamsRequest],
    ["/arkeo.arkeo.ProtoInt64", ProtoInt64],
    ["/arkeo.arkeo.Provider", Provider],
    ["/arkeo.arkeo.MsgOpenContract", MsgOpenContract],
    ["/arkeo.arkeo.MsgOpenContractResponse", MsgOpenContractResponse],
    ["/arkeo.arkeo.QueryAllProviderResponse", QueryAllProviderResponse],
    ["/arkeo.arkeo.QueryActiveContractRequest", QueryActiveContractRequest],
    ["/arkeo.arkeo.MsgModProvider", MsgModProvider],
    ["/arkeo.arkeo.MsgSetVersionResponse", MsgSetVersionResponse],
    ["/arkeo.arkeo.MsgModProviderResponse", MsgModProviderResponse],
    ["/arkeo.arkeo.EventCloseContract", EventCloseContract],
    ["/arkeo.arkeo.ProtoBools", ProtoBools],
    ["/arkeo.arkeo.MsgBondProvider", MsgBondProvider],
    ["/arkeo.arkeo.MsgCloseContract", MsgCloseContract],
    ["/arkeo.arkeo.EventModProvider", EventModProvider],
    ["/arkeo.arkeo.ProtoUint64", ProtoUint64],
    ["/arkeo.arkeo.GenesisState", GenesisState],
    ["/arkeo.arkeo.ContractSet", ContractSet],
    ["/arkeo.arkeo.Params", Params],
    ["/arkeo.arkeo.MsgSetVersion", MsgSetVersion],
    ["/arkeo.arkeo.QueryAllContractResponse", QueryAllContractResponse],
    ["/arkeo.arkeo.ContractExpirationSet", ContractExpirationSet],
    ["/arkeo.arkeo.QueryFetchProviderRequest", QueryFetchProviderRequest],
    ["/arkeo.arkeo.QueryActiveContractResponse", QueryActiveContractResponse],
    ["/arkeo.arkeo.ProtoStrings", ProtoStrings],
    ["/arkeo.arkeo.MsgCloseContractResponse", MsgCloseContractResponse],
    ["/arkeo.arkeo.QueryFetchContractResponse", QueryFetchContractResponse],
    ["/arkeo.arkeo.MsgClaimContractIncome", MsgClaimContractIncome],
    ["/arkeo.arkeo.ProtoAccAddresses", ProtoAccAddresses],
    ["/arkeo.arkeo.Contract", Contract],
    ["/arkeo.arkeo.QueryAllProviderRequest", QueryAllProviderRequest],
    ["/arkeo.arkeo.UserContractSet", UserContractSet],
    ["/arkeo.arkeo.MsgBondProviderResponse", MsgBondProviderResponse],
    ["/arkeo.arkeo.EventBondProvider", EventBondProvider],
    ["/arkeo.arkeo.EventOpenContract", EventOpenContract],
    ["/arkeo.arkeo.QueryFetchProviderResponse", QueryFetchProviderResponse],
    ["/arkeo.arkeo.EventSettleContract", EventSettleContract],
    ["/arkeo.arkeo.QueryParamsResponse", QueryParamsResponse],
    ["/arkeo.arkeo.QueryFetchContractRequest", QueryFetchContractRequest],
    ["/arkeo.arkeo.QueryAllContractRequest", QueryAllContractRequest],
    
];

export { msgTypes }