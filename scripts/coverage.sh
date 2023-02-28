#!/bin/bash
mkdir -p scripts/sresult_cover
pkgs=$(go list ./...)
cov_result="./scripts/result_cover/coverage.out"
cov_html="./scripts/result_cover/cover.html"
deps=`echo ${pkgs} | tr ' ' ","`
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
