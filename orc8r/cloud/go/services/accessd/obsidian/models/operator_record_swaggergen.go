// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"strconv"

	"github.com/go-openapi/errors"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// OperatorRecord Operator Record
// swagger:model operator_record
type OperatorRecord struct {

	// certificate sns
	CertificateSns []CertificateSn `json:"certificate_sns"`

	// entities
	Entities ACLType `json:"entities,omitempty"`
}

// Validate validates this operator record
func (m *OperatorRecord) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCertificateSns(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateEntities(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *OperatorRecord) validateCertificateSns(formats strfmt.Registry) error {

	if swag.IsZero(m.CertificateSns) { // not required
		return nil
	}

	for i := 0; i < len(m.CertificateSns); i++ {

		if err := m.CertificateSns[i].Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("certificate_sns" + "." + strconv.Itoa(i))
			}
			return err
		}

	}

	return nil
}

func (m *OperatorRecord) validateEntities(formats strfmt.Registry) error {

	if swag.IsZero(m.Entities) { // not required
		return nil
	}

	if err := m.Entities.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("entities")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *OperatorRecord) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OperatorRecord) UnmarshalBinary(b []byte) error {
	var res OperatorRecord
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
