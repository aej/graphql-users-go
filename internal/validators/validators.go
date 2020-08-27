package validators

import (
	"errors"
	"fmt"
	"github.com/go-ozzo/ozzo-validation"
)

func MinLength(minLength int) validation.RuleFunc {
	return func(value interface{}) error {
		s, _ := value.(string)
		if len(s) < minLength {
			message := fmt.Sprintf("must be more than %v characters", minLength)
			return errors.New(message)
		}
		return nil
	}
}
