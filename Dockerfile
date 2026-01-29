FROM golang:1.25.3 AS compile
WORKDIR /app
COPY . /app
RUN go mod download
RUN CGO_ENABLED=0 go build -v -o /app/app ./cmd

FROM scratch AS execute
EXPOSE 8000
WORKDIR /app
COPY --from=compile /app/app /app/app
ENTRYPOINT [ "/app/app" ]
