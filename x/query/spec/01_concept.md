<!--
order: 1
-->

# Concepts

The queries module is a little different from the other modules. It does
not maintain a state. Its purpose is to have one place to manage all queries.
A lot of queries require interaction with multiple modules and often do not
belong to a single module. 

Most queries align with the cosmos convention. Api documentation can be found
in the generated swagger file or in the proto files.

## Finalized Bundles

Finalized bundles are one of the main features of KYVE. This query will also be
the main query for people building applications on top of KYVE and using
KYVE's data.

### Bundles query

The basic structure of the bundles query works as follows.
For the field `finalized_bundles` always the latest schema version is 
returned. The different version are explained below.

**Query**: `/kyve/v1/bundles/{poolId}`

**Params**:

| Name              | Type    | Description                                           |
|-------------------|---------|-------------------------------------------------------|
| pagination.limit  | number  | Defines the amount of bundles returned                |
| pagination.offset | number  | The amount of bundles to skip                         |
| pagination.key    | string  | Define key if next_key iteration should be used.      |
| pagination.revers | boolean | Reverse order                                         |
| index             | number  | Filters for the bundle which contains the given index |


**Response**:
```yaml
{
  "finalized_bundles": "[]FinalizedBundle",
  "pagination": {
    next_key: "string",
    total: number
  }
}
```


#### Version 1

```yaml
{
  "pool_id": "number",
  "id": "number",
  "storage_id": "string",
  "uploader": "string",
  "from_index": "number",
  "to_index": "number",
  "to_key": "number",
  "bundle_summary": "string",
  "data_hash": "string",
  "finalized_at": {
    "height": "number",
    "timestamp": "number"
  },
  "from_key": "number",
  "storage_provider_id": "number",
  "compression_id": "number",
}
```

#### Version 2

For version 2 the field `stake_security` was added. Bundles which 
were finalized before the field existed return null.

```yaml
{
  "stake_security": "number"|null
}
```

### Query by ID
To obtain a specific bundle specified by its ID use

**Query**: `/kyve/v1/bundles/{poolId}/{id}`
