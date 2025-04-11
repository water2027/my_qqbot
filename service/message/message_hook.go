package message

type BeforeSendFunc func(msg *Message) error

type BeforeSendHook struct {
	Fn       BeforeSendFunc
	Priority uint
}