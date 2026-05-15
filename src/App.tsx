import { ConfigProvider, type ThemeConfig } from "antd";
import { BrowserRouter } from "react-router-dom";
import { Router } from "./router/Router";

const theme: ThemeConfig = {
  token: {
    "colorBgBase": "#fbfbfb",
    "colorSuccess": "#52c41a",
    "borderRadius": 8,
    "fontSize": 14,
    "colorLink": "#40a2e3",
    "colorError": "#ed0a28",
    "colorWarning": "#fa541c"
  }
}


function App() {
  return (
    <ConfigProvider theme={theme}>
      <BrowserRouter>
        <Router />
      </BrowserRouter>
    </ConfigProvider>
  )
}

export default App;
