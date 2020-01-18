FROM golang:alpine AS build

WORKDIR /go/src/app

ENV USER=go \
    UID=1000 \
    GID=1000 \
    GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0

COPY . .

RUN go build -ldflags="-w -s" -o awsnews && \
  addgroup --gid "$GID" "$USER" && \
  adduser \
  --disabled-password \
  --gecos "" \
  --home "$(pwd)" \
  --ingroup "$USER" \
  --no-create-home \
  --uid "$UID" \
  "$USER" && \
  chown "$UID":"$GID" /go/src/app/awsnews

FROM scratch
COPY --from=build /etc/passwd /etc/group /etc/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /go/src/app/awsnews /go/src/app/config.yaml /
USER 1000
ENTRYPOINT ["/awsnews"]
