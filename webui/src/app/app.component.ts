import { CommonModule } from '@angular/common';
import { Component, OnInit, OnDestroy } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

import { Message, SocketEventType, SocketEvent } from "./domain.message"
import { SocketService } from "./socket.service";

enum ChatState {
  Registrating,
  Connecting,
  Connected,
}

@Component({
  standalone: true,
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
  imports: [
    CommonModule,
    FormsModule,
    MatButtonModule,
    MatInputModule,
    MatFormFieldModule,
    MatProgressSpinnerModule,
  ]
})
export class AppComponent implements OnInit, OnDestroy {

  public messages: Array<string>;
  public chatBox: string;
  public state: ChatState;
  public ChatState: typeof ChatState;
  public name: string;

  public constructor(private socket: SocketService) {
    this.messages = [];
    this.chatBox = "";
    this.name = "";

    this.state = ChatState.Registrating;
    this.ChatState = ChatState;
  }

  public ngOnInit() {
    let handler: (event: SocketEvent) => void = event => {
      switch (event.type) {
        case SocketEventType.Open:
          this.state = ChatState.Connected
          if (this.messages.length == 0) {
            this.socket.setName(this.name);
          } else {
            this.messages.push("_system_: You have been reconnected");
          }
          break;
        case SocketEventType.Close:
          let prevState = this.state
          this.state = ChatState.Connecting
          if (prevState == ChatState.Connected) {
            this.connect();
          }
          break;
        case SocketEventType.Message:
          const data = `${event.message.Sender}: ${event.message.Content}`;
          this.messages.push(data);
          break;
      }
    }
    this.socket.getEventListener().subscribe(handler)
  }

  public ngOnDestroy() {
    this.socket.close();
  }

  public connect() {
    if (this.state == ChatState.Connecting) {
      this.socket.connect();
      setTimeout(() => {
        this.connect();
      }, 5000);
    }
  }

  public enterName() {
    this.name = this.chatBox;
    this.chatBox = "";
    this.state = ChatState.Connecting;
    this.connect();
  }

  public send() {
    if (!this.chatBox) {
      return;
    }
    this.socket.send({Sender: this.name, Content: this.chatBox});
    this.chatBox = "";
  }

  public isSystemMessage(message: string) {
    return message.startsWith("_system_: ") ? "<strong>" + message.substring(9) + "</strong>" : message;
  }
}
