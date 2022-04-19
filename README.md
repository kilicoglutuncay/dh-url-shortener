# dh-url-shortener
[![codecov](https://codecov.io/github/kilicoglutuncay/dh-url-shortener/branch/master/graph/badge.svg?token=Lc1XnvB6YE)](https://codecov.io/github/kilicoglutuncay/dh-url-shortener)
[![pipeline](https://github.com/kilicoglutuncay/dh-url-shortener/actions/workflows/main.yml/badge.svg?branch=master)](https://github.com/kilicoglutuncay/dh-url-shortener/actions/workflows/main.yml)
## Getting Started
Url shortener makes it easy to shorten long urls.  
Stores data in memory and supports snapshotting of data to file periodically. Also restores data from snapshot.


## Running 

You should have installed [Docker](https://www.docker.com/)

```
docker run -p 8080:8080 -it tujix/url-shortener:latest
```


TODO:
- Continues deployment to pipeline.
- Env variables for deployment.