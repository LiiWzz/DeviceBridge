FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends gcc gcc-arm-linux-gnueabihf libc6-dev libc6-dev-armhf-cross busybox-static && rm -rf /var/lib/apt/lists/*
