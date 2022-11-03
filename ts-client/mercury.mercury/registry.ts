import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgOpenContract } from "./types/mercury/tx";
import { MsgBondProvider } from "./types/mercury/tx";
import { MsgModProvider } from "./types/mercury/tx";
import { MsgCloseContract } from "./types/mercury/tx";
import { MsgClaimContractIncome } from "./types/mercury/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/mercury.mercury.MsgOpenContract", MsgOpenContract],
    ["/mercury.mercury.MsgBondProvider", MsgBondProvider],
    ["/mercury.mercury.MsgModProvider", MsgModProvider],
    ["/mercury.mercury.MsgCloseContract", MsgCloseContract],
    ["/mercury.mercury.MsgClaimContractIncome", MsgClaimContractIncome],
    
];

export { msgTypes }