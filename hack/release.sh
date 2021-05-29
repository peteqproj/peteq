
#!/bin/bash
# uses github.com/davidrjonas/semver-cli

set -e

version=$(cat VERSION)
echo "Previous version: $version"
minor=$(semver-cli inc minor $version)
echo $minor > VERSION

echo "Releasing version $minor"
git add VERSION
git commit -m "chore(release): version $minor"
git push
fqrn="v$minor"
git checkout -b release-$fqrn
git tag $fqrn
git push --tags 
git checkout master