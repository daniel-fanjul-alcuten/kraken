#!/bin/bash

dir=$(dirname $0)

# bare
if [ -d ./hooks/ ]; then
  cp "$dir"/post-receive ./hooks/
fi
