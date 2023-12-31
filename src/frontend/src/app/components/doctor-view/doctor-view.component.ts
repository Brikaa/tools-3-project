import { CommonModule } from '@angular/common';
import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { DoctorAppointment, Slot, UserContext } from '../../types';
import { setEntities, withPromptValues } from '../common/common';
import { isSuccessResponse, sendRequest } from '../../httpClient';

@Component({
  selector: 'doctor-view',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './doctor-view.component.html',
  styleUrls: ['../common/table.css']
})
export class DoctorViewComponent implements OnInit, OnDestroy {
  @Input({ required: true }) ctx!: UserContext;
  slots: Slot[] = [];
  appointments: DoctorAppointment[] = [];
  ws: WebSocket | null = null;

  setSlots() {
    setEntities(this.ctx, this.slots, 'slots', (body) => {
      this.slots = body['slots'];
    });
  }

  setAppointments() {
    setEntities(this.ctx, this.appointments, 'doctor-appointments', (body) => {
      this.appointments = body['appointments'];
    });
  }

  putSlot(endpoint: string, alertMessage: string) {
    withPromptValues(
      async (start, end) => {
        const res = await sendRequest(this.ctx.token, 'PUT', endpoint, {
          start,
          end
        });
        if (!isSuccessResponse(res)) {
          return;
        }
        this.setSlots();
        this.setAppointments();
        alert(alertMessage);
      },
      'Start date (RFC 3339 format)',
      'End date (RFC 3339 format)'
    );
  }

  createSlot() {
    this.putSlot('slots', 'Slot has been added!');
  }

  editSlot(id: string) {
    this.putSlot(`slots/${id}`, 'Slot has been edited!');
  }

  async deleteSlot(id: string) {
    const res = await sendRequest(this.ctx.token, 'DELETE', `slots/${id}`);
    if (!isSuccessResponse(res)) {
      return;
    }
    alert('Slot has been deleted!');
    this.setSlots();
    this.setAppointments();
  }

  ngOnInit() {
    this.setSlots();
    this.setAppointments();
    const protocol = window.location.protocol === 'http:' ? 'ws' : 'wss';
    this.ws = new WebSocket(
      `${protocol}://${location.host}/api/doctor-appointments/ws?token=${this.ctx.token}`
    );
    this.ws.addEventListener('message', (message) => {
      if (typeof message.data === 'string') {
        alert(message.data);
      }
      this.setSlots();
      this.setAppointments();
    });
  }

  ngOnDestroy(): void {
    this.ws?.close();
  }
}
