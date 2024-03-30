import { EventEmitter, Injectable } from '@angular/core';

import { Message, SocketEventType } from './domain.message'

@Injectable({
  providedIn: 'root'
})
export class SocketService {
    private socket?: WebSocket;
    private listener: EventEmitter<any> = new EventEmitter();

  public connect() {
    this.socket = new WebSocket("ws://localhost:8000/ws");
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
    this.socket?.send(JSON.stringify(msg));
  }

  public setName(name: string) {
    this.socket?.send(name)
  }

  public close() {
    if (this.socket) {
      this.socket.onclose = null;
      this.socket.close();
    }
  }

  public getEventListener() {
    return this.listener;
  }
}
