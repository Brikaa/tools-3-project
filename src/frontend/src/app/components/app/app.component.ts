import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { User, UserContext } from '../../types';
import { LOCAL_STORAGE_TOKEN } from '../../constants';
import { isSuccessResponse, sendRequest } from '../../httpClient';
import { GuestNavbarComponent } from '../guest-navbar/guest-navbar.component';
import { UserNavbarComponent } from '../user-navbar/user-navbar.component';
import { DoctorViewComponent } from '../doctor-view/doctor-view.component';
import { PatientViewComponent } from '../patient-view/patient-view.component';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    CommonModule,
    GuestNavbarComponent,
    UserNavbarComponent,
    DoctorViewComponent,
    PatientViewComponent
  ],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {
  ctx: UserContext | null = null;

  async setUserCtx(token: string | null | undefined) {
    if (token) {
      const res = await sendRequest(token, 'GET', 'user');
      if (isSuccessResponse(res)) {
        const user: User = await res.json();
        this.ctx = {
          id: user.id,
          role: user.role,
          token,
          username: user.username
        };
      }
      localStorage.setItem(LOCAL_STORAGE_TOKEN, token);
    } else {
      this.ctx = null;
      localStorage.removeItem(LOCAL_STORAGE_TOKEN);
    }
  }

  ngOnInit(): void {
    this.setUserCtx(localStorage.getItem(LOCAL_STORAGE_TOKEN));
  }
}
