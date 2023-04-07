/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-02-27 15:51:55
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-31 15:53:19
 * @FilePath: /sidetube/src/layout/appbarMemu.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

import * as React from "react";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";
import TextField from "@mui/material/TextField";
import axios from "axios";
import { Form } from "react-router-dom";
import LoadingButton from "@mui/lab/LoadingButton";
import Snackbar from "@mui/material/Snackbar";
import MuiAlert from "@mui/material/Alert";
import { AuthContext } from "common/context/authcontext";
import { Box } from "@mui/system";
import Button from "@mui/material/Button";
import getHost from "common/axios";

const Alert = React.forwardRef(function Alert(props, ref) {
  return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />;
});

export default function LoginBox() {
  const [open, setOpen] = React.useState(false);
  const [openSnackbar, setOpenSnackbar] = React.useState(false);
  const [emailInput, setEmailInput] = React.useState("");
  const [passwordInput, setPasswordInput] = React.useState("");
  const [repeatPasswordInput, setRepeatPasswordInput] = React.useState("");
  const [nameInput, setNameInput] = React.useState("");
  const [register, setRegister] = React.useState(false);
  const [submitLoading, setSubmitLoading] = React.useState(false);
  const [errorPassword, setErrorPassword] = React.useState(false);
  const [errorResponse, setErrorResponse] = React.useState(false);
  const [errorResponseMesg, setErrorResponseMesg] = React.useState("");
  const { setIsUserLoggedIn, setShowLoginBox } = React.useContext(AuthContext);

  const formId = "userLoginForm";
  const handleDialogClose = () => {
    setOpen(false);
  };

  const handleDialogOpen = () => {
    setOpen(true);
  };

  const handleSnackbarClose = () => {
    setOpenSnackbar(false);
  };

  const handleRegisterClick = () => {
    setRegister((pre) => !pre);
    setEmailInput("");
    setPasswordInput("");
    setRepeatPasswordInput("");
    setNameInput("");
    setErrorPassword(false);
  };

  const handelEmailOnChange = (e) => {
    setEmailInput(e.target.value);
  };

  const handelPasswordOnChange = (e) => {
    setPasswordInput(e.target.value);
    setErrorPassword(false);
  };

  const handelRepeatPasswordOnChange = (e) => {
    setRepeatPasswordInput(e.target.value);
    setErrorPassword(false);
  };

  const handelNameInputOnChange = (e) => {
    setNameInput(e.target.value);
  };

  const handleDialogSubmit = (e) => {
    e.preventDefault();

    async function login() {
      try {
        setSubmitLoading(true);
        const response = await UserAxios()
          .post(
            "/user/login",
            {
              email: emailInput,
              passWord: passwordInput,
            },
            { withCredentials: true }
          )
          .then(() => {
            handleDialogClose();
            setIsUserLoggedIn(true);
            setErrorResponse(false);
          });
      } catch (e) {
        console.log(e);
        if (e.response && e.response.status === 400) {
          setOpenSnackbar(true);
          setErrorResponse(true);
          setErrorResponseMesg(e.response.data ?? "Error");
        }
      } finally {
        setSubmitLoading(false);
      }
    }

    async function Register() {
      try {
        setSubmitLoading(true);
        const response = await UserAxios().post("/user/register", {
          email: emailInput,
          passWord: passwordInput,
          name: nameInput,
        });

        handleRegisterClick();
        setErrorResponse(false);
      } catch (e) {
        console.log(e);
        if (e.response && e.response.status === 400) {
          setErrorResponse(true);
          setErrorResponseMesg(e.response.data ?? "Error");
        }
      } finally {
        setOpenSnackbar(true);
        setSubmitLoading(false);
      }
    }

    if (register) {
      if (passwordInput !== repeatPasswordInput) {
        setErrorPassword(true);
        return;
      }
      Register();
    } else {
      login();
    }
  };

  React.useEffect(() => {
    setShowLoginBox(() => handleDialogOpen);
  }, []);

  return (
    <Box sx={{ display: { xs: "none", md: "flex" } }}>
      <Button
        variant="contained"
        onClick={handleDialogOpen}
        size="small"
        sx={{
          borderRadius: 10,
          bgcolor: "primary.contrastText",
          color: "primary.main",
          fontWeight: "bold",
          "&:hover": {
            bgcolor: "primary.contrastText",
            color: "primary.main",
          },
        }}
      >
        Login
      </Button>
      <Dialog open={open} onClose={handleDialogClose}>
        <Form onSubmit={handleDialogSubmit} id={formId}>
          <DialogTitle>
            {register ? "Register an account" : "Log In"}
          </DialogTitle>
          <DialogContent sx={{ pd: 0, pb: 0 }}>
            <TextField
              margin="dense"
              id="email"
              label="Email Address"
              type="email"
              fullWidth
              variant="standard"
              value={emailInput}
              onChange={handelEmailOnChange}
              sx={{ mt: 0 }}
              required
              autoFocus
            />
            <TextField
              margin="dense"
              id="password"
              label="Password"
              type="password"
              fullWidth
              variant="standard"
              value={passwordInput}
              onChange={handelPasswordOnChange}
              required
              error={errorPassword}
            />
            {register && (
              <TextField
                margin="dense"
                id="repeatPassword"
                label="Repeat Password"
                type="password"
                fullWidth
                variant="standard"
                value={repeatPasswordInput}
                onChange={handelRepeatPasswordOnChange}
                required
                error={errorPassword}
                helperText={
                  errorPassword ? "Password not equal Repeat Password." : ""
                }
              />
            )}
            {register && (
              <TextField
                margin="dense"
                id="name"
                label="Name"
                type="text"
                fullWidth
                variant="standard"
                value={nameInput}
                onChange={handelNameInputOnChange}
                sx={{ mt: 0 }}
                required
              />
            )}
            <DialogContentText sx={{ pt: 1, fontSize: "15px", color: "#888" }}>
              {register ? "Click" : "Still not have an account? Click"}
              <Button
                sx={{ minWidth: "50px", pt: "4px" }}
                onClick={handleRegisterClick}
              >
                here
              </Button>
              for {register ? "login" : "register"}
            </DialogContentText>
          </DialogContent>

          <DialogActions>
            <Button onClick={handleDialogClose}>Cancel</Button>
            <LoadingButton
              variant="contained"
              size="small"
              loading={submitLoading}
              sx={{ float: "right", ml: 1 }}
              type="submit"
            >
              {register ? "Register" : "LogIn"}
            </LoadingButton>
          </DialogActions>
        </Form>
        <Snackbar
          open={openSnackbar}
          autoHideDuration={3000}
          onClose={handleSnackbarClose}
        >
          <Alert
            severity={errorResponse ? "error" : "success"}
            sx={{ width: "100%", zIndex: 1501 }}
          >
            {errorResponse ? errorResponseMesg : "Register successfully!!"}
          </Alert>
        </Snackbar>
      </Dialog>
    </Box>
  );
}

function UserAxios() {
  // Set config defaults when creating the instance
  const instance = axios.create({
    baseURL: getHost() + "/api",
  });

  instance.defaults.headers.post["Content-Type"] = "application/json";
  // defaultHeaders(instance)
  return instance;
}
