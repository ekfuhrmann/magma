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
import type {subscriber} from '../../../../../fbcnms-packages/fbcnms-magma-api';
import type {KPIRows} from '../../components/KPIGrid';

import Button from '@material-ui/core/Button';
import CardHeader from '@material-ui/core/CardHeader';
import Collapse from '@material-ui/core/Collapse';
import Divider from '@material-ui/core/Divider';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
import KPIGrid from '../../components/KPIGrid';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import Text from '../../theme/design-system/Text';
import TextField from '@material-ui/core/TextField';

import {colors, typography} from '../../theme/default';
import {CardTitleFilterRow} from '../../components/layout/CardTitleRow';
import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  kpiHeaderBlock: {
    display: 'flex',
    alignItems: 'center',
    padding: 0,
  },
  kpiHeaderContent: {
    display: 'flex',
    alignItems: 'center',
  },
  kpiHeaderIcon: {
    fill: colors.primary.comet,
    marginRight: theme.spacing(1),
  },
  kpiBlock: {
    boxShadow: `0 0 0 1px ${colors.primary.concrete}`,
  },
  kpiLabel: {
    color: colors.primary.comet,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
  },
  kpiValue: {
    color: colors.primary.brightGray,
    whiteSpace: 'nowrap',
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    width: '100%',
  },
  kpiBox: {
    width: '100%',
    '& > div': {
      width: '100%',
    },
  },
}));

export default function SubscriberDetailConfig({
  subscriberInfo,
}: {
  subscriberInfo: subscriber,
}) {
  const classes = useStyles();

  function ConfigFilter() {
    return <Button variant="text">Edit</Button>;
  }

  function TrafficFilter() {
    return <Button variant="text">Edit</Button>;
  }

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <Grid container spacing={4}>
            <Grid item xs={12} md={6} alignItems="center">
              <CardTitleFilterRow
                icon={SettingsIcon}
                label="Config"
                filter={ConfigFilter}
              />
              <SubscriberInfoConfig
                readOnly={true}
                subscriberInfo={subscriberInfo}
              />
            </Grid>

            <Grid item xs={12} md={6} alignItems="center">
              <CardTitleFilterRow
                icon={GraphicEqIcon}
                label="Traffic Policy"
                filter={TrafficFilter}
              />
              <SubscriberConfigTrafficPolicy
                readOnly={true}
                subscriberInfo={subscriberInfo}
              />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

function SubscriberConfigTrafficPolicy({
  subscriberInfo,
  readOnly,
}: {
  subscriberInfo: subscriber,
  readOnly: boolean,
}) {
  const [open, setOpen] = useState({
    activeAPN: true,
    baseNames: true,
    activePolicies: true,
  });
  const handleCollapse = (config: string) => {
    setOpen({
      ...open,
      [config]: !open[config],
    });
  };
  const classes = useStyles();

  return (
    <List component={Paper} elevation={0}>
      <ListItem button onClick={() => handleCollapse('activeAPN')}>
        <TextField
          fullWidth={true}
          value={subscriberInfo.active_apns?.length || 0}
          label="Active APNs"
          InputProps={{disableUnderline: true, readOnly: readOnly}}
        />
        {open['activeAPN'] ? <ExpandLess /> : <ExpandMore />}
      </ListItem>
      <Divider />
      <Collapse
        key="activeAPN"
        in={open['activeAPN']}
        timeout="auto"
        unmountOnExit>
        {subscriberInfo.active_apns?.map(data => (
          <>
            <ListItem>
              <ListItemText
                primary={data}
                primaryTypographyProps={{variant: 'body3'}}
              />
            </ListItem>
            <Divider />
          </>
        )) || null}
      </Collapse>
      <ListItem button onClick={() => handleCollapse('baseNames')}>
        <TextField
          fullWidth={true}
          value={subscriberInfo.active_base_names?.length || 0}
          label="Base Names"
          InputProps={{disableUnderline: true, readOnly: readOnly}}
        />
        {open['baseNames'] ? <ExpandLess /> : <ExpandMore />}
      </ListItem>
      <Divider />
      <Collapse
        key="baseNames"
        in={open['baseNames']}
        timeout="auto"
        unmountOnExit>
        <ListItem>
          <TextField
            fullWidth={true}
            value={subscriberInfo.active_base_names?.join(', ') || 0}
            InputProps={{disableUnderline: true, readOnly: readOnly}}
          />
        </ListItem>
        <Divider />
      </Collapse>
      <ListItem button onClick={() => handleCollapse('activePolicies')}>
        <TextField
          fullWidth={true}
          value={subscriberInfo.active_policies?.length || 0}
          label="Active Policies"
          InputProps={{disableUnderline: true, readOnly: readOnly}}
        />
        {open['activePolicies'] ? <ExpandLess /> : <ExpandMore />}
      </ListItem>
      <Collapse
        key="activePolicies"
        in={open['activePolicies']}
        timeout="auto"
        unmountOnExit>
        <ListItem>
          <TextField
            fullWidth={true}
            value={subscriberInfo.active_policies?.join(', ') || 0}
            InputProps={{disableUnderline: true, readOnly: readOnly}}
          />
        </ListItem>
      </Collapse>
    </List>
  );
}

function SubscriberInfoConfig({
  subscriberInfo,
  readOnly,
}: {
  subscriberInfo: subscriber,
  readOnly: boolean,
}) {
  const [authKey, setAuthKey] = useState(subscriberInfo.lte.auth_key);
  const [authOPC, setAuthOPC] = useState(subscriberInfo.lte.auth_opc ?? false);
  const [dataPlan, setDataPlan] = useState(subscriberInfo.lte.sub_profile);

  const kpiData: KPIRows[] = [
    [
      {
        category: 'LTE Network Access',
        value: subscriberInfo.lte.state,
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Data plan',
        value: dataPlan,
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Auth Key',
        value: authKey,
        statusCircle: false,
      },
    ],
  ];

  if (authOPC) {
    kpiData.push([
      {
        category: 'Auth OPC',
        value: authOPC,
        statusCircle: false,
      },
    ]);
  }

  return <KPIGrid data={kpiData} />;
}
