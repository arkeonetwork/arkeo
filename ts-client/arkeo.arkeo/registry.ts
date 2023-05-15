import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgBondProvider } from "./types/arkeo/arkeo/tx";
import { MsgSetVersion } from "./types/arkeo/arkeo/tx";
import { MsgClaimContractIncome } from "./types/arkeo/arkeo/tx";
import { MsgModProvider } from "./types/arkeo/arkeo/tx";
import { MsgOpenContract } from "./types/arkeo/arkeo/tx";
import { MsgCloseContract } from "./types/arkeo/arkeo/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/arkeo.arkeo.MsgBondProvider", MsgBondProvider],
    ["/arkeo.arkeo.MsgSetVersion", MsgSetVersion],
    ["/arkeo.arkeo.MsgClaimContractIncome", MsgClaimContractIncome],
    ["/arkeo.arkeo.MsgModProvider", MsgModProvider],
    ["/arkeo.arkeo.MsgOpenContract", MsgOpenContract],
    ["/arkeo.arkeo.MsgCloseContract", MsgCloseContract],
    
];

export { msgTypes }