/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-03-03 19:43:13
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-04-04 20:25:42
 * @FilePath: /sidetube/src/common/timeAgo.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
export default function timeAgo(timestamp) {
  const now = new Date().getTime();
  const timeZoneOffset = new Date().getTimezoneOffset();
  const date = new Date(timestamp) - timeZoneOffset * 60 * 1000;

  const diff = now - date;

  const year = Math.floor(diff / (1000 * 60 * 60 * 24 * 365));
  if (year > 0) {
    return `${year} year${year > 1 ? "s" : ""} ago`;
  }

  const month = Math.floor(diff / (1000 * 60 * 60 * 24 * 30));
  if (month > 0) {
    return `${month} month${month > 1 ? "s" : ""} ago`;
  }

  const week = Math.floor(diff / (1000 * 60 * 60 * 24 * 7));
  if (week > 0) {
    return `${week} week${week > 1 ? "s" : ""} ago`;
  }

  const days = Math.floor(diff / (1000 * 60 * 60 * 24));
  if (days > 0) {
    return `${days} day${days > 1 ? "s" : ""} ago`;
  }

  const hours = Math.floor(diff / (1000 * 60 * 60));
  if (hours > 0) {
    return `${hours} hour${hours > 1 ? "s" : ""} ago`;
  }

  const minutes = Math.floor(diff / (1000 * 60));
  return `${minutes} minute${minutes > 1 ? "s" : ""} ago`;
}
