FROM golang:1.18.9

COPY ./ /src/
WORKDIR /src/
RUN ls 
RUN CGO_ENABLED=0 GOOS=linux go build -o generateMPD

FROM ubuntu:22.04
# Update the package list and install necessary packages
RUN apt-get update \
    && apt-get install -y ffmpeg
RUN mkdir /var/videos
COPY --from=0 /src/generateMPD /
ENTRYPOINT ["/generateMPD"]