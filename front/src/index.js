/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-13 20:51:41
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-03-30 15:55:51
 * @FilePath: /sidetube/src/index.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import App from "./App";
import reportWebVitals from "./reportWebVitals";
import { createBrowserRouter, RouterProvider } from "react-router-dom";
import VideoOwnerList from "./pages/channel/owner/video";
import VideoPlayer from "./pages/video/video";
import HomeVideoRecommend from "pages/home/home";
import Channel from "pages/channel/viewer/channel";
import SearchPage from "pages/search/search";
import HistoryPage from "pages/history/history";

const root = ReactDOM.createRoot(document.getElementById("root"));

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    children: [
      {
        path: "channel/videoUpload",
        element: <VideoOwnerList />,
      },
      {
        path: "channel/:user_id",
        element: <Channel />,
      },
      {
        path: "/",
        element: <HomeVideoRecommend />,
      },
      {
        path: "/history",
        element: <HistoryPage />,
      },
      {
        path: "video/:id",
        element: <VideoPlayer />,
      },
      {
        path: "search/:query",
        element: <SearchPage />,
      },
    ],
  },
]);
root.render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
