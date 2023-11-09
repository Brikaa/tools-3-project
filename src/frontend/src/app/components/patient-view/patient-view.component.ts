import { Component, Input, OnInit } from '@angular/core';
import { isSuccessResponse, sendRequest } from '../../httpClient';
import { Doctor, PatientAppointment, Slot, UserContext } from '../../types';

@Component({
  selector: 'patient-view',
  standalone: true,
  templateUrl: './patient-view.component.html'
})
export class PatientViewComponent implements OnInit {
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

  #refreshAppointments() {
    this.setSelectedDoctorId(this.selectedDoctorId);
    this.setAppointments();
  }

  async scheduleAppointment(slotId: string) {
    const res = await sendRequest(this.ctx, 'PUT', '/appointments', { slotId });
    if (!isSuccessResponse(res)) return;
    alert('Appointment scheduled!');
    this.#refreshAppointments();
  }

  async cancelAppointment(id: string) {
    const res = await sendRequest(this.ctx, 'DELETE', `/appointments/${id}`);
    if (!isSuccessResponse(res)) return;
    alert('Appointment cancelled!');
    this.#refreshAppointments();
  }

  async editAppointmentSlot(id: string) {
    const slotId = prompt('Enter the new slot id');
    const res = await sendRequest(this.ctx, 'PUT', `/appointments/${id}`, { slotId });
    if (!isSuccessResponse(res)) return;
    alert('Appointment modified!');
    this.#refreshAppointments();
  }

  ngOnInit() {
    this.setDoctors();
    this.setAppointments();
  }
}
