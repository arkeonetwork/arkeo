import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgRegisterProvider } from "./types/mercury/tx";
import { MsgBondProvider } from "./types/mercury/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/mercury.mercury.MsgRegisterProvider", MsgRegisterProvider],
    ["/mercury.mercury.MsgBondProvider", MsgBondProvider],
    
];

export { msgTypes }