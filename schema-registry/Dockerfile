FROM golang:alpine3.11
# Setup python environment with its needed packages
RUN apk add --no-cache \
    py-pip \
    && pip install genson
# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /app
 #Copy and download dependency using go mod
COPY go.mod .
RUN go mod download
# Copy the code into the container
COPY . .
# Build the application
RUN go build -o server ./main
# Export necessary port
EXPOSE 8080
# Command to run when starting the container
CMD ["/app/server"]
