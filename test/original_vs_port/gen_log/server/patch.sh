#!/usr/bin/env bash

# Replace strict version constraint with a flexible ">= version" requirement
sed -i 's/s\.add_dependency("json", \["\([0-9.]*\)"\])/s.add_dependency("json", ">= \1")/' mjai.gemspec
sed -i 's/s\.add_dependency("nokogiri", \["\([0-9.]*\)"\])/s.add_dependency("nokogiri", ">= \1")/' mjai.gemspec
sed -i 's/s\.add_dependency("bundler", \["\([0-9.]*\)"\])/s.add_dependency("bundler", ">= \1")/' mjai.gemspec
sed -i 's/s\.add_dependency("sass", \["\([0-9.]*\)"\])/s.add_dependency("sass", ">= \1")/' mjai.gemspec

# Replace removed URI.decode with URI.decode_www_form_component
sed -i 's/URI\.decode/URI\.decode_www_form_component/' lib/mjai/tenhou_archive.rb
