import * as React from "react";
import { styled, alpha } from "@mui/material/styles";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Toolbar from "@mui/material/Toolbar";
import IconButton from "@mui/material/IconButton";
import Typography from "@mui/material/Typography";
import InputBase from "@mui/material/InputBase";
import Badge from "@mui/material/Badge";

import SearchIcon from "@mui/icons-material/Search";
import AccountCircle from "@mui/icons-material/AccountCircle";
import NotificationsIcon from "@mui/icons-material/Notifications";
import MoreIcon from "@mui/icons-material/MoreVert";
import { YouTube } from "@mui/icons-material";
import MenuIcon from "@mui/icons-material/Menu";
import { AuthContext } from "common/context/authcontext";
import { Link, useNavigate } from "react-router-dom";
import AppBarUserMenu from "./component/appbarUserMemu";
import LoginButton from "component/loginBox";
import { useMediaQuery } from "@mui/material";
import AppBarMsgsMenu from "./component/appbarMsgMenu";

const drawerWidth = 240;

const Search = styled("div")(({ theme }) => ({
  position: "relative",
  borderRadius: theme.shape.borderRadius,
  backgroundColor: alpha(theme.palette.common.white, 0.15),
  "&:hover": {
    backgroundColor: alpha(theme.palette.common.white, 0.25),
  },
  marginRight: theme.spacing(2),
  marginLeft: 0,
  width: "100%",
  border: "1px solid",
  borderColor: "#aaa",
  [theme.breakpoints.up("sm")]: {
    marginLeft: theme.spacing(3),
    width: "100%",
    maxWidth: 500,
  },
}));

const SearchIconWrapper = styled("div")(({ theme }) => ({
  padding: theme.spacing(0, 2),
  height: "100%",
  position: "absolute",
  pointerEvents: "none",
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
}));

const StyledInputBase = styled(InputBase)(({ theme }) => ({
  // color: "inherit",
  "& .MuiInputBase-input": {
    padding: theme.spacing(1, 1, 1, 0),
    // vertical padding + font size from searchIcon
    paddingLeft: `calc(1em + ${theme.spacing(4)})`,
    transition: theme.transitions.create("width"),
    width: "100%",
  },
}));

export default function PrimarySearchAppBar({ handleDrawerToggle }) {
  const [anchorEl, setAnchorEl] = React.useState(null);
  const isMenuOpen = Boolean(anchorEl);

  const [msgAnchorEl, setMsgAnchorEl] = React.useState(null);
  const isMsgMenuOpen = Boolean(msgAnchorEl);
  const [messageRead, setMessageRead] = React.useState(
    localStorage.getItem("menu_index") ?? 0
  );

  const [search, setSearch] = React.useState("");
  const [mobileMoreAnchorEl, setMobileMoreAnchorEl] = React.useState(null);

  const { isUserLoggedIn, setIsUserLoggedIn } = React.useContext(AuthContext);
  const isPhoneScreen = useMediaQuery("(max-width: 425px)");

  const handleProfileMenuOpen = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMessageMenuOpen = (event) => {
    localStorage.setItem("menu_index", 1);
    setMessageRead(1);
    setMsgAnchorEl(event.currentTarget);
  };

  const handleMobileMenuClose = () => {
    setMobileMoreAnchorEl(null);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
    handleMobileMenuClose();
  };

  const handleMsgMenuClose = () => {
    setMsgAnchorEl(null);
    handleMobileMenuClose();
  };

  const handleMenuButtonOnClick = () => {
    handleDrawerToggle();
  };

  const handleSearchChange = (e) => {
    setSearch(e.target.value);
  };

  const history = useNavigate();

  function handleSubmit(event) {
    event.preventDefault(); // Prevent the default form submission
    history(`search/${encodeURIComponent(search)}`); // Navigate to the search results page
  }

  const menuId = "primary-search-account-menu";
  const msgMenuId = "primary-search-messsage-menu";
  return (
    <Box sx={{ display: "flex" }}>
      <AppBar
        position="fixed"
        sx={{
          ml: { sm: `${drawerWidth}px` },
          zIndex: (theme) => theme.zIndex.drawer + 1,
          color: "primary.main",
          bgcolor: "primary.contrastText",
        }}
      >
        <Toolbar>
          <IconButton
            size="large"
            edge="start"
            color="inherit"
            aria-label="menu"
            sx={{ mr: isPhoneScreen ? 0 : 2 }}
            onClick={handleMenuButtonOnClick}
          >
            <MenuIcon />
          </IconButton>
          <Link
            to={`/`}
            style={{
              textDecoration: "none",
              color: "inherit",
            }}
          >
            <Typography
              variant="h6"
              noWrap
              component="div"
              sx={{
                display: "flex",
                alignItems: "center",
                mx: isPhoneScreen ? 2 : 0,
              }}
            >
              <YouTube />
              {!isPhoneScreen && "SideTube"}
            </Typography>
          </Link>

          <Box sx={{ flexGrow: 1, display: "flex", justifyContent: "center" }}>
            <Search>
              <form onSubmit={handleSubmit}>
                <SearchIconWrapper>
                  <SearchIcon />
                </SearchIconWrapper>
                <StyledInputBase
                  value={search}
                  onChange={handleSearchChange}
                  placeholder="Search…"
                  inputProps={{ "aria-label": "search" }}
                  sx={{ width: "100%" }}
                />
              </form>
            </Search>
          </Box>
          {isUserLoggedIn && (
            <Box sx={{ display: { xs: "none", md: "flex" } }}>
              <IconButton
                size="large"
                aria-label="show 1 new notifications"
                color="inherit"
                aria-controls={msgMenuId}
                aria-haspopup="true"
                onClick={handleMessageMenuOpen}
              >
                <Badge badgeContent={messageRead > 0 ? 0 : 1} color="error">
                  <NotificationsIcon />
                </Badge>
              </IconButton>
              <IconButton
                size="large"
                edge="end"
                aria-label="account of current user"
                aria-controls={menuId}
                aria-haspopup="true"
                onClick={handleProfileMenuOpen}
                color="inherit"
              >
                <AccountCircle />
              </IconButton>
            </Box>
          )}
          {!isUserLoggedIn && <LoginButton />}
          <Box sx={{ display: { xs: "flex", md: "none" } }}>
            <IconButton
              size="large"
              aria-label="show more"
              aria-controls={menuId}
              aria-haspopup="true"
              onClick={handleProfileMenuOpen}
              color="inherit"
            >
              <MoreIcon />
            </IconButton>
          </Box>
        </Toolbar>
      </AppBar>
      {isUserLoggedIn && (
        <AppBarUserMenu
          isMenuOpen={isMenuOpen}
          menuId={menuId}
          anchorEl={anchorEl}
          handleMenuClose={handleMenuClose}
        />
      )}
      {isUserLoggedIn && (
        <AppBarMsgsMenu
          isMenuOpen={isMsgMenuOpen}
          menuId={msgMenuId}
          anchorEl={msgAnchorEl}
          handleMenuClose={handleMsgMenuClose}
        />
      )}
    </Box>
  );
}
