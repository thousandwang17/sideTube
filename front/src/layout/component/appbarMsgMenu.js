/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-27 15:51:55
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-04-02 15:56:04
 * @FilePath: /sidetube/src/layout/appbarMemu.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

import * as React from "react";
import MenuItem from "@mui/material/MenuItem";
import Menu from "@mui/material/Menu";

export default function AppBarMsgsMenu({
  isMenuOpen,
  handleMenuClose,
  anchorEl,
  menuId,
}) {
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
      <MenuItem>Welcome to SideTube</MenuItem>
    </Menu>
  );
}
