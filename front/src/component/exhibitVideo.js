/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-03-02 18:49:58
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-27 16:21:00
 * @FilePath: /sidetube/src/compmnent/exhibitVideo.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Grid } from "@mui/material";
import { Box } from "@mui/system";
import PropTypes from "prop-types";
import { styled } from "@mui/material/styles";
import timeAgo from "common/timeAgo";
import { Link } from "react-router-dom";
import ColorAvatar from "./avatar";
import { durationFormat } from "common/durationFormat";
import getHost from "common/axios";

const EditTittleBox = styled(Box)`
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  font-weight: bold;
  margin-bottom: 5px;
  line-height: 21px;
  font-size: 14px;
`;

const EditNameBox = styled(Box)`
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  font-size: 12px;
`;

const EditInfoBox = styled(Box)`
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  font-size: 12px;
`;

const SpanDot = styled("span")`
  &::before {
    content: "•";
    margin-right: 0.5em;
  }
`;
const DurationBox = styled("div")(({ theme }) => ({
  position: "absolute",
  bottom: 10,
  right: 2,
  backgroundColor: theme.palette.primary.main,
  color: theme.palette.primary.contrastText,
  paddingLeft: 3,
  paddingRight: 3,
  paddingTop: 0.2,
  paddingBottom: 0.2,
  borderRadius: 1,
  fontSize: 12,
  fontWeight: "bold",
}));

export default function ExhibitVideo({ meta, mini, showAvater }) {
  return (
    <Grid container sx={{ marginBottom: { xs: 3, md: 0 } }}>
      <Grid
        item
        xs={12}
        md={mini ? 5 : 12}
        sx={{ justifyContent: "right", display: "flex" }}
      >
        <Link to={"/video/" + meta.video_id}>
          <div style={{ position: "relative", display: "inline-block" }}>
            <img
              style={{
                maxWidth: "100%",
                borderRadius: "7px",
              }}
              src={
                meta?.png
                  ? getHost() + `/picture/video/` + meta.png
                  : "https://images.unsplash.com/photo-1512917774080-9991f1c4c750?auto=format&w=350&dpr=2"
              }
              alt={meta?.title ?? ""}
            />
            <DurationBox>
              {meta?.duration ? durationFormat(meta.duration) : "00:00"}
            </DurationBox>
          </div>
        </Link>
      </Grid>
      <Grid
        item
        xs={12}
        md={mini ? 7 : 12}
        sx={{ display: "flex", pl: mini ? 1 : 0 }}
      >
        <Box sx={{ pr: 1.5, display: showAvater === false ? "none" : "flex" }}>
          <Link
            to={"/channel/" + meta.user_id}
            style={{ textDecoration: "none", color: "inherit" }}
          >
            <ColorAvatar userName={meta.user_name} />
          </Link>
        </Box>
        <Box sx={{ display: "flex", flexDirection: "column" }}>
          <EditTittleBox>
            <Link
              to={"/video/" + meta.video_id}
              style={{ textDecoration: "none", color: "inherit" }}
            >
              {meta.title ?? ""}
            </Link>
          </EditTittleBox>

          <EditNameBox>
            <Link
              to={
                showAvater
                  ? "/channel/" + meta.user_id
                  : "/video/" + meta.video_id
              }
              style={{ textDecoration: "none", color: "inherit" }}
            >
              {meta.user_name}
            </Link>
          </EditNameBox>

          <EditInfoBox>
            <Link
              to={"/video/" + meta.video_id}
              style={{ textDecoration: "none", color: "inherit" }}
            >
              <span>views : {meta.views ?? 0} </span>
              <SpanDot>{timeAgo(meta.createTime)}</SpanDot>
            </Link>
          </EditInfoBox>
        </Box>
      </Grid>
    </Grid>
  );
}

ExhibitVideo.propTypes = {
  meta: PropTypes.shape({
    video_id: PropTypes.string.isRequired,
    user_name: PropTypes.string.isRequired,
    user_id: PropTypes.string.isRequired,
    png: PropTypes.string.isRequired,
    view: PropTypes.string,
    createTime: PropTypes.string.isRequired,
    duration: PropTypes.string.isRequired,
  }).isRequired,
  mini: PropTypes.bool.isRequired,
  showAvater: PropTypes.bool.isRequired,
};

ExhibitVideo.defaultProps = {
  mini: false, // Set default prop value to true
  showAvater: false,
};
