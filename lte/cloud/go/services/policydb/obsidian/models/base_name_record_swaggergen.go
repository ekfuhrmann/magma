// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// BaseNameRecord base name record
// swagger:model base_name_record
type BaseNameRecord struct {

	// name
	Name BaseName `json:"name,omitempty"`

	// rule names
	RuleNames RuleNames `json:"rule_names,omitempty"`
}

// Validate validates this base name record
func (m *BaseNameRecord) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRuleNames(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *BaseNameRecord) validateName(formats strfmt.Registry) error {

	if swag.IsZero(m.Name) { // not required
		return nil
	}

	if err := m.Name.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("name")
		}
		return err
	}

	return nil
}

func (m *BaseNameRecord) validateRuleNames(formats strfmt.Registry) error {

	if swag.IsZero(m.RuleNames) { // not required
		return nil
	}

	if err := m.RuleNames.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("rule_names")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *BaseNameRecord) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *BaseNameRecord) UnmarshalBinary(b []byte) error {
	var res BaseNameRecord
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
