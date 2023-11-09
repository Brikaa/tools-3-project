import { Component, Input, OnInit } from '@angular/core';
import { isSuccessResponse, sendRequest } from '../../httpClient';
import { Doctor, PatientAppointment, Slot, UserContext } from '../../types';
import { CommonModule } from '@angular/common';
import { setEntities, withPromptValues } from '../common/common';

@Component({
  selector: 'patient-view',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './patient-view.component.html'
})
export class PatientViewComponent implements OnInit {
  @Input({ required: true }) ctx!: UserContext;
  doctors: Doctor[] = [];
  slots: Slot[] = [];
  appointments: PatientAppointment[] = [];
  selectedDoctorId: string = '';

  setDoctors() {
    setEntities(this.ctx, this.doctors, 'doctors', (body) => {
      this.doctors = body['doctors'];
    });
  }

  #setSelectedDoctorId(id: string) {
    this.selectedDoctorId = id;
    setEntities(this.ctx, this.slots, `doctors/${id}/slots`, (body) => {
      this.slots = body['slots'];
    });
  }

  onSelectedDoctorChange(target: EventTarget) {
    this.#setSelectedDoctorId((target as HTMLOptionElement).value);
  }

  setAppointments() {
    setEntities(this.ctx, this.appointments, 'appointments', (body) => {
      this.appointments = body['appointments'];
    });
  }

  #refreshAppointments() {
    this.#setSelectedDoctorId(this.selectedDoctorId);
    this.setAppointments();
  }

  async scheduleAppointment(slotId: string) {
    const res = await sendRequest(this.ctx, 'PUT', 'appointments', { slotId });
    if (!isSuccessResponse(res)) return;
    alert('Appointment scheduled!');
    this.#refreshAppointments();
  }

  async cancelAppointment(id: string) {
    const res = await sendRequest(this.ctx, 'DELETE', `appointments/${id}`);
    if (!isSuccessResponse(res)) return;
    alert('Appointment cancelled!');
    this.#refreshAppointments();
  }

  async editAppointmentSlot(id: string) {
    withPromptValues(async (slotId) => {
      const res = await sendRequest(this.ctx, 'PUT', `appointments/${id}`, { slotId });
      if (!isSuccessResponse(res)) return;
      alert('Appointment modified!');
      this.#refreshAppointments();
    }, 'Enter the new slot id');
  }

  ngOnInit() {
    this.setDoctors();
    this.setAppointments();
  }
}
