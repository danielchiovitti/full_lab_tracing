package business_error

import "fmt"

type BusinessError struct {
	Message string
}

func (b *BusinessError) Error() string {
	return fmt.Sprintf("Error: %s", b.Message)
}
