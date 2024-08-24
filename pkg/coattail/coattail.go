package coattail

/* ====== Singleton Instance ====== */

var local *Peer

func Manage() *Peer {
	if local == nil {
		local = newPeer(
			PeerDetails{
				PeerID: "local",
			},
			&localPeerAdapter{
				units: []anyUnit{},
				peers: []PeerDetails{},
			},
		)
	}

	return local
}
