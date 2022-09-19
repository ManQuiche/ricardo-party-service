package boot

func ListenEvents() {
	_, _ = natsEncConn.Subscribe(natsRegisterTopic, asyncHandler.OnAccountCreated)
	_, _ = natsEncConn.Subscribe(natsUserDeleted, asyncHandler.OnUserDelete)
}
