import {
  AccountBookOutlined,
  CaretLeftOutlined,
  CaretRightOutlined,
  HomeOutlined,
  UserOutlined,
} from "@ant-design/icons";

import { Avatar, Button, Dropdown, Flex, Layout, Menu, Typography, theme, type MenuProps } from "antd";
import type { ItemType, MenuItemType } from "antd/es/menu/interface";
import { useState } from "react";
import { FaArrowTrendUp } from "react-icons/fa6";
import { Navigate, Outlet, useNavigate } from "react-router-dom";
import { useUserStore } from "../../store/User";

const { Header, Sider, Content } = Layout;
const { Title } = Typography;

export const ContentLayout = () => {
  const userId = useUserStore((state) => state.userId);
  const navigate = useNavigate();
  const {
    token: { colorBorderSecondary, colorPrimary },
  } = theme.useToken();

  if (!userId) {
    return <Navigate to="/" replace />;
  }
  const [collapsed, setCollapsed] = useState(false);
  const SidebarItems: ItemType<MenuItemType>[] = [
    {
      key: "1",
      icon: <HomeOutlined />,
      label: "Home",
      onClick: () => { navigate("/home") }
    },
    {
      key: "2",
      icon: <AccountBookOutlined />,
      label: "Accounts",
      onClick: () => { navigate("/accounts") }
    }
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
      extra: '⌘P',
    },
    {
      key: '3',
      label: 'Settings',
      extra: '⌘S',
    },
    {
      type: 'divider',
    },
    {
      key: '4',
      label: 'Logout',
      danger: true,
      extra: '⌘L'
    }
  ];

  const ProfileOnClick: MenuProps['onClick'] = (e) => {
    switch (e.key) {
      case '2':
        console.log("clicked Profile");
        break;
      case '3':
        console.log("clicked settings");
        break;

      case '4':
        console.log("clicked Logout");
        useUserStore.getState().clearUser()
        navigate("/")
        break;

      default:
        console.log("Unknown value");
    }
  }

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        width={240}
        theme="light"
        style={{
          borderRight: `1px solid ${colorBorderSecondary}`,
        }}
      >
        <Flex
          align="center"
          gap={12}
          style={{
            height: 64,
            paddingInline: 20,
          }}
        >
          <FaArrowTrendUp size={20} color={colorPrimary} />
          {!collapsed && (
            <Title
              level={4}
              style={{
                margin: 0,
              }}
            >
              Takara
            </Title>
          )}
        </Flex>

        <Menu
          mode="inline"
          items={SidebarItems}
          style={{
            borderInlineEnd: 0,
            paddingInline: 8,
          }}
        />

      </Sider>

      <Layout>
        <Header
          style={{
            display: "flex",
            padding: "16px",
            justifyContent: "space-between",
          }}
        >
          <Button
            type="text"
            icon={collapsed ? <CaretRightOutlined /> : <CaretLeftOutlined />}
            onClick={() => setCollapsed(!collapsed)}
          />

          <Dropdown menu={{ items: ProfileItems, onClick: ProfileOnClick }}>
            <Avatar
              icon={<UserOutlined />}
              style={{ cursor: "pointer", backgroundColor: colorPrimary }}
            />
          </Dropdown>
        </Header>

        <Content
          style={{
            padding: 24,
          }}
        >
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
};
