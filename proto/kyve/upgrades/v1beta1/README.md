# Upgrade Types

When upgrading from `v1.0.0` to `v1.1.0`, a significant amount of internal
Protobuf types changed. This is mainly the change from `string` to `sdk.Dec`.

Unfortunately, we can't decode strings that don't have the correct zero 
padding. Therefore, all affected fields need to be updated manually. This
directory contains all the old Protobuf definitions so that we can correctly
decode the original values.

These Protobuf definitions are directly built to `app/upgrades/v1_1/types`.

After the `v1.1.0` upgrade, these definitions can be removed.
