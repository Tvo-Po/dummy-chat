export interface Message {
  Sender: string;
  Content: string;
};

export enum SocketEventType {
  Open,
  Message,
  Close,
}

export type SocketEvent = {type: SocketEventType.Open | SocketEventType.Close} | 
  {type: SocketEventType.Message, message: Message};
