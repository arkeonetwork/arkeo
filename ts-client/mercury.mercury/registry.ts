import { GeneratedType } from "@cosmjs/proto-signing";
import { MsgBondProvider } from "./types/mercury/tx";
import { MsgRegisterProvider } from "./types/mercury/tx";

const msgTypes: Array<[string, GeneratedType]>  = [
    ["/mercury.mercury.MsgBondProvider", MsgBondProvider],
    ["/mercury.mercury.MsgRegisterProvider", MsgRegisterProvider],
    
];

export { msgTypes }