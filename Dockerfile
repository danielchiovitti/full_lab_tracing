FROM golang:1.21-alpine as build
WORKDIR /app
COPY . .
EXPOSE 3500
RUN apk update && apk add --no-cache ca-certificates

RUN update-ca-certificates
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o cep ./cmd/full_cep/main.go


FROM scratch
WORKDIR /app
COPY --from=build /app/cep .
EXPOSE 3500
ENTRYPOINT ["./cep"]
