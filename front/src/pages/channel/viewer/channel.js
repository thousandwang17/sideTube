/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-14 17:50:53
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-31 17:02:01
 * @FilePath: /sidetube/src/channel/videoUpload/upload.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import React, { useState, useEffect } from "react";
import { Box } from "@mui/system";

import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
import VideoList from "./videoList";
import ColorAvatar from "component/avatar";
import { useParams } from "react-router-dom";
import { userInfoAxios } from "./axios";
import { AuthContext } from "common/context/authcontext";
import NotLogIn from "pages/notLogin";
import { useLocation } from "react-router-dom";
import { getUserID } from "common/jwt";

const uAxios = userInfoAxios();
function TabPanel(props) {
  const { children, value, index, ...other } = props;
  let { user_id } = useParams();

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`simple-tabpanel-${index}`}
      aria-labelledby={`simple-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ pt: 3 }}>
          <div>{children}</div>
        </Box>
      )}
    </div>
  );
}

function a11yProps(index) {
  return {
    id: `simple-tab-${index}`,
    "aria-controls": `simple-tabpanel-${index}`,
  };
}

export default function Channel() {
  const [value, setValue] = React.useState(0);
  const [videoCount, setVideoCount] = React.useState(0);
  const [data, setData] = React.useState({
    user_name: null,
    videos: 0,
    user_id: null,
  });
  const { isUserLoggedIn } = React.useContext(AuthContext);
  const location = useLocation();
  const isOwner = location.pathname.includes("/myChannel");
  const [showPage, setShowPage] = React.useState(
    isOwner && !isUserLoggedIn ? false : true
  );

  let { user_id } = useParams();

  const fetchUserInfo = () => {
    async function fetch() {
      const uid = isOwner ? getUserID() : user_id;
      try {
        uAxios
          .post("/info", {
            user_id: uid,
          })
          .then((resp) => {
            if (resp?.data == null) {
              throw new Error("user info data is missing");
            }
            setData({
              user_name: resp?.data?.user_name ?? null,
              videos: resp?.data?.videos ?? 0,
              user_id: resp?.data?.user_id ?? null,
            });
          });
      } catch (e) {
        console.log(e);
      }
    }
    fetch();
  };

  const handleChange = (event, newValue) => {
    setValue(newValue);
  };

  useEffect(() => {
    if (isOwner && !isUserLoggedIn) {
      setShowPage(false);
      return;
    }
    setShowPage(true);
  }, [isOwner, isUserLoggedIn]);

  const firstLoading = React.useRef(true);
  const firstRefreshLoading = React.useRef(true);

  // when id change ,  set firstLoading to ture to get new video data
  // used for refresh page and fetch new data
  useEffect(() => {
    if (firstRefreshLoading.current) {
      return;
    }
    window.scrollTo(0, 0);
    firstLoading.current = true;
  }, [user_id]);

  React.useEffect(() => {
    if (!firstLoading.current) {
      //  unlock `firstRefreshLoading` when Updating of life cycle
      firstRefreshLoading.current = false;
      return;
    }
    firstLoading.current = false;
    fetchUserInfo();
  }, [user_id]);

  return (
    <Box>
      {showPage && (
        <Box key={data.user_id}>
          <Box sx={{ display: "flex", alignItems: "center" }}>
            <Box sx={{ flexGrow: 1, display: "flex", mb: 3 }}>
              <ColorAvatar
                userName={data.user_name}
                sx={{
                  width: "80px",
                  height: "80px",
                  fontSize: "40px",
                  display: "flex",
                }}
              />
              <Box sx={{ display: "flex", ml: 2, flexDirection: "column" }}>
                <Box sx={{ fontSize: "24px" }}>{data.user_name ?? " "}</Box>
                <Box sx={{ fontSize: "13px" }}>
                  {videoCount > 0 ? videoCount + ` videos` : "no videos"}
                </Box>
              </Box>
            </Box>
            <Box></Box>
          </Box>

          <Box sx={{ borderBottom: 1, borderColor: "divider" }}>
            <Tabs
              value={value}
              onChange={handleChange}
              aria-label="basic tabs example"
            >
              <Tab label="videos" {...a11yProps(0)} />
              {/* <Tab label="play list" {...a11yProps(1)} /> */}
            </Tabs>
          </Box>

          <TabPanel value={value} index={0}>
            <VideoList
              userId={isOwner ? getUserID() : user_id}
              setVideoCount={setVideoCount}
            />
          </TabPanel>
        </Box>
      )}
      {!showPage && <NotLogIn />}
    </Box>
  );
}
