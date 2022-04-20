# URL Shortener API

---
[![codecov](https://codecov.io/github/kilicoglutuncay/dh-url-shortener/branch/master/graph/badge.svg?token=Lc1XnvB6YE)](https://codecov.io/github/kilicoglutuncay/dh-url-shortener)
[![pipeline](https://github.com/kilicoglutuncay/dh-url-shortener/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/kilicoglutuncay/dh-url-shortener/actions/workflows/main.yml)

URL Shortener API shortens long url to 7 character hash. Encodes URL in base-36 and store them in memory. 
Also, it periodically writes stored data in memory to file.

## How To Use 

You should have installed [Docker](https://www.docker.com/)

Run the following command to start the shortener api in container:

```
docker run -p 8080:8080 -it tujix/url-shortener:latest
```

You can also change the app address and short domain url to your liking thank to passing them as env variable.
```
docker run -p 8080:8090 -it -e APP_ADDR=":8090" -e SHORT_URL_DOMAIN=https://tujix.me tujix/url-shortener:latest
```

Shorten URL request:

```
curl -X POST -H "Content-Type: application/json" -d '{"url":"https://github.com/kilicoglutuncay/dh-url-shortener"}' http://localhost:8080/shorten
```
Shorten URL response:

```
{
  "url": "http://localhost:8080/a89145c"
}
```

Expand URL request, redirects you (302) to the original URL:

```
curl -X GET http://localhost:8080/a89145c
```

List all URLs request, shows all stored URLs with their hits:

```
curl -X GET http://localhost:8080/list
```

List response:

```
[
    {
        "hash": "a89145c",
        "original_url": "https://github.com/kilicoglutuncay/dh-url-shortener",
        "hits": 42
    }
]
```

Latest version of url shortener api is available on [tujix.me](http://tujix.me/list)