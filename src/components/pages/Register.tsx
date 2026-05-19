import { ArrowLeftOutlined } from "@ant-design/icons";
import {
  Alert,
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
  type FormProps,
} from "antd";
import { FcGoogle } from "react-icons/fc";
import { GrGithub } from "react-icons/gr";
import { useNavigate } from "react-router-dom";
import {
  createUser,
  type CreateUserRequest,
  type CreateUserResponse,
} from "../../hooks/user";
import { useMutation } from "@tanstack/react-query";
import Link from "antd/es/typography/Link";

const { Title, Text } = Typography;

type FormFields = {
  username: string;
  password: string;
  confirmPassword: string;
  altName: string;
  email: string;
  altEmail: string
};

export const Register = () => {
  const { mutate, isPending, isError, isSuccess } = useMutation<
    CreateUserResponse,
    Error,
    CreateUserRequest
  >({
    mutationFn: createUser
  });

  const onFinish: FormProps<FormFields>["onFinish"] = (values) => {
    const userData: CreateUserRequest = {
      userName: values.username,
      altEmail: values.altEmail,
      email: values.email,
      altName: values.altName,
      password: values.password,

    };

    mutate(userData);
  };

  return (
    <Row justify="center" align="middle">
      <Col span={20}>
        <Card variant="borderless" style={{ padding: 16 }}>
          <Row gutter={[48, 32]} align="middle">
            <Col span={12}>
              <Flex vertical align="center" gap={0} style={{ padding: 16 }}>
                <Title level={2} style={{ margin: 0 }}>
                  Takara
                </Title>
                <Text type="secondary">Your Personal Finance Dashboard</Text>
              </Flex>


              <Image
                src="/signup.svg"
                alt="signup"
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
                <Space align="center">
                  <Button icon={<ArrowLeftOutlined />} type="text" href="/" />
                  <Text italic style={{ fontSize: 20 }}>
                    Sign Up
                  </Text>
                </Space>

                {isError && (
                  <Alert
                    description="User name already exists"
                    type="error"
                    showIcon
                  />
                )}

                {isSuccess && (
                  <Alert
                    title="User Created Successfully"
                    description={<Text>Please Login with your credentials <Link href="/">here</Link></Text>}
                    type="success"
                    showIcon
                  />
                )}

                <Form<FormFields>
                  layout="vertical"
                  size="large"
                  onFinish={onFinish}
                  requiredMark="optional"
                >
                  <Form.Item
                    label="Username"
                    name="username"
                    required
                    rules={[
                      { required: true, message: "Enter a username" },
                      { min: 3, message: "Minimum 3 characters" },
                    ]}
                  >
                    <Input placeholder="Enter username" />
                  </Form.Item>

                  <Form.Item
                    label="Nickname"
                    name="altName"
                    rules={[
                      { min: 1, message: "Minimum 1 characters" },
                    ]}
                  >
                    <Input placeholder="Enter Nickname" />
                  </Form.Item>

                  <Form.Item
                    label="Email"
                    name="email"
                  >
                    <Input placeholder="Enter Email" />
                  </Form.Item>

                  <Form.Item
                    label="Personal Email"
                    name="altEmail"
                  >
                    <Input placeholder="Enter Personal Email" />
                  </Form.Item>

                  <Form.Item
                    label="Password"
                    name="password"
                    required
                    rules={[
                      { required: true, message: "Enter your password" },
                      { min: 6, message: "Minimum 6 characters" },
                    ]}
                  >
                    <Input.Password placeholder="Enter password" />
                  </Form.Item>

                  <Form.Item
                    label="Confirm Password"
                    name="confirmPassword"
                    required
                    dependencies={["password"]}
                    rules={[
                      { required: true, message: "Confirm your password" },
                      ({ getFieldValue }) => ({
                        validator(_, value) {
                          if (!value || getFieldValue("password") === value) {
                            return Promise.resolve();
                          }
                          return Promise.reject(
                            new Error("Passwords do not match"),
                          );
                        },
                      }),
                    ]}
                  >
                    <Input.Password placeholder="Confirm password" />
                  </Form.Item>

                  <Row justify="center" style={{ paddingBottom: 12 }}>
                    <Button
                      type="primary"
                      htmlType="submit"
                      size="large"
                      loading={isPending}
                    >
                      Sign Up
                    </Button>
                  </Row>

                  {/*
                    TODO: Add SSO                  
                  */}
                  {/* <Divider>
                    <Text type="secondary">Or Sign Up With</Text>
                  </Divider>

                  <Row
                    justify="space-evenly"
                    align="middle"
                    style={{ padding: 14 }}
                  >
                    <Button icon={<FcGoogle />}>Google</Button>
                    <Button icon={<GrGithub />}>Github</Button>
                  </Row> */}
                </Form>
              </Space>
            </Col>
          </Row>
        </Card>
      </Col>
    </Row>
  );
};
