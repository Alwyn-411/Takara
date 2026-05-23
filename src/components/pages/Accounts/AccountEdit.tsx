import { Alert, Breadcrumb, Button, Card, Col, Form, Input, InputNumber, message, Radio, Row, Select, Typography, type FormProps } from 'antd';
import { useParams, useNavigate } from 'react-router-dom';
import { useMutation, useQuery } from '@tanstack/react-query';
import { editAccountWithAccountId, getAccountWithAccountId } from '../../../hooks/accounts';
import { currencies, type Accounts } from '../../../types/Accounts';
import { useUserStore } from '../../../store/User';
import { options, type Formfields } from './AccountCreate';
import { useEffect } from 'react';

const { TextArea } = Input;
const { Text } = Typography;

export const AccountEdit = () => {
    const { accountId } = useParams<{ accountId: string }>();
    const { userId } = useUserStore();

    const navigate = useNavigate();
    const [form] = Form.useForm();

    const { data, isFetching } = useQuery({
        queryFn: () => getAccountWithAccountId(accountId!!),
        queryKey: ['accounts', accountId],
        enabled: !!accountId && !!userId,
    });

    const { mutate, isError, isPending } = useMutation({
        mutationFn: editAccountWithAccountId,
        onSuccess: () => {
            message.success('Edited Successfully');
            navigate(-1);
        },
    });

    const selectedCurrency: string = Form.useWatch('currency', form);

    const selectedType: string = Form.useWatch('type', form);

    const currencyObj = currencies.find((c) => c.value === selectedCurrency);

    const onFinish: FormProps<Formfields>['onFinish'] = (values) => {
        const AccountData: Partial<Accounts> = {
            accountId: accountId!!,
            type: values.type,
            name: values.name,
            accountNumber: values.accountNumber,
            description: values.description,
            currency: values.currency,
        };

        if (!!values.currency) {
            AccountData.balance = values.balance?.toString();
        }

        if (values.type === 'Savings') {
            AccountData.interest = values.interest?.toString();
        }

        if (values.type === 'Current') {
            AccountData.interest = '0';
        }

        mutate({ accountId: accountId!!, account: AccountData });
    };

    useEffect(() => {
        if (!!data) form.setFieldsValue(data);
    }, [data]);

    return (
        <Col span={24}>
            <Row>
                <Breadcrumb
                    items={[
                        {
                            title: 'Accounts',
                            href: '/accounts',
                        },
                        {
                            title: 'Edit Bank Account',
                        },
                    ]}
                />
            </Row>
            <Row justify="center">
                <Card style={{ width: '100%' }}>
                    {isError && (
                        <Alert
                            title="Account Edit Failed"
                            description={<Text>Error occured while processing this request</Text>}
                            type="error"
                            showIcon
                        />
                    )}

                    <Form<Formfields> form={form} layout="vertical" size="large" onFinish={onFinish} requiredMark="optional">
                        <Row gutter={16}>
                            <Col span={12}>
                                <Form.Item
                                    label="Name"
                                    name="name"
                                    required
                                    rules={[
                                        { required: true, message: 'Enter a username' },
                                        { min: 3, message: 'Minimum 3 characters' },
                                    ]}
                                >
                                    <Input placeholder="Enter the Name of the Bank Account" />
                                </Form.Item>
                            </Col>
                            <Col span={12}>
                                <Form.Item
                                    label="Account Number"
                                    name="accountNumber"
                                    rules={[
                                        { required: true, message: 'Enter Your Account Number' },
                                        { min: 15, message: 'Minimum 15 characters' },
                                    ]}
                                >
                                    <Input placeholder="Enter the Bank Account Number" />
                                </Form.Item>
                            </Col>
                        </Row>

                        <Form.Item label="Description" name="description">
                            <TextArea rows={4} placeholder="Enter a short description about this account" />
                        </Form.Item>

                        <Row gutter={16}>
                            <Col span={12}>
                                <Form.Item label="Currency" name="currency" rules={[{ required: true, message: 'Select a Currency' }]}>
                                    <Select options={currencies} placeholder="Select a Currency" />
                                </Form.Item>
                            </Col>
                            <Col span={12}>
                                <Form.Item
                                    label="Balance"
                                    name="balance"
                                    dependencies={['currency']}
                                    rules={[{ required: true, message: 'Enter Your Bank Balance' }]}
                                >
                                    <InputNumber
                                        style={{ width: '100%' }}
                                        disabled={!selectedCurrency}
                                        prefix={currencyObj?.symbol}
                                        placeholder="Enter Your Bank Balance"
                                        suffix={currencyObj?.value}
                                    />
                                </Form.Item>
                            </Col>
                        </Row>
                        <Row gutter={16}>
                            <Col span={12}>
                                <Form.Item label="Account Type" name="type" rules={[{ required: true, message: 'Choose your account Type' }]}>
                                    <Radio.Group block options={options} defaultValue="Savings" optionType="button" buttonStyle="solid" />
                                </Form.Item>
                            </Col>
                            <Col span={12}>
                                <Form.Item
                                    label="Interest"
                                    name="interest"
                                    dependencies={['type']}
                                    rules={[
                                        ({}) => ({
                                            validator(_, value) {
                                                if (selectedType === 'Savings' && !value) return Promise.reject(new Error('Interest is Required'));

                                                return Promise.resolve();
                                            },
                                        }),
                                    ]}
                                >
                                    <InputNumber
                                        disabled={selectedType === 'Current'}
                                        style={{ width: '100%' }}
                                        placeholder="Enter Your Bank's Interest Rate"
                                        suffix="%"
                                    />
                                </Form.Item>
                            </Col>
                        </Row>
                        <Row justify="end">
                            <Button type="primary" htmlType="submit" size="large" loading={isPending || isFetching}>
                                Edit Bank Account
                            </Button>
                        </Row>
                    </Form>
                </Card>
            </Row>
        </Col>
    );
};
