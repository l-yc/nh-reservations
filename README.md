# nh-reservations

This is a web app for [New House](https://newhouse.mit.edu/) facility reservations. Made using Go, some ChatGPT, SIPB's [Petrock](https://petrock.mit.edu/), and sqlite3.

## Config
Insert `CLIENT_ID`, `CLIENT_SECRET`, and `SESSION_KEY` in a `.env` file placed in the same working directory. 

## Dev
Add a mapping from `localhost` to `nh.xvm.mit.edu` in your `/etc/hosts` file. You will need to generate a local SSL Cert, e.g. like so

```sh
openssl req -x509 -newkey rsa:4096 -sha256 -days 3650 \
  -nodes -keyout server.key -out server.crt -subj "/CN=nh.xvm.mit.edu" \
  -addext "subjectAltName=DNS:example.com,DNS:*.example.com,IP:10.0.0.1"
```

## Build

The following will build a single binary `nh-reservations`.

```sh
CGO_ENABLED=0 go build
```

## List of Webmasters
- Yue Chen Li (yuecli), 2026
