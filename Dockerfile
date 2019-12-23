FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o TIKSN.Lionize.SettingsService .

EXPOSE 8080

CMD ["./TIKSN.Lionize.SettingsService"]