import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgModProvider } from "./types/arkeo/arkeo/tx";
import { MsgClaimContractIncome } from "./types/arkeo/arkeo/tx";
import { MsgCloseContract } from "./types/arkeo/arkeo/tx";
import { MsgBondProvider } from "./types/arkeo/arkeo/tx";
import { MsgOpenContract } from "./types/arkeo/arkeo/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/arkeo.arkeo.MsgModProvider", MsgModProvider],
    ["/arkeo.arkeo.MsgClaimContractIncome", MsgClaimContractIncome],
    ["/arkeo.arkeo.MsgCloseContract", MsgCloseContract],
    ["/arkeo.arkeo.MsgBondProvider", MsgBondProvider],
    ["/arkeo.arkeo.MsgOpenContract", MsgOpenContract],
    
];

export { msgTypes }