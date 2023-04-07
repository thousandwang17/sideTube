/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-17 17:29:17
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-04-05 13:56:58
 * @FilePath: /sidetube/src/common/axios.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import axios from "axios";
import getHost from "common/axios";
import getJWT from "../../../common/jwt";

export function videoUploadAxios() {
  // Set config defaults when creating the instance
  const instance = axios.create({
    baseURL: getHost() + "/api/videoUpload",
  });

  instance.interceptors.request.use((config) => {
    // set jwt token to  Authorization of header before all request send
    config.headers["Authorization"] = getJWT();

    return config;
  });
  instance.defaults.headers.post["Content-Type"] = "application/json";

  // defaultHeaders(instance)
  return instance;
}

export function videoListAxios() {
  // Set config defaults when creating the instance
  const instance = axios.create({
    baseURL: getHost() + "/api/studio/video",
  });

  instance.interceptors.request.use((config) => {
    // set jwt token to  Authorization of header before all request send
    config.headers["Authorization"] = getJWT();
    return config;
  });
  instance.defaults.headers.post["Content-Type"] = "application/json";
  // defaultHeaders(instance)
  return instance;
}

export function videMetaUpdateAxios() {
  // Set config defaults when creating the instance
  const instance = axios.create({
    baseURL: getHost() + "/api/studio/video",
  });

  instance.interceptors.request.use((config) => {
    // set jwt token to  Authorization of header before all request send
    config.headers["Authorization"] = getJWT();
    return config;
  });
  instance.defaults.headers.post["Content-Type"] = "application/json";
  // defaultHeaders(instance)
  return instance;
}
