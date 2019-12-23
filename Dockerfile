FROM golang:latest AS build-env

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o TIKSN.Lionize.SettingsService .

FROM golang:latest

WORKDIR /app

COPY --from=build-env /app/TIKSN.Lionize.SettingsService /app/TIKSN.Lionize.SettingsService

ENTRYPOINT ["/app/TIKSN.Lionize.SettingsService"]
