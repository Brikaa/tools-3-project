import { Component, EventEmitter, Input, Output } from '@angular/core';
import { UserContext } from '../../types';

@Component({
  selector: 'user-navbar',
  standalone: true,
  templateUrl: './user-navbar.component.html'
})
export class UserNavbarComponent {
  @Input({ required: true }) ctx!: UserContext;
  @Output() userCtxEvent = new EventEmitter<string>();

  logout() {
    this.userCtxEvent.emit();
  }
}
