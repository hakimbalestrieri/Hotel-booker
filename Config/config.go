/**
File: Config.go
Authors: Hakim Balestrieri
Date: 23.12.2021
*/

package Config

var Debug = true
var RoomsNumber = 30
var DayNumber = 31
var ServerPorts = map[int]int{0: 3100, 1: 3101, 2: 3102, 3: 3103, 4: 3104, 5: 3105}
var ServerNumber = len(ServerPorts)
var ServerWithoutParent = -1
var ClientPorts = map[int]int{0: 4000, 1: 4001, 2: 4002, 3: 4003, 4: 4004, 5: 4005}

type Structure struct {
	Root     int
	Children []int
}

var ServerSchema = map[int]*Structure{
	0: {
		Children: []int{3},
		Root:     -1,
	},
	1: {
		Children: []int{2},
		Root:     0,
	},
	2: {
		Children: []int{},
		Root:     1,
	},
	3: {
		Children: []int{4},
		Root:     0,
	},
	4: {
		Children: []int{},
		Root:     3,
	},
}
