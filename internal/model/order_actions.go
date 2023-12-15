package model

import (
	"errors"
	"fmt"
)

func (o *OrderActions) stateToState(action func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()
	action()
	return err
}
