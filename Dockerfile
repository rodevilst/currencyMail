FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /currencyMail cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=0 /currencyMail .

COPY migrations ./migrations

ENV DB_USER=root
ENV DB_PASSWORD=rootroot
ENV DB_HOST=mysql
ENV DB_PORT=3306
ENV DB_NAME=currency

ENV SMTP_HOST=smtp.example.com
ENV SMTP_PORT=587
ENV SMTP_USER=your-email@example.com
ENV SMTP_PASSWORD=your-email-password
ENV FROM_EMAIL=your-email@example.com

CMD ["./currencyMail"]
