FROM resin/raspberry-pi-golang AS build-env
ADD . /src
RUN cd /src && go get -u github.com/eclipse/paho.mqtt.golang && go build -ldflags "-linkmode external -extldflags -static" -x -o mqtt-external-ip .

FROM hypriot/rpi-alpine-scratch
WORKDIR /app
COPY --from=build-env /src/mqtt-external-ip /app/
ENTRYPOINT ["./mqtt-external-ip"]
