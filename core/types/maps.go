package types

// for cache maps.

type AuthMaps struct {
	// user <-> auth channel
	Data map[int64]int64
	// user <-> 10 group ids
	GroupIds map[int64][]int64
	// user <-> 10 group message ids
	GroupMessages map[int64][]int64
	// user <-> challenge steps
	Steps map[int64]int64
	// user <-> if open chat tunnel
	ChatOpened map[int64]bool
}
