#!/bin/bash

dir=$(dirname $0)

# non bare
if [ -d ./.git/hooks/ ]; then
  cp "$dir"/hooks/post-{checkout,commit,rewrite} ./.git/hooks/
fi

# bare
if [ -d ./hooks/ ]; then
  cp "$dir"/hooks/post-receive ./hooks/
fi
