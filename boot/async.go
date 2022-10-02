package boot

func ListenEvents() {
	_, _ = natsEncConn.Subscribe(natsUserCreated, asyncHandler.Created)
	_, _ = natsEncConn.Subscribe(natsUserUpdated, asyncHandler.Updated)
	_, _ = natsEncConn.Subscribe(natsUserDeleted, asyncHandler.Deleted)
}
