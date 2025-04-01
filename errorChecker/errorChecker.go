package errorChecker

import "log"

type ErrorChecker struct {
	message        string
	defaultMessage string
}

func NewErrorChecker() *ErrorChecker {
	defaultMessage := "error"
	return &ErrorChecker{
		message:        defaultMessage,
		defaultMessage: defaultMessage,
	}
}

func (ec *ErrorChecker) WithMessage(message string) *ErrorChecker {
	ec.message = message + ": %s"
	return ec
}

func (ec *ErrorChecker) ResetMessage() {
	ec.message = ec.defaultMessage
}

func (ec *ErrorChecker) CheckErr(e error) {
	if e != nil {
		log.Fatalf(ec.message, e)
	}
	ec.ResetMessage()
}
