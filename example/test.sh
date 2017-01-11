#!/bin/bash -eu

node=`find ./external/nodejs* -name node | head -1`


once_output=`$node ./tool/js-squish/example/once.squished.js`
if [ "$once_output" != "Hello World" ]; then
  echo "Expected 'Hello World', got $once_output"
  exit 2
fi


twice_output=`$node ./tool/js-squish/example/twice.squished.js`
if [ "$twice_output" != "Hello World" ]; then
  echo "Expected 'Hello World', got $twice_output"
  exit 2
fi
