package firstfailure

import (
	"context"
	"errors"
	"fmt"
	"go.temporal.io/sdk/temporal"
	"math/rand"
)

const (
	if_chance   = 0.10 // InsufficientFundsError
	auth_chance = 0.05 // AuthFailure, e.g., 403
	ce_chance   = 0.10 // generic ConnectionError
)

type Activities struct{}

type InsufficientFundsError struct {
	Message string
}

func (e InsufficientFundsError) Error() string {
	return "InsufficientFundsError should NOT be retried"
}

type ConnectionError struct {
	code    int
	message string
}

func (e ConnectionError) Error() string {
	return fmt.Sprintf("Error code %d: %s", e.code, e.message)
}

func (a *Activities) CheckInventory(ctx context.Context, order OrderInfo) error {
	return nil
}

func (a *Activities) Charge(ctx context.Context, order OrderInfo) error {
	/* ... attempt to charge CC ... */
	err := chargeCreditCard("", 0.0)

	if err != nil {
		var connErr *ConnectionError
		if errors.As(err, &connErr) {
			if connErr.code == 403 {
				return temporal.NewNonRetryableApplicationError(
					connErr.message,
					"ConnectionError",
					err,
				)
			}
		}
	}

	return err
}

func (a *Activities) FulfillOrder(ctx context.Context, order OrderInfo) error {
	return nil
}

func chargeCreditCard(acct string, amt float64) error {
	r := rand.Float32()
	if r <= if_chance {
		return &InsufficientFundsError{}
	} else if r <= if_chance+auth_chance {
		return &ConnectionError{code: 403, message: "Authentication Failure"}
	} else if r <= if_chance+auth_chance+ce_chance {
		return &ConnectionError{code: 500, message: "Unknown connection failure"}
	}

	return nil
}
