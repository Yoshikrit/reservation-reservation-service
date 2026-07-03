package utility

import (
	"fmt"
	"github.com/Yoshikrit/reservation/internal/pkg/apperror"

	gonanoid "github.com/matoous/go-nanoid"
)

func CreateUUID(name string) (string, *apperror.AppError) {
	alphabet := "23456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ"

	// สุ่มความยาว 10 ตัวอักษร
	randomStr, err := gonanoid.Generate(alphabet, 10)
	if err != nil {
		return "", apperror.NewError(50000000, err)
	}

	return fmt.Sprintf("%s-%s", name, randomStr), nil
}
