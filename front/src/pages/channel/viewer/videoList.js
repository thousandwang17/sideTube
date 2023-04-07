/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-03-03 14:06:08
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-04-04 19:43:26
 * @FilePath: /sidetube/src/pages/video/component/recomend.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

import * as React from "react";
import InfiniteScroll from "react-infinite-scroll-component";
import { channelAxios } from "./axios";
import ExhibitVideo from "component/exhibitVideo";
import Grid from "@mui/material/Grid";
import { Box } from "@mui/system";

const pageSize = 10;

const cAxios = channelAxios();
export default function VideoList({ userId, setVideoCount }) {
  const [data, setData] = React.useState([]);
  const [skip, setSkip] = React.useState(0);
  const [videoListLength, setVideoListLength] = React.useState(30);

  const fetchChannelVideoData = async () => {
    return await cAxios
      .post("/publicList", {
        user_id: userId,
        skip: skip,
        limit: pageSize,
      })
      .then((resp) => {
        if (resp?.data?.list == null) {
          throw new Error("channel data is missing");
        }

        if (resp.data.list.length === 0) {
          setVideoListLength(data.length);
        }

        if (skip === 0 && resp?.data?.count) {
          setVideoCount(resp.data.count);
        }

        const fetch_datas = resp.data.list.map((d) => {
          if (
            d.user_id &&
            d.user_name &&
            d.duration &&
            d.png &&
            d.video_id &&
            // d.title &&
            d.uploadTime
          ) {
            return {
              user_id: d.user_id,
              user_name: d.user_name,
              duration: d.duration,
              png: d.png,
              video_id: d.video_id,
              title: d.title ?? "",
              views: d.views ?? 0,
              createTime: d.uploadTime,
            };
          }
          throw new Error("data missing required field");
        });

        return fetch_datas;
      });
  };

  // invoke by InfiniteScroll
  const fetchMoreData = () => {
    setSkip((pre) => pre + pageSize);
  };

  // when 'skip' be changed , it will get more  data and append to the list.
  React.useEffect(() => {
    let ignore = false;

    try {
      fetchChannelVideoData().then((resp) => {
        if (!ignore) {
          setData((pre) => [...pre, ...resp]);
        }
      });
    } catch (e) {
      console.log(e);
    }

    return () => {
      ignore = true;
    };
  }, [skip]);

  return (
    <>
      <InfiniteScroll
        dataLength={data.length}
        next={fetchMoreData}
        hasMore={skip + pageSize <= videoListLength}
      >
        <Grid container spacing={2}>
          {data.map((d) => (
            <Grid item xs={12} sm={4} md={3} pl={0} key={d.video_id}>
              <ExhibitVideo meta={d} />
            </Grid>
          ))}
        </Grid>
      </InfiniteScroll>

      {videoListLength === 0 && (
        <Box
          sx={{
            display: "flex",
            justifyContent: "center",
            fontSize: "18px",
            color: "#666",
          }}
        >
          <p>This channel has not uploaded a video yet</p>
        </Box>
      )}
    </>
  );
}
