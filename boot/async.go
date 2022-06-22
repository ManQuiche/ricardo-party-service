package boot

func ListenEvents() {
	_, _ = natsEncConn.Subscribe(natsUserDeleted, asyncHandler.OnUserDelete)
}
