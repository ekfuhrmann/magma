// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// SubscriberState MME state for a subscriber
// swagger:model subscriber_state
type SubscriberState struct {

	// mme
	Mme UntypedMmeState `json:"mme,omitempty"`

	// s1ap
	S1ap UntypedMmeState `json:"s1ap,omitempty"`

	// spgw
	Spgw UntypedMmeState `json:"spgw,omitempty"`
}

// Validate validates this subscriber state
func (m *SubscriberState) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SubscriberState) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SubscriberState) UnmarshalBinary(b []byte) error {
	var res SubscriberState
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}