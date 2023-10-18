# set arg1 to "ts" to generate typescript code

cd proto || exit 1

if [ "$1" = "ts" ]; then
  echo "Generating typescript code"
  template="buf.gen.ts.yaml"
else
  echo "Generate go code and docs"
  template="buf.gen.yaml"
fi

buf generate --template $template
cd ..

if [ "$1" != "ts" ]; then
  cp -r github.com/KYVENetwork/chain/* ./
  swagger-combine ./docs/config.json -o ./docs/swagger.yml
  rm -rf tmp-swagger-gen
fi

rm -rf github.com
