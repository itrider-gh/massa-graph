#!/bin/bash

set -e

# Met à jour la version et le changelog
npx standard-version

# Extrait la version pour le tag
VERSION=$(node -p "require('./package.json').version")

# Écrit la version dans un fichier si besoin
echo "$VERSION" > VERSION

# Commit le fichier VERSION (et les autres modifiés par standard-version)
git add CHANGELOG.md VERSION package.json package-lock.json
git commit -m "chore(release): $VERSION [skip ci]" || echo "Aucun commit à faire"

# Pousse les changements
git push origin main
git push origin "v$VERSION"
