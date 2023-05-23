cd proto
buf generate
buf generate --template buf.gen.pulsar.yaml
cd ..

cp -r github.com/KYVENetwork/chain/* ./
rm -rf github.com
rm -rf kyve

swagger-combine ./docs/config.json -o ./docs/swagger.yml
rm -rf tmp-swagger-gen
