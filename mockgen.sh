set -x

rm -rf .mocks/*

$GOPATH/bin/mockgen -destination=./.mocks/mock_url_handler.go -source=./handler/url.go -package=mocks &&
$GOPATH/bin/mockgen -destination=./.mocks/mock_shortener_service.go -source=./service/shortener.go -package=mocks