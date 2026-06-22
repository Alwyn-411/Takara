import { EditOutlined, HomeOutlined } from '@ant-design/icons';
import { useMutation } from '@tanstack/react-query';
import {
    Avatar,
    Badge,
    Breadcrumb,
    Button,
    Card,
    Col,
    Flex,
    Form,
    Input,
    message,
    Row,
    Select,
    Space,
    Typography,
    Upload,
    type FormProps,
} from 'antd';
import { useState } from 'react';
import { updateUser, updateUserPref, uploadUserAvatar } from '../../../api/user';
import { useUserPrefStore } from '../../../store/Preferences';
import { useUserStore } from '../../../store/User';
import { currencies } from '../../../types/Accounts';
import type { User, UserPref } from '../../../types/User';

const { Title, Text } = Typography;

type UserFormFields = Partial<Pick<User, 'userName' | 'altName' | 'email' | 'altEmail'>>;
interface PasswordFormFields extends Pick<User, 'password'> {
    oldPassword: string;
}
type UserPrefFormFields = Partial<Pick<UserPref, 'currency' | 'theme'>>;

export const appTheme = [
    {
        label: 'Light',
        value: 'light',
    },
    {
        label: 'Dark',
        value: 'dark',
    },
];

export const Profile = () => {
    const [editable, setEditable] = useState(false);
    const [form] = Form.useForm<PasswordFormFields>();

    const updateUserStore = useUserStore((s) => s.updateUser);
    const updateUserPrefStore = useUserPrefStore((s) => s.updateUserPref);

    const createUserPrefApi = useMutation({
        mutationFn: updateUserPref,
        onSuccess: (_, req) => {
            updateUserPrefStore({
                currency: req?.currency,
                theme: req?.theme,
            });

            message.success('User preferances updated successfully');
            setEditable(false);
        },
    });

    const updateUserApi = useMutation({
        mutationFn: updateUser,
        onSuccess: (_, req) => {
            updateUserStore({
                userName: req.userName,
                altName: req.altName,
                email: req.email,
                altEmail: req.altEmail,
            });

            message.success('User updated successfully');
            setEditable(false);
        },
    });

    const uploadAvatarMutation = useMutation({
        mutationFn: uploadUserAvatar,
        onSuccess: () => {
            message.success('User Avatar updated successfully');

            setEditable(false);
        },
    });

    const onFinishUser: FormProps<UserFormFields>['onFinish'] = (values) => {
        const userData: UserFormFields = {
            userName: values.userName,
            altEmail: values.altEmail,
            email: values.email,
            altName: values.altName,
        };

        updateUserApi.mutate(userData);
    };

    const onFinishPassword: FormProps<PasswordFormFields>['onFinish'] = (values) => {
        const userData: PasswordFormFields = {
            password: values.password,
            oldPassword: values.oldPassword,
        };

        updateUserApi.mutate(userData);
    };

    const onFinishPref: FormProps<UserPrefFormFields>['onFinish'] = (values) => {
        const prefs: UserPrefFormFields = {
            currency: values.currency,
            theme: values.theme,
        };

        createUserPrefApi.mutate(prefs);
    };

    const filledOldPassword = Form.useWatch('oldPassword', form);
    const filledPassword = Form.useWatch('confirmPassword', form);

    return (
        <>
            <Row justify="space-between" gutter={16} style={{ padding: 12 }}>
                <Breadcrumb
                    items={[
                        {
                            href: '/home',
                            title: <HomeOutlined />,
                        },
                        {
                            title: `${useUserStore.getState().userName}'s Profile`,
                        },
                    ]}
                />
            </Row>
            <Card
                styles={{
                    body: {
                        position: 'relative',
                        padding: 32,
                    },
                }}
            >
                <Space orientation="vertical" size="large" style={{ width: '100%' }}>
                    <Flex vertical align="center" gap={12}>
                        <Badge
                            count={<Button shape="circle" icon={<EditOutlined />} onClick={() => setEditable((prev) => !prev)} />}
                            offset={[-10, 10]}
                        >
                            {editable ? (
                                <Upload
                                    showUploadList={false}
                                    accept="image/*"
                                    beforeUpload={(file) => {
                                        uploadAvatarMutation.mutate(file);
                                        return false;
                                    }}
                                >
                                    <Avatar size={120} src={`/v1/user/avatar/${useUserStore.getState().userId}`} style={{ cursor: 'pointer' }} />
                                </Upload>
                            ) : (
                                <Avatar size={120} src={`/v1/user/avatar/${useUserStore.getState().userId}`} />
                            )}
                        </Badge>

                        <Title level={2} style={{ margin: 0 }}>
                            My Profile
                        </Title>

                        <Text type="secondary">Manage your account details and preferences</Text>
                    </Flex>

                    <Card type="inner" title="Profile Information">
                        <Form<UserFormFields>
                            layout="vertical"
                            onFinish={onFinishUser}
                            disabled={!editable}
                            initialValues={{
                                userName: useUserStore.getState().userName,
                                altName: useUserStore.getState().altName,
                                email: useUserStore.getState().email,
                                altEmail: useUserStore.getState().altEmail,
                            }}
                        >
                            <Row gutter={24}>
                                <Col span={6}>
                                    <Form.Item label="Username" name="userName">
                                        <Input />
                                    </Form.Item>
                                </Col>

                                <Col span={6}>
                                    <Form.Item label="Nickname" name="altName">
                                        <Input placeholder="Enter Your Nickname" />
                                    </Form.Item>
                                </Col>

                                <Col span={6}>
                                    <Form.Item label="Email" name="email">
                                        <Input placeholder="Enter Your Email" />
                                    </Form.Item>
                                </Col>

                                <Col span={6}>
                                    <Form.Item label="Personal Email" name="altEmail">
                                        <Input placeholder="Enter Your Alt Email" />
                                    </Form.Item>
                                </Col>
                            </Row>

                            <Flex justify="flex-end">
                                <Button type="primary" htmlType="submit" loading={updateUserApi.isPending}>
                                    Save Profile
                                </Button>
                            </Flex>
                        </Form>
                    </Card>

                    <Card type="inner" title="Application Preferences">
                        <Form<UserPrefFormFields>
                            layout="vertical"
                            onFinish={onFinishPref}
                            disabled={!editable}
                            initialValues={{
                                currency: useUserPrefStore.getState().currency,
                                theme: useUserPrefStore.getState().theme,
                            }}
                        >
                            <Row gutter={24}>
                                <Col span={12}>
                                    <Form.Item label="Currency" name="currency">
                                        <Select options={currencies} />
                                    </Form.Item>
                                </Col>

                                <Col span={12}>
                                    <Form.Item label="Theme" name="theme">
                                        <Select options={appTheme} />
                                    </Form.Item>
                                </Col>
                            </Row>

                            <Flex justify="flex-end">
                                <Button type="primary" htmlType="submit" loading={createUserPrefApi.isPending}>
                                    Save Preferences
                                </Button>
                            </Flex>
                        </Form>
                    </Card>

                    <Card type="inner" title=" Change Password">
                        <Form<PasswordFormFields> layout="vertical" onFinish={onFinishPassword} form={form} disabled={!editable}>
                            <Row gutter={24}>
                                <Col span={8}>
                                    <Form.Item label="Old Password" name="oldPassword" required>
                                        <Input.Password placeholder="Enter your old password" />
                                    </Form.Item>
                                </Col>

                                <Col span={8}>
                                    <Form.Item
                                        label="New Password"
                                        name="confirmPassword"
                                        dependencies={['oldPassword']}
                                        rules={[
                                            {
                                                required: true,
                                                validator(_, value) {
                                                    if (value === filledOldPassword) {
                                                        return Promise.reject('Old Password is same as new Password');
                                                    }

                                                    return Promise.resolve();
                                                },
                                            },
                                        ]}
                                    >
                                        <Input.Password placeholder="Enter your new password" />
                                    </Form.Item>
                                </Col>

                                <Col span={8}>
                                    <Form.Item
                                        label="Confirm Password"
                                        name="password"
                                        dependencies={['confirmPassword', 'oldPassword']}
                                        rules={[
                                            {
                                                required: true,
                                                validator(_, value) {
                                                    if (filledOldPassword === value || filledOldPassword === filledPassword) {
                                                        return Promise.reject('Old Password is same as new Password');
                                                    }

                                                    if (!filledPassword) {
                                                        return Promise.reject('password field is not filled');
                                                    }

                                                    if (value !== filledPassword) {
                                                        return Promise.reject('password does not match');
                                                    }

                                                    return Promise.resolve();
                                                },
                                            },
                                        ]}
                                    >
                                        <Input.Password placeholder="Confirm your new password" />
                                    </Form.Item>
                                </Col>
                            </Row>

                            <Flex justify="flex-end">
                                <Button type="primary" htmlType="submit">
                                    Change Password
                                </Button>
                            </Flex>
                        </Form>
                    </Card>
                </Space>
            </Card>
        </>
    );
};
