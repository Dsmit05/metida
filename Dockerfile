ARG GITHUB_PATH=github.com/Dsmit05/metida

FROM golang:1.16 AS builder
WORKDIR /home/${GITHUB_PATH}
RUN apt-get install make -y
COPY . .
RUN make build

FROM alpine:3.16.0
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/metida .
COPY --from=builder /home/${GITHUB_PATH}/config.yml .
RUN chown root:root metida

EXPOSE 8080
CMD ["./metida", "prod"]