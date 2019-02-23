// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// DiameterClientConfigs Diameter Configuration of The Client
// swagger:model diameter_client_configs
type DiameterClientConfigs struct {

	// address
	// Pattern: [^\:]+(:[0-9]{1,5})?
	Address string `json:"address,omitempty"`

	// dest host
	DestHost string `json:"dest_host,omitempty"`

	// dest realm
	DestRealm string `json:"dest_realm,omitempty"`

	// host
	// Min Length: 1
	Host string `json:"host,omitempty"`

	// local address
	// Pattern: [0-9a-f\:\.]*(:[0-9]{1,5})?
	LocalAddress string `json:"local_address,omitempty"`

	// product name
	// Min Length: 1
	ProductName string `json:"product_name,omitempty"`

	// protocol
	// Enum: [tcp tcp4 tcp6 sctp sctp4 sctp6]
	Protocol string `json:"protocol,omitempty"`

	// realm
	// Min Length: 1
	Realm string `json:"realm,omitempty"`

	// retransmits
	Retransmits uint32 `json:"retransmits,omitempty"`

	// retry count
	RetryCount uint32 `json:"retry_count,omitempty"`

	// watchdog interval
	WatchdogInterval uint32 `json:"watchdog_interval,omitempty"`
}

// Validate validates this diameter client configs
func (m *DiameterClientConfigs) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHost(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateLocalAddress(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProductName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProtocol(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRealm(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *DiameterClientConfigs) validateAddress(formats strfmt.Registry) error {

	if swag.IsZero(m.Address) { // not required
		return nil
	}

	if err := validate.Pattern("address", "body", string(m.Address), `[^\:]+(:[0-9]{1,5})?`); err != nil {
		return err
	}

	return nil
}

func (m *DiameterClientConfigs) validateHost(formats strfmt.Registry) error {

	if swag.IsZero(m.Host) { // not required
		return nil
	}

	if err := validate.MinLength("host", "body", string(m.Host), 1); err != nil {
		return err
	}

	return nil
}

func (m *DiameterClientConfigs) validateLocalAddress(formats strfmt.Registry) error {

	if swag.IsZero(m.LocalAddress) { // not required
		return nil
	}

	if err := validate.Pattern("local_address", "body", string(m.LocalAddress), `[0-9a-f\:\.]*(:[0-9]{1,5})?`); err != nil {
		return err
	}

	return nil
}

func (m *DiameterClientConfigs) validateProductName(formats strfmt.Registry) error {

	if swag.IsZero(m.ProductName) { // not required
		return nil
	}

	if err := validate.MinLength("product_name", "body", string(m.ProductName), 1); err != nil {
		return err
	}

	return nil
}

var diameterClientConfigsTypeProtocolPropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["tcp","tcp4","tcp6","sctp","sctp4","sctp6"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		diameterClientConfigsTypeProtocolPropEnum = append(diameterClientConfigsTypeProtocolPropEnum, v)
	}
}

const (

	// DiameterClientConfigsProtocolTCP captures enum value "tcp"
	DiameterClientConfigsProtocolTCP string = "tcp"

	// DiameterClientConfigsProtocolTcp4 captures enum value "tcp4"
	DiameterClientConfigsProtocolTcp4 string = "tcp4"

	// DiameterClientConfigsProtocolTcp6 captures enum value "tcp6"
	DiameterClientConfigsProtocolTcp6 string = "tcp6"

	// DiameterClientConfigsProtocolSctp captures enum value "sctp"
	DiameterClientConfigsProtocolSctp string = "sctp"

	// DiameterClientConfigsProtocolSctp4 captures enum value "sctp4"
	DiameterClientConfigsProtocolSctp4 string = "sctp4"

	// DiameterClientConfigsProtocolSctp6 captures enum value "sctp6"
	DiameterClientConfigsProtocolSctp6 string = "sctp6"
)

// prop value enum
func (m *DiameterClientConfigs) validateProtocolEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, diameterClientConfigsTypeProtocolPropEnum); err != nil {
		return err
	}
	return nil
}

func (m *DiameterClientConfigs) validateProtocol(formats strfmt.Registry) error {

	if swag.IsZero(m.Protocol) { // not required
		return nil
	}

	// value enum
	if err := m.validateProtocolEnum("protocol", "body", m.Protocol); err != nil {
		return err
	}

	return nil
}

func (m *DiameterClientConfigs) validateRealm(formats strfmt.Registry) error {

	if swag.IsZero(m.Realm) { // not required
		return nil
	}

	if err := validate.MinLength("realm", "body", string(m.Realm), 1); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *DiameterClientConfigs) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DiameterClientConfigs) UnmarshalBinary(b []byte) error {
	var res DiameterClientConfigs
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
