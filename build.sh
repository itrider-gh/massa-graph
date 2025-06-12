#!/bin/bash

set -e

if [ -z "$1" ]; then
  echo "Usage: $0 <graph.json URL>"
  exit 1
fi

GRAPH_URL="$1"
VERSION=$(cat VERSION)

mkdir -p dist

# Minify, replace fetch URL, and inject version into the banner
cat index.html \
  | sed "s|fetch('graph.json')|fetch('$GRAPH_URL')|" \
  | sed "s|<h1>Massa Graph</h1>|<h1>Massa Graph <small style='font-size:12px;color:#ccc'>v$VERSION</small></h1>|" \
  | tr '\n' ' ' \
  | sed 's/  */ /g' \
  | sed 's/> </></g' \
  > dist/index.html

echo "✅ Build done: dist/index.html"
