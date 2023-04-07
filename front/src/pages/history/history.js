/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-03-03 14:06:08
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-04-04 18:43:46
 * @FilePath: /sidetube/src/pages/video/component/recomend.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

import * as React from "react";

import InfiniteScroll from "react-infinite-scroll-component";
import historyAxios from "./axios";
import { Box } from "@mui/system";
import DayHistory from "./component/dayHistory";
import { Grid } from "@mui/material";
import { AuthContext } from "common/context/authcontext";
import NotLogIn from "pages/notLogin";
import { useInfiniteQuery } from "react-query";

const pageSize = 10;

const hAxios = historyAxios();
export default function HistoryPage() {
  const [objectData, setObjectData] = React.useState({});
  const { isUserLoggedIn } = React.useContext(AuthContext);
  const { data, fetchNextPage, hasNextPage, isFetching } = useInfiniteQuery(
    ["history", isUserLoggedIn],
    async ({ pageParam = 0 }) => {
      if (!isUserLoggedIn) {
        return [];
      }
      const resp = await hAxios.post("/history", {
        skip: pageParam,
        limit: pageSize,
      });
      if (!resp?.data?.list) {
        throw new Error(" history is missing");
      }

      const fetch_datas = resp.data.list.map((d) => {
        if (
          d.user_id &&
          d.user_name &&
          d.duration &&
          d.png &&
          d.video_id &&
          d.title &&
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
            desc: d.desc ?? "",
            historyTime: d.viewTime,
            createTime: d.uploadTime,
          };
        }
        throw new Error("data missing required field");
      });
      return fetch_datas;
    },
    {
      getNextPageParam: (lastPage) =>
        lastPage.length === pageSize ? lastPage.length : false,
      keepPreviousData: true,
      refetchOnWindowFocus: false,
      retry: false,
      cacheTime: 0,
    }
  );

  const handleFetchMore = React.useCallback(() => {
    if (hasNextPage) {
      fetchNextPage();
    }
  }, [hasNextPage, fetchNextPage]);

  const handleflatMapData = (data) => {
    const res = data.pages.flatMap((p) => p);
    return res;
  };

  const dataToRender = data ? handleflatMapData(data) : [];

  React.useEffect(() => {
    let obj = {};
    if (dataToRender == 0 && objectData != {}) {
      setObjectData({});
      return;
    }

    dataToRender.map((d) => {
      const t = new Date(d.historyTime);
      const key = t.getMonth() + 1 + "-" + (t.getDate() + 1);
      if (obj[key] === undefined) {
        obj[key] = new Array();
      }
      obj[key].push(d);
    });

    const sortedObj = Object.fromEntries(
      Object.entries(obj).sort(([keyA], [keyB]) => keyB.localeCompare(keyA))
    );

    console.log(dataToRender, sortedObj);
    setObjectData(sortedObj);
  }, [data]);

  return (
    <>
      {isUserLoggedIn && (
        <Box>
          <InfiniteScroll
            dataLength={dataToRender.length}
            next={handleFetchMore}
            hasMore={hasNextPage}
          >
            <Grid container spacing={0}>
              {Object.keys(objectData).map((key) => (
                <DayHistory key={key} data={objectData[key]} />
              ))}
            </Grid>
          </InfiniteScroll>

          {dataToRender.length === 0 && (
            <Box
              sx={{
                display: "flex",
                justifyContent: "center",
                fontSize: "18px",
                color: "#666",
              }}
            >
              <p>No history found</p>
            </Box>
          )}
        </Box>
      )}

      {!isUserLoggedIn && <NotLogIn />}
    </>
  );
}
