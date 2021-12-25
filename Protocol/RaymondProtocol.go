package Protocol

/**
File: RaymondProtocol.go
Authors: Hakim Balestrieri
Date: 25.12.2021
*/

const (
	RAYMOND_PRO_REQ RaymondType = iota
	RAYMOND_WAIT
	RAYMOND_END
	RAYMOND_REQ
	RAYMOND_TOKEN
)

const (
	RAY_NO RaymondStatusType = iota
	RAY_ASK
	RAY_SC
)

type RaymondType int
type RaymondStatusType int
type RaymondProtocol struct {
	ReqType  RaymondType
	ServerId int
}