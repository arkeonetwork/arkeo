import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgClaimContractIncome } from "./types/arkeo/tx";
import { MsgBondProvider } from "./types/arkeo/tx";
import { MsgModProvider } from "./types/arkeo/tx";
import { MsgOpenContract } from "./types/arkeo/tx";
import { MsgCloseContract } from "./types/arkeo/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/arkeo.arkeo.MsgClaimContractIncome", MsgClaimContractIncome],
    ["/arkeo.arkeo.MsgBondProvider", MsgBondProvider],
    ["/arkeo.arkeo.MsgModProvider", MsgModProvider],
    ["/arkeo.arkeo.MsgOpenContract", MsgOpenContract],
    ["/arkeo.arkeo.MsgCloseContract", MsgCloseContract],
    
];

export { msgTypes }