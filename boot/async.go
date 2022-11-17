package boot

func ListenEvents() {
	_, _ = natsEncConn.Subscribe(natsUserCreated, userAsync.Created)
	_, _ = natsEncConn.Subscribe(natsUserUpdated, userAsync.Updated)
	_, _ = natsEncConn.Subscribe(natsUserDeleted, userAsync.Deleted)

	_, _ = natsEncConn.Subscribe(natsPartyRequested)
	_, _ = natsEncConn.Subscribe(natsPartyJoined)
}
