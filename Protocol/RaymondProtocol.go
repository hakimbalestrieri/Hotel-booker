package Protocol

/**
File: RaymondProtocol.go
Authors: Hakim Balestrieri
Date: 27.11.2021
*/

type RaymondStatus int

const (
	RAY_REQ RaymondStatus = iota
	RAY_WAIT
	RAY_END
	RAY_REQ_SENDER
	RAY_TOKEN

	RAY_NO
	RAY_SC
)

type RaymondProtocol struct {
	ReqType  RaymondStatus
	ServerId int
	ParentId int
}
