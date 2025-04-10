# Use official Golang base image
FROM golang:1.23-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the rest of the app
COPY . .

# Build the Go app
RUN go build -o main .

# Expose the port your app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
