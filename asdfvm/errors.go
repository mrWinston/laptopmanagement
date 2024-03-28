package main

import "fmt"

type InvalidStateError struct {
  StateName string
}

func (ise *InvalidStateError) Error() string {
  return fmt.Sprintf("'%s' is not a valid state", ise.StateName) 
}

func NewInvalidStateError(stateName string) error{
  return &InvalidStateError{
    StateName: stateName,
  }
} 
