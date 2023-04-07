/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-10 13:37:29
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-25 20:40:25
 * @FilePath: /sidetube/src/pages/video/rely.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import * as React from "react";
import ListItem from "@mui/material/ListItem";
import ListItemText from "@mui/material/ListItemText";
import ListItemAvatar from "@mui/material/ListItemAvatar";
import ListItemSecondaryAction from "@mui/material/ListItemSecondaryAction";
import { getUserID } from "common/jwt";
import ColorAvatar from "component/avatar.js";
import PropTypes from "prop-types";
import MoreVertIcon from "@mui/icons-material/MoreVert";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import IconButton from "@mui/material/IconButton";
import MessageFiled from "./messageInput";
import { replyAxios } from "../axios";
import timeAgo from "common/timeAgo";

const rAxios = replyAxios();

export default function Relpy({ meta, onEditMessage }) {
  const [isHovered, setIsHovered] = React.useState(false);
  const [anchorEl, setAnchorEl] = React.useState(null);
  const [editMode, setEditMode] = React.useState(false);
  const openMenu = Boolean(anchorEl);
  const [updateLoading, setUpdateLoading] = React.useState(false);

  // Transfer to edit mode
  const handleEditOnClick = () => {
    setEditMode(true);
    handleMenuClose();
  };

  const handleMenuClick = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const handleCancleOnClick = () => {
    setEditMode(false);
  };

  const handleUpdate = async (new_message) => {
    if (new_message === "") {
      console.log("message is empty");
      return;
    }

    setUpdateLoading(true);
    try {
      await rAxios
        .post("/edit", {
          reply_id: meta.id,
          Message: new_message,
        })
        .then((resp) => {
          if (typeof onEditMessage === "function") {
            onEditMessage(meta.id, new_message);
          }
        });
      setEditMode(false);
    } catch (e) {
      console.error(e);
    } finally {
      setUpdateLoading(false);
    }
  };

  return editMode ? (
    // edit mode
    <MessageFiled
      onSubmit={handleUpdate}
      onCancel={handleCancleOnClick}
      defaultMessage={meta.message}
      mini
    />
  ) : (
    // watch mode
    <ListItem
      alignItems="flex-start"
      key={meta.id}
      sx={{ p: 0, pb: 1 }}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      <ListItemAvatar sx={{ minWidth: 45 }}>
        <ColorAvatar userName={meta.user_name} sx={{ width: 30, height: 30 }} />
      </ListItemAvatar>
      <ListItemText
        primary={
          <React.Fragment>
            <span> {meta.user_name} </span>
            <span style={{ color: "#aaa", paddingLeft: 2, fontSize: 13 }}>
              {timeAgo(meta.message_time)} {meta?.update_time && " (edited) "}
            </span>
          </React.Fragment>
        }
        secondary={
          <React.Fragment>
            <span>{meta.message}</span>
          </React.Fragment>
        }
      />
      {getUserID() === meta.user_id && isHovered && !editMode && (
        <ListItemSecondaryAction sx={{ top: 5, right: 50, transform: "unset" }}>
          <IconButton
            edge="end"
            aria-label="more"
            id={"rely_more_vertIcon_" + meta.id}
            aria-controls={openMenu ? "rely_menu_" + meta.id : undefined}
            aria-haspopup="true"
            aria-expanded={openMenu ? "true" : undefined}
            onClick={handleMenuClick}
          >
            <MoreVertIcon />
          </IconButton>
          <Menu
            id={"rely_menu_" + meta.id}
            aria-labelledby={"rely_more_vertIcon_" + meta.id}
            anchorEl={anchorEl}
            open={openMenu}
            onClose={handleMenuClose}
            anchorOrigin={{
              vertical: "top",
              horizontal: "left",
            }}
            transformOrigin={{
              vertical: "top",
              horizontal: "left",
            }}
          >
            <MenuItem onClick={handleEditOnClick}>edit</MenuItem>
          </Menu>
        </ListItemSecondaryAction>
      )}
    </ListItem>
  );
}

Relpy.propTypes = {
  meta: PropTypes.shape({
    message_id: PropTypes.string.isRequired,
    id: PropTypes.string.isRequired,
    user_name: PropTypes.string.isRequired,
    message: PropTypes.string.isRequired,
    message_time: PropTypes.string.isRequired,
  }).isRequired,
};
