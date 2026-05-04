FROM alpine:3.21

RUN apk add --no-cache openssh \
    && ssh-keygen -A \
    && adduser -D -s /sbin/nologin -u 1000 ssh-bastion \
    && passwd -u ssh-bastion \
    && mkdir -p /home/ssh-bastion/.ssh \
    && chmod 700 /home/ssh-bastion/.ssh \
    && chown -R ssh-bastion:ssh-bastion /home/ssh-bastion/.ssh \
    && rm -f /bin/sh /bin/ash /usr/bin/env

COPY sshd_config /etc/ssh/sshd_config

EXPOSE 22

CMD ["/usr/sbin/sshd", "-D", "-e"]
