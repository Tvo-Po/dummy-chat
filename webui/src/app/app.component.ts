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

  public messages: Array<Message>;
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
          this.socket.setName(this.name);
          break;
        case SocketEventType.Close:
          let prevState = this.state
          this.state = ChatState.Connecting
          if (prevState == ChatState.Connected) {
            this.connect();
          }
          break;
        case SocketEventType.Message:
          this.messages.push(event.message);
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

  public send(last: HTMLElement) {
    if (!this.chatBox) {
      return;
    }
    this.socket.send({Sender: this.name, Content: this.chatBox});
    this.chatBox = "";
    setTimeout(() => {
      last.scrollIntoView({behavior: 'smooth'});
    }, 100);
  }

  public getMessageClass(message: Message) {
    return message.Sender == "_system_" ? "chat__message_system" : "chat__message";
  }
}
