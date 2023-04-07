/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-18 15:23:51
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-04-05 14:05:50
 * @FilePath: /sidetube/src/common/jwt.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import axios from "axios";
import getHost from "./axios";

let JwtToken = null;
let intervalId = null; // variable to store the interval ID for clearing later
let lastRequestTime = 0; // variable to store the interval ID for clearing later
const REQUEST_INTERVAL = 5 * 60 * 1000; // 5 minutes in milliseconds
let invokeSetIsLoginIn = null;

export default function getJWT() {
  return JwtToken;
}

function requestRefreshToken() {
  tokenAxios()
    .post("/token", {}, { withCredentials: true })
    .then((response) => {
      lastRequestTime = Date.now();
      JwtToken = getCookie("RefreshToken");
      invokeSetIsLoginIn(true);
    })
    .catch((error) => {
      if (error.response?.status === 403) {
        clear();
        invokeSetIsLoginIn(false);
      }
    });
}

export function requestJwtTokenEveryfiveMins(setIsUserLoggedIn) {
  JwtToken = getCookie("RefreshToken");
  invokeSetIsLoginIn = setIsUserLoggedIn;
  if (intervalId === null && checkUserLoggedIn()) {
    intervalId = setTimeout(requestRefreshToken, REQUEST_INTERVAL);
    // add an event listener for when the computer wakes up from sleep
    document.addEventListener("visibilitychange", visibilitychange);
  } else if (!checkUserLoggedIn() && intervalId !== null) {
    clear();
  }
}

function clear() {
  clearInterval(intervalId);
  intervalId = null;
  document.removeEventListener("visibilitychange", visibilitychange);
}

function visibilitychange() {
  if (document.visibilityState === "visible") {
    if (Date.now() - lastRequestTime < 10 * 60 * 1000) {
      return;
    }

    if (intervalId) {
      clearInterval(intervalId);
    }
    intervalId = setTimeout(requestRefreshToken, REQUEST_INTERVAL);
    requestRefreshToken(); // make a request immediately if the computer wakes up from sleep and there is a pending interval
  }
}

function tokenAxios() {
  // Set config defaults when creating the instance
  const instance = axios.create({
    baseURL: getHost() + "/api",
  });

  instance.defaults.headers.post["Content-Type"] = "application/json";
  // defaultHeaders(instance)
  return instance;
}

function getCookie(name) {
  const cookieString = document.cookie;
  const cookies = cookieString.split(";");

  for (let i = 0; i < cookies.length; i++) {
    const cookie = cookies[i].trim();
    if (cookie.startsWith(name + "=")) {
      return cookie.substring(name.length + 1);
    }
  }

  return null;
}

export function getUserName() {
  return getUserInfoFromJwt()?.name ?? "";
}

export function getUserID() {
  return getUserInfoFromJwt()?.uid ?? "";
}

export function getUserInfoFromJwt() {
  if (JwtToken === null) {
    return null;
  }

  const parts = JwtToken.split(".");
  if (parts.length < 3) {
    console.log("token format error");
    return null;
  }
  const header = parts[0];
  const payload = parts[1];
  const signature = parts[2];
  const decodedPayload = atob(payload);
  return JSON.parse(decodedPayload);
}

export function checkUserLoggedIn() {
  JwtToken = getCookie("RefreshToken");
  return !!JwtToken;
}
