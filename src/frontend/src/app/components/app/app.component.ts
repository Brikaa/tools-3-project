import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { User, UserContext } from '../../types';
import { LOCAL_STORAGE_TOKEN } from '../../constants';
import { isSuccessResponse, sendRequest } from '../../httpClient';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  ctx: UserContext | null = null;

  async setUserCtx() {
    const token = localStorage.getItem(LOCAL_STORAGE_TOKEN);
    if (token) {
      const res = await sendRequest(this.ctx, 'GET', 'user');
      if (isSuccessResponse(res)) {
        const user: User = await res.json();
        this.ctx = {
          id: user.id,
          role: user.role,
          token,
          username: user.username
        };
      }
    }
  }

  ngOnInit(): void {
    this.setUserCtx();
  }
}
