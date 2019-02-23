"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
from enum import Enum

# Register names
IMSI_REG = 'metadata'
DIRECTION_REG = 'reg1'
SCRATCH_REGS = ['reg0']

# Register values
REG_ZERO_VAL = 0x0


class Direction(Enum):
    """
    Direction bits for direction reg
    """
    OUT = 0x01
    IN = 0x10


def load_direction(parser, direction: Direction):
    """
    Wrapper for loading the direction register
    """
    if not is_valid_direction(direction):
        raise Exception("Invalid direction")
    return parser.NXActionRegLoad2(dst=DIRECTION_REG, value=direction.value)


def is_valid_direction(direction: Direction):
    return isinstance(direction, Direction)
