FROM golang:1.21-alpine3.19 as builder

# The latest alpine images don't have some tools like (`git` and `bash`).
# Adding git, bash and openssh to the image
RUN apk update && apk upgrade && apk add --no-cache bash git openssh make autoconf gcc libc-dev sudo procps curl

# Set the Current Working Directory inside the container
WORKDIR /var/www/backend

COPY . .

RUN go build -v

# Expose port 3000 to the outside world
EXPOSE 3000

CMD ["./backend","start"]