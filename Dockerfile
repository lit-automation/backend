ARG base_img

FROM ${base_img} as builder

WORKDIR /app/src/slr-api
RUN GOOS=linux go build -v -o bin/main

FROM frolvlad/alpine-glibc

RUN apk --no-cache add ca-certificates tzdata && update-ca-certificates

COPY --from=builder /app/src/slr-api/bin/main /
EXPOSE 9001
STOPSIGNAL SIGTERM
CMD ["./main"]
