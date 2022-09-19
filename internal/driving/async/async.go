package async

type Handler interface {
	OnUserDelete(userID uint)
	OnAccountCreated(userID uint)
}
