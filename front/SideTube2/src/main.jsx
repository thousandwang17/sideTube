/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-16 13:55:45
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-01-16 14:22:14
 * @FilePath: /SideTube2/src/main.jsx
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import './index.css'
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import VideoOwnerList from './channel/owner/list';

const root = ReactDOM.createRoot(document.getElementById('root'));

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    children: [
      {
        path: "channel/",
        element: <VideoOwnerList />,
      },
    ],
  },


]);
root.render(
  <React.StrictMode>
     <RouterProvider router={router} />
  </React.StrictMode>
);