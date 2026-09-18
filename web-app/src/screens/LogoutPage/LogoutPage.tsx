// This file is part of MinIO Console Server
// Copyright (c) 2022 MinIO, Inc.
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

import React, { useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useAppDispatch } from "../../store";
import { ErrorResponseHandler } from "../../common/types";
import { clearSession } from "../../common/utils";
import { userLogged } from "../../systemSlice";
import { resetSession } from "../Console/consoleSlice";
import request from "superagent";
import LoadingComponent from "../../common/LoadingComponent";

const LogoutPage = () => {
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  useEffect(() => {
    // postLogoutUrl, when the server supplies one, is where the browser has to
    // go next so that the session is also ended at the identity provider.
    // Clearing local state alone is not a logout: the IdP session survives and
    // the next login returns with no credential prompt.
    const deleteSession = (postLogoutUrl?: string) => {
      dispatch(userLogged(false));
      // Disconnect OB Websocket
      dispatch({ type: "socket/OBDisconnect" });
      localStorage.setItem("userLoggedIn", "");
      localStorage.setItem("redirect-path", "");
      dispatch(resetSession());
      clearSession();

      if (postLogoutUrl) {
        // replace() rather than assign() so that Back does not return to a
        // page belonging to the session just ended.
        window.location.replace(postLogoutUrl);
        return;
      }

      navigate("/login");
      window.location.reload(); //reset-all redux states etc. by force reloading.
    };

    const logout = () => {
      const state = localStorage.getItem("auth-state");
      // superagent directly rather than api.invoke, because invoke resolves
      // with res.body and discards the response, and the redirect target
      // arrives as a header. Changing invoke would affect every caller.
      request
        .post("api/v1/logout")
        .send({ state })
        .then((res) => {
          deleteSession(res.headers["x-console-post-logout-url"]);
        })
        .catch((err: ErrorResponseHandler) => {
          console.error(err);
          // Still tear down what we can: a failed IdP logout must not leave
          // the user logged in locally.
          deleteSession();
        });
    };
    logout();
  }, [dispatch, navigate]);
  return <LoadingComponent />;
};

export default LogoutPage;
