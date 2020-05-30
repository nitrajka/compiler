#!/bin/bash

run_test() {
  test_name=$1
  echo "Running test: $test_name"
  ./.tmp/compiler -o ./.tmp/$test_name.go $test_name.compilang
  go build -o ./.tmp/$test_name ./.tmp/$test_name.go
  ./.tmp/$test_name > ./.tmp/$test_name.output
  diff ./.tmp/$test_name.output $test_name.output
  error=$?
  if [ $error -eq 0 ]
  then
     echo "OK"
  else
     echo "Failed"
  fi
}

# Test setup
rm -Rf .tmp
mkdir -p .tmp
go build -o .tmp/compiler ../cmd/compilang/main.go

# Run tests
ALL_TESTS="./*.compilang"
for f in $ALL_TESTS;
do
  base=$(basename $f)
  test_name="${base%.*}"
  run_test $test_name
done

# Cleanup
rm -Rf .tmp
