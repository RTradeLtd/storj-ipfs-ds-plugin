package addrutil

import (
	ma "gx/ipfs/QmT4U94DnD8FRfqr21obWY32HLM5VExccPKMjQHofeYqr9/go-multiaddr"
	mafmt "gx/ipfs/QmcASJiigZMaU3dKghcQo6FtcbqwfZ7UFp5W4fp62f8KJE/mafmt"
)

// SubtractFilter returns a filter func that filters all of the given addresses
func SubtractFilter(addrs ...ma.Multiaddr) func(ma.Multiaddr) bool {
	addrmap := make(map[string]bool, len(addrs))
	for _, a := range addrs {
		addrmap[string(a.Bytes())] = true
	}

	return func(a ma.Multiaddr) bool {
		return !addrmap[string(a.Bytes())]
	}
}

// IsFDCostlyTransport returns true for transports that require a new file
// descriptor per connection created
func IsFDCostlyTransport(a ma.Multiaddr) bool {
	return mafmt.TCP.Matches(a)
}

// FilterNeg returns a negated version of the passed in filter
func FilterNeg(f func(ma.Multiaddr) bool) func(ma.Multiaddr) bool {
	return func(a ma.Multiaddr) bool {
		return !f(a)
	}
}
