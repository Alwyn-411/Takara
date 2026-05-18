import { useMutation } from "@tanstack/react-query";
import {
  Button,
  Card,
  Checkbox,
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
import { AuthUser } from "../../hooks/auth";
import { getUser } from "../../hooks/user";
import { useUserStore } from "../../store/User";

const { Title, Text, Link } = Typography;

type FormFields = {
  username: string;
  password: string;
  remember: boolean;
};

export const Login = () => {
  const navigate = useNavigate();

  const fetchUserMutation = useMutation({
    mutationFn: getUser,
    onSuccess: (data) => {
      useUserStore.getState().updateUser({
        userId: data.userId,
        userName: data.userName,
        altName: data.altName,
        email: data.email,
        altEmail: data.altEmail,
      });
    },
    onError: (error) => {
      console.error("Failed to Fetch user:", error);
    },
  });

  const loginMutation = useMutation({
    mutationFn: AuthUser,
    onSuccess: (data) => {
      fetchUserMutation.mutate(data.id);
      navigate("/home");
    },
    onError: (error) => {
      console.error("Login failed:", error);
    },
  });

  const onFinish: FormProps<FormFields>["onFinish"] = (values) => {
    loginMutation.mutate({
      userName: values.username!,
      password: values.password!,
      remember: values.remember,
    });
  };

  const onFinishFailed: FormProps<FormFields>["onFinishFailed"] = (
    errorInfo,
  ) => {
    console.log("Failed:", errorInfo);
  };

  return (
    <Row justify="center" align="middle">
      <Col span={20}>
        <Card variant="borderless" style={{ padding: 16 }}>
          <Row gutter={[48, 32]} align="middle">
            <Col span={12}>
              <Space
                orientation="vertical"
                size="middle"
                style={{ width: "100%" }}
              >
                <Flex vertical align="center" gap={0} style={{ padding: 16 }}>
                  <Title level={2} style={{ margin: 0 }}>
                    Welcome to Takara
                  </Title>

                  <Text type="secondary">Your Personal Finance Dashboard</Text>
                </Flex>

                <Form<FormFields>
                  layout="vertical"
                  initialValues={{ remember: true }}
                  size="large"
                  onFinish={onFinish}
                  onFinishFailed={onFinishFailed}
                >
                  <Form.Item<FormFields>
                    label="Username"
                    name="username"
                    required={false}
                    rules={[
                      {
                        required: true,
                        message: "Enter your username",
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
                        message: "Enter your password",
                      },
                    ]}
                  >
                    <Input.Password placeholder="Enter password" />
                  </Form.Item>

                  <Row
                    justify="space-between"
                    align="middle"
                    style={{ paddingInline: 12 }}
                  >
                    <Form.Item<FormFields>
                      name="remember"
                      valuePropName="checked"
                      noStyle
                    >
                      <Checkbox>Remember me</Checkbox>
                    </Form.Item>

                    <Link>Forgot password?</Link>
                  </Row>

                  <Row justify="center" style={{ paddingBottom: 12 }}>
                    <Button type="primary" htmlType="submit" size="large">
                      Login in
                    </Button>
                  </Row>

                  <Divider>
                    <Text type="secondary">Or Login With</Text>
                  </Divider>

                  <Row
                    justify="space-evenly"
                    align="middle"
                    style={{ padding: 14 }}
                  >
                    <Button icon={<FcGoogle />}>Google</Button>
                    <Button icon={<GrGithub />}>Github</Button>
                  </Row>

                  <Row justify="center">
                    <Space size="small">
                      <Text type="secondary">Don't have an account ?</Text>
                      <Link href="/signup">Register Now.</Link>
                    </Space>
                  </Row>
                </Form>
              </Space>
            </Col>

            <Col span={12}>
              <Image
                src="/login.svg"
                alt="login"
                preview={false}
                width="100%"
              />
            </Col>
          </Row>
        </Card>
      </Col>
    </Row>
  );
};
