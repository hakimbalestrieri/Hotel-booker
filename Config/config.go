package Config

/**
File: Config.go
Authors: Hakim Balestrieri
Date: 27.11.2021
*/

var Debug = false
var RoomsNumber = 30
var DayNumber = 31
var ServerPorts = map[int]int{0: 3020, 1: 3000, 2: 3001}
var ServerNumber = len(ServerPorts)
