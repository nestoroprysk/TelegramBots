package mock

import "github.com/nestoroprysk/TelegramBots/internal/responder"

type Responder struct {
	onSucceed func(interface{}) error
	onFail    func(error) error
	onError   func(error) error
}

var _ responder.Responder = &Responder{}

func NoOpResponder() responder.Responder {
	return &Responder{}
}

func ForwardFailErrResponder() responder.Responder {
	f := func(err error) error {
		return err
	}

	return &Responder{
		onSucceed: func(interface{}) error {
			panic("shouldn't call Succeed")
		},
		onFail:  f,
		onError: f,
	}
}

func SucceedResponder(f func(interface{}) error) responder.Responder {
	return &Responder{
		onSucceed: f,
		onFail: func(error) error {
			panic("shouldn't call Fail")
		},
		onError: func(error) error {
			panic("shouldn't call Error")
		},
	}
}

func (r Responder) Succeed(b interface{}) error {
	if r.onSucceed != nil {
		return r.onSucceed(b)
	}

	return nil
}

func (r Responder) Fail(err error) error {
	if r.onFail != nil {
		return r.onFail(err)
	}

	return nil
}

func (r Responder) Error(err error) error {
	if r.onError != nil {
		return r.onError(err)
	}

	return nil
}

func (Responder) Close() error {
	return nil
}
