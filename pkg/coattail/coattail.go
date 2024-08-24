package coattail

/* ====== Singleton Instance ====== */

var local *Peer

func Manage() *Peer {
	if local == nil {
		local = newPeer(
			PeerDetails{
				PeerID: LocalPeerId,
			},
			&localPeerAdapter{
				units: []anyUnit{},
				peers: []PeerDetails{},
			},
		)
	}

	return local
}
