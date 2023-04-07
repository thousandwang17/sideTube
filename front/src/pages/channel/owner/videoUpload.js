/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-14 17:50:53
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-21 17:25:36
 * @FilePath: /sidetube/src/channel/videoUpload/upload.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import React, { useState, useEffect } from "react";
import { Button } from "@mui/material";
import { videoUploadAxios } from "./axios";
import { Box } from "@mui/system";
import CircularProgress from "@mui/material/CircularProgress";

// Define the chunk size (e.g. 5MB)
const CHUNK_SIZE = 5 * 1024 * 1024;
const VUAxios = videoUploadAxios();

export default function VideoUpload({
  setProgress,
  setUploadFileName,
  setRefreshList,
}) {
  const [file, setFile] = useState();
  const [loading, setLoading] = useState(false);

  const handleFileChange = (e) => {
    if (e.target.files) {
      setFile(e.target.files[0]);
    }
  };

  useEffect(() => {
    handleUpload();
  }, [file]);

  // Before uploading the video, must invoke this function to inform the API server to prepare
  const startUpload = async () => {
    if (!file) {
      throw new Error("file is miss");
    }
    setLoading(true);
    setUploadFileName(file.name);
    let totalChunks = Math.ceil(file.size / CHUNK_SIZE);
    const response = await VUAxios.post("/start", {
      totalChunks: totalChunks,
    });

    if (response.data.videoId == null) {
      throw new Error("video Id is missing");
    }
    return response.data;
  };

  // This function will upload video by chunk until all part has been sent
  const handleUploadByChunk = async (response) => {
    // Create a FileReader object
    const fileReader = new FileReader();

    // Define the current chunk and total number of chunks
    let currentChunk = 1;
    let totalChunks = Math.ceil(file.size / CHUNK_SIZE);

    // Handle the 'load' event
    fileReader.onload = async (e) => {
      // Get the current chunk of data
      const chunk = e.target.result;

      // Send the current chunk of data to the server
      // Create a FormData object
      const formData = new FormData();

      // Append the current chunk to the form data
      formData.append("video_id", response.videoId);
      formData.append("part_id", currentChunk);
      formData.append("streaming_data", new Blob([chunk]));
      const res = await VUAxios.post("/updatePart", formData, {
        headers: {
          "Content-Type": "multipart/form-data",
        },
        onUploadProgress: (progressEvent) => {
          console.log(
            progressEvent.loaded,
            CHUNK_SIZE * (currentChunk - 1),
            file.size
          );
          setProgress(
            Math.round(
              ((progressEvent.loaded + CHUNK_SIZE * (currentChunk - 1)) /
                file.size) *
                100
            )
          );
        },
      });

      if (currentChunk === totalChunks) {
        await finishUpload(response);
      } else {
        fileReader.readAsArrayBuffer(
          file.slice(currentChunk * CHUNK_SIZE, (currentChunk + 1) * CHUNK_SIZE)
        );
        currentChunk++;
      }
    };

    fileReader.onerror = (err) => {
      console.log(err);
    };

    fileReader.onabort = () => {
      console.log("file reading was aborted");
    };

    // Read the file as a array buffer , after readAsArrayBuffer have been loaded , will trigger "fileReader.onload" function
    fileReader.readAsArrayBuffer(file.slice(0, currentChunk * CHUNK_SIZE));
  };

  // Inform the API server that the video has been uploaded, starting to encode the video to multi-type of video
  const finishUpload = async (response) => {
    await VUAxios.post("/finish", {
      video_id: response.videoId,
    });

    setProgress(0);
    setFile(null);
    setRefreshList((prev) => prev + 1);
  };

  // The main logic that handles the video uploading
  const handleUpload = () => {
    startUpload()
      .then((response) => {
        return handleUploadByChunk(response);
      })
      .catch((err) => {
        console.log(err.message);
      })
      .finally((err) => {
        setLoading(false);
      });
  };
  return (
    <>
      <Box sx={{ position: "relative" }}>
        <input
          accept="video/mp4"
          style={{ display: "none" }}
          id="raised-button-file"
          onChange={handleFileChange}
          type="file"
        />
        <label htmlFor="raised-button-file">
          <Button variant="contained" component="span" disabled={loading}>
            Upload
          </Button>
          {loading && (
            <CircularProgress
              size={24}
              sx={{
                // color: green[500],
                position: "absolute",
                top: "50%",
                left: "50%",
                marginTop: "-12px",
                marginLeft: "-12px",
              }}
            />
          )}
        </label>
      </Box>
      {/* {progress}% */}
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
