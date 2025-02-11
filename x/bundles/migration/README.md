## Bundles Migration

This migration is used to update the bundle summaries by adding the created Merkle roots that are
required by the [Trustless API](https://docs.kyve.network/access-data-sets/trustless-api/overview).
To create the correct Merkle proofs for the archived bundles, an approach was used that downloads
the bundle and computes the correct Merkle proof. To validate the correct proof calculation, it was
implemented in [Go and Python](https://github.com/KYVENetwork/merkle-script) , whereas it was tested with the TypeScript runtime implementation.

These scripts created the `files/merkle_roots_pool_X` binary files, that are used for the v2 migration.
Both Python and Go implementations were used to compute the hashes, which provide identical results:

```
merkle_roots_pool_0
da4bb9bf0a60c5c79e399d8bb54ae4cf916f6c1dbdd5cdae45cb991f4e56158f

merkle_roots_pool_1
3c4eeb915cd01c6adea3241ea3536dfce5cec87017557b7e43d92c6ceec3096e  

merkle_roots_pool_2
754eb4680fe550cd3a7277ab0fc12c8f7ce794d18ca71d247561e40b05629c39  

merkle_roots_pool_3
df26b886928dbec03e84eca9b41c02b15ae7c5e7cf39ab540fcf381d3e1d27cc  

merkle_roots_pool_5
051efd6e44d7ac5bca41abb20aaf79d34dd095b5d6797d536bf13face7e397f9  

merkle_roots_pool_7
303d5ccaa18cc9e23298d599e3ba4c5bcf46f44d0fb5dd2cfdebcd02dcd8dc95 

merkle_roots_pool_9
e2f1c174350e5925d3f61b7adfb077f38507aec1562900b79c645099809ae617 
```