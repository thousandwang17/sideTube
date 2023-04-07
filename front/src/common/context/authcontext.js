/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-03-07 15:19:05
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-04-04 18:07:12
 * @FilePath: /sidetube/src/common/cotext/authcontext.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
// AuthContext.js
import React, { createContext, useState } from "react";
import { requestJwtTokenEveryfiveMins, checkUserLoggedIn } from "common/jwt";

export const AuthContext = createContext();

export function AuthProvider({ children }) {
  const [isUserLoggedIn, setIsUserLoggedIn] = useState(checkUserLoggedIn());
  const [showLoginBox, setShowLoginBox] = useState(
    () => () => console.log(" showLoginBox ")
  );
  requestJwtTokenEveryfiveMins(setIsUserLoggedIn);

  return (
    <AuthContext.Provider
      value={{
        isUserLoggedIn,
        setIsUserLoggedIn,
        showLoginBox,
        setShowLoginBox,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}
