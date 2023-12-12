# Author

Omar Adel Abdel Hamid Ahmed Brikaa - 20206043 - S7 - brikaaomar@gmail.com

# GitHub repo

https://github.com/Brikaa/tools-3-project

# Prerequisites

- Docker CLI (even if you are using Podman backend)
- Docker Compose CLI (even if you are using Podman backend)
- Podman (optional. If you want to use Docker instead, read on)
- A Linux environment with systemd
- GNU Make (pre-installed on most Linux environments)

# Choosing between Docker and Podman

Podman backend is used by default. If you wish to use Docker, create a file called `.no-podman` in the `src` directory.

# Running

To run the containers:

```bash
cd src
make dev
```

Access the app at `clinic.localhost:4000` (or whatever you configured frontend host to).

# Configuration

In order to configure different parameters like ports and hosts create a file called `.env.overrides` overriding
the variables you want to change in `.env.sample`.

# Tech stack

- MySQL database
- Redis PubSub messaging
- Go backend
- Angular frontend
