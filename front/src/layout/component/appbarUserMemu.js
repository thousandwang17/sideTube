/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-27 15:51:55
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-04-05 13:56:39
 * @FilePath: /sidetube/src/layout/appbarMemu.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

import * as React from "react";
import MenuItem from "@mui/material/MenuItem";
import Menu from "@mui/material/Menu";
import axios from "axios";
import MuiAlert from "@mui/material/Alert";
import { AuthContext } from "common/context/authcontext";
import { getUserID } from "common/jwt";
import { Link } from "react-router-dom";
import getHost from "common/axios";

const Alert = React.forwardRef(function Alert(props, ref) {
  return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />;
});

export default function AppBarUserMenu({
  isMenuOpen,
  handleMenuClose,
  anchorEl,
  menuId,
}) {
  const { isUserLoggedIn, setIsUserLoggedIn } = React.useContext(AuthContext);

  const handleLogout = (e) => {
    e.preventDefault();

    async function Logout() {
      try {
        const response = await UserAxios()
          .post("/logout", {}, { withCredentials: true })
          .then(() => {
            setIsUserLoggedIn(false);
          });
      } catch (e) {
        console.log(e);
      }
    }
    Logout();
    handleMenuClose();
  };

  return (
    <Menu
      anchorEl={anchorEl}
      anchorOrigin={{
        vertical: "top",
        horizontal: "right",
      }}
      id={menuId}
      keepMounted
      transformOrigin={{
        vertical: "top",
        horizontal: "right",
      }}
      open={isMenuOpen}
      onClose={handleMenuClose}
    >
      <MenuItem>
        <Link
          to={`channel/myChannel`}
          style={{
            textDecoration: "none",
            color: "inherit",
            width: "100%",
          }}
          onClick={handleMenuClose}
        >
          My Channel
        </Link>
      </MenuItem>
      <MenuItem onClick={handleLogout}>log out</MenuItem>
    </Menu>
  );
}

function UserAxios() {
  // Set config defaults when creating the instance
  const instance = axios.create({
    baseURL: getHost() + "/api/user",
  });

  instance.defaults.headers.post["Content-Type"] = "application/json";
  // defaultHeaders(instance)
  return instance;
}
