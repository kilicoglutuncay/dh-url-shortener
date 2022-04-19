set -x

rm -rf .mocks/*

$GOPATH/bin/mockgen -destination=./.mocks/mock_url_handler.go -source=./internal/api/handler/url.go -package=mocks &&
$GOPATH/bin/mockgen -destination=./.mocks/mock_shortener_service.go -source=./internal/api/service/shortener.go -package=mocks