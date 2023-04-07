/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-03-03 14:06:08
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-04-03 16:55:13
 * @FilePath: /sidetube/src/pages/video/component/recomend.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

import React, { useState, useEffect } from "react";
import List from "@mui/material/List";
import CircularProgress from "@mui/material/CircularProgress";
import InfiniteScroll from "react-infinite-scroll-component";
import recommendAxios from "./axios";
import { Stack } from "@mui/material";
import ExhibitVideo from "component/exhibitVideo";
import PropTypes from "prop-types";
import { useInfiniteQuery } from "react-query";
const pageSize = 10;

const RAxios = recommendAxios();
export default function RelationVideoRecommend({ videoID }) {
  const { data, fetchNextPage, hasNextPage, isFetching } = useInfiniteQuery(
    ["recommendations", videoID],
    async ({ pageParam = 0 }) => {
      const resp = await RAxios.post("/relation/videos", {
        video_id: videoID,
        skip: pageParam,
        limit: pageSize,
      });
      if (!resp?.data?.list) {
        throw new Error("recommend is missing");
      }

      const timeZoneOffset = new Date().getTimezoneOffset();

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
            createTime:
              new Date(d.uploadTime).getTime() + timeZoneOffset * 60 * 1000,
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

  const dataToRender = data ? data.pages.flatMap((p) => p) : [];

  const handleFetchMore = React.useCallback(() => {
    if (hasNextPage) {
      fetchNextPage();
    }
  }, [hasNextPage, fetchNextPage]);

  React.useEffect(() => {
    // Reset the query when the videoID changes
    fetchNextPage({ pageParam: 0, resetError: true });
  }, [videoID, fetchNextPage]);

  return (
    <List sx={{ width: "100%", bgcolor: "background.paper", pt: 0 }}>
      <InfiniteScroll
        dataLength={dataToRender.length}
        next={handleFetchMore}
        hasMore={hasNextPage}
        loader={
          <Stack alignItems="center">
            <CircularProgress size={30} />
          </Stack>
        }
      >
        {dataToRender.map((d) => (
          <ExhibitVideo meta={d} key={d.video_id} mini={true} />
        ))}
      </InfiniteScroll>
    </List>
  );
}

RelationVideoRecommend.propTypes = {
  videoID: PropTypes.string.isRequired,
};
