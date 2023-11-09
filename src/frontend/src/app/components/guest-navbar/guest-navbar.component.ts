import { Component, EventEmitter, Output } from '@angular/core';
import { isSuccessResponse, sendRequest } from '../../httpClient';
import { withPromptValues } from '../common/common';

@Component({
  selector: 'guest-navbar',
  standalone: true,
  templateUrl: './guest-navbar.component.html'
})
export class GuestNavbarComponent {
  @Output() userCtxEvent = new EventEmitter<string>();

  login() {
    withPromptValues(
      async (username, password) => {
        const response = await sendRequest(null, 'POST', 'login', { username, password });
        if (!isSuccessResponse(response)) {
          return;
        }
        const body: { token: string } = await response.json();
        this.userCtxEvent.emit(body.token);
      },
      'Username',
      'Password'
    );
  }

  register() {
    withPromptValues(
      async (username, password) => {
        const response = await sendRequest(null, 'POST', 'signup', { username, password });
        if (isSuccessResponse(response)) {
          alert('Registered successfully! You can now log in with your account.');
        }
      },
      'Username',
      'Password'
    );
  }
}
