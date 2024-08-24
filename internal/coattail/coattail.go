package coattail

/* ====== Singleton Instance ====== */

var local *Peer

func Manage() *Peer {
	if local == nil {
		local = newPeer(&localPeerAdapter{
			units: []AnyUnit{},
			peers: []PeerDetails{},
		})
	}

	return local
}
