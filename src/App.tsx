import type { ThemeConfig } from 'antd';
import { ConfigProvider, theme } from 'antd';
import { BrowserRouter } from 'react-router-dom';
import { Router } from './router/Router';

const appTheme: ThemeConfig = {
    algorithm: theme.defaultAlgorithm,

    token: {
        colorPrimary: '#1677ff',
        colorSuccess: '#73d13d',
        colorWarning: '#fa541c',

        borderRadius: 10,

        colorBgLayout: '#f5f7fa',
        colorBgContainer: '#ffffff',

        colorText: '#262626',
        colorTextSecondary: '#595959',

        fontSize: 14,
        fontFamily: 'Inter, system-ui, sans-serif',
    },

    components: {
        Layout: {
            headerBg: '#ffffff',
            siderBg: '#ffffff',
            bodyBg: '#f5f7fa',
        },

        Menu: {
            itemBorderRadius: 8,
            itemHeight: 42,
        },

        Card: {
            borderRadiusLG: 14,
        },
    },
};

function App() {
    return (
        <ConfigProvider theme={appTheme}>
            <BrowserRouter>
                <Router />
            </BrowserRouter>
        </ConfigProvider>
    );
}

export default App;
