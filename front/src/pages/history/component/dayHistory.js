/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-03-30 20:21:54
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-31 15:38:38
 * @FilePath: /sidetube/src/pages/history/component/dayHistory.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Grid } from "@mui/material";
import ExhibitVideo from "component/exhibitVideo";
import ExhibitVideoForPC from "component/exhibitVideoForPC";
import * as React from "react";
import { useMediaQuery } from "@mui/material";
import { format, isToday, isYesterday } from "date-fns";
import PropTypes from "prop-types";
import { Box } from "@mui/system";

export default function DayHistory({ data }) {
  const isScreenSmall = useMediaQuery("(max-width: 960px)");
  const _date = new Date(data[0].historyTime);
  const formattedDate = format(_date, "MMMM d");

  const key = format(_date, "yyyy-MM-dd");
  const [title, setTitle] = React.useState(formattedDate);

  React.useEffect(() => {
    if (isToday(_date)) {
      setTitle("Today");
    } else if (isYesterday(_date)) {
      setTitle("YesterDay");
    }
  }, [data]);
  return (
    <>
      <Box component="h4" key={key}>
        {title}
      </Box>
      <Grid container spacing={2}>
        {data.map((d) => (
          <Grid item key={d.video_id}>
            {isScreenSmall ? (
              <ExhibitVideo meta={d} showAvater={true} />
            ) : (
              <ExhibitVideoForPC meta={d} />
            )}
          </Grid>
        ))}
      </Grid>
    </>
  );
}

DayHistory.propTypes = {
  data: PropTypes.oneOfType([
    PropTypes.arrayOf(
      PropTypes.shape({
        video_id: PropTypes.string.isRequired,
        user_name: PropTypes.string.isRequired,
        user_id: PropTypes.string.isRequired,
        png: PropTypes.string.isRequired,
        view: PropTypes.string,
        viewTime: PropTypes.string.isRequired,
        duration: PropTypes.string.isRequired,
      })
    ),
    PropTypes.array,
  ]).isRequired,
};
