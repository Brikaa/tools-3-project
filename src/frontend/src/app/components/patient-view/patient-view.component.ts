import { Component, Input } from '@angular/core';
import { isSuccessResponse, sendRequest } from '../../httpClient';
import { Doctor, PatientAppointment, Slot, UserContext } from '../../types';

@Component({
  selector: 'patient-view',
  standalone: true,
  templateUrl: './patient-view.component.html'
})
export class PatientViewComponent {
  @Input({ required: true }) ctx!: UserContext;
  doctors: Doctor[] = [];
  slots: Slot[] = [];
  appointments: PatientAppointment[] = [];
  selectedDoctorId: string = '';

  async #setEntities<T extends Doctor | Slot | PatientAppointment>(
    entities: T[],
    endpoint: string,
    setEntities: (body: { [key: string]: T[] }) => void
  ) {
    const res = await sendRequest(this.ctx, 'GET', endpoint);
    if (!isSuccessResponse(res)) {
      entities.length = 0;
      return;
    }
    const body: { [key: string]: T[] } = await res.json();
    setEntities(body);
  }

  setDoctors() {
    this.#setEntities(this.doctors, '/doctors', (body) => {
      this.doctors = body['doctors'];
    });
  }

  setSelectedDoctorId(id: string) {
    this.selectedDoctorId = id;
    this.#setEntities(this.slots, `/doctors/${id}/slots`, (body) => {
      this.slots = body['slots'];
    });
  }

  setAppointments() {
    this.#setEntities(this.appointments, '/appointments', (body) => {
      this.appointments = body['appointments'];
    });
  }

  async scheduleAppointment() {
    const res = await sendRequest(this.ctx, 'PUT', '/appointments');
    if (!isSuccessResponse(res)) return;
    this.setSelectedDoctorId(this.selectedDoctorId);
    this.setAppointments();
    alert('Appointment scheduled!');
  }
}
