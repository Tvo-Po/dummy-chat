import { CommonModule } from '@angular/common';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { FormsModule } from '@angular/forms';

import { Message, SocketEventType, SocketEvent } from "./domain.message"
import { SocketService } from "./socket.service";

@Component({
  standalone: true,
  selector: 'app-root',
  templateUrl: './app.component.html',
  imports: [CommonModule, FormsModule]
})
export class AppComponent implements OnInit, OnDestroy {

  public messages: Array<string>;
  public chatBox: string;
  private name: string;

  public constructor(private socket: SocketService) {
    this.messages = [];
    this.chatBox = "";
    this.name = "";
  }

  public ngOnInit() {
    let handler: (event: SocketEvent) => void = event => {
      switch (event.type) {
        case SocketEventType.Open:
          this.messages.push("/The socket connection has been established. Enter your name.");
          break;
        case SocketEventType.Close:
          this.messages.push("/The socket connection has been closed.");
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
    return message.startsWith("/") ? "<strong>" + message.substring(1) + "</strong>" : message;
  }
}
