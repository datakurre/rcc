package fail

import "fmt"

func Around(err *error) {
	original := recover()
	if original == nil {
		return
	}

	catch, ok := original.(delimited)
	if !ok {
		panic(original)
	}

	*err = catch()
}

func On(condition bool, form string, details ...interface{}) {
	if condition {
		panic(failure(form, details...))
	}
}

func failure(form string, details ...interface{}) delimited {
	err := fmt.Errorf(form, details...)
	return func() error {
		return err
	}
}

type delimited func() error
