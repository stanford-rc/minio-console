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

import React from "react";

export const promoHeader = (
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

export const promoInfo = (
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
      For assistance or any questions regarding your storage, please contact
      us at{" "}
      <a href="mailto:srcc-stanford@stanford.edu">
        srcc-stanford@stanford.edu
      </a>
      .
    </p>
    <p
      style={{
        fontSize: "0.8em",
        opacity: 0.8,
        marginTop: "10px",
        marginBottom: 0,
      }}
    >
      This service utilizes MinIO software. MinIO® is a registered trademark
      of MinIO, Inc. This customized console is provided by Stanford Research
      Computing and is not affiliated with or endorsed by MinIO, Inc. The
      source code for our modified MinIO server is available at{" "}
      <a
        href="https://github.com/stanford-rc/elm-minio"
        target="_blank"
        rel="noopener noreferrer"
      >
        github.com/stanford-rc/elm-minio
      </a>
      . Our modifications are released under the GNU Affero General Public
      License v3, the same license as MinIO itself.
    </p>
  </>
);
