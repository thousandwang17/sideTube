/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-09 15:16:39
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-04-04 20:29:18
 * @FilePath: /sidetube/src/pages/video/message.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import * as React from "react";
import List from "@mui/material/List";
import CircularProgress from "@mui/material/CircularProgress";
import { Box } from "@mui/system";
import { messageAxios } from "./axios";
import { getUserName, getUserID } from "common/jwt";
import InfiniteScroll from "react-infinite-scroll-component";
import MessageFiled from "./component/messageInput";
import Message from "./component/message";
import { Stack } from "@mui/material";
import PropTypes from "prop-types";

const MAxios = messageAxios();
const pageSize = 10;

export default function MessageItemsList({
  videoID,
  messageLength,
  incMessageLength,
}) {
  const [messages, setMessages] = React.useState([]);
  const [skip, setSkip] = React.useState(0);

  // after adding new data, it will cause getting replica message, because skip is not increase correctly
  // To avoid that, sholud count of new Message to get current skip number (skip += newMessageCount )
  const [newMessageCount, setNewMessageCount] = React.useState(0);

  // this function will send to the "MessageFiled" component for invoke,
  // it  will store a new message by API and refresh the dom
  const handleSubmitOnClick = async (message) => {
    await MAxios.post("/create", {
      message: message,
      video_id: videoID,
    }).then((resp) => {
      resp = resp?.data;
      if (resp?.MessageId == null) {
        throw new Error("message Id is missing");
      }
      const timeZoneOffset = new Date().getTimezoneOffset();

      setMessages([
        {
          id: resp.MessageId,
          user_id: getUserID(),
          message: message,
          user_name: getUserName(),
          message_time: Date.now() + timeZoneOffset * 60 * 1000,
          reply_count: 0,
          update_time: null,
        },
        ...messages,
      ]);

      incMessageLength();
      setNewMessageCount((pre) => pre + 1);
    });
  };

  const handleIncMessageReply = (message_id) => {
    setMessages((prevMessageMeta) => {
      const message = prevMessageMeta.find((m) => m.id === message_id);
      if (message) {
        message.reply_count += 1;
      }
      return [...prevMessageMeta];
    });
  };

  const handleEditMessage = (message_id, new_message) => {
    setMessages((prevMessageMeta) => {
      const message = prevMessageMeta.find((m) => m.id === message_id);
      if (message) {
        message.message = new_message;
        const date = new Date(Date.now());
        const formattedDate = date.toLocaleString("default", {
          year: "numeric",
          month: "2-digit",
          day: "2-digit",
          hour: "2-digit",
          minute: "2-digit",
          second: "2-digit",
        });
        message.update_time = formattedDate;
      }
      return [...prevMessageMeta];
    });
  };

  // This function retrieves messages from an API by specifying
  // the range of messages to be obtained using the parameters "skip" and "limit"
  const fetchMessageData = async () => {
    return await MAxios.post("/list", {
      video_id: videoID,
      skip: skip,
      limit: pageSize,
    }).then((resp) => {
      if (resp?.data?.list == null) {
        throw new Error("message Id is missing");
      }

      const fetch_messages = resp.data.list.map((d) => {
        if (d.id && d.user_id && d.user_name && d.message && d.create_time) {
          return {
            id: d.id,
            user_id: d.user_id,
            user_name: d.user_name,
            message: d.message,
            message_time: d.create_time,
            update_time:
              d.update_time === "0001-01-01 00:00:00" ? null : d.update_time,
            reply_count: d.replies ?? 0,
          };
        }
        throw new Error("data missing required field");
      });

      return fetch_messages;
    });
  };

  // invoke by InfiniteScroll
  const fetchMoreData = () => {
    setSkip((pre) => pre + newMessageCount + pageSize);
    setNewMessageCount(0);
  };

  React.useEffect(() => {
    setSkip(0);
  }, [videoID]);

  // when 'skip' be changed , it will get more older message and append to the list.
  React.useEffect(() => {
    let ignore = false;
    fetchMessageData().then((res) => {
      if (!ignore) {
        setMessages((pre) => [...pre, ...res]);
      }
    });
    return () => {
      ignore = true;
    };
  }, [skip]);

  return (
    <>
      <Box sx={{ pl: 1, pt: 4 }}>
        {messageLength === 0 ? "no message" : messageLength + "  messages"}
      </Box>
      <MessageFiled onSubmit={handleSubmitOnClick} />
      <List sx={{ width: "100%", bgcolor: "background.paper" }}>
        <InfiniteScroll
          dataLength={messages.length}
          next={fetchMoreData}
          hasMore={skip + pageSize <= messageLength}
          loader={
            <Stack alignItems="center">
              <CircularProgress size={30} />
            </Stack>
          }
        >
          {messages.map((d) => (
            <Message
              meta={d}
              key={d.id}
              onReply={handleIncMessageReply}
              onEditMessage={handleEditMessage}
            />
          ))}
        </InfiniteScroll>
        {messageLength === 0 && (
          <Box
            sx={{
              display: "flex",
              justifyContent: "center",
              fontSize: "18px",
              color: "#666",
            }}
          >
            No message yet
          </Box>
        )}
      </List>
    </>
  );
}

MessageItemsList.propTypes = {
  videoID: PropTypes.string.isRequired,
};
