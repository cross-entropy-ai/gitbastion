FROM golang:1.24-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /gitbastion .

FROM alpine:3.21

RUN apk add --no-cache openssh \
    && ssh-keygen -A \
    && adduser -D -s /sbin/nologin -u 1000 git \
    && passwd -u git \
    && mkdir -p /home/git/.ssh \
    && chmod 700 /home/git/.ssh \
    && chown -R git:git /home/git/.ssh

COPY sshd_config /etc/ssh/sshd_config
COPY --from=build /gitbastion /usr/local/bin/gitbastion

EXPOSE 22

CMD ["gitbastion"]
