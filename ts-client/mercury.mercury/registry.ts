import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgBondProvider } from "./types/mercury/tx";
import { MsgModProvider } from "./types/mercury/tx";
import { MsgOpenContract } from "./types/mercury/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/mercury.mercury.MsgBondProvider", MsgBondProvider],
    ["/mercury.mercury.MsgModProvider", MsgModProvider],
    ["/mercury.mercury.MsgOpenContract", MsgOpenContract],
    
];

export { msgTypes }