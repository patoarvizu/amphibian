#!/bin/bash

helm package helm/amphibian/
version=$(cat helm/amphibian/Chart.yaml | yaml2json | jq -r '.version')
mv amphibian-$version.tgz docs/
helm repo index docs --url https://patoarvizu.github.io/amphibian
helm-docs
mv helm/amphibian/README.md docs/index.md
git add docs/