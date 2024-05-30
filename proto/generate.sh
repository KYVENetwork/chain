#!/bin/sh

workspace_dir="/workspace"
cosmos_proto_dir="/cosmos-sdk/proto"

# Find all proto files in the kyve directory (except for module.proto's) and make a comma separated list
proto_list=$(cd proto && find kyve -name '*.proto' -not -name 'module.proto' -print0 | xargs -0 -n1 | tr '\n' ',' | sed 's/,$//')

# Find all module.proto files in the kyve directory and make a comma separated list
module_list=$(cd proto && find kyve -name 'module.proto' -print0 | xargs -0 -n1 | tr '\n' ',' | sed 's/,$//')

# Cleanup
cleanup() {
  rm -rf tmp-gen
  rm -rf tmp-swagger-gen
}

# Cleanup and error
cleanup_and_error() {
  cleanup
  exit 1
}

# Generate gogo proto files
generate_gogo_proto() {
  # Generate proto files
  (cd proto && buf generate -v --template buf.gen.gogo.yaml --path "$proto_list") || cleanup_and_error

  # Copy the generated proto files to the x/ directory
  cp -r tmp-gen/github.com/KYVENetwork/chain/* ./
}

# Generate module proto files
generate_pulsar_proto() {
  (cd proto && buf generate --template buf.gen.pulsar.yaml --path "$module_list") || cleanup_and_error
}

# Generate openapi docs
generate_docs() {
  # Generate Cosmos-SDK swagger files (this part is mostly copied from the cosmos-sdk repo)
  cosmos_sdk_proto_dirs=$(cd $cosmos_proto_dir && find ./cosmos -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
  for dir in $cosmos_sdk_proto_dirs; do
    # Filter query files
    query_file=$(cd $cosmos_proto_dir && find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
    if [ -n "$query_file" ]; then
      # Don't generate swagger files for certain problematic files
      if [ "$query_file" != "./cosmos/app/v1alpha1/query.proto" ] && [ "$query_file" != "./cosmos/orm/query/v1alpha1/query.proto" ]; then
        (cd $cosmos_proto_dir && buf generate --template $workspace_dir/proto/buf.gen.swagger.yaml "$query_file" --output $workspace_dir/tmp-swagger-gen) || cleanup_and_error
      fi
    fi
  done

  # Generate Kyve swagger files
  (cd proto && buf generate --template buf.gen.swagger.yaml --path "$proto_list") || cleanup_and_error

  # Combine the swagger files
  swagger-combine ./docs/config.json -o ./docs/static/openapi.yml
}

generate_gogo_proto
generate_pulsar_proto
generate_docs
cleanup
