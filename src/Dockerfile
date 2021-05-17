FROM golang:latest
COPY ./ /synche/src

WORKDIR /synche/src
RUN apt update && \
    apt install jq -y && \
    download_url=$(curl -s https://api.github.com/repos/go-swagger/go-swagger/releases/latest | \
    jq -r '.assets[] | select(.name | contains("'"$(uname | tr '[:upper:]' '[:lower:]')"'_amd64")) | .browser_download_url') && \
    curl -o /usr/local/bin/swagger -L'#' "$download_url" && \
    chmod +x /usr/local/bin/swagger && \
    mkdir /synche/build && \
    make generate && \
    make build

EXPOSE 9449 2121

CMD ["/synche/build/synche-server", "serve"]