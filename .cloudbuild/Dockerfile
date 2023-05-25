FROM golang:1.20 as build
WORKDIR /go/src/app
COPY ./ .
RUN go build -v -o /app .

# Now copy it into our base image.
FROM gcr.io/distroless/cc

COPY --from=build /app/ /app
COPY --from=build /go/src/app/.private /private

CMD ["/app"]
