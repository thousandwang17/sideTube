/*
 * @Author: dennyWang thousandwang17@gmail.com
 * @Date: 2023-01-13 20:51:41
 * @LastEditors: dennyWang thousandwang17@gmail.com
 * @LastEditTime: 2023-04-03 15:02:05
 * @FilePath: /sidetube/src/App.js
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import "./App.css";
import MainContainer from "./layout/main";
import { QueryClientProvider, QueryClient } from "react-query";
import { AuthProvider } from "common/context/authcontext";

const queryClient = new QueryClient();
function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <MainContainer></MainContainer>
      </AuthProvider>
    </QueryClientProvider>
  );
}

export default App;
