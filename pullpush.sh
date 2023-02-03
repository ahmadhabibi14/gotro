#!/usr/bin/env bash

set -x # print all command
set -e # exit on error

if [ $# -eq 0 ] ; then
  echo "Usage: 
  ./pullpush.sh 'the commit message'"
  exit
fi

# generate documentation
for x in A B F I L M T X; do # remove S C until https://github.com/robertkrimen/godocdown/issues/17 fixed
  godocdown $x > $x/README.md
  replacer gotro/$x github.com/kokizzu/gotro/$x '# ' $x/README.md
done

# format indentation
go install golang.org/x/tools/cmd/goimports@latest
goimports -w **/*.go
echo "codes formatted.."

go get -u -v github.com/ory/dockertest/v3@latest
go get -u -v github.com/kokizzu/id64@latest
go get -u -v github.com/kokizzu/lexid@latest
go mod tidy

# testing if has error
go build loader.go || ( echo 'has error, exiting..' ; kill 0 )
echo "codes tested.."

# testing if has "gokil" included
ag gokil **/*.go && ( echo 'echo should not import previous gokil library..' ; kill 0 )
echo "imports checked.."

# run linter
#golangci-lint run 

# run tests
go test ./...

# add and commit all files
git add .
git status
read -p "Press Ctrl+C to exit, press any enter key to check the diff..
"

# recheck again
git diff --staged
echo 'Going to commit with message: '\""$*"\"
read -p "Press Ctrl+C to exit, press any enter key to really commit..
"

git commit -m "$*" && git pull && git push origin master

TAG=$(ruby -e 't = Time.now; print "v1.#{t.month+(t.year-2021)*12}%02d.#{t.hour}%02d" % [t.day, t.min]')
git tag -a "$TAG" -m "$*"
git push --tags 

echo "# to undo this release: 
git tag -d $TAG
git push -d origin $TAG"
