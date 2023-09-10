/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-03-03 14:06:08
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-06-27 10:05:43
 * @FilePath: /sidetube/src/pages/video/component/recomend.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

import * as React from "react";

import InfiniteScroll from "react-infinite-scroll-component";
import searchAxios from "./axios";
import ExhibitVideo from "component/exhibitVideo";
import Grid from "@mui/material/Grid";
import { useMediaQuery } from "@mui/material";
import ExhibitVideoForPC from "component/exhibitVideoForPC";
import { useParams } from "react-router-dom";
import { Box } from "@mui/system";
import { useInfiniteQuery } from "react-query";

const pageSize = 10;

const RAxios = searchAxios();
export default function SearchPage() {
  const isScreenSmall = useMediaQuery("(max-width: 960px)");
  let { query } = useParams();

  const { data, fetchNextPage, hasNextPage, isFetching } = useInfiniteQuery(
    ["search", query],
    async ({ pageParam = 0 }) => {
      const resp = await RAxios.post("", {
        query: decodeURIComponent(query),
        skip: pageParam,
        limit: pageSize,
      });

      console.log(resp);
      if (!resp?.data?.list) {
        throw new Error("search is missing");
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
    window.scrollTo(0, 0);
    fetchNextPage({ pageParam: 0, resetError: true });
  }, [query, fetchNextPage]);

  return (
    <>
      <InfiniteScroll
        dataLength={dataToRender.length}
        next={handleFetchMore}
        hasMore={hasNextPage}
      >
        <Grid container spacing={2}>
          {dataToRender.map((d) => (
            <Grid item key={d.video_id}>
              {isScreenSmall ? (
                <ExhibitVideo meta={d} showAvater={true} />
              ) : (
                <ExhibitVideoForPC meta={d} />
              )}
            </Grid>
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
          <p>
            No videos were found matching the search term '
            {decodeURIComponent(query)}'.
          </p>
        </Box>
      )}
    </>
  );
}
