#!/bin/bash

# build go plugins
# pushd ../go-plugins
# sub_packages=$(find . -maxdepth 1 -type d ! -name '.')
# for sub_package in $sub_packages; do
#     if [ -d $sub_package ] && [ -f $sub_package/go.mod ]; then
#         echo "build $sub_package"
#         go mod tidy
#         go build -o ../bin/$(basename $sub_package) $sub_package
#     fi
# done

# popd

chmod 640 bin/*
chown 3000:3000 bin/*

# build cos docker
export http_proxy="http://proxyxa.zte.com.cn:80"
export https_proxy="https://proxyxa.zte.com.cn:80"
docker build --no-cache --network host -f build/dockerfiles/alpine.Dockerfile -t uapgateway:v0.0.1 .