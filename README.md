# DNSPod DDNS client

DDNS client for Dnspod, you can run it as system service

## How to get token?

Dndpod API doc:

- https://www.dnspod.cn/docs/info.html
- https://support.dnspod.cn/Kb/showarticle/tsid/227/

## Test

make test

## Build

~~~
go build
~~~

## Run

~~~
./ddns -token="xxx,xxxxx" -domain="example.com" -record="test.ddns" -interval=60
~~~

## Install as service

~~~
./ddns -token="xxx,xxxxx" -domain="example.com" -record="test.ddns" -interval=60 -install
~~~

## Uninstall

~~~
./ddns -uninstall
~~~