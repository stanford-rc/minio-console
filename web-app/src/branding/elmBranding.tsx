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

// Central source of truth for Elm-specific console branding: promo panel
// content, the license blurb, and the login footer. Shared by the LoginPage
// and Console/License screens so branding lives in one neutral location.

import React from "react";
import { Box } from "mds";

// Single source of truth for the modified MinIO server's source link.
//
// This points at the fork itself, which is where the modifications are. The
// build repository, elm-minio, is a checkout and a compile and no longer
// carries any source changes of its own.
export const MINIO_SOURCE_URL = "https://github.com/stanford-rc/minio";

// The console is separately modified AGPL software that the user is
// interacting with over the network, so it needs its own source offer rather
// than being covered by the server's.
export const CONSOLE_SOURCE_URL =
  "https://github.com/stanford-rc/minio-console";

// Documentation for the service the user is actually logging in to.
//
// This replaced a link to docs.min.io. MinIO folded the community object-store
// docs into their commercial AIStor product, so both
// docs.min.io/community/minio-object-store/index.html and the bare docs.min.io
// now redirect to docs.min.io/aistor/, a landing page for software this is
// not. Sending an Elm user there answers no question they have and implies
// they are running AIStor.
//
// Exported because the console's help menu offers the same link and should
// not carry its own copy of the URL.
export const ELM_DOCS_URL = "https://docs.elm.stanford.edu/";

export const PromoHeader = (
  <span
    style={{
      fontSize: "clamp(6px, 4vw, 56px)",
      lineHeight: 1,
      display: "inline-block",
      width: "100%",
    }}
  >
    Welcome to Elm
  </span>
);

export const PromoLicense = () => (
  <p
    style={{
      fontSize: "0.8em",
      opacity: 0.8,
      marginTop: "10px",
      marginBottom: 0,
    }}
  >
    This service utilizes MinIO software. MinIO® is a registered trademark of
    MinIO, Inc. This customized console is provided by Stanford Research
    Computing and is not affiliated with or endorsed by MinIO, Inc. The source
    code for our modified MinIO server is available at{" "}
    <a href={MINIO_SOURCE_URL} target="_blank" rel="noopener noreferrer">
      github.com/stanford-rc/minio
    </a>
    , and the source code for this modified console at{" "}
    <a href={CONSOLE_SOURCE_URL} target="_blank" rel="noopener noreferrer">
      github.com/stanford-rc/minio-console
    </a>
    . Our modifications to both are released under the GNU Affero General Public
    License v3, the same license as MinIO itself.
  </p>
);

export const PromoInfo = (
  <>
    <p style={{ marginTop: 0 }}>
      Welcome to the MinIO management console for <strong>Elm</strong>, the
      object storage service provided by{" "}
      <strong>Stanford Research Computing</strong>.
    </p>
    <p style={{ marginTop: "0.5em", marginBottom: "0.5em" }}>
      This user interface is a customized version of the open-source MinIO
      Console. Our build is based on the community-maintained{" "}
      <a
        href="https://github.com/georgmangold/console"
        target="_blank"
        rel="noopener noreferrer"
      >
        georgmangold/console
      </a>{" "}
      fork and is tailored to support the needs of the Stanford research
      community.
    </p>
    <p style={{ marginTop: "0.5em", marginBottom: "0.5em" }}>
      For assistance or any questions regarding your storage, please contact us
      at{" "}
      <a href="mailto:srcc-stanford@stanford.edu">
        srcc-stanford@stanford.edu
      </a>
      .
    </p>

    <PromoLicense />
  </>
);

export const PromoFooter = () => (
  <Box
    sx={{
      "& .separator": {
        marginLeft: 4,
        marginRight: 4,
      },
    }}
  >
    <a href={ELM_DOCS_URL} target="_blank" rel="noopener">
      Elm Documentation
    </a>
    <span className={"separator"}>|</span>
    <a href={MINIO_SOURCE_URL} target="_blank" rel="noopener">
      GitHub
    </a>
  </Box>
);
