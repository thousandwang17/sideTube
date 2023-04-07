/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-03-03 14:15:49
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-04-05 13:57:25
 * @FilePath: /sidetube/src/common/recommendAxios.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import axios from "axios";
import getHost from "common/axios";
import getJWT from "common/jwt";

export default function recommendAxios() {
  // Set config defaults when creating the instance
  const instance = axios.create({
    baseURL: getHost() + "/api/recommend",
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
