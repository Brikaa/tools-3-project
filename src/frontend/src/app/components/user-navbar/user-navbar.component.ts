import { Component, EventEmitter, Input, Output } from '@angular/core';
import { UserContext } from '../../types';
import { LOCAL_STORAGE_TOKEN } from '../../constants';

@Component({
  selector: 'user-navbar',
  standalone: true,
  templateUrl: './user-navbar.component.html'
})
export class UserNavbarComponent {
  @Input({ required: true }) ctx!: UserContext;
  @Output() userCtxEvent = new EventEmitter<string>();

  logout() {
    localStorage.removeItem(LOCAL_STORAGE_TOKEN);
    this.userCtxEvent.emit();
  }
}
