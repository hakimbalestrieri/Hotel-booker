package Protocol
/**
File: UpdateProtocol.go
Authors: Hakim Balestrieri & Alexandre Mottier
Date: 27.11.2021
*/

type UpdateType int

const (
	UPD_CLIENT UpdateType = iota
	UPD_ROOM
)

type UpdateProtocol struct {
	ReqType  UpdateType
	Arguments []string
	clientId int
	ServerIdFrom int
	ServerIdTo []int
}