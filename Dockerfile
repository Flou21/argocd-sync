FROM golang:1.18-alpine3.16 as builder

WORKDIR /app

COPY . .

RUN go build .

CMD /app/main


FROM alpine:3.16


RUN apk update

RUN apk --no-cache add curl

RUN curl -sSL -o /usr/local/bin/argocd https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64
RUN chmod +x /usr/local/bin/argocd

COPY --from=builder /app/argocd-sync /usr/local/bin/argocd-sync

CMD ["argocd-sync", "--app", "$APP_NAME"]
