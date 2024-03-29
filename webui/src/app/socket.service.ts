import { EventEmitter, Injectable } from '@angular/core';

import { Message, SocketEventType } from './domain.message'

@Injectable({
  providedIn: 'root'
})
export class SocketService {
    private socket: WebSocket;
    private listener: EventEmitter<any> = new EventEmitter();

  constructor() {
    this.socket = new WebSocket("ws://0.0.0.0:8000/ws");
    this.socket.onopen = event => {
      this.listener.emit({"type": SocketEventType.Open});
    }
    this.socket.onclose = event => {
      this.listener.emit({"type": SocketEventType.Close});
    }
    this.socket.onmessage = event => {
      this.listener.emit({"type": SocketEventType.Message, "message": JSON.parse(event.data)});
    }
  }

  public send(msg: Message) {
    this.socket.send(JSON.stringify(msg));
  }

  public setName(name: string) {
    this.socket.send(name)
  }

  public close() {
    this.socket.close();
  }

  public getEventListener() {
    return this.listener;
  }
}
