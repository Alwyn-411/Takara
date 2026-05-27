import { useMutation } from '@tanstack/react-query';
import { Button, Card, Checkbox, Col, Divider, Flex, Form, Image, Input, message, Row, Space, Typography, type FormProps } from 'antd';
import { FcGoogle } from 'react-icons/fc';
import { GrGithub } from 'react-icons/gr';
import { useNavigate } from 'react-router-dom';
import { AuthUser } from '../../../api/auth';
import { getUser } from '../../../api/user';
import { useUserStore } from '../../../store/User';

const { Title, Text, Link } = Typography;

type FormFields = {
    username: string;
    password: string;
    remember: boolean;
};

export const Login = () => {
    const navigate = useNavigate();

    const updateUser = useUserStore((s) => s.updateUser);
    const setToken = useUserStore((s) => s.setToken);

    const { mutate } = useMutation({
        mutationFn: async (values: FormFields) => {
            const auth = await AuthUser({
                userName: values.username,
                password: values.password,
                remember: values.remember,
            });

            setToken(auth.token);

            const user = await getUser(auth.id);
            return user;
        },
        onSuccess: (user) => {
            updateUser({
                userId: user.userId,
                userName: user.userName,
                altName: user.altName,
                email: user.email,
                altEmail: user.altEmail,
            });
            navigate('/home');
        },
        onError: (error) => {
            console.error(error);
            message.error('Invalid Credentials');
        },
    });

    const onFinish: FormProps<FormFields>['onFinish'] = (values) => {
        mutate(values);
    };

    const onFinishFailed: FormProps<FormFields>['onFinishFailed'] = (errorInfo) => {
        console.log('Failed:', errorInfo);
    };

    return (
        <Row justify="center" align="middle">
            <Col span={20}>
                <Card variant="borderless" style={{ padding: 16 }}>
                    <Row gutter={[48, 32]} align="middle">
                        <Col span={12}>
                            <Space orientation="vertical" size="middle" style={{ width: '100%' }}>
                                <Flex vertical align="center" gap={0} style={{ padding: 16 }}>
                                    <Title level={2} style={{ margin: 0 }}>
                                        Welcome to Takara
                                    </Title>

                                    <Text type="secondary">Your Personal Finance Dashboard</Text>
                                </Flex>

                                <Row justify="center">
                                    <Text italic style={{ fontSize: 20 }}>
                                        Login In
                                    </Text>
                                </Row>

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
                                                message: 'Enter your username',
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
                                                message: 'Enter your password',
                                            },
                                        ]}
                                    >
                                        <Input.Password placeholder="Enter password" />
                                    </Form.Item>

                                    <Row justify="space-between" align="middle" style={{ paddingInline: 12 }}>
                                        <Form.Item<FormFields> name="remember" valuePropName="checked" noStyle>
                                            <Checkbox>Remember me</Checkbox>
                                        </Form.Item>

                                        <Link>Forgot password?</Link>
                                    </Row>

                                    <Row justify="center" style={{ paddingBottom: 12 }}>
                                        <Button type="primary" htmlType="submit" size="large">
                                            Login in
                                        </Button>
                                    </Row>

                                    {/*
                    TODO: Add SSO                  
                  */}
                                    {/* <Divider>
                    <Text type="secondary">Or Log in With</Text>
                  </Divider>

                  <Row
                    justify="space-evenly"
                    align="middle"
                    style={{ padding: 14 }}
                  >
                    <Button icon={<FcGoogle />}>Google</Button>
                    <Button icon={<GrGithub />}>Github</Button>
                  </Row> */}

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
                            <Image src="/login.svg" alt="login" preview={false} width="100%" />
                        </Col>
                    </Row>
                </Card>
            </Col>
        </Row>
    );
};
