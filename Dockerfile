FROM golang:1.20-buster as builder

WORKDIR /app

COPY go.* ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o service

FROM mcr.microsoft.com/playwright:v1.35.0-jammy as runner

RUN mv /ms-playwright/firefox-1408 /ms-playwright/firefox-1319

RUN rm -rf /ms-playwright/firefox-1408 && \
    rm -rf /ms-playwright/chromium* && \
    rm -rf /ms-playwright/webkit* && \
    apt-get autoremove && \
    apt-get clean && \
    apt-get autoclean

FROM runner

WORKDIR /app

COPY --from=builder /app/service .

CMD ["./service"]