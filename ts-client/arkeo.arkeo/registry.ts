import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgOpenContract } from "./types/arkeo/tx";
import { MsgModProvider } from "./types/arkeo/tx";
import { MsgCloseContract } from "./types/arkeo/tx";
import { MsgBondProvider } from "./types/arkeo/tx";
import { MsgClaimContractIncome } from "./types/arkeo/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/arkeo.arkeo.MsgOpenContract", MsgOpenContract],
    ["/arkeo.arkeo.MsgModProvider", MsgModProvider],
    ["/arkeo.arkeo.MsgCloseContract", MsgCloseContract],
    ["/arkeo.arkeo.MsgBondProvider", MsgBondProvider],
    ["/arkeo.arkeo.MsgClaimContractIncome", MsgClaimContractIncome],
    
];

export { msgTypes }