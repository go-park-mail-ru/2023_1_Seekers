#!/bin/bash

res_dir="./scripts/result_cover"
mkdir -p $res_dir
pkgs=$(go list ./... | grep -v "mocks")
cov_result="${res_dir}/coverage.out"
cov_html="${res_dir}/cover.html"
deps=`echo ${pkgs} | tr ' ' ","`
echo deps
echo "mode: atomic" > $cov_result

for pkg in $pkgs; do
    set -e
    go test -v -cover -coverpkg "$deps" -coverprofile=coverage.tmp $pkg
    set +e

    if [ -f coverage.tmp ]; then
        tail -n +2 coverage.tmp >> $cov_result
        rm coverage.tmp
    fi
done;
# get total cov
go tool cover -func=$cov_result
# get html
go tool cover -html $cov_result -o $cov_html
