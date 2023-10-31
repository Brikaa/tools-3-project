import requests
import json
import datetime

BASE = "http://backend_runner:8000"
token = ""
headers = {"Content-Type": "application/json"}
slots = {}


def create_url(endpoint):
    return BASE + "/" + endpoint


def send_request(method, endpoint, payload):
    url = create_url(endpoint)
    return requests.request(method, url, headers=headers, data=json.dumps(payload))


def action(message, function):
    print(message)
    result = function()
    print(result + "\n")


def signup(username, password, role):
    res = send_request(
        "POST", "signup", {"username": username, "password": password, "role": role}
    )
    return res.text, res.status_code


def login(username, password):
    res = send_request(
        "POST", "login", {"username": username, "password": password, "role": role}
    )
    status = res.status_code
    if status == 200:
        token = res.json()["token"]
    return res.text, status


def create_slot(start, end):
    res = send_request(
        "PUT", "slots", {"start": start.isoformat(), "end": end.isoformat()}
    )
    return res.text, res.status_code


def get_slots():
    res = send_request("GET", "slots", None)
    return res.text, res.status_code


def get_doctor_appointments():
    res = send_request("GET", "doctor-appointments", None)
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

    d1s3_start = d1s2_end + datetime.timedelta(hours=2)
    d1s3_end = d1s3_start + datetime.timedelta(hours=2)

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
        lambda: create_slot("2011-10-31T18:30:16.320Z", "2011-10-31T19:30:16.320Z"),
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
    action("Get slots ([d1s1, d1s2, d1s3])", get_slots)
    action("Get appointments ([])", get_doctor_appointments)
