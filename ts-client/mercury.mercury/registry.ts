import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgCloseContract } from "./types/mercury/tx";
import { MsgClaimContractIncome } from "./types/mercury/tx";
import { MsgBondProvider } from "./types/mercury/tx";
import { MsgOpenContract } from "./types/mercury/tx";
import { MsgModProvider } from "./types/mercury/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/mercury.mercury.MsgCloseContract", MsgCloseContract],
    ["/mercury.mercury.MsgClaimContractIncome", MsgClaimContractIncome],
    ["/mercury.mercury.MsgBondProvider", MsgBondProvider],
    ["/mercury.mercury.MsgOpenContract", MsgOpenContract],
    ["/mercury.mercury.MsgModProvider", MsgModProvider],
    
];

export { msgTypes }