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

// PromqlData promql data
// swagger:model promql_data
type PromqlData struct {

	// result
	// Required: true
	Result PromqlResult `json:"result"`

	// result type
	// Required: true
	ResultType *string `json:"resultType"`
}

// Validate validates this promql data
func (m *PromqlData) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateResult(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResultType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *PromqlData) validateResult(formats strfmt.Registry) error {

	if err := validate.Required("result", "body", m.Result); err != nil {
		return err
	}

	if err := m.Result.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("result")
		}
		return err
	}

	return nil
}

func (m *PromqlData) validateResultType(formats strfmt.Registry) error {

	if err := validate.Required("resultType", "body", m.ResultType); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *PromqlData) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *PromqlData) UnmarshalBinary(b []byte) error {
	var res PromqlData
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
