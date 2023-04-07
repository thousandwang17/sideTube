/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-03-02 18:49:58
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-31 16:41:37
 * @FilePath: /sidetube/src/compmnent/exhibitVideo.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Grid } from "@mui/material";
import { Box } from "@mui/system";
import PropTypes from "prop-types";
import { styled } from "@mui/material/styles";
import timeAgo from "common/timeAgo";
import { Link } from "react-router-dom";
import ColorAvatar from "component/avatar";
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
  line-height: 1.6rem;
  font-size: 1.2rem;
`;

const EditNameBox = styled(Box)`
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  font-size: 15px;
  padding-top: 8px;
  padding-bottom: 8px;
`;

const EditInfoBox = styled(Box)`
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  font-size: 14px;
  color: #666666;
`;

const EditDescBox = styled(Box)`
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  font-size: 14px;
  color: #666666;
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

export default function ExhibitVideoForPC({ meta }) {
  return (
    <Grid container>
      <Grid
        item
        sx={{
          justifyContent: "right",
          display: "flex",
          maxWidth: 360,
          minWidth: 240,
        }}
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
      <Grid item sx={{ pl: 2, flex: 1 }}>
        <Box sx={{ display: "flex", flexDirection: "column" }}>
          <EditTittleBox>
            <Link
              to={"/video/" + meta.video_id}
              style={{ textDecoration: "none", color: "inherit" }}
            >
              {meta.title ?? ""}
            </Link>
          </EditTittleBox>

          <EditInfoBox>
            <Link
              to={"/video/" + meta.video_id}
              style={{ textDecoration: "none", color: "inherit" }}
            >
              <span>views : {meta.views ?? 0} </span>
              <SpanDot>{timeAgo(meta.createTime)}</SpanDot>
            </Link>
          </EditInfoBox>

          <EditNameBox>
            <Link
              to={"/channel/" + meta.user_id}
              style={{
                textDecoration: "none",
                color: "inherit",
                display: "flex",
              }}
            >
              <ColorAvatar
                userName={meta.user_name}
                sx={{ width: 21, height: 21, fontSize: 14, mr: 1 }}
              />
              <span>{meta.user_name}</span>
            </Link>
          </EditNameBox>

          <EditDescBox>
            <Link
              to={"/video/" + meta.video_id}
              style={{ textDecoration: "none", color: "inherit" }}
            >
              <span style={{ whiteSpace: "break-spaces" }}>
                {meta.desc ?? ""}
              </span>
            </Link>
          </EditDescBox>
        </Box>
      </Grid>
    </Grid>
  );
}

ExhibitVideoForPC.propTypes = {
  meta: PropTypes.shape({
    video_id: PropTypes.string.isRequired,
    user_name: PropTypes.string.isRequired,
    user_id: PropTypes.string.isRequired,
    png: PropTypes.string.isRequired,
    view: PropTypes.string,
    desc: PropTypes.string,
    createTime: PropTypes.string.isRequired,
    duration: PropTypes.string.isRequired,
  }).isRequired,
};
