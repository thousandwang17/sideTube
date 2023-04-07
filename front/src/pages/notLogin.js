/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-03-24 16:11:15
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-31 16:58:56
 * @FilePath: /sidetube/src/pages/notLogin.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Button } from "@mui/material";
import { Box } from "@mui/system";
import React from "react";
import { AuthContext } from "common/context/authcontext";

export default function NotLogIn() {
  const { showLoginBox } = React.useContext(AuthContext);

  const handleButtonOnClick = () => {
    showLoginBox();
  };

  return (
    <Box
      sx={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        height: "calc(100vh - 100px)",
      }}
    >
      <Button size="large" variant="outlined" onClick={handleButtonOnClick}>
        Log in and continue
      </Button>
    </Box>
  );
}
