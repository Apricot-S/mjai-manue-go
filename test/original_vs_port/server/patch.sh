#!/usr/bin/env bash

# Replace strict version constraint with a flexible ">= version" requirement
sed -i 's/s\.add_dependency("json", \["\([0-9.]*\)"\])/s.add_dependency("json", ">= \1")/' mjai.gemspec
sed -i 's/s\.add_dependency("nokogiri", \["\([0-9.]*\)"\])/s.add_dependency("nokogiri", ">= \1")/' mjai.gemspec
sed -i 's/s\.add_dependency("bundler", \["\([0-9.]*\)"\])/s.add_dependency("bundler", ">= \1")/' mjai.gemspec
sed -i 's/s\.add_dependency("sass", \["\([0-9.]*\)"\])/s.add_dependency("sass", ">= \1")/' mjai.gemspec
