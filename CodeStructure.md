# Code Structure

### Getters
`getters` belong to the keeper directory. Getter files are prefixed with 
`getters_`. All methods there always succeed. In the worst case nothing happens.
These methods never throw an error. This is the only place where methods are allowed
to write to the KV-Store. Also, all aggregation variables are updated here.


### Logic-Files
These files are prefixed with `logic_` and handle complex tasks. 
They are allowed and encouraged to emit events and call the getters functions.
All logic happens here.


### Msg-Server
Handle transactions on a high level. As much logic as possible should be forwarded
to the logic files. This file should always be easy to read.
