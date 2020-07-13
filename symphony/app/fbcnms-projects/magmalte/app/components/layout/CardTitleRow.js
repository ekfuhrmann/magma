/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {ComponentType} from 'react';

import Grid from '@material-ui/core/Grid';
import React from 'react';
import Text from '../../theme/design-system/Text';

import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  cardTitleRow: {
    marginBottom: props => (props.noMargin ? '0' : theme.spacing(1)),
    minHeight: '36px',
  },
  cardTitleIcon: {
    fill: colors.primary.comet,
    marginRight: theme.spacing(1),
  },
}));

export type CardTitleRowProps = {
  icon: ComponentType<SvgIconExports>,
  label: string,
  noMargin: boolean,
};

export const CardTitleRow = (props: CardTitleRowProps) => {
  const classes = useStyles(props);
  const Icon = props.icon;

  return (
    <Grid container alignItems="center" className={classes.cardTitleRow}>
      <Icon className={classes.cardTitleIcon} />
      <Text variant="body1">{props.label}</Text>
    </Grid>
  );
};

export type CardTitleFilterRowProps = {
  icon: ComponentType<SvgIconExports>,
  label: string,
  filter: () => React$Node,
};

export const CardTitleFilterRow = (props: CardTitleFilterRowProps) => {
  const classes = useStyles();
  const Filters = props.filter;

  return (
    <Grid
      container
      alignItems="center"
      justify="space-between"
      className={classes.cardTitleRow}>
      <Grid item>
        <CardTitleRow icon={props.icon} label={props.label} noMargin={true} />
      </Grid>
      <Grid item>
        <Filters />
      </Grid>
    </Grid>
  );
};
