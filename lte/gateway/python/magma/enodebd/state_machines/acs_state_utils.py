"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging
from typing import Any, Optional, Dict, Type, List
from magma.enodebd.data_models.data_model import DataModel
from magma.enodebd.data_models.data_model_parameters import ParameterName
from magma.enodebd.device_config.enodeb_configuration import \
    EnodebConfiguration
from magma.enodebd.devices.device_utils import get_device_name, \
    EnodebDeviceName
from magma.enodebd.exceptions import Tr069Error, ConfigurationError, \
    IncorrectDeviceHandlerError
from magma.enodebd.tr069 import models
from magma.enodebd.tr069.models import Tr069ComplexModel


def process_inform_message(
    inform: Any,
    device_name: EnodebDeviceName,
    data_model: DataModel,
    device_cfg: EnodebConfiguration,
) -> None:
    """
    Modifies the device configuration based on what is received in the Inform
    message. Will raise an error if it turns out that the data model we are
    using is incorrect. This is decided based on the device OUI and
    software-version that is reported in the Inform message.

    Args:
        inform: Inform Tr069 message
        device_handler: The state machine we are using for our device
    """
    logging.info('Processing Inform message')

    def _get_param_value_from_path_suffix(
        suffix: str,
        path_list: List[str],
        param_values_by_path: Dict[str, Any],
    ) -> Any:
        for path in path_list:
            if path.endswith(suffix):
                return param_values_by_path[path]
        raise ConfigurationError('Did not receive expected info in Inform')

    if hasattr(inform, 'ParameterList') and \
            hasattr(inform.ParameterList, 'ParameterValueStruct'):
        param_values_by_path = {}
        for param_value in inform.ParameterList.ParameterValueStruct:
            path = param_value.Name
            value = param_value.Value.Data
            logging.debug('(Inform msg) Received parameter: %s = %s', path,
                          value)
            param_values_by_path[path] = value

        # Check the OUI and version number to see if the data model matches
        path_list = list(param_values_by_path.keys())
        device_oui = _get_param_value_from_path_suffix(
            'DeviceInfo.ManufacturerOUI',
            path_list,
            param_values_by_path,
        )
        sw_version = _get_param_value_from_path_suffix(
            'DeviceInfo.SoftwareVersion',
            path_list,
            param_values_by_path,
        )
        logging.info('OUI: %s, Software: %s', device_oui, sw_version)
        correct_device_name = get_device_name(device_oui, sw_version)
        if device_name is not correct_device_name:
            raise IncorrectDeviceHandlerError(
                device_name=correct_device_name)

        param_name_list = data_model.get_parameter_names()
        name_to_val = {}
        for name in param_name_list:
            path = data_model.get_parameter(name).path
            if path in param_values_by_path:
                value = param_values_by_path[path]
                name_to_val[name] = value

        for name, val in name_to_val.items():
            device_cfg.set_parameter(name, val)


def are_tr069_params_equal(param_a: Any, param_b: Any, type_: str) -> bool:
    """
    Compare two parameters in TR-069 format.
    The following differences are ignored:
    - Leading and trailing whitespace, commas and quotes
    - Capitalization, for booleans (true, false)
    Returns:
        True if params are the same
    """
    # Cast booleans to integers
    cmp_a, cmp_b = param_a, param_b
    if type_ == 'boolean' and cmp_b in ('0', '1') or cmp_a in ('0', '1'):
        cmp_a, cmp_b = map(int, (cmp_a, cmp_b))
    cmp_a, cmp_b = map(str, (cmp_a, cmp_b))
    cmp_a, cmp_b = map(lambda s: s.strip(', \'"'), (cmp_a, cmp_b))
    if cmp_a.lower() in ['true', 'false']:
        cmp_a, cmp_b = map(lambda s: s.lower(), (cmp_a, cmp_b))
    return cmp_a == cmp_b


def get_all_objects_to_add(
    desired_cfg: EnodebConfiguration,
    device_cfg: EnodebConfiguration,
) -> List[ParameterName]:
    """
    Find a ParameterName that needs to be added to the eNB configuration,
    if any
    """
    desired = desired_cfg.get_object_names()
    current = device_cfg.get_object_names()
    return list(set(desired).difference(set(current)))


def get_all_objects_to_delete(
    desired_cfg: EnodebConfiguration,
    device_cfg: EnodebConfiguration,
) -> List[ParameterName]:
    """
    Find a ParameterName that needs to be deleted from the eNB configuration,
    if any
    """
    desired = desired_cfg.get_object_names()
    current = device_cfg.get_object_names()
    return list(set(current).difference(set(desired)))


def get_params_to_get(
    device_cfg: EnodebConfiguration,
    data_model: DataModel,
) -> List[ParameterName]:
    """
    Returns the names of params not belonging to objects that are added/removed
    """
    desired_names = data_model.get_present_params()
    known_names = device_cfg.get_parameter_names()
    names = list(set(desired_names) - set(known_names))
    return names


def get_object_params_to_get(
    desired_cfg: Optional[EnodebConfiguration],
    device_cfg: EnodebConfiguration,
    data_model: DataModel,
) -> List[ParameterName]:
    """
    Returns a list of parameter names for object parameters we don't know the
    current value of
    """
    names = []
    # TODO: This might a string for some strange reason, investigate why
    num_plmns = \
        int(device_cfg.get_parameter(ParameterName.NUM_PLMNS))
    for i in range(1, num_plmns + 1):
        obj_name = ParameterName.PLMN_N % i
        if not device_cfg.has_object(obj_name):
            device_cfg.add_object(obj_name)
        obj_to_params = data_model.get_numbered_param_names()
        desired = obj_to_params[obj_name]
        current = []
        if desired_cfg is not None:
            current = desired_cfg.get_parameter_names_for_object(obj_name)
        names_to_add = list(set(desired) - set(current))
        names = names + names_to_add
    return names


# We don't attempt to set these parameters on the eNB configuration
READ_ONLY_PARAMETERS = [
    ParameterName.OP_STATE,
    ParameterName.RF_TX_STATUS,
    ParameterName.GPS_STATUS,
    ParameterName.PTP_STATUS,
    ParameterName.MME_STATUS,
    ParameterName.GPS_LAT,
    ParameterName.GPS_LONG,
]


def get_param_values_to_set(
    desired_cfg: EnodebConfiguration,
    device_cfg: EnodebConfiguration,
    data_model: DataModel,
    exclude_admin: bool = False,
) -> Dict[ParameterName, Any]:
    """
    Get a map of param names to values for parameters that we will
    set on the eNB's configuration, excluding parameters for objects that can
    be added/removed.

    Also exclude special parameters like admin state, since it may be set at
    a different time in the provisioning process than most parameters.
    """
    param_values = {}
    # Get the parameters we might set
    params = set(desired_cfg.get_parameter_names()) - set(READ_ONLY_PARAMETERS)
    if exclude_admin:
        params = set(params) - {ParameterName.ADMIN_STATE}
    # Values of parameters
    for name in params:
        new = desired_cfg.get_parameter(name)
        old = device_cfg.get_parameter(name)
        _type = data_model.get_parameter(name).type
        if not are_tr069_params_equal(new, old, _type):
            param_values[name] = new

    return param_values


def get_obj_param_values_to_set(
    desired_cfg: EnodebConfiguration,
    device_cfg: EnodebConfiguration,
    data_model: DataModel,
) -> Dict[ParameterName, Dict[ParameterName, Any]]:
    """ Returns a map from object name to (a map of param name to value) """
    param_values = {}
    objs = desired_cfg.get_object_names()
    for obj_name in objs:
        param_values[obj_name] = {}
        params = desired_cfg.get_parameter_names_for_object(obj_name)
        for name in params:
            new = desired_cfg.get_parameter_for_object(name, obj_name)
            old = device_cfg.get_parameter_for_object(name, obj_name)
            _type = data_model.get_parameter(name).type
            if not are_tr069_params_equal(new, old, _type):
                param_values[obj_name][name] = new
    return param_values


def get_all_param_values_to_set(
    desired_cfg: EnodebConfiguration,
    device_cfg: EnodebConfiguration,
    data_model: DataModel,
    exclude_admin: bool = False,
) -> Dict[ParameterName, Any]:
    """ Returns a map of param names to values that we need to set """
    param_values = get_param_values_to_set(desired_cfg, device_cfg,
                                           data_model, exclude_admin)
    obj_param_values = get_obj_param_values_to_set(desired_cfg, device_cfg,
                                                   data_model)
    for _obj_name, param_map in obj_param_values.items():
        for name, val in param_map.items():
            param_values[name] = val
    return param_values


def parse_get_parameter_values_response(
    data_model: DataModel,
    message: models.GetParameterValuesResponse,
) -> Dict[ParameterName, Any]:
    """ Returns a map of ParameterName to the value read from the response """
    assert_is_msg_type_and_not_fault(message,
                                     models.GetParameterValuesResponse)
    param_values_by_path = {}
    for param_value_struct in message.ParameterList.ParameterValueStruct:
        param_values_by_path[param_value_struct.Name] = \
            param_value_struct.Value.Data

    param_name_list = data_model.get_parameter_names()
    name_to_val = {}
    for name in param_name_list:
        path = data_model.get_parameter(name).path
        if path in param_values_by_path:
            value = param_values_by_path[path]
            name_to_val[name] = value

    return name_to_val


def get_optional_param_to_check(
    data_model: DataModel,
) -> Optional[ParameterName]:
    """
    If there is a parameter which is optional in the data model, and we do not
    know if it exists or not, then return it so we can check for its presence.
    """
    params = data_model.get_names_of_optional_params()
    for param in params:
        try:
            data_model.is_parameter_present(param)
        except KeyError:
            return param
    return None


def assert_is_msg_type_and_not_fault(
    message: Tr069ComplexModel,
    type_: Type[Tr069ComplexModel],
) -> None:
    if type(message) == models.Fault:
        raise Tr069Error(
            'Received Fault in response to GetParameterValues '
            '(faultstring = %s)' % message.FaultString)
    elif type(message) != type_:
        raise Tr069Error(
            'Unexpected response type: %s' % type(message))

