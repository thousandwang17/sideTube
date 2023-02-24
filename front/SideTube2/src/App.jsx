/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-16 13:55:45
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-01-16 14:31:42
 * @FilePath: /SideTube2/src/App.jsx
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import './App.css';
import './layout/SideBar'
import { Outlet } from "react-router-dom";

import SearchAppBar from './layout/appbar'

function App() {
  return (
    <div>
    <SearchAppBar>
    </SearchAppBar>
    <Outlet />
    </div>
    
  );
}

export default App
