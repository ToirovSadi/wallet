package wallet

import "errors"

// Errors that can occur in these functions
var ErrAccountNotFound = errors.New("account that you want doesn't exist")
var ErrPhoneRegistred = errors.New("phone already registred")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment that you asked not found")
var ErrFavoriteNotFound = errors.New("favorite payment that you ask not found")
