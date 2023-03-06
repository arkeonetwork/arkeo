import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgOpenContract } from "./types/arkeo/arkeo/tx";
import { MsgClaimContractIncome } from "./types/arkeo/arkeo/tx";
import { MsgBondProvider } from "./types/arkeo/arkeo/tx";
import { MsgCloseContract } from "./types/arkeo/arkeo/tx";
import { MsgModProvider } from "./types/arkeo/arkeo/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/arkeo.arkeo.MsgOpenContract", MsgOpenContract],
    ["/arkeo.arkeo.MsgClaimContractIncome", MsgClaimContractIncome],
    ["/arkeo.arkeo.MsgBondProvider", MsgBondProvider],
    ["/arkeo.arkeo.MsgCloseContract", MsgCloseContract],
    ["/arkeo.arkeo.MsgModProvider", MsgModProvider],
    
];

export { msgTypes }