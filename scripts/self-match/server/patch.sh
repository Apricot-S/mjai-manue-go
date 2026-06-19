#!/usr/bin/env bash

set -euo pipefail

# Remove deprecated `s.has_rdoc`.
sed -i 's/s\.has_rdoc = true//' mjai.gemspec

# Replace strict version constraints with flexible ">= version" requirements.
sed -i 's/s\.add_dependency("json", \["\([0-9.]*\)"\])/s.add_dependency("json", ">= \1")/' mjai.gemspec
sed -i 's/s\.add_dependency("nokogiri", \["\([0-9.]*\)"\])/s.add_dependency("nokogiri", ">= \1")/' mjai.gemspec
sed -i 's/s\.add_dependency("bundler", \["\([0-9.]*\)"\])/s.add_dependency("bundler", ">= \1")/' mjai.gemspec
sed -i 's/s\.add_dependency("sass", \["\([0-9.]*\)"\])/s.add_dependency("sass", ">= \1")/' mjai.gemspec

# Replace removed `URI.decode` with `URI.decode_www_form_component`.
sed -i 's/URI\.decode/URI\.decode_www_form_component/' lib/mjai/tenhou_archive.rb

# Ruby 4 removed the 3-argument form of `ERB.new`.
sed -i '
/html = ERB\.new(File\.read("#{res_dir}\/views\/archive_player\.erb"), nil, "<>")\./c\
html = ERB.new(File.read("#{res_dir}/views/archive_player.erb"), trim_mode: "<>").
' lib/mjai/file_converter.rb
