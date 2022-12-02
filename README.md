# bomberman-dom

## Prerequisites

- [pnpm](https://pnpm.io/installation)
- [golang](https://go.dev/doc/install)
- Access to [mini-framework](https://01.kood.tech/git/roosarula/mini-framework) and [bomberman-dom](https://01.kood.tech/git/AS/bomberman-dom) repos

## Getting started

1. `git clone https://01.kood.tech/git/AS/bomberman-dom`

### Backend

1. `cd backend`
2. `go run .`

### Frontend

1. `cd frontend`
2. `pnpm install`
3. `pnpm start`


### Using docker compose

1. Fill out your git credentials in docker-compose.yml GIT_USER & GIT_PASS
2. `docker-compose up --build`
3. visit `localhost:8081`