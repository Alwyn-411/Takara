import {
  CaretLeftOutlined,
  CaretRightOutlined,
  HomeOutlined,
} from "@ant-design/icons";

import { Button, Flex, Layout, Menu, Typography, theme } from "antd";
import type { ItemType, MenuItemType } from "antd/es/menu/interface";
import { useState } from "react";
import { FaArrowTrendUp } from "react-icons/fa6";
import { Navigate, Outlet } from "react-router-dom";
import { useUserStore } from "../../store/User";

const { Header, Sider, Content } = Layout;
const { Title } = Typography;

export const ContentLayout = () => {
  const userId = useUserStore((state) => state.userId);
  if (!userId) {
    return <Navigate to="/" replace />;
  }

  const [collapsed, setCollapsed] = useState(false);
  const items: ItemType<MenuItemType>[] = [
    {
      key: "1",
      icon: <HomeOutlined />,
      label: "Home",
    },
  ];

  const {
    token: { colorBorderSecondary, colorPrimary },
  } = theme.useToken();

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
          defaultSelectedKeys={["1"]}
          items={items}
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
            alignItems: "center",
            paddingInline: 16,
            borderBottom: `1px solid ${colorBorderSecondary}`,
          }}
        >
          <Button
            type="text"
            icon={collapsed ? <CaretRightOutlined /> : <CaretLeftOutlined />}
            onClick={() => setCollapsed(!collapsed)}
          />
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
