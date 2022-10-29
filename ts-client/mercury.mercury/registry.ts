import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgModProvider } from "./types/mercury/tx";
import { MsgBondProvider } from "./types/mercury/tx";
import { MsgOpenContract } from "./types/mercury/tx";
import { MsgCloseContract } from "./types/mercury/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/mercury.mercury.MsgModProvider", MsgModProvider],
    ["/mercury.mercury.MsgBondProvider", MsgBondProvider],
    ["/mercury.mercury.MsgOpenContract", MsgOpenContract],
    ["/mercury.mercury.MsgCloseContract", MsgCloseContract],
    
];

export { msgTypes }