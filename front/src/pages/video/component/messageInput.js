/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-12 20:13:55
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-24 14:16:06
 * @FilePath: /sidetube/src/pages/video/component/message.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import * as React from "react";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemText from "@mui/material/ListItemText";
import ListItemAvatar from "@mui/material/ListItemAvatar";
import LoadingButton from "@mui/lab/LoadingButton";
import TextField from "@mui/material/TextField";
import { Button } from "@mui/material";
import AccountCircle from "@mui/icons-material/AccountCircle";
import { AuthContext } from "common/context/authcontext";

// Message input filed for message and rely
export default function MessageFiled({
  onSubmit,
  sx,
  onCancel,
  mini,
  defaultMessage,
}) {
  const [message, setMessage] = React.useState(defaultMessage ?? "");
  const [showButtons, setShowButtons] = React.useState(false);
  const [submitLoading, setSubmitLoading] = React.useState(false);
  const { isUserLoggedIn } = React.useContext(AuthContext);

  const handleMessageOnChange = (e) => {
    setMessage(e.target.value);
  };

  const handleCancleOnClick = () => {
    setMessage("");
    setShowButtons(false);
    if (typeof onCancel === "function") {
      onCancel();
    }
  };

  const handleInputOnFocus = () => {
    setShowButtons(true);
  };

  const handleSubmitOnClick = async () => {
    if (message === "") {
      console.log("message is empty");
      return;
    }

    setSubmitLoading(true);
    try {
      await onSubmit(message).then(() => setMessage(""));
    } catch (e) {
      console.error(e);
    } finally {
      setSubmitLoading(false);
    }
  };

  return (
    <List
      sx={{
        width: "100%",
        p: 0,
        mb: mini ? "0" : "16px",
        display: "flex",
        alignItems: "flex-end",
        bgcolor: "background.paper",
        ...sx,
      }}
    >
      <ListItem
        sx={{ pl: 0, pb: "1px", pt: mini ? 0 : 1, alignItems: "flex-start" }}
      >
        <ListItemAvatar sx={{ minWidth: mini ? 35 : 56 }}>
          <AccountCircle
            sx={{
              width: mini ? 30 : 40,
              height: mini ? 30 : 40,
              color: "primary.main",
              marginTop: "15px",
            }}
          />
        </ListItemAvatar>
        <ListItemText
          primary={
            <TextField
              variant="standard"
              sx={{ width: "100%", p: 1 }}
              multiline
              onChange={handleMessageOnChange}
              placeholder={
                isUserLoggedIn ? "type some message..." : "please login first"
              }
              disabled={!isUserLoggedIn}
              value={message}
              onFocus={handleInputOnFocus}
            />
          }
          secondary={
            <React.Fragment>
              <span
                style={{
                  width: "100%",
                  alignItems: "flex-end",
                  justifyContent: "end",
                  display: showButtons ? "flex" : "none",
                }}
              >
                <Button
                  size="small"
                  sx={{ float: "right" }}
                  onClick={handleCancleOnClick}
                >
                  cancel
                </Button>
                <LoadingButton
                  variant="contained"
                  size="small"
                  loading={submitLoading}
                  disabled={message === ""}
                  onClick={handleSubmitOnClick}
                  sx={{ float: "right", ml: 1 }}
                >
                  Send
                </LoadingButton>
              </span>
            </React.Fragment>
          }
        />
      </ListItem>
    </List>
  );
}
