import {
    Alert,
    AutoComplete,
    Breadcrumb,
    Button,
    Card,
    Col,
    DatePicker,
    Form,
    Input,
    InputNumber,
    message,
    Row,
    Segmented,
    Select,
    Space,
    Typography,
} from 'antd';
import { useMutation, useQuery } from '@tanstack/react-query';
import { useNavigate, useParams } from 'react-router-dom';
import { getAccountWithAccountId } from '../../../api/accounts';
import { useUserStore } from '../../../store/User';
import { HomeOutlined } from '@ant-design/icons';
import { TransactionType, type Transaction } from '../../../types/Transactions';
import { currencies } from '../../../types/Accounts';
import { createTransaction, listCategories, listMerchants, listTag } from '../../../api/transactions';
import { useState } from 'react';
import { useDebounce } from '../../core/Debounce';

export interface Formfields extends Partial<Omit<Transaction, 'accountId' | 'transactionId' | 'userId' | 'active' | 'updatedAt' | 'createdAt'>> {}

const { Text } = Typography;
const { TextArea } = Input;

export const TransactionsCreate = () => {
    const navigate = useNavigate();
    const [form] = Form.useForm();
    const userId = useUserStore.getState().userId;
    const { accountId } = useParams<{ accountId: string }>();

    const { data } = useQuery({
        queryFn: () => getAccountWithAccountId(accountId!!),
        queryKey: ['accounts', accountId],
        enabled: !!accountId && !!userId,
    });

    const [tagSearch, setTagSearch] = useState('');
    const debouncedTagSearch = useDebounce(tagSearch, 150);

    const tagOptionsQuery = useQuery({
        queryFn: () => listTag(debouncedTagSearch),
        queryKey: ['tags', userId, debouncedTagSearch],
        enabled: !!userId,
    });

    const [categorySearch, setCategorySearch] = useState('');
    const debouncedCategorySearch = useDebounce(categorySearch, 150);

    const categoryOptionsQuery = useQuery({
        queryFn: () => listCategories(debouncedCategorySearch),
        queryKey: ['categories', userId, debouncedCategorySearch],
        enabled: !!userId,
    });

    const [merchantSearch, setMerchantSearch] = useState('');
    const debouncedMerchantSearch = useDebounce(merchantSearch, 150);

    const merchantOptionsQuery = useQuery({
        queryFn: () => listMerchants(debouncedMerchantSearch),
        queryKey: ['merchants', userId, debouncedMerchantSearch],
        enabled: !!userId,
    });

    const { mutate, isPending, isError } = useMutation({
        mutationFn: createTransaction,
        onSuccess: () => {
            message.success('Creation Successfully');
            navigate(-1);
        },
    });

    const onFinish = (values: Formfields) => {
        const payload = {
            userId,
            accountId,
            type: values.type,
            settledAmount: String(values.settledAmount),
            settledCurrency: values.settledCurrency,
            merchantName: values.merchantName ?? '',
            categoryName: values.categoryName ?? '',
            description: values.description ?? '',
            tagNames: values.tags ?? [],
            transactionAt: (values.transactionAt as any)?.unix(),
        };

        mutate(payload);
    };

    const selectedType: string = Form.useWatch('type', form);
    const selectedCurrency: string = Form.useWatch('settledCurrency', form);
    const currencyObj = currencies.find((c) => c.value === selectedCurrency);

    const options = [
        {
            label: <Text type="danger">Debit</Text>,
            value: TransactionType.Debit,
        },
        {
            label: <Text type="success">Credit</Text>,
            value: TransactionType.Credit,
        },
    ];

    const merchantOptions = (merchantOptionsQuery.data?.records ?? []).map((m) => ({
        value: m.merchantName,
        label: m.merchantName,
    }));
    const categoryOptions = (categoryOptionsQuery.data?.records ?? []).map((m) => ({
        value: m.categoryName,
        label: m.categoryName,
    }));
    const tagOptions = (tagOptionsQuery.data?.records ?? []).map((m) => ({
        value: m.tagName,
        label: m.tagName,
    }));

    return (
        <Col span={24}>
            <Row justify="space-between" gutter={16} style={{ padding: 12 }}>
                <Breadcrumb
                    items={[
                        {
                            href: '/home',
                            title: <HomeOutlined />,
                        },
                        {
                            title: 'Accounts',
                            href: '/accounts',
                        },
                        {
                            title: data?.name,
                            href: '../details',
                        },
                        {
                            title: 'Create New Transaction',
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
                            style={{ padding: 8 }}
                        />
                    )}
                    <Form<Formfields>
                        form={form}
                        layout="vertical"
                        size="large"
                        onFinish={onFinish}
                        requiredMark="optional"
                        initialValues={{
                            type: TransactionType.Debit,
                        }}
                    >
                        <Form.Item name="type" required>
                            <Segmented options={options} block />
                        </Form.Item>

                        {selectedType === TransactionType.Debit && (
                            <Row gutter={16} align="middle">
                                <Col span={4}>
                                    <Form.Item
                                        label="Currency"
                                        name="settledCurrency"
                                        required
                                        rules={[{ required: true, message: 'Select a Currency' }]}
                                    >
                                        <Select options={currencies} placeholder="Currency" />
                                    </Form.Item>
                                </Col>
                                <Col span={5}>
                                    <Form.Item label="Amount" name="settledAmount" required rules={[{ required: true, message: 'Enter an Amount' }]}>
                                        <InputNumber
                                            disabled={!selectedCurrency}
                                            prefix={currencyObj?.symbol}
                                            placeholder="Amount"
                                            suffix={currencyObj?.value}
                                            style={{ width: '100%' }}
                                        />
                                    </Form.Item>
                                </Col>
                                <Col span={7}>
                                    <Form.Item
                                        label="Paid to"
                                        name="merchantName"
                                        required
                                        rules={[{ required: true, message: 'Enter merchant name' }]}
                                    >
                                        <AutoComplete
                                            options={merchantOptions}
                                            placeholder="Reciever Name"
                                            showSearch={{ onSearch: setMerchantSearch }}
                                        />
                                    </Form.Item>
                                </Col>
                                <Col span={8}>
                                    <Form.Item
                                        label="On"
                                        name="transactionAt"
                                        required
                                        rules={[{ required: true, message: 'Select Transaction Time' }]}
                                    >
                                        <DatePicker showTime use12Hours style={{ width: '100%' }} />
                                    </Form.Item>
                                </Col>
                            </Row>
                        )}

                        {selectedType !== TransactionType.Debit && (
                            <Row gutter={16} align="middle">
                                <Col span={7}>
                                    <Form.Item
                                        label="Sent By"
                                        name="merchantName"
                                        required
                                        rules={[
                                            { required: true, message: 'Enter merchant name' },
                                            { min: 3, message: 'Minimum 3 characters' },
                                        ]}
                                    >
                                        <AutoComplete
                                            options={merchantOptions}
                                            placeholder="Sender Name"
                                            showSearch={{
                                                onSearch: setMerchantSearch,
                                            }}
                                        />
                                    </Form.Item>
                                </Col>
                                <Col span={4}>
                                    <Form.Item
                                        label="Currency"
                                        name="settledCurrency"
                                        required
                                        rules={[{ required: true, message: 'Select a Currency' }]}
                                    >
                                        <Select options={currencies} placeholder="Currency" />
                                    </Form.Item>
                                </Col>
                                <Col span={5}>
                                    <Form.Item label="Amount" name="settledAmount" required rules={[{ required: true, message: 'Enter an Amount' }]}>
                                        <InputNumber
                                            disabled={!selectedCurrency}
                                            prefix={currencyObj?.symbol}
                                            placeholder="Amount"
                                            suffix={currencyObj?.value}
                                            style={{ width: '100%' }}
                                        />
                                    </Form.Item>
                                </Col>
                                <Col span={8}>
                                    <Form.Item
                                        label="On"
                                        name="transactionAt"
                                        required
                                        rules={[{ required: true, message: 'Select Transaction Time' }]}
                                    >
                                        <DatePicker showTime use12Hours style={{ width: '100%' }} />
                                    </Form.Item>
                                </Col>
                            </Row>
                        )}

                        <Row>
                            <Col span={24}>
                                <Form.Item label="Description" name="description">
                                    <TextArea rows={4} placeholder="Enter a short description about this transaction" />
                                </Form.Item>
                            </Col>
                        </Row>

                        <Row gutter={16}>
                            <Col span={12}>
                                <Form.Item label="Category" name="categoryName">
                                    <AutoComplete
                                        options={categoryOptions}
                                        placeholder="Select a category"
                                        showSearch={{
                                            onSearch: setCategorySearch,
                                        }}
                                    />
                                </Form.Item>
                            </Col>
                            <Col span={12}>
                                <Form.Item label="Tags" name="tags">
                                    <Select
                                        mode="tags"
                                        allowClear
                                        placeholder="Select transaction tags"
                                        onChange={() => {}}
                                        showSearch={{
                                            onSearch: setTagSearch,
                                        }}
                                        options={tagOptions}
                                    />
                                </Form.Item>
                            </Col>
                        </Row>
                        <Row justify="end">
                            <Space>
                                <Button onClick={() => form.resetFields()}>Reset</Button>
                                <Button type="primary" htmlType="submit" size="large" loading={isPending}>
                                    Create Transaction
                                </Button>
                            </Space>
                        </Row>
                    </Form>
                </Card>
            </Row>
        </Col>
    );
};
