// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// AccessGatewayRecord access gateway record
// swagger:model access_gateway_record
type AccessGatewayRecord struct {

	// hw id
	// Required: true
	HwID *HwGatewayID `json:"hw_id" magma_alt_name:"HwId"`

	// key
	// Required: true
	Key *ChallengeKey `json:"key"`

	// name
	// Min Length: 1
	Name string `json:"name,omitempty"`
}

// Validate validates this access gateway record
func (m *AccessGatewayRecord) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateHwID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateKey(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *AccessGatewayRecord) validateHwID(formats strfmt.Registry) error {

	if err := validate.Required("hw_id", "body", m.HwID); err != nil {
		return err
	}

	if m.HwID != nil {
		if err := m.HwID.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("hw_id")
			}
			return err
		}
	}

	return nil
}

func (m *AccessGatewayRecord) validateKey(formats strfmt.Registry) error {

	if err := validate.Required("key", "body", m.Key); err != nil {
		return err
	}

	if m.Key != nil {
		if err := m.Key.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("key")
			}
			return err
		}
	}

	return nil
}

func (m *AccessGatewayRecord) validateName(formats strfmt.Registry) error {

	if swag.IsZero(m.Name) { // not required
		return nil
	}

	if err := validate.MinLength("name", "body", string(m.Name), 1); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *AccessGatewayRecord) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AccessGatewayRecord) UnmarshalBinary(b []byte) error {
	var res AccessGatewayRecord
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
