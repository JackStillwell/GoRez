package interfaces

type ListeningService interface {
	AddListener()
	RemoveListener()
	Send(interface{})
}
