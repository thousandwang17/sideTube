/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-29 16:04:42
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-27 16:21:06
 * @FilePath: /sidetube/src/pages/channel/owner/editVideoMeta.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import React, { useEffect } from "react";
import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import Dialog from "@mui/material/Dialog";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import AppBar from "@mui/material/AppBar";
import Toolbar from "@mui/material/Toolbar";
import IconButton from "@mui/material/IconButton";
import Typography from "@mui/material/Typography";
import CloseIcon from "@mui/icons-material/Close";
import Slide from "@mui/material/Slide";
import { videMetaUpdateAxios } from "./axios";
import ImageWithFallback from "component/imageWithfallback";
import { maxWidth } from "@mui/system";
import getHost from "common/axios";

const VMUAxios = videMetaUpdateAxios();

const Transition = React.forwardRef(function Transition(props, ref) {
  return <Slide direction="up" ref={ref} {...props} />;
});

export default function FormDialog({ data }) {
  const [open, setOpen] = React.useState(false);

  const [videoId, setVideoId] = React.useState(null);
  const [title, setTitle] = React.useState("");
  const [desc, setDesc] = React.useState("");
  const [file, setFile] = React.useState(null);
  const [imgSrc, setImageSrc] = React.useState(null);

  console.log(data?.row?.png);
  const handleFileChange = (e) => {
    if (e.target.files) {
      setFile(e.target.files[0]);

      const reader = new FileReader();
      reader.readAsDataURL(e.target.files[0]);
      reader.onload = () => {
        console.log(123);
        setImageSrc(reader.result);
      };
    }
  };

  const handleTitleOnChange = (e) => {
    setTitle(e.target.value);
  };

  const handleDescOnChange = (e) => {
    setDesc(e.target.value);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const handleSave = (e) => {
    e.preventDefault();
    console.log(videoId);
    if (!videoId) {
      return;
    }
    (async () => {
      try {
        const formData = new FormData();

        formData.append("video_id", videoId);
        formData.append("title", title);
        formData.append("desc", desc);
        if (file !== null) {
          formData.append("png", file);
        }

        const response = await VMUAxios.post("/setInfo", formData, {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        });

        data.row.refreshList();
      } catch (e) {
        console.log(e);
        return;
      }
      setOpen(false);
    })();
  };

  useEffect(() => {
    if (data?.row) {
      setVideoId(data.row.video_id);
      setTitle(data.row.title);
      setDesc(data.row.desc);
      setFile(null);
      setOpen(true);
      setImageSrc(getHost() + `/picture/video/` + data?.row?.png);
    }
  }, [data]);

  return (
    <div>
      <Dialog
        fullScreen
        open={open}
        onClose={handleClose}
        TransitionComponent={Transition}
      >
        <form onSubmit={handleSave}>
          <AppBar sx={{ position: "relative" }}>
            <Toolbar>
              <IconButton
                edge="start"
                color="inherit"
                onClick={handleClose}
                aria-label="close"
              >
                <CloseIcon />
              </IconButton>
              <Typography sx={{ ml: 2, flex: 1 }} variant="h6" component="div">
                Sound
              </Typography>
              <Button autoFocus color="inherit" onClick={handleClose}>
                cancel
              </Button>
              <Button type="submit" color="inherit">
                save
              </Button>
            </Toolbar>
          </AppBar>
          <DialogTitle>video</DialogTitle>

          <DialogContent>
            <DialogContentText>Title*</DialogContentText>
            <TextField
              autoFocus
              margin="dense"
              id="title"
              type="text"
              fullWidth
              variant="standard"
              required
              defaultValue={data?.row?.title || ""}
              onChange={handleTitleOnChange}
              style={{ marginBottom: 30 }}
            />
            <DialogContentText>Description*</DialogContentText>
            <TextField
              placeholder="..."
              autoFocus
              margin="dense"
              id="description"
              required
              multiline
              rows={10}
              fullWidth
              defaultValue={data?.row?.desc || ""}
              onChange={handleDescOnChange}
              style={{ marginBottom: 30 }}
            />

            <DialogContentText>showing</DialogContentText>
            <input
              accept="image/jpg,image/png"
              style={{ display: "none" }}
              id="upload-showing"
              onChange={handleFileChange}
              type="file"
            />
            <label htmlFor="upload-showing">
              <ImageWithFallback
                style={{
                  maxWidth: "500px",
                  minWidth: "240px",
                  cursor: "pointer",
                }}
                primarySrc={imgSrc}
                fallbackSrc={
                  "https://images.unsplash.com/photo-1512917774080-9991f1c4c750?auto=format&w=350&dpr=2"
                }
              />
            </label>
          </DialogContent>
        </form>
      </Dialog>
    </div>
  );
}
