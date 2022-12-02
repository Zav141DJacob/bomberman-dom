FROM node:19-alpine3.16 AS builder

WORKDIR /build

ARG GIT_USER
ARG GIT_PASS

RUN npm install -g pnpm
RUN apk add git python3 make g++

COPY package.json .
RUN pnpm remove mini-framework
RUN git clone https://$GIT_USER:$GIT_PASS@01.kood.tech/git/roosarula/mini-framework.git
RUN pnpm add file:mini-framework
RUN pnpm install

COPY . .
RUN pnpm build

FROM caddy:2.6-alpine
EXPOSE 80 443

COPY --from=builder /build/dist/ /srv
COPY Caddyfile /etc/caddy/Caddyfile