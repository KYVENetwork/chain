cd proto
buf generate --template buf.gen.gogo.yaml
buf generate --template buf.gen.pulsar.yaml
cd ..

cp -r github.com/KYVENetwork/chain/* ./
rm -rf github.com
rm -rf kyve

exit 0

# TODO: fix docs

swagger-combine ./docs/config.json -o ./docs/swagger.yml
rm -rf tmp-swagger-gen