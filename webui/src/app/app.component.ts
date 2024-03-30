import { CommonModule } from '@angular/common';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

import { Message, SocketEventType, SocketEvent } from "./domain.message"
import { SocketService } from "./socket.service";

enum ConnectionState {
  Establishing,
  Established,
  Reconnecting,
}

@Component({
  standalone: true,
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
  imports: [CommonModule, FormsModule, MatProgressSpinnerModule]
})
export class AppComponent implements OnInit, OnDestroy {

  public messages: Array<string>;
  public chatBox: string;
  public state: ConnectionState;
  public ConnectionState: typeof ConnectionState;
  private name: string;

  public constructor(private socket: SocketService) {
    this.messages = [];
    this.chatBox = "";
    this.name = "";

    this.state = ConnectionState.Establishing;
    this.ConnectionState = ConnectionState;
  }

  public ngOnInit() {
    this.socket.connect()
    let handler: (event: SocketEvent) => void = event => {
      switch (event.type) {
        case SocketEventType.Open:
          switch (this.state) {
            case ConnectionState.Establishing:
              this.messages.push("_system_: The socket connection has been established. Enter your name.");
              break;
            case ConnectionState.Reconnecting:
              this.socket.setName(this.name)
              break;
          }
          this.state = ConnectionState.Established
          break;
        case SocketEventType.Close:
          this.state = ConnectionState.Reconnecting
          break;
        case SocketEventType.Message:
          const data = `${event.message.Sender}: ${event.message.Content}`;
          this.messages.push(data);
          break;
      }
    }
    this.socket.getEventListener().subscribe(handler);
  }

  public ngOnDestroy() {
    this.socket.close();
  }

  public send() {
    if (!this.chatBox) {
      return
    }
    if (this.name === "") {
      this.name = this.chatBox !== "" ? this.chatBox : "unnamed"
      this.socket.setName(this.chatBox);
      this.chatBox = "";
    } else {
      this.socket.send({Sender: this.name, Content: this.chatBox});
      this.chatBox = "";
    }
  }

  public isSystemMessage(message: string) {
    return message.startsWith("_system_: ") ? "<strong>" + message.substring(9) + "</strong>" : message;
  }
}
