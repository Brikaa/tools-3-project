import requests
import json

BASE = "http://backend_runner:8000"
token = ""
headers = {"Content-Type": "application/json"}


def create_url(endpoint):
    return BASE + "/" + endpoint


def send_request(method, url, payload):
    res = requests.request(
        method, url, headers=headers, data=json.dumps(payload)
    )
    return res.text, res.status_code


def action(message, function):
    print(message)
    result = function()
    print(result)
    print()


def signup(username, password, role):
    return send_request(
        "POST",
        create_url("signup"),
        {"username": username, "password": password, "role": role},
    )


if __name__ == "__main__":
    action("Signup p1", lambda: signup("p1", "p1123", "patient"))
    action("Signup p2", lambda: signup("p2", "p2123", "patient"))
    action("Signup d1", lambda: signup("d1", "d1123", "doctor"))
    action("Signup d2", lambda: signup("d2", "d2123", "doctor"))
    action("Signup d2 (invalid)", lambda: signup("d2", "d2123", "doctor"))
    action("Signup i1 (invalid role)", lambda: signup("i1", "i123", "doctors"))
    action("Signup i2 (invalid username)", lambda: signup("i-2", "i123", "doctor"))
