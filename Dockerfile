FROM golang:1.16

WORKDIR /app
COPY go.mod go.sum ./

# Download all dependencies.
RUN go mod download

# Copy the source from the current directory to
# the Working Directory inside the container
COPY . .
RUN go build -o main .
#RUN go run .
EXPOSE 8080
# Command to run the executable
CMD ["./main"]