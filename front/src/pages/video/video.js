/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-31 20:12:44
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-04-05 17:16:21
 * @FilePath: /sidetube/src/pages/video/video.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import React, { useState, useEffect } from "react";
import { videoAxios, videoMpdUrl } from "./axios";
import { useParams } from "react-router-dom";
import "shaka-player/dist/controls.css";
import MessageItemsList from "./messages";
import Grid from "@mui/material/Grid";
import { Box } from "@mui/system";
import Typography from "@mui/material/Typography";
import Accordion from "@mui/material/Accordion";
import AccordionSummary from "@mui/material/AccordionSummary";
import AccordionDetails from "@mui/material/AccordionDetails";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import RelationVideoRecommend from "./component/recommend/recomend";
import ColorAvatar from "component/avatar";
import { AuthContext } from "common/context/authcontext";
import { styled } from "@mui/material/styles";
import { Link } from "react-router-dom";
import { useMediaQuery } from "@mui/material";

const shaka = require("shaka-player/dist/shaka-player.ui.js");

const ErrorVideoDiv = styled("div")`
  position: absolute;
  top: 0;
  bottom: 0;
  right: 0;
  left: 0;
  justify-content: center;
  display: flex;
  align-items: center;
  font-size: 21px;
  color: white;
  z-index: -1;
`;

const VAxios = videoAxios();
export default function VideoPlayer() {
  // video id from react-router-dom.useParams
  let { id } = useParams();
  const [watchedFor15Seconds, setWatchedFor15Seconds] = useState(false);
  const { isUserLoggedIn } = React.useContext(AuthContext);

  const [error, setError] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");
  const [player, setPlayer] = useState(null);
  const [videoMeta, setVideoMeta] = useState({
    title: null,
    describe: null,
    user_id: null,
    user_name: null,
    publish_time: null,
    messages: 0,
    views: 0,
  });

  // use for check first loading
  const firstLoading = React.useRef(true);
  const firstRefreshLoading = React.useRef(true);
  const videoRef = React.useRef(null);
  const playerRef = React.useRef(null);

  // when id change ,  set firstLoading to ture to get new video data
  // used for refresh page and fetch new data
  useEffect(() => {
    if (firstRefreshLoading.current) {
      return;
    }
    window.scrollTo(0, 0);
    firstLoading.current = true;
    setWatchedFor15Seconds(false);
    setError(false);
  }, [id]);

  useEffect(() => {
    let ignore = false;
    async function fetchListData() {
      const resp = await VAxios.post("/video", {
        video_id: id,
      }).catch((error) => {
        // Handle errors
        if (error?.response?.status) {
          const status = error.response.status;
          switch (status) {
            case 400:
              setErrorMessage("the video is not exists");
              break;
            case 403:
              setErrorMessage("the video is private");
              break;
            default:
              setErrorMessage("Temporary error, try again later");
              break;
          }
        } else {
          setErrorMessage("Temporary error, try again later");
        }
        setError(true);
      });
      return resp.data;
    }

    try {
      fetchListData().then((resp) => {
        if (!ignore) {
          firstRefreshLoading.current = false;
          initApp(videoMpdUrl(resp?.mpd));
          const timeZoneOffset = new Date().getTimezoneOffset();
          const uploadTime =
            new Date(resp.uploadTime).getTime() + timeZoneOffset * 60 * 1000;
          setVideoMeta({
            title: resp.title,
            describe: resp.desc,
            user_id: resp.user_id,
            user_name: resp.user_name,
            publish_time: new Date(uploadTime).toLocaleString("en-US", {
              timeZone: "UTC",
            }),
            messages: resp.messages ?? 0,
            views: resp.views ?? 0,
          });
        }
      });
    } catch (e) {
      console.error("Error code", e);
    }

    return () => {
      ignore = true;
    };
  }, [id]);

  useEffect(() => {
    const p = playerRef.current;
    if (p !== player) {
      if (p !== null) {
        p.unload();
        closePiP();
      }

      playerRef.current = player;
    }

    return () => {
      if (playerRef.current !== null) {
        const p = playerRef.current;
        p.unload();
        closePiP();
      }
    };
  }, [player]);

  useEffect(() => {
    if (!isUserLoggedIn) {
      return;
    }
    const video = videoRef.current;

    const handleTimeUpdate = () => {
      if (video.currentTime >= 15 && !watchedFor15Seconds) {
        setWatchedFor15Seconds(true);
        try {
          incViews();
        } catch (e) {
          console.error("Error code", e);
        }
      }
    };

    async function incViews() {
      const resp = await VAxios.post("/video/incViews", {
        video_id: id,
      });
    }

    video.addEventListener("timeupdate", handleTimeUpdate);

    return () => {
      video.removeEventListener("timeupdate", handleTimeUpdate);
    };
  }, [watchedFor15Seconds]);

  const closePiP = () => {
    if (document.pictureInPictureElement) {
      document
        .exitPictureInPicture()
        .then(() => {
          console.log("Exited PiP mode.");
        })
        .catch((error) => {
          console.error("Error exiting PiP mode:", error);
        });
    }
  };

  const initApp = async (mpaUrl) => {
    if (mpaUrl == null) {
      throw new Error("invalid mpd url");
    }

    // Install built-in polyfills to patch browser incompatibilities.
    shaka.polyfill.installAll();

    // Check to see if the browser supports the basic APIs Shaka needs.
    if (shaka.Player.isBrowserSupported()) {
      // Everything looks good!
      initPlayer(mpaUrl);
    } else {
      // This browser does not have the minimum set of APIs we need.
      console.error("Browser not supported!");
    }
  };

  const incMessageLength = () => {
    setVideoMeta((prevVideoMeta) => ({
      ...prevVideoMeta,
      messages: prevVideoMeta.messages + 1,
    }));
  };

  async function initPlayer(mpaUrl) {
    // Create a Player instance.
    const video = document.getElementById("MainVideoPlayer");

    const player = new shaka.Player(video);
    const videoContainer = document.getElementById("video-container");
    const ui = new shaka.ui.Overlay(player, videoContainer, video);
    const uiConfig = {};

    setPlayer(player);

    // Configuring elements to be displayed on video player control panel
    uiConfig["overflowMenuButtons"] = [
      "quality",
      "picture_in_picture",
      "playback_rate",
      "airplay",
      "cast",
    ];
    ui.configure(uiConfig);

    // Listen for error events.
    player.addEventListener("error", onErrorEvent);

    // Try to load a manifest.
    // This is an asynchronous process.

    await player.load(mpaUrl);
    // This runs if the asynchronous load is successful.
    console.log("The video has now been loaded!");
  }

  const onErrorEvent = (event) => {
    // Extract the shaka.util.Error object from the event.
    onError(event.detail);
  };

  const onError = (error) => {
    // Log the error.
    console.error("Error code", error.code, "object", error);
  };

  return (
    <Grid container mt={0.2} key={id}>
      <Grid item xs={12} sm={12} md={8} px={3} sx={{ px: { xs: 0, md: 3 } }}>
        <Box
          id="video-container"
          sx={{
            position: "relative",
            backgroundColor: "#555",
            mx: { xs: -2, md: 0 },
          }}
        >
          <ErrorVideoDiv sx={{ zIndex: 1, display: error ? "flex" : "none" }}>
            {errorMessage}
          </ErrorVideoDiv>
          <video id="MainVideoPlayer" width="100%" autoPlay ref={videoRef} />
        </Box>

        {videoMeta.title && !error && (
          <Typography
            variant="h6"
            sx={{
              maxHeight: "4rem",
              overflow: "hidden",
              textOverflow: "ellipsis",
              WebkitLineClamp: 2,
              WebkitBoxOrient: "vertical",
              display: "-webkit-box",
            }}
          >
            {videoMeta.title}
          </Typography>
        )}

        {!error && (
          <Box>
            <Link
              to={"/channel/" + videoMeta.user_id}
              style={{ textDecoration: "none", color: "inherit" }}
            >
              <Box sx={{ display: "flex", alignItems: "center" }}>
                <ColorAvatar userName={videoMeta.user_name ?? ""} />
                <Box sx={{ ml: 1.5, fontSize: 18, my: 2 }}>
                  {videoMeta.user_name ?? ""}
                </Box>
              </Box>
            </Link>

            <Box sx={{ mt: 1, mb: 3 }}>
              <Accordion>
                <AccordionSummary
                  expandIcon={videoMeta.describe ? <ExpandMoreIcon /> : ""}
                  aria-controls="panel1bh-content"
                  id="panel1bh-header"
                >
                  <Typography sx={{ width: "33%", flexShrink: 0 }}>
                    views : {videoMeta.views ?? 0}
                  </Typography>
                  <Typography sx={{ color: "text.secondary" }}>
                    publish : {videoMeta.publish_time}
                  </Typography>
                </AccordionSummary>
                {videoMeta.describe && (
                  <AccordionDetails>
                    <Typography sx={{ whiteSpace: "break-spaces" }}>
                      {videoMeta.describe}
                    </Typography>
                  </AccordionDetails>
                )}
              </Accordion>
            </Box>
            <MessageItemsList
              videoID={id}
              messageLength={videoMeta.messages}
              incMessageLength={incMessageLength}
            />
          </Box>
        )}
      </Grid>
      <Grid item xs={12} sm={12} md={4} pl={0}>
        <RelationVideoRecommend videoID={id} />
      </Grid>
    </Grid>
  );
}
