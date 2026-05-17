import { ArrowLeftOutlined } from "@ant-design/icons";
import {
  Button,
  Card,
  Col,
  Divider,
  Flex,
  Form,
  Image,
  Input,
  Row,
  Space,
  Typography,
} from "antd";
import { FcGoogle } from "react-icons/fc";
import { GrGithub } from "react-icons/gr";
import { useNavigate } from "react-router-dom";

const { Title, Text } = Typography;

type FormFields = {
  username?: string;
  password?: string;
  confirmPassword?: string;
  remember?: boolean;
};

export const Register = () => {
  const navigate = useNavigate();
  const onLogin = async () => {
    navigate("/home");
  };
  return (
    <Row justify="center" align="middle">
      <Col span={20}>
        <Card variant="borderless" style={{ padding: 16 }}>
          <Row gutter={[48, 32]} align="middle">
            <Col span={12}>
              <Image
                src="/signup.svg"
                alt="login"
                preview={false}
                width="100%"
              />
            </Col>
            <Col span={12}>
              <Space
                orientation="vertical"
                size="middle"
                style={{ width: "100%" }}
              >
                <Flex vertical align="center" gap={0} style={{ padding: 16 }}>
                  <Title level={2} style={{ margin: 0 }}>
                    Takara
                  </Title>

                  <Text type="secondary">Your Personal Finance Dashboard</Text>
                </Flex>

                <Space align="center">
                  <Button icon={<ArrowLeftOutlined />} type="text" href="/" />
                  <Text strong italic style={{ fontSize: 20 }}>
                    Sign up
                  </Text>
                </Space>

                <Form<FormFields>
                  layout="vertical"
                  initialValues={{ remember: true }}
                  size="large"
                >
                  <Form.Item<FormFields>
                    label="Username"
                    name="username"
                    required={false}
                    rules={[
                      {
                        required: true,
                        message: "Enter a username",
                      },
                    ]}
                  >
                    <Input placeholder="Enter username" />
                  </Form.Item>

                  <Form.Item<FormFields>
                    label="Password"
                    name="password"
                    required={false}
                    rules={[
                      {
                        required: true,
                        message: "Enter your new password",
                      },
                    ]}
                  >
                    <Input.Password placeholder="Enter new password" />
                  </Form.Item>

                  <Form.Item<FormFields>
                    label="Confirm Password"
                    name="confirmPassword"
                    required={false}
                    rules={[
                      {
                        required: true,
                        message: "Enter your new password again",
                      },
                    ]}
                  >
                    <Input.Password placeholder="Enter new password again" />
                  </Form.Item>

                  <Row justify="center" style={{ paddingBottom: 12 }}>
                    <Button
                      type="primary"
                      htmlType="submit"
                      size="large"
                      onClick={onLogin}
                    >
                      Sign Up
                    </Button>
                  </Row>

                  <Divider>
                    <Text type="secondary">Or Sign Up With</Text>
                  </Divider>

                  <Row
                    justify="space-evenly"
                    align="middle"
                    style={{ padding: 14 }}
                  >
                    <Button icon={<FcGoogle />}>Google</Button>
                    <Button icon={<GrGithub />}>Github</Button>
                  </Row>
                </Form>
              </Space>
            </Col>
          </Row>
        </Card>
      </Col>
    </Row>
  );
};
