<!--
order: 8
-->

# Parameters

The bundles module contains the following parameters:

| Key           | Type                                                      | Example                                |
|---------------|-----------------------------------------------------------|----------------------------------------|
| UploadTimeout | uint64 (time s)                                           | 600                                    |
| StorageCosts  | []StorageCost (storageProviderId, cost in tkyve per byte) | ["storage_provider_id": 1, "cost": 25] |
| NetworkFee    | sdk.Dec (%)                                               | "0.01"                                 |
| MaxPoints     | uint64                                                    | 5                                      |
