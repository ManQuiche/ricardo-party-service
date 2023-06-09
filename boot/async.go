package boot

func ListenEvents() {
	_, _ = natsEncConn.Subscribe(natsUserCreated, userAsync.Created)
	_, _ = natsEncConn.Subscribe(natsUserUpdated, userAsync.Updated)
	_, _ = natsEncConn.Subscribe(natsUserDeleted, userAsync.Deleted)

	_, _ = natsEncConn.Subscribe(natsPartyRequested, partyAsync.Requested)
	_, _ = natsEncConn.Subscribe(natsPartyJoined, partyAsync.Joined)
}
