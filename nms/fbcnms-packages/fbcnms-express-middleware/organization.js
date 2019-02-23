/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {Request, Response, NextFunction} from 'express';

function getSubdomainList(host: ?string): Array<string> {
  if (!host) {
    return [];
  }
  const subdomainList = host.split('.');
  if (subdomainList) {
    subdomainList.splice(-1, 1);
  }
  return subdomainList;
}

export async function getOrganization(
  req: Request,
  OrganizationModel: any,
): any {
  const subDomains = getSubdomainList(req.get('host'));
  if (subDomains.length != 1 && subDomains.length != 2) {
    throw new Error('Invalid organization!');
  }
  const org = await OrganizationModel.findOne({
    where: {
      name: subDomains[0],
    },
  });
  if (!org) {
    throw new Error('Invalid organization!');
  }
  return org;
}

export function organizationMiddleware({model}: {model: any}) {
  return (req: Request, res: Response, next: NextFunction) => {
    try {
      req.organization = () => getOrganization(req, model);
      next();
    } catch (err) {
      res.status(404).send();
      next(err);
    }
  };
}
