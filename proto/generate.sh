#!/bin/sh

# Find all proto files in the kyve directory except for module.proto's
protos=$(cd proto && find kyve -name '*.proto' -not -name 'module.proto')
# Transform the proto files into a comma separated list
proto_list=$(echo $protos | sed 's/ /,/g')

# Find all module.proto files in the kyve directory
modules=$(cd proto && find kyve -name module.proto)
# Transform the module.proto files into a comma separated list
module_list=$(echo $modules | sed 's/ /,/g')

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
  $(cd proto && buf generate -v --template buf.gen.gogo.yaml --path $proto_list) || cleanup_and_error

  # Copy the generated proto files to the x/ directory
  cp -r tmp-gen/github.com/KYVENetwork/chain/* ./
}

# Generate module proto files
generate_pulsar_proto() {
  $(cd proto && buf generate --template buf.gen.pulsar.yaml --path $module_list) || cleanup_and_error
}

# Generate openapi docs
generate_docs() {
  # Generate cosmos-sdk swagger files (this part is mostly copied from the cosmos-sdk repo)
  cosmos_sdk_proto_dirs=$(cd /cosmos-sdk/proto && find ./cosmos -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
  for dir in $cosmos_sdk_proto_dirs; do
    # Filter query files
    query_file=$(cd /cosmos-sdk/proto && find "${dir}" -maxdepth 1 \( -name 'query.proto' -o -name 'service.proto' \))
    if [ -n "$query_file" ]; then
      # Don't generate swagger files for certain problematic files
      if [ "$query_file" != "./cosmos/app/v1alpha1/query.proto" ] && [ "$query_file" != "./cosmos/orm/query/v1alpha1/query.proto" ]; then
        $(cd  /cosmos-sdk/proto && buf generate --template /workspace/proto/buf.gen.swagger.yaml $query_file --output /workspace/tmp-swagger-gen) || cleanup_and_error
      fi
    fi
  done

  # Generate Kyve swagger files
  $(cd proto && buf generate --template buf.gen.swagger.yaml --path $proto_list) || cleanup_and_error

  # Combine the swagger files
  swagger-combine ./docs/config.json -o ./docs/static/openapi.yml
}

generate_gogo_proto
generate_pulsar_proto
generate_docs
cleanup
