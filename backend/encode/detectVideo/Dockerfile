FROM golang:1.18.9

COPY ./ /src/
WORKDIR /src/
RUN ls 
RUN CGO_ENABLED=0 GOOS=linux go build -o detectVideo

FROM ubuntu:22.04
# Update the package list and install necessary packages
RUN apt-get update \
    && apt-get install -y ffmpeg

RUN mkdir /var/videos 
# Copy the binary from the builder stage
COPY --from=0 /src/detectVideo /

ENTRYPOINT ["/detectVideo"]