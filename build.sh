#!/bin/bash

set -e

if [ -z "$1" ]; then
  echo "Usage: $0 <graph.json URL>"
  exit 1
fi

GRAPH_URL="$1"
VERSION=$(cat VERSION)

mkdir -p dist

# Crée une version temporaire avec URL et version injectées
sed \
  -e "s|fetch('graph.json')|fetch('$GRAPH_URL')|" \
  -e "s|<h1>Massa Graph</h1>|<h1>Massa Graph <small style='font-size:12px;color:#ccc'>v$VERSION</small></h1>|" \
  index.html > dist/index.tmp.html

# Minifie avec html-minifier-terser
html-minifier-terser dist/index.tmp.html \
  --collapse-whitespace \
  --remove-comments \
  --remove-optional-tags \
  --minify-css true \
  --minify-js true \
  --output dist/index.html

# Supprime le temporaire
rm dist/index.tmp.html

echo "✅ Build done: dist/index.html"
