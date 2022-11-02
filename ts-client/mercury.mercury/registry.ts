import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgClaimContractIncome } from "./types/mercury/tx";
import { MsgOpenContract } from "./types/mercury/tx";
import { MsgModProvider } from "./types/mercury/tx";
import { MsgBondProvider } from "./types/mercury/tx";
import { MsgCloseContract } from "./types/mercury/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/mercury.mercury.MsgClaimContractIncome", MsgClaimContractIncome],
    ["/mercury.mercury.MsgOpenContract", MsgOpenContract],
    ["/mercury.mercury.MsgModProvider", MsgModProvider],
    ["/mercury.mercury.MsgBondProvider", MsgBondProvider],
    ["/mercury.mercury.MsgCloseContract", MsgCloseContract],
    
];

export { msgTypes }