FROM gcr.io/distroless/static-debian12
ARG BIN=server
COPY bin/${BIN}-linux /app/server
EXPOSE 9092
ENTRYPOINT ["/app/server"]
