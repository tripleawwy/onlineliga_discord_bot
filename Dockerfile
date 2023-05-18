FROM golang:1.20 AS builder

WORKDIR /app

ENV CGO_ENABLED=0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Add Packages like curl, jq and tar and fontconfig
#RUN apk add --no-cache curl jq tar fontconfig
RUN apt-get update && apt-get install -y curl jq tar fontconfig unzip

# Download Monospace Cascadio font from https://github.com/microsoft/cascadia-code/releases
# Get zip of latest release
RUN mkdir -p /usr/share/fonts
RUN curl -s https://api.github.com/repos/microsoft/cascadia-code/releases/latest \
  | jq -r '.assets[] | select(.name | test(".zip$")) | .browser_download_url' \
  # the expected output is something like https://github.com/microsoft/cascadia-code/releases/download/v2111.01/CascadiaCode-2111.01.zip
  | xargs curl -L -o /tmp/cascadia.zip \
  && unzip -o /tmp/cascadia.zip -d /usr/share/fonts \
  && fc-cache -f -v

# Build binaries \
RUN ./scripts/build.sh

#From gcr.io/distroless/static as final
FROM alpine:latest as final

COPY --from=builder /app/bin/ol_discord_bot-linux-amd64 /app/bin/
COPY --from=builder /usr/share/fonts/ /usr/share/fonts/

ENV FONT_PATH=/usr/share/fonts/ttf/static/CascadiaCode-Bold.ttf

CMD ["/app/bin/ol_discord_bot-linux-amd64"]
