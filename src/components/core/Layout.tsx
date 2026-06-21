import { AccountBookOutlined, BankOutlined, CaretLeftOutlined, CaretRightOutlined, HomeOutlined, UserOutlined } from '@ant-design/icons';

import { Avatar, Button, Dropdown, Flex, Layout, Menu, Typography, theme, type MenuProps } from 'antd';

import type { ItemType, MenuItemType } from 'antd/es/menu/interface';
import { useState } from 'react';
import { FaArrowTrendUp } from 'react-icons/fa6';
import { Navigate, Outlet, useLocation, useNavigate } from 'react-router-dom';

import { useUserStore } from '../../store/User';

import './Layout.css';

const { Header, Sider, Content } = Layout;
const { Title, Text } = Typography;

export const ContentLayout = () => {
    const userId = useUserStore((state) => state.userId);

    const navigate = useNavigate();
    const location = useLocation();

    const [collapsed, setCollapsed] = useState(false);

    const {
        token: { colorPrimary, colorBorderSecondary },
    } = theme.useToken();

    if (!userId) {
        return <Navigate to="/" replace />;
    }

    const SidebarItems: ItemType<MenuItemType>[] = [
        {
            key: '/home',
            icon: <HomeOutlined />,
            label: 'Home',
            onClick: () => navigate('/home'),
        },
        {
            key: '/accounts',
            icon: <AccountBookOutlined />,
            label: 'Accounts',
            onClick: () => navigate('/accounts'),
        },
        {
            key: '/holdings',
            icon: <BankOutlined />,
            label: 'Holdings',
            onClick: () => navigate('/holdings'),
        },
    ];

    const ProfileItems: MenuProps['items'] = [
        {
            key: '1',
            label: 'My Account',
            disabled: true,
        },
        {
            type: 'divider',
        },
        {
            key: '2',
            label: 'Profile',
        },
        {
            key: '3',
            label: 'Settings',
        },
        {
            type: 'divider',
        },
        {
            key: '4',
            label: 'Logout',
            danger: true,
        },
    ];

    const ProfileOnClick: MenuProps['onClick'] = (e) => {
        switch (e.key) {
            case '4':
                useUserStore.getState().clearUser();
                navigate('/');
                break;
        }
    };

    return (
        <Layout
            className="layout-root"
            style={
                {
                    '--primary-color': colorPrimary,
                    '--border-color': colorBorderSecondary,
                } as React.CSSProperties
            }
        >
            <Sider
                trigger={null}
                collapsible
                collapsed={collapsed}
                width={240}
                theme="light"
                className="layout-sider"
                style={{
                    borderRight: `1px solid ${colorBorderSecondary}`,
                }}
            >
                <Flex align="center" gap={12} className="logo-container">
                    <div className="logo-icon">
                        <FaArrowTrendUp size={16} />
                    </div>

                    <Flex vertical gap={0} className={`logo-text ${collapsed ? 'logo-text-collapsed' : 'logo-text-expanded'}`}>
                        <Title level={4} className="logo-title">
                            Takara
                        </Title>
                    </Flex>
                </Flex>

                <Menu mode="inline" selectedKeys={[location.pathname]} items={SidebarItems} className="sidebar-menu" />
            </Sider>

            <Layout>
                <Header
                    className="layout-header"
                    style={{
                        borderBottom: `1px solid ${colorBorderSecondary}`,
                    }}
                >
                    <Button
                        type="default"
                        shape="circle"
                        icon={collapsed ? <CaretRightOutlined /> : <CaretLeftOutlined />}
                        onClick={() => setCollapsed(!collapsed)}
                    />

                    <Dropdown
                        placement="bottomRight"
                        menu={{
                            items: ProfileItems,
                            onClick: ProfileOnClick,
                        }}
                    >
                        <Avatar icon={<UserOutlined />} className="profile-avatar" />
                    </Dropdown>
                </Header>

                <Content className="layout-content">
                    <Outlet />
                </Content>
            </Layout>
        </Layout>
    );
};
