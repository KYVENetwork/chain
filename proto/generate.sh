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
  # Generate proto files
  $(cd proto && buf generate --template buf.gen.pulsar.yaml --path $module_list) || cleanup_and_error

  # Copy the generated proto files to the x/ directory
#  cp -r tmp-gen-pulsar/kyve/* ./x
}

# Generate openapi docs
generate_docs() {
  # Generate swagger files
  $(cd proto && buf generate --template buf.gen.swagger.yaml --path $proto_list) || cleanup_and_error

  # Combine the swagger files
  swagger-combine ./docs/config.json -o ./docs/static/openapi.yml
}

generate_gogo_proto
generate_pulsar_proto
generate_docs
cleanup
