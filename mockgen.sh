set -x

arr=(shortener)

rm -rf .mocks/* || true

for i in ${arr[@]}
do
  $GOPATH/bin/mockgen -destination=./.mocks/mock_${i}_handler.go -source=./handler/${i}.go -package=mocks &&
  #$GOPATH/bin/mockgen -destination=./.mocks/mock_${i}_service.go -source=./service/${i}.go -package=mocks &&
  #$GOPATH/bin/mockgen -destination=./.mocks/mock_${i}_client.go -source=./client/${i}.go -package=mocks
  echo $?
#$GOPATH/bin/mockgen -destination=./.mocks/mock_user_segment_client.go -source=./client/user_segment.go -package=mocks
done
