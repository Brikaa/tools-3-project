import requests
import json
import datetime
import os

BASE = f"http://backend_runner:{os.environ['BACKEND_PORT']}"
headers = {"Content-Type": "application/json"}


def create_url(endpoint):
    return BASE + "/" + endpoint


def send_request(method, endpoint, payload):
    url = create_url(endpoint)
    return requests.request(method, url, headers=headers, data=json.dumps(payload))


def action(message, function):
    print(message)
    result = function()
    print(result, "\n")
    return result


def signup(username, password, role):
    res = send_request(
        "POST", "signup", {"username": username, "password": password, "role": role}
    )
    return res.text, res.status_code


def login(username, password):
    res = send_request("POST", "login", {"username": username, "password": password})
    status = res.status_code
    if status == 200:
        headers["Authorization"] = "Basic " + res.json()["token"]
    return res.text, status


def create_slot(start, end):
    res = send_request(
        "PUT",
        "slots",
        {
            "start": datetime.datetime.isoformat(start),
            "end": datetime.datetime.isoformat(end),
        },
    )
    return res.text, res.status_code


def update_slot(id, start, end):
    res = send_request(
        "PUT",
        f"slots/{id}",
        {
            "start": datetime.datetime.isoformat(start),
            "end": datetime.datetime.isoformat(end),
        },
    )
    return res.text, res.status_code


def get_slots():
    res = send_request("GET", "slots", None)
    return res.text, res.status_code


def get_doctor_appointments():
    res = send_request("GET", "doctor-appointments", None)
    return res.text, res.status_code


def get_doctors():
    res = send_request("GET", "doctors", None)
    return res.text, res.status_code


def get_available_slots_for_doctor(id):
    res = send_request("GET", f"doctor/{id}/slots", None)
    return res.text, res.status_code


def create_appointment(slot_id):
    res = send_request("PUT", f"appointments", {"slotId": slot_id})
    return res.text, res.status_code


def get_slot_id_by_start_date(slots, start):
    return list(
        filter(lambda slot: datetime.datetime.fromisoformat(slot["start"]) == start),
        slots,
    )[0]["id"]


def get_doctor_id_by_username(doctors, username):
    return list(filter(lambda doctor: doctor["username"] == username))[0]["id"]


def get_appointments():
    res = send_request("GET", "appointments", None)
    return res.text, res.status_code


def get_appointment_id_by_slot_id(appointments, slot_id):
    return list(filter(lambda appointment: appointment["slotId"] == slot_id))[0]["id"]


def delete_appointment(id):
    res = send_request("DELETE", f"appointments/{id}", None)
    return res.text, res.status_code


def delete_slot(id):
    res = send_request("DELETE", f"slots/{id}", None)
    return res.text, res.status_code


def update_appointment(id, slot_id):
    res = send_request("PUT", f"appointments/{id}", {"slotId": slot_id})
    return res.text, res.status_code


if __name__ == "__main__":
    p1_username = "p1"
    p1_password = "p1123"
    p2_username = "p2"
    p2_password = "p2123"

    d1_username = "d1"
    d1_password = "d1123"
    d2_username = "d2"
    d2_password = "d2123"

    d1s1_start = datetime.datetime.now()
    d1s1_delta = datetime.timedelta(hours=1)
    d1s1_end = d1s1_start + d1s1_delta

    d1s2_start = d1s1_end + datetime.timedelta(hours=2)
    d1s2_end = d1s2_start + datetime.timedelta(hours=2)

    d1s3_start = d1s2_end + datetime.timedelta(hours=1)
    d1s3_end = d1s3_start + datetime.timedelta(hours=1)

    d2s1_start = datetime.datetime.now() + datetime.timedelta(minutes=30)
    d2s1_end = d2s1_start + datetime.timedelta(hours=2)

    d2s2_start = d2s1_end + datetime.timedelta(hours=3)
    d2s2_end = d2s2_start + datetime.timedelta(hours=3)

    d2s3_start = d2s2_end + datetime.timedelta(hours=4)
    d2s3_end = d2s3_start + datetime.timedelta(hours=4)

    action("Signup p1", lambda: signup(p1_username, p1_password, "patient"))
    action("Signup p2", lambda: signup(p2_username, p2_password, "patient"))
    action("Signup d1", lambda: signup(d1_username, d1_password, "doctor"))
    action("Signup d2", lambda: signup(d2_username, d2_password, "doctor"))
    action("Signup d2 (invalid)", lambda: signup(d2_username, d2_password, "doctor"))
    action("Signup i1 (invalid role)", lambda: signup("i1", "i123", "doctors"))
    action("Signup i2 (invalid username)", lambda: signup("i-2", "i123", "doctor"))

    action("Login p2 (invalid password)", lambda: login(p2_username, p2_password + "a"))
    action("Login p2", lambda: login(p2_username, p2_password))
    action("Create valid slot (forbidden)", lambda: create_slot(d1s1_start, d1s1_end))

    action("Login d1", lambda: login(d1_username, d1_password))
    action("Get slots ([])", get_slots)
    action(
        "Create slot in the past (invalid)",
        lambda: create_slot(
            datetime.datetime.now() - datetime.timedelta(hours=1), d1s1_end
        ),
    )
    action(
        "Create slot with end before start (invalid)",
        lambda: create_slot(d1s1_end, d1s1_start),
    )
    action("Create slot d1s1", lambda: create_slot(d1s1_start, d1s1_end))
    action(
        "Create slot overlapping with d1s1 at target start (invalid)",
        lambda: create_slot(d1s1_start + d1s1_delta / 2, d1s1_end + d1s1_delta / 2),
    )
    action(
        "Create slot overlapping with d1s1 at target end (invalid)",
        lambda: create_slot(d1s1_start - d1s1_delta / 2, d1s1_end - d1s1_delta / 2),
    )
    action("Create slot d1s2", lambda: create_slot(d1s2_start, d1s2_end))
    action("Create slot d1s3", lambda: create_slot(d1s3_start, d1s3_end))

    slots_text, _ = action(
        f"Get slots ([{d1s1_start}, {d1s2_start}, {d1s3_start}])", get_slots
    )
    slots = json.loads(slots_text)["slots"]
    d1s1 = get_slot_id_by_start_date(slots, d1s1_start)
    d1s2 = get_slot_id_by_start_date(slots, d1s2_start)
    d1s3 = get_slot_id_by_start_date(slots, d1s3_start)

    action("Get appointments ([])", get_doctor_appointments)

    action("Login d2", lambda: login(d2_username, d2_password))
    action("Get slots ([])", get_slots)
    action("Create slot d2s1", lambda: create_slot(d2s1_start, d2s1_end))
    action("Create slot d2s2", lambda: create_slot(d2s2_start, d2s2_end))
    action("Create slot d2s3", lambda: create_slot(d2s3_start, d2s3_end))

    slots_text, _ = action(
        f"Get slots ([{d2s1_start}, {d2s2_start}, {d2s3_start}])", get_slots
    )
    slots = json.loads(slots_text)["slots"]
    d2s1 = get_slot_id_by_start_date(slots, d2s1_start)
    d2s2 = get_slot_id_by_start_date(slots, d2s2_start)
    d2s3 = get_slot_id_by_start_date(slots, d2s3_start)

    action("Login p1", lambda: login(p1_username, p1_password))

    doctors_text, _ = action("Get doctors ([d1, d2])", get_doctors)
    doctors = json.loads(doctors)["doctors"]
    d1 = get_doctor_id_by_username(doctors, "d1")
    d2 = get_doctor_id_by_username(doctors, "d2")

    action(
        f"Get slots for d1 ([{d1s1}, {d1s2}, {d1s3}])",
        lambda: get_available_slots_for_doctor(d1),
    )
    action(
        f"Get slots for d2 ([{d2s1}, {d2s2}, {d2s3}])",
        lambda: get_available_slots_for_doctor(d2),
    )

    action(f"Reserve slot d1s1 (p1a1d1s1)", lambda: create_appointment(d1s1))
    action(f"Reserve slot d2s2 (p1a2d2s2)", lambda: create_appointment(d2s2))
    appointments_text, _ = action(
        "Get appointments ([p1a1d1s1, p1a2d2s2])", get_appointments
    )
    appointments = json.loads(appointments_text)["appointments"]
    p1a1d1s1 = get_appointment_id_by_slot_id(appointments, d1s1)
    p1a2d2s2 = get_appointment_id_by_slot_id(appointments, d2s2)

    action(
        f"Get slots for d1 ([{d1s2}, {d1s3}])",
        lambda: get_available_slots_for_doctor(d1),
    )
    action(
        f"Get slots for d2 ([{d2s1}, {d2s3}])",
        lambda: get_available_slots_for_doctor(d2),
    )

    action("Login p2", lambda: login(p2_username, p2_password))
    action(
        f"Get slots for d1 ([{d1s2}, {d1s3}])",
        lambda: get_available_slots_for_doctor(d1),
    )
    action(
        f"Get slots for d2 ([{d2s1}, {d2s3}])",
        lambda: get_available_slots_for_doctor(d2),
    )
    action(f"Reserve slot d1s1 (invalid)", lambda: create_appointment(d1s1))
    action(f"Reserve slot d1s3 (p2a1d1s3)", lambda: create_appointment(d1s3))
    action(f"Reserve slot d2s3 (p2a2d2s3)", lambda: create_appointment(d2s3))

    appointments_text, _ = action(
        "Get appointments ([p2a1d1s3, p2a2d2s3])", get_appointments
    )
    appointments = json.loads(appointments_text)["appointments"]
    p2a1d1s3 = get_appointment_id_by_slot_id(appointments, d1s3)
    p2a2d2s3 = get_appointment_id_by_slot_id(appointments, d2s3)

    action(f"Get slots for d1 ([{d1s2}])", lambda: get_available_slots_for_doctor(d1))
    action(f"Get slots for d2 ([{d2s1}])", lambda: get_available_slots_for_doctor(d2))
    action(
        f"Delete appointment p1a1d1s1 (invalid)", lambda: delete_appointment(p1a1d1s1)
    )
    action("Delete appointment p2a1d1s3 (valid)", lambda: delete_appointment(p2a1d1s3))
    action(
        f"Get slots for d1 ([{d1s2}, {d1s3}])",
        lambda: get_available_slots_for_doctor(d1),
    )
    action(f"Get slots for d2 ([{d2s1}])", lambda: get_available_slots_for_doctor(d2))
    action(f"Get appointments ([{p2a2d2s3}])", get_appointments)

    action(
        "Update appointment p2a2d2s3 -> p2a2d2s1", update_appointment(p2a2d2s3, d2s1)
    )
    p2a2d2s1 = p2a2d2s3
    p2a2d2s3 = None
    action(f"Get slots for d2 ([{d2s3}])", lambda: get_available_slots_for_doctor(d2))
    action(
        "Update appointment p2a2d2s1 -> p2a2d1s3", update_appointment(p2a2d2s1, d1s3)
    )
    p2a2d1s3 = p2a2d2s1
    p2a2d2s1 = None
    action(f"Get slots for d1 ([{d1s2}])", lambda: get_available_slots_for_doctor(d1))
    action(
        f"Get slots for d2 ([{d2s1, d2s3}])", lambda: get_available_slots_for_doctor(d2)
    )

    action("Login d2", lambda: login(d2_username, d2_password))
    action(f"Get appointments ([{p1a2d2s2}])", get_doctor_appointments)
    action(f"Update slot d2s2", lambda: update_slot(d2s2, d1s1_start, d1s1_end))
    action(f"Delete slot d2s2", lambda: delete_slot(d2s2))
    action(f"Get appointments ([])", get_appointments)
