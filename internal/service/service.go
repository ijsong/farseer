package service

import fmt "fmt"

type serviceError struct {
	api    string
	reason string
}

func NewNotInitiatedMessageError(api string) error {
	return serviceError{
		api:    api,
		reason: "not initiated field",
	}
}

func (e serviceError) Error() string {
	return fmt.Sprintf("service error (%s): %s", e.api, e.reason)
}
