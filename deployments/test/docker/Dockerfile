# Build the application from source
FROM golang:1.20 AS build-stage

WORKDIR /app

# Download go dependencies. Doing this separately rather than simply building
# the image every time should not download the dependencies if there isn't
# a change within the go modules.
COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./pkg ./pkg


RUN CGO_ENABLED=0 GOOS=linux go build -o /app/build/photoflux cmd/main.go

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a small image
# When you run the backend, you don't need the ability in the container to run "go" code
# For this reason, you just copy the executable to a small image that allows running executables
# This small image is the one below, a linux environment to run an executable
FROM gcr.io/distroless/base-debian11 AS build-release-stage

# ^^ all above was build stage

WORKDIR /

COPY --from=build-stage /app/build/photoflux /bin/photoflux

EXPOSE 8033

USER nonroot:nonroot

ENTRYPOINT ["/bin/photoflux"]
