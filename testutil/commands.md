Create Funder

./build/kyved tx funders create-funder funder1 --from alice --keyring-backend test --yes --fees 5000tkyve --gas 200000 --home ~/.chain --chain-id kyve-local
./build/kyved tx funders fund-pool 0 100000000acoin 1000000acoin --from alice --keyring-backend test --yes --fees 5000tkyve --gas 200000 --home ~/.chain --chain-id kyve-local

Create Validator

```shell
./build/kyved tx staking create-validator validator_bob.json --from bob --keyring-backend test --yes --fees 5000tkyve --gas 200000 --home ~/.chain --chain-id kyve-local
```

```shell
./build/kyved tx staking create-validator validator_charlie.json --from charlie --keyring-backend test --yes --fees 5000tkyve --gas 200000 --home ~/.chain --chain-id kyve-local
```

Join pools

# bob - dummy
```shell
./build/kyved tx stakers join-pool 0 kyve137v27tfyegc083w5kj9zhhrfk34n8vhjma73gq 1000000 --yes --from bob --home ~/.chain --keyring-backend test --chain-id kyve-local --fees 4000tkyve
```

# charlie - faucet kyve1kahmjds2rxj2qzamdvy5m8ljnkqrf5xhetes7q
```shell
./build/kyved tx stakers join-pool 0 kyve1kahmjds2rxj2qzamdvy5m8ljnkqrf5xhetes7q 1000000 --yes --from charlie --home ~/.chain --keyring-backend test --chain-id kyve-local --fees 4000tkyve
```

Submit Bundle proposal

Charlie submit bundle: 
staker pool_id storage_id ... hash, from_index, bundle_size ...

# Submit Charlie kyve1ay22rr3kz659fupu0tcswlagq4ql6rwm4nuv0s faucet
```shell
./build/kyved tx bundles submit-bundle-proposal kyve1ay22rr3kz659fupu0tcswlagq4ql6rwm4nuv0s 0 "not-empty" 1024 hash 0 100 0 99 summary --from faucet --keyring-backend test --chain-id kyve-local --home ~/.chain --yes --fees 8000tkyve --gas 400000
```

# Submit Bob kyve1hvg7zsnrj6h29q9ss577mhrxa04rn94h7zjugq dummy
```shell
./build/kyved tx bundles submit-bundle-proposal kyve1hvg7zsnrj6h29q9ss577mhrxa04rn94h7zjugq 0 "not-empty" 1024 hash 500 100 0 99 summary --from dummy --keyring-backend test --chain-id kyve-local --home ~/.chain --yes --fees 8000tkyve --gas 400000
```

./build/kyved tx bundles claim-uploader-role kyve1ay22rr3kz659fupu0tcswlagq4ql6rwm4nuv0s 0 --from faucet --keyring-backend test --chain-id kyve-local --home ~/.chain --yes --fees 4000tkyve
./build/kyved tx bundles claim-uploader-role kyve1hvg7zsnrj6h29q9ss577mhrxa04rn94h7zjugq 0 --from bob --keyring-backend test --chain-id kyve-local --home ~/.chain --yes --fees 4000tkyve

./build/kyved tx bundles vote-bundle-proposal kyve1ay22rr3kz659fupu0tcswlagq4ql6rwm4nuv0s 0 "not-empty" 1 --from faucet --keyring-backend test --chain-id kyve-local --home ~/.chain --yes --fees 4000tkyve
./build/kyved tx bundles vote-bundle-proposal kyve1hvg7zsnrj6h29q9ss577mhrxa04rn94h7zjugq 0 "not-empty" 1 --from dummy --keyring-backend test --chain-id kyve-local --home ~/.chain --yes --fees 4000tkyve





# Join Pool mynode Alice foundation
```
./build/kyved tx stakers join-pool 0 kyve1fd4qu868n7arav8vteghcppxxa0p2vna5f5ep8 1000000 --yes --from alice --home ~/.chain --keyring-backend test --chain-id kyve-alpha --fees 4000tkyve --node https://rpc.alpha.kyve.network:443
```


# Submit Charlie kyve1ay22rr3kz659fupu0tcswlagq4ql6rwm4nuv0s faucet
```shell
```
./build/kyved tx bundles submit-bundle-proposal kyve1eka2hngntu5r2yeuyz5pd45a0fadarp3zue8gd 0 "not-empty" 1024 hash 300 100 0 99 summary --from dummy --keyring-backend test --chain-id kyve-alpha --home ~/.chain --yes --fees 8000tkyve --gas 400000 --node https://rpc.alpha.kyve.network:443
./build/kyved tx bundles submit-bundle-proposal kyve1s7j6ccd4ule2cwtxsecqvfjmfm0u40g5drx8zl 0 "not-empty" 1024 hash 400 100 0 99 summary --from faucet --keyring-backend test --chain-id kyve-alpha --home ~/.chain --yes --fees 8000tkyve --gas 400000 --node https://rpc.alpha.kyve.network:443
./build/kyved tx bundles submit-bundle-proposal kyve1jq304cthpx0lwhpqzrdjrcza559ukyy3zsl2vd 0 "not-empty" 1024 hash 200 100 0 99 summary --from dummy --keyring-backend test --chain-id kyve-alpha --home ~/.chain --yes --fees 8000tkyve --gas 400000 --node https://rpc.alpha.kyve.network:443


# Submit Bob kyve1hvg7zsnrj6h29q9ss577mhrxa04rn94h7zjugq dummy
```shell
```

./build/kyved tx bundles claim-uploader-role kyve1eka2hngntu5r2yeuyz5pd45a0fadarp3zue8gd 0 --from bob --keyring-backend test --chain-id kyve-alpha --home ~/.chain --yes --fees 4000tkyve --node https://rpc.alpha.kyve.network:443
./build/kyved tx bundles claim-uploader-role kyve1s7j6ccd4ule2cwtxsecqvfjmfm0u40g5drx8zl 0 --from faucet --keyring-backend test --chain-id kyve-alpha --home ~/.chain --yes --fees 4000tkyve --node https://rpc.alpha.kyve.network:443

./build/kyved tx bundles vote-bundle-proposal kyve1eka2hngntu5r2yeuyz5pd45a0fadarp3zue8gd 0 "not-empty" 1 --from dummy --keyring-backend test --chain-id kyve-alpha --home ~/.chain --yes --fees 4000tkyve --node https://rpc.alpha.kyve.network:443
./build/kyved tx bundles vote-bundle-proposal kyve1s7j6ccd4ule2cwtxsecqvfjmfm0u40g5drx8zl 0 "not-empty" 1 --from faucet --keyring-backend test --chain-id kyve-alpha --home ~/.chain --yes --fees 4000tkyve --node https://rpc.alpha.kyve.network:443
./build/kyved tx bundles vote-bundle-proposal kyve1jq304cthpx0lwhpqzrdjrcza559ukyy3zsl2vd 0 "not-empty" 2 --from foundation --keyring-backend test --chain-id kyve-alpha --home ~/.chain --yes --fees 4000tkyve --node https://rpc.alpha.kyve.network:443


# Create Funder

./build/kyved tx funders create-funder funder1 --from alice --keyring-backend test --yes --fees 5000tkyve --gas 200000 --home ~/.chain --chain-id kyve-alpha --node https://rpc.alpha.kyve.network
./build/kyved tx funders fund-pool 0 100000000acoin 1000000acoin --from alice --keyring-backend test --yes --fees 5000tkyve --gas 200000 --home ~/.chain --chain-id kyve-alpha --node https://rpc.alpha.kyve.network
