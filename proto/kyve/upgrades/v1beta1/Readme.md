# Upgrade Types

For the upgrade from v1.0.0 to v1.1.0 a lot of internal proto string
types were changed to sdk.Dec. 

Unfortunately Cosmos proto can not decode strings which do not have to
correct zero padding. Therefore, all affected fields need to be updated.
This directory contains all (old) proto definitions which got changed.

The proto files are built directly to `app/upgrades/types`.

After the `v1.1.0` upgrade was performed, this directory can be deleted.