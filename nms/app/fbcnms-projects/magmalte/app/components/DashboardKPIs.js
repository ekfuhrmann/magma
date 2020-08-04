/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */

import EnodebKPIs from './EnodebKPIs';
import GatewayKPIs from './GatewayKPIs';
import Grid from '@material-ui/core/Grid';
import React from 'react';

import {CardTitleRow} from './layout/CardTitleRow';
import {GpsFixed} from '@material-ui/icons';

export default function () {
  return (
    <>
      <CardTitleRow icon={GpsFixed} label="Events" />
      <Grid container item zeroMinWidth alignItems="center" spacing={4}>
        <Grid item xs={12} lg={6}>
          <GatewayKPIs />
        </Grid>
        <Grid item xs={12} lg={6}>
          <EnodebKPIs />
        </Grid>
      </Grid>
    </>
  );
}
