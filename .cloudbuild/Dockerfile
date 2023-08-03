FROM golang:1.20 as build
WORKDIR /go/src/app
COPY ./ .
RUN go build -v -o /app .

# Note: I moved off the official "distroless" image, due to this issue: https://github.com/GoogleContainerTools/distroless/issues/1342
#       I have no clue of the differences, but this works *shrug*
FROM cgr.dev/chainguard/glibc-dynamic

COPY --from=build /app/ /app
COPY --from=build /go/src/app/.private /private

CMD ["/app"]
