// This file is part of MinIO Console Server
// Copyright (c) 2021 MinIO, Inc.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

import { PromoLicense } from "../../LoginPage/promoContent";

import React, { Fragment } from "react";
import {
  AGPLV3DarkLogo,
  ApplicationLogo,
  Box,
  Grid,
  HelpBox,
  LicenseIcon,
  PageLayout,
} from "mds";
import PageHeaderWrapper from "../Common/PageHeaderWrapper/PageHeaderWrapper";
import { version } from "version";

const License = () => {
  return (
    <Fragment>
      <PageHeaderWrapper label={"License & About"} />
      <PageLayout>
        <Grid item xs={12}>
          <Box
            sx={{
              display: "flex",
              flexFlow: "column",
              alignItems: "center",
              justifyContent: "center",
              textAlign: "center",
              "& a": {
                color: "#2781B0",
                fontWeight: 600,
              },
            }}
          >
            <Box sx={{ my: "30px", fontSize: "30px" }}>
              Version: v{version}
            </Box>

            <Box sx={{ maxWidth: "800px", px: 2 }}>
              <PromoLicense />
            </Box>

            <Box sx={{ my: "30px", borderTop: '1px solid #ddd', width: '80%', pt: '30px' }}>
              <p>
                This software is licensed under the GNU Affero General Public License
                (AGPL) Version 3.0.
              </p>
              <a
                href="https://www.gnu.org/licenses/agpl-3.0.en.html"
                target="_blank"
                rel="noopener noreferrer"
              >
                <AGPLV3DarkLogo />
              </a>
              <p style={{ marginTop: '1em' }}>
                For more information, please refer to the full license text at{" "}
                <a
                  href="https://www.gnu.org/licenses/agpl-3.0.en.html"
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  www.gnu.org/licenses
                </a>
                .
              </p>
            </Box>
          </Box>
        </Grid>
      </PageLayout>
    </Fragment>
  );
};

export default License;
