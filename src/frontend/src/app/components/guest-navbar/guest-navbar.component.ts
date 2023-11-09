import { Component, EventEmitter, Output } from '@angular/core';
import { isSuccessResponse, sendRequest } from '../../httpClient';
import { LOCAL_STORAGE_TOKEN } from '../../constants';

@Component({
  selector: 'guest-navbar',
  standalone: true,
  templateUrl: './guest-navbar.component.html'
})
export class GuestNavbarComponent {
  @Output() userCtxEvent = new EventEmitter<string>();

  #withUsernameAndPassword(fn: (username: string, password: string) => void) {
    const username = prompt('Username');
    if (!username) {
      return;
    }
    const password = prompt('Password');
    if (!password) {
      return;
    }
    fn(username, password);
  }

  login() {
    this.#withUsernameAndPassword(async (username, password) => {
      const response = await sendRequest(null, 'POST', 'login', { username, password });
      if (!isSuccessResponse(response)) {
        return;
      }
      const body = await response.json();
      localStorage.setItem(LOCAL_STORAGE_TOKEN, body.token);
      this.userCtxEvent.emit();
    });
  }

  register() {
    this.#withUsernameAndPassword(async (username, password) => {
      const response = await sendRequest(null, 'POST', 'signup', { username, password });
      if (isSuccessResponse(response)) {
        alert('Registered successfully! You can now log in with your account.');
      }
    });
  }
}
