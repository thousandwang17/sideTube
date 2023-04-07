/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-14 17:50:53
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-24 16:13:44
 * @FilePath: /sidetube/src/channel/videoUpload/upload.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import React, { useState, useEffect } from "react";
import { Box } from "@mui/system";
import LinearProgress from "@mui/material/LinearProgress";
import Typography from "@mui/material/Typography";
import Tabs from "@mui/material/Tabs";
import Tab from "@mui/material/Tab";
import VideoUpload from "./videoUpload";
import ControlledSelectionServerPaginationGrid from "./videoList";
import { AuthContext } from "common/context/authcontext";
import NotLogIn from "pages/notLogin";

function TabPanel(props) {
  const { children, value, index, ...other } = props;

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

export default function VideoOwnerList() {
  const [value, setValue] = React.useState(0);
  const [progress, setProgress] = React.useState(0);
  const [uploadFileName, setUploadFileName] = React.useState("");
  const [refreshList, setRefreshList] = React.useState(0);
  const { isUserLoggedIn } = React.useContext(AuthContext);

  const handleChange = (event, newValue) => {
    setValue(newValue);
  };
  return (
    <>
      {isUserLoggedIn && (
        <Box>
          <Box sx={{ display: "flex", alignItems: "center" }}>
            <Box sx={{ flexGrow: 1 }}>
              <h2>Channel Content</h2>
            </Box>
            <Box>
              <VideoUpload
                setProgress={setProgress}
                setUploadFileName={setUploadFileName}
                setRefreshList={setRefreshList}
              />
            </Box>
          </Box>
          <Box>
            {progress > 0 && (
              <LinearProgressWithLabel
                value={progress}
                uploadfilename={uploadFileName}
              />
            )}
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
            <ControlledSelectionServerPaginationGrid
              propRefresh={refreshList}
            />
          </TabPanel>
          {/* <TabPanel value={value} index={1} >
                Item Two
            </TabPanel> */}
        </Box>
      )}

      {!isUserLoggedIn && <NotLogIn />}
    </>
  );
}

function LinearProgressWithLabel(props) {
  return (
    <>
      <Box>{props.uploadfilename}</Box>
      <Box sx={{ display: "flex", alignItems: "center" }}>
        <Box sx={{ width: "100%", mr: 1 }}>
          <LinearProgress variant="determinate" {...props} />
        </Box>
        <Box sx={{ minWidth: 35 }}>
          <Typography variant="body2" color="text.secondary">{`${Math.round(
            props.value
          )}%`}</Typography>
        </Box>
      </Box>
    </>
  );
}

// LinearProgressWithLabel.propTypes = {
//   /**
//    * The value of the progress indicator for the determinate and buffer variants.
//    * Value between 0 and 100.
//    */
//   value: PropTypes.number.isRequired,
// };
