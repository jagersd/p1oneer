FROM debian:latest
WORKDIR /init
RUN apt update && apt install -y python3 curl procps
ENV P1ONEER_CONFIG_DIR=/init/p1-configs
COPY ./examples/*.json /init/p1-configs/
COPY --from=p1oneer:latest p1oneer /init/bin/p1oneer
ENV PATH="$PATH:/init/bin"
ENTRYPOINT [ "p1oneer" ]

