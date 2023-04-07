/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-03-03 14:06:08
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-04-05 17:02:14
 * @FilePath: /sidetube/src/pages/video/component/recomend.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

import * as React from "react";

import InfiniteScroll from "react-infinite-scroll-component";
import recommendAxios from "pages/home/axios";
import ExhibitVideo from "component/exhibitVideo";
import Grid from "@mui/material/Grid";
import { useInfiniteQuery } from "react-query";

const pageSize = 12;

const RAxios = recommendAxios();
export default function HomeVideoRecommend() {
  const { data, fetchNextPage, hasNextPage, isFetching } = useInfiniteQuery(
    ["home"],
    async ({ pageParam = 0 }) => {
      const resp = await RAxios.post("/home/videos", {
        skip: pageParam,
        limit: pageSize,
      });
      if (!resp?.data?.list) {
        throw new Error("home recommends is missing");
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

  return (
    <InfiniteScroll
      dataLength={dataToRender.length}
      next={handleFetchMore}
      hasMore={hasNextPage}
    >
      <Grid container spacing={2}>
        {dataToRender.map((d) => (
          <Grid item xs={12} sm={4} md={3} pl={0} key={d.video_id}>
            <ExhibitVideo meta={d} showAvater={true} />
          </Grid>
        ))}
      </Grid>
    </InfiniteScroll>
  );
}
