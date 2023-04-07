import * as React from "react";
import PropTypes from "prop-types";
import Box from "@mui/material/Box";
import CssBaseline from "@mui/material/CssBaseline";
import Divider from "@mui/material/Divider";
import Drawer from "@mui/material/Drawer";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Toolbar from "@mui/material/Toolbar";
import { Outlet, Link } from "react-router-dom";
import VideoLibraryIcon from "@mui/icons-material/VideoLibrary";
import PrimarySearchAppBar from "./appbar";
import { Container } from "@mui/system";
import { useLocation } from "react-router-dom";
import { styled, useTheme } from "@mui/material/styles";
import HomeIcon from "@mui/icons-material/Home";
import { useMediaQuery } from "@mui/material";
import SubscriptionsIcon from "@mui/icons-material/Subscriptions";
import HistoryIcon from "@mui/icons-material/History";
const drawerWidth = 240;

const Main = styled("main", { shouldForwardProp: (prop) => prop !== "open" })(
  ({ theme, open }) => ({
    transition: theme.transitions.create(["margin", "width"], {
      easing: theme.transitions.easing.sharp,
      duration: theme.transitions.duration.leavingScreen,
    }),
    ...(open && {
      width: `calc(100% - ${drawerWidth}px)`,
      marginLeft: `${drawerWidth}px`,
      transition: theme.transitions.create(["margin", "width"], {
        easing: theme.transitions.easing.easeOut,
        duration: theme.transitions.duration.enteringScreen,
      }),
    }),
  })
);

function ResponsiveDrawer(props) {
  const { window } = props;
  const location = useLocation();
  const isVideoPage = location.pathname.includes("/video/");
  const [mobileOpen, setMobileOpen] = React.useState(!isVideoPage);
  const [sideBarType, setSideBarType] = React.useState("temporary");
  const isPhoneScreen = useMediaQuery("(max-width: 425px)");

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  React.useEffect(() => {
    if (isPhoneScreen) {
      setMobileOpen(false);
      setSideBarType("temporary");
      return;
    }

    if (isVideoPage === true) {
      setMobileOpen(false);
      setSideBarType("temporary");
    } else {
      setSideBarType("persistent");
      setMobileOpen(true);
    }
  }, [isVideoPage, isPhoneScreen]);

  const baseList = [
    {
      name: "Home",
      icon: <HomeIcon />,
      linkTo: `/`,
    },
    {
      name: "History",
      icon: <HistoryIcon />,
      linkTo: `/history`,
    },
  ];

  const channelList = [
    {
      name: "My Channel",
      icon: <VideoLibraryIcon />,
      linkTo: `channel/myChannel`,
    },
    {
      name: "Channel Studio",
      icon: <SubscriptionsIcon />,
      linkTo: `channel/videoUpload`,
    },
  ];

  const drawer = (
    <Box>
      <Toolbar />
      <Box sx={{ paddingTop: 1 }} />
      {baseList.map((obj, index) => (
        <ListItem key={obj.name} disablePadding>
          <Link
            to={obj.linkTo}
            style={{
              textDecoration: "none",
              color: "inherit",
              width: "100%",
            }}
          >
            <ListItemButton>
              <ListItemIcon sx={{ color: "primary.main" }}>
                {obj.icon}
              </ListItemIcon>
              <ListItemText primary={obj.name} sx={{ color: "#555555" }} />
            </ListItemButton>
          </Link>
        </ListItem>
      ))}
      <Divider />
      <List>
        {channelList.map((obj, index) => (
          <ListItem key={obj.name} disablePadding>
            <Link
              to={obj.linkTo}
              style={{
                textDecoration: "none",
                color: "inherit",
                width: "100%",
              }}
            >
              <ListItemButton>
                <ListItemIcon sx={{ color: "primary.main" }}>
                  {obj.icon}
                </ListItemIcon>
                <ListItemText primary={obj.name} sx={{ color: "#555555" }} />
              </ListItemButton>
            </Link>
          </ListItem>
        ))}
      </List>
    </Box>
  );

  const container =
    window !== undefined ? () => window().document.body : undefined;

  return (
    <Box sx={{ display: "flex" }}>
      <CssBaseline />
      <PrimarySearchAppBar handleDrawerToggle={handleDrawerToggle} />
      <Box
        component="nav"
        sx={{ flexShrink: { sm: 0 }, paddingTop: 5 }}
        aria-label="mailbox folders"
      >
        {/* The implementation can be swapped with js to avoid SEO duplication of links. */}
        <Drawer
          container={container}
          variant={sideBarType}
          open={mobileOpen}
          onClose={handleDrawerToggle}
          ModalProps={{
            keepMounted: true, // Better open performance on mobile.
          }}
          sx={{
            "& .MuiDrawer-paper": {
              boxSizing: "border-box",
              width: drawerWidth,
            },
          }}
        >
          {drawer}
        </Drawer>
      </Box>
      <Main
        open={!isVideoPage && !isPhoneScreen && mobileOpen}
        sx={{ flexGrow: 1, padding: isPhoneScreen ? 0 : 3 }}
      >
        <Toolbar />
        <Container maxWidth="xl">
          <Outlet />
        </Container>
      </Main>
    </Box>
  );
}

ResponsiveDrawer.propTypes = {
  /**
   * Injected by the documentation to work in an iframe.
   * You won't need it on your project.
   */
  window: PropTypes.func,
};

export default ResponsiveDrawer;
