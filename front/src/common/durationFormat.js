/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-03-09 14:25:21
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-09 14:26:21
 * @FilePath: /sidetube/src/common/durationFormat.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
export function durationFormat(duration) {
  const sec = parseInt(duration, 10);
  if (isNaN(sec)) {
    return "00:00";
  }

  const hours = parseInt(sec / 3600, 10);
  let mins = parseInt((sec % 3600) / 60);
  let secs = parseInt((sec % 3600) % 60);

  if (mins < 10) {
    mins = `0` + mins;
  }

  if (hours >= 1) {
    return hours + `:` + mins + `:` + secs;
  }

  if (secs < 10) {
    secs = `0` + secs;
  }
  return mins + `:` + secs;
}
