package types

import "time"

type OracleAcknowledgement struct {
	IBCAcknowledgement []byte
	OracleResponse     []byte
	// TODO(@john): Is it better to use a unix timestamp?
	Timestamp time.Time
}
