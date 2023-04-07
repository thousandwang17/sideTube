/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-25 19:36:32
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-04-04 19:13:36
 * @FilePath: /sidetube/src/pages/channel/owner/video.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import * as React from "react";
import { DataGrid } from "@mui/x-data-grid";
import { videoListAxios } from "./axios";
import FormDialog from "./editVideoMeta";
import PropTypes from "prop-types";
import { IconButton } from "@mui/material";
import EditIcon from "@mui/icons-material/Edit";
import Box from "@mui/material/Box";
import Grid from "@mui/material/Grid";
import FormGroup from "@mui/material/FormGroup";
import FormControlLabel from "@mui/material/FormControlLabel";
import Switch from "@mui/material/Switch";
import { styled } from "@mui/material/styles";
import { durationFormat } from "common/durationFormat";
import Tooltip from "@mui/material/Tooltip";
import HelpIcon from "@mui/icons-material/Help";
import getHost from "common/axios";

const VLAxios = videoListAxios();

const EditIconBox = styled(Box)`
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
`;

const RenderDate = (props) => {
  const { hasFocus, value, row } = props;
  const buttonElement = React.useRef(null);
  const rippleRef = React.useRef(null);

  React.useLayoutEffect(() => {
    if (hasFocus) {
      const input = buttonElement.current?.querySelector("input");
      input?.focus();
    } else if (rippleRef.current) {
      // Only available in @mui/material v5.4.1 or later
      rippleRef.current.stop({});
    }
  }, [hasFocus]);

  const showFormDialog = () => {
    row.showFormDialog(props);
  };

  return (
    <Grid container spacing={2} sx={{ padding: 1 }} onClick={showFormDialog}>
      <Grid item xs={4}>
        <div style={{ position: "relative", display: "inline-block" }}>
          <img
            style={{
              maxHeight: 100,
              maxWidth: 100,
            }}
            src={
              row?.png
                ? getHost() + `/picture/video/` + row.png
                : "https://images.unsplash.com/photo-1512917774080-9991f1c4c750?auto=format&w=350&dpr=2"
            }
            alt={row?.title ?? ""}
          />
          <Box
            sx={{
              position: "absolute",
              bottom: 0,
              right: 0,
              bgcolor: "primary.main",
              color: "primary.contrastText",
              px: 0.4,
              py: 0.2,
              borderRadius: 1,
              fontSize: 12,
              fontWeight: "bold",
            }}
          >
            {row?.duration ? durationFormat(row.duration) : "00:00"}
          </Box>
        </div>
      </Grid>
      <Grid item xs={8} sx={{ display: "flex" }}>
        <EditIconBox> {row?.title ?? ""}</EditIconBox>
        <Box
          sx={{
            marginLeft: "auto",
            display: "flex",
          }}
        >
          <IconButton color="primary" aria-label="edit video">
            <EditIcon />
          </IconButton>
        </Box>
      </Grid>
    </Grid>
  );
};

const Encoding = 0;
const UnPublish = 0;
const Publish = 1;
const permission2String = {
  0: `private`,
  1: `public`,
};

const RenderSwitch = (props) => {
  const { row } = props;

  const showFormDialog = () => {
    row.showFormDialog(props);
  };

  const [permission, setPermission] = React.useState(
    row.permission === Publish ? Publish : UnPublish
  );

  const switchOnChangeHandler = (e) => {
    (async () => {
      const resp = await VLAxios.post("/setState", {
        video_id: row.video_id,
        state: e.target.checked ? Publish : UnPublish,
      });
    })();

    console.log(e.target.checked);
    setPermission(e.target.checked ? Publish : UnPublish);
  };

  return (
    <FormGroup>
      <FormControlLabel
        control={
          <Switch
            checked={permission === Publish}
            onChange={switchOnChangeHandler}
            disabled={!row?.title && !row?.desc}
          />
        }
        label={permission2String[permission]}
      />
    </FormGroup>
  );
};

const columns = [
  {
    field: "video",
    headerName: "video",
    width: 450,
    sortable: false,
    renderCell: RenderDate,
  },
  {
    field: "permission",
    renderHeader: () => (
      <div>
        {"permission "}
        <Tooltip
          title="After setting the title and description, then able to set permission"
          placement="top"
        >
          <IconButton sx={{ padding: 0 }}>
            <HelpIcon />
          </IconButton>
        </Tooltip>
      </div>
    ),
    width: 150,
    sortable: false,
    renderCell: RenderSwitch,
  },
  {
    field: "uploadTime",
    headerName: "Upload Time",
    width: 160,
    sortable: false,
  },
  {
    field: "views",
    headerName: "views",
    type: "number",
    width: 120,
    sortable: false,
    valueGetter: (params) =>
      `${params.row.views ? params.row.views.toLocaleString("en-US") : "0"} `,
  },
  {
    field: "messages",
    headerName: "messages",
    description: "This column has a value getter and is not sortable.",
    type: "number",
    sortable: false,
    width: 120,
    valueGetter: (params) =>
      `${
        params.row.messages ? params.row.messages.toLocaleString("en-US") : "0"
      } `,
  },
];

const pageSize = 5;

export default function ControlledSelectionServerPaginationGrid({
  propRefresh,
}) {
  const [page, setPage] = React.useState(0);
  const [rows, setRows] = React.useState([]);
  const [rowCount, setRowCount] = React.useState(0);
  const [loading, setLoading] = React.useState(false);
  const [formDialog, setFormDialog] = React.useState(null);
  const [refresh, setRefresh] = React.useState(0);

  const SetRefreshCallback = () => {
    setRefresh(refresh + 1);
  };

  // Callback for FormDialog
  const FormDialogCallback = (data) => {
    setFormDialog(data);
  };

  const fetchVidoeList = async (page) => {
    return await VLAxios.post("/list", {
      skip: page * pageSize,
      limit: pageSize,
    }).then((resp) => {
      if (page === 0 && resp?.data?.count) {
        setRowCount(resp.data.count);
      }

      if (resp?.data?.list) {
        const rows = resp.data.list.map((element) => {
          return {
            ...element,
            showFormDialog: FormDialogCallback,
            refreshList: SetRefreshCallback,
          };
        });

        return rows;
      }

      return [];
    });
  };

  React.useEffect(() => {
    let ignore = false;

    setLoading(true);

    try {
      fetchVidoeList(page).then((resp) => {
        if (!ignore) {
          setRows(resp);
        }
      });
    } catch (e) {
    } finally {
      setLoading(false);
    }

    return () => {
      ignore = true;
    };
  }, [page, refresh, propRefresh]);

  return (
    <div style={{ height: 400, width: "100%" }}>
      <DataGrid
        rows={rows}
        getRowId={(row) => row.video_id}
        columns={columns}
        pagination
        pageSize={pageSize}
        rowsPerPageOptions={[5]}
        rowCount={rowCount}
        paginationMode="server"
        onPageChange={(newPage) => {
          setPage(newPage);
        }}
        loading={loading}
        keepNonExistentRowsSelected
        disableColumnMenu
        disableSelectionOnClick
        getRowHeight={() => "auto"}
      />
      <FormDialog data={formDialog} />
    </div>
  );
}
