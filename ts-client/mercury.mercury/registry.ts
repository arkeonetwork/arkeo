import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgModProvider } from "./types/mercury/tx";
import { MsgBondProvider } from "./types/mercury/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/mercury.mercury.MsgModProvider", MsgModProvider],
    ["/mercury.mercury.MsgBondProvider", MsgBondProvider],
    
];

export { msgTypes }