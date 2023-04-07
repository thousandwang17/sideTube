/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-12 20:50:25
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-07 12:53:18
 * @FilePath: /sidetube/src/compmnent/Avatar.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

import Avatar from "@mui/material/Avatar";
import PropTypes from "prop-types";

function stringToColor(string) {
  let hash = 0;
  let i;

  /* eslint-disable no-bitwise */
  for (i = 0; i < string.length; i += 1) {
    hash = string.charCodeAt(i) + ((hash << 5) - hash);
  }

  let colour = "#";

  for (i = 0; i < 3; i += 1) {
    const value = (hash >> (i * 8)) & 0xff;
    colour += `00${value.toString(16)}`.substr(-2);
  }
  /* eslint-enable no-bitwise */

  return colour;
}

const ColorAvatar = (props) => {
  return (
    <Avatar
      sx={{ bgcolor: stringToColor(props.userName ?? "user"), ...props.sx }}
    >
      {props.userName?.charAt(0)?.toUpperCase()}
    </Avatar>
  );
};

ColorAvatar.propTypes = {
  userName: PropTypes.string.isRequired,
};

export default ColorAvatar;
