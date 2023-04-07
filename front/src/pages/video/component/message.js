/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-12 20:13:55
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-04-04 20:28:49
 * @FilePath: /sidetube/src/pages/video/component/message.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import * as React from "react";
import ListItem from "@mui/material/ListItem";
import ListItemText from "@mui/material/ListItemText";
import ListItemAvatar from "@mui/material/ListItemAvatar";
import ListItemSecondaryAction from "@mui/material/ListItemSecondaryAction";
import { Box } from "@mui/system";
import { Button } from "@mui/material";
import { ExpandLess, ExpandMore } from "@mui/icons-material";
import Relpy from "./reply";
import ColorAvatar from "component/avatar.js";
import PropTypes from "prop-types";
import MessageFiled from "./messageInput";
import { replyAxios, messageAxios } from "../axios";
import { getUserName, getUserID } from "common/jwt";
import List from "@mui/material/List";
import Collapse from "@mui/material/Collapse";
import IconButton from "@mui/material/IconButton";
import MoreVertIcon from "@mui/icons-material/MoreVert";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import timeAgo from "common/timeAgo";

const pageSize = 10;
const rAxios = replyAxios();
const mAxios = messageAxios();

// Message input filed for message and reply
export default function Message({ meta, onReply, onEditMessage }) {
  const [showReplyInput, SetShowReplyInput] = React.useState(false);
  const [replies, setReplies] = React.useState([]);
  const [openReplies, setOpenReplies] = React.useState(false);
  const [skip, setSkip] = React.useState(0);
  const [isHovered, setIsHovered] = React.useState(false);
  const [anchorEl, setAnchorEl] = React.useState(null);
  const [message, setMessage] = React.useState(meta.message);
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

  const handleMessageOnChange = (e) => {
    setMessage(e.target.value);
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
      await mAxios
        .post("/edit", {
          message_id: meta.id,
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

  // after adding new data, it will cause getting replica message, because skip is not increase correctly
  // To avoid that, sholud count of new Reply to get current skip number (skip += newReplyCount )
  const [newReplyCount, setNewReplyCount] = React.useState(0);

  // if input cancel button on Click , it will hide the inputFiled
  const handlereplyOnCancel = () => {
    SetShowReplyInput(false);
  };

  // if reply button on Click , it  will display the inputFiled for typing
  const handlereplyOnClick = () => {
    SetShowReplyInput(true);
  };

  // if ShowReplies button on Click, if first time, it will trigger to get data by api
  // , also toggle to  display Replies
  const handleShowReplies = () => {
    setOpenReplies((pre) => !pre);
  };

  // This function retrieves messages from an API by specifying
  // the range of messages to be obtained using the parameters "skip" and "limit"
  const fetchMessageData = async () => {
    try {
      return await rAxios
        .post("/list", {
          message_id: meta.id,
          skip: skip,
          limit: pageSize,
        })
        .then((resp) => {
          if (resp?.data?.list == null) {
            throw new Error("message Id is missing");
          }

          const fetch_replies = resp.data.list.map((d) => {
            if (d.id && d.user_id && d.user_name && d.message && d.time) {
              return {
                id: d.id,
                user_id: d.user_id,
                user_name: d.user_name,
                message: d.message,
                message_time: d.time,
                update_time:
                  d.update_time === "0001-01-01 00:00:00"
                    ? null
                    : d.update_time,
                message_id: meta.id,
              };
            }
            throw new Error("data missing required field");
          });

          return fetch_replies;
        });
    } catch (e) {
      console.error(e);
    }
  };

  // this function will send to the "MessageFiled" component for invoke,
  // it  will store a new message by API and refresh the dom
  const handlereplyOnSubmit = async (message) => {
    await rAxios
      .post("/create", {
        message: message,
        message_id: meta.id,
      })
      .then((resp) => {
        resp = resp?.data;
        if (resp?.ReplyId == null) {
          throw new Error("Reply Id is missing");
        }

        // to avoid show replica data , only after firstLoading that can append new message,
        // otherwise get new reply from api
        if (initListFirstLoading.current) {
          const timeZoneOffset = new Date().getTimezoneOffset();

          setReplies((pre) => [
            {
              id: resp.ReplyId,
              user_id: getUserID(),
              message: message,
              user_name: getUserName(),
              message_time: Date.now() + timeZoneOffset * 60 * 1000,
              update_time: null,
            },
            ...pre,
          ]);
        }

        SetShowReplyInput(false);
        // replyLength = replyLength +1
        setNewReplyCount((pre) => pre + 1);
        if (typeof onReply === "function") {
          onReply(meta.id);
        }
      });
  };

  const handleEditReply = (reply_id, new_message) => {
    setReplies((prevMessageMeta) => {
      const reply = prevMessageMeta.find((m) => m.id === reply_id);
      if (reply) {
        reply.message = new_message;
        const date = new Date(Date.now());
        const formattedDate = date.toLocaleString("default", {
          year: "numeric",
          month: "2-digit",
          day: "2-digit",
          hour: "2-digit",
          minute: "2-digit",
          second: "2-digit",
        });
        reply.update_time = formattedDate;
      }
      return [...prevMessageMeta];
    });
  };

  // if openReplies is true and first loading , will invode fetchMessageData to get data
  const initListFirstLoading = React.useRef(false);
  React.useEffect(() => {
    if (!openReplies) return;

    if (initListFirstLoading.current) {
      return;
    }
    fetchMessageData().then((res) => setReplies((pre) => [...pre, ...res]));
    initListFirstLoading.current = true;
  }, [openReplies]);

  // invoke by InfiniteScroll
  const fetchMoreReplies = () => {
    setSkip((pre) => pre + newReplyCount + pageSize);
    // after setting skip , reset count of New Reply
    setNewReplyCount(0);
  };

  // when 'skip' be changed , it will get more older message and append to the list.
  // const skipFirstLoading = React.useRef(true);
  React.useEffect(() => {
    if (!openReplies) return;
    fetchMessageData().then((res) => setReplies((pre) => [...pre, ...res]));
  }, [skip]);

  return (
    <Box
      sx={{ p: 0, mb: "16px" }}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {editMode ? (
        // edit mode
        <MessageFiled
          onSubmit={handleUpdate}
          onCancel={handleCancleOnClick}
          defaultMessage={meta.message}
        />
      ) : (
        // watch mode
        <ListItem alignItems="flex-start" sx={{ p: 0 }}>
          <ListItemAvatar>
            <ColorAvatar userName={meta.user_name} />
          </ListItemAvatar>
          <ListItemText
            primary={
              <React.Fragment>
                <span> {meta.user_name} </span>
                <span style={{ color: "#aaa", paddingLeft: 2, fontSize: 13 }}>
                  {timeAgo(meta.message_time)}
                  {meta?.update_time && " (edited) "}
                </span>
              </React.Fragment>
            }
            secondary={
              <React.Fragment>
                <span style={{ whiteSpace: "pre-line" }}> {meta.message} </span>
                <span style={{ display: "block" }}>
                  <Button
                    size="small"
                    onClick={handlereplyOnClick}
                    sx={{ ml: -2, zIndex: 2 }}
                  >
                    reply
                  </Button>
                </span>
              </React.Fragment>
            }
          />
          {getUserID() === meta.user_id && isHovered && !editMode && (
            <ListItemSecondaryAction
              sx={{ top: 10, right: 30, transform: "unset" }}
            >
              <IconButton
                edge="end"
                aria-label="more"
                id={"MoreVertIcon_" + meta.id}
                aria-controls={openMenu ? "menu_" + meta.id : undefined}
                aria-haspopup="true"
                aria-expanded={openMenu ? "true" : undefined}
                onClick={handleMenuClick}
              >
                <MoreVertIcon />
              </IconButton>
              <Menu
                id={"menu_" + meta.id}
                aria-labelledby={"MoreVertIcon_" + meta.id}
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
      )}
      <Box sx={{ ml: 7, width: "calc(100% - 56px)" }}>
        {showReplyInput && (
          <MessageFiled
            mini
            onSubmit={handlereplyOnSubmit}
            onCancel={handlereplyOnCancel}
            sx={{ mt: -2 }}
          />
        )}
        {meta.reply_count > 0 && (
          <Button
            startIcon={openReplies ? <ExpandLess /> : <ExpandMore />}
            sx={{ mt: -2 }}
            style={{ textTransform: "none" }}
            onClick={handleShowReplies}
          >
            {meta.reply_count + ` replies `}
          </Button>
        )}
        {meta.reply_count > 0 && (
          <Collapse in={openReplies} timeout="auto" unmountOnExit>
            <List sx={{ width: "100%", bgcolor: "background.paper", pd: 0 }}>
              {replies.map((d) => (
                <Relpy meta={d} key={d.id} onEditMessage={handleEditReply} />
              ))}

              {skip + pageSize + newReplyCount < meta.reply_count && (
                <Button
                  sx={{ textTransform: "none" }}
                  onClick={fetchMoreReplies}
                >
                  See more replies
                </Button>
              )}
            </List>
          </Collapse>
        )}
      </Box>
    </Box>
  );
}

Message.propTypes = {
  meta: PropTypes.shape({
    id: PropTypes.string.isRequired,
    user_name: PropTypes.string.isRequired,
    message: PropTypes.string.isRequired,
    message_time: PropTypes.string.isRequired,
    reply_count: PropTypes.number.isRequired,
  }).isRequired,
};
