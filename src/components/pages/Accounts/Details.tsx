import { useNavigate, useParams } from 'react-router-dom';
import { useUserStore } from '../../../store/User';
import { Breadcrumb, Button, Card, Row, Space, Statistic, Timeline, Typography } from 'antd';
import { EditOutlined, HomeOutlined, LeftOutlined, PlusOutlined, RightOutlined } from '@ant-design/icons';
import { getAccountWithAccountId } from '../../../api/accounts';
import { useQuery } from '@tanstack/react-query';
import { currencies, type Accounts } from '../../../types/Accounts';
import type { valueType } from 'antd/es/statistic/utils';
import { listTransactionsByAccountId } from '../../../api/transactions';

const { Title, Text } = Typography;

export const AccountDetails = () => {
    const { accountId } = useParams<{ accountId: string }>();
    const { userId } = useUserStore();
    const navigate = useNavigate();

    const { data, isLoading, isSuccess } = useQuery({
        queryFn: () => getAccountWithAccountId(accountId!!),
        queryKey: ['accounts', accountId],
        enabled: !!accountId && !!userId,
    });

    // listTransactionsByAccountId
    const transactionsQuery = useQuery({
        queryFn: () => listTransactionsByAccountId(accountId!!),
        queryKey: ['listTransactions', accountId],
        enabled: !!accountId && !!userId,
    });

    const timelineItems = transactionsQuery.data?.records.forEach((value) => {
        return {
            title: value.transactionAt,
            content: `${value.settledAmount} ${value.type === 'Debit' ? 'Paid to' : 'Credited from'} ${value.merchant}`,
        };
    });

    const EditAccount = () => {
        navigate(`../../${accountId}/edit`, { relative: 'path' });
    };

    const AddTransaction = () => {
        navigate(`../transactions/create`, { relative: 'path' });
    };

    const getAccountCurrency = (data: Accounts) => {
        return currencies.find(({ value }) => {
            return value === data.currency;
        });
    };

    // const unmaskedCardNumber = (value: valueType) => {
    //     if (!value) return 'N/A';

    //     const str = String(value).replace(/\D/g, '');

    //     return str.match(/.{1,4}/g)?.join(' ') || str;
    // };

    const maskedCardNumber = (value: valueType) => {
        if (!value) return 'N/A';

        const str = String(value).replace(/\D/g, '');

        if (str.length <= 4) return str;

        const last4 = str.slice(-4);

        return `**** **** **** ${last4}`;
    };

    return (
        <>
            <Row justify="space-between" gutter={16} style={{ padding: 12 }}>
                <Space>
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
                                title: 'View Account Details',
                            },
                        ]}
                    />
                </Space>
                <Space>
                    <Button icon={<LeftOutlined />} />
                    <Button icon={<RightOutlined />} />
                </Space>
            </Row>
            {isSuccess && !!data && (
                <Row gutter={16} style={{ padding: 12 }}>
                    <Card
                        variant="borderless"
                        loading={isLoading}
                        title={
                            <Title level={3} italic>
                                {data.name}
                            </Title>
                        }
                        style={{ width: '100%' }}
                        extra={<Button onClick={EditAccount} type="link" icon={<EditOutlined />} />}
                    >
                        <Row justify="space-between" style={{ padding: 8 }}>
                            <Statistic title="Account Number" value={data.accountNumber} formatter={maskedCardNumber} />
                            <Statistic title="Account Type" value={data.type} />
                            <Statistic title="Currency" value={`${getAccountCurrency(data)?.value} - ${getAccountCurrency(data)?.symbol}`} />
                            {data.type === 'Savings' && <Statistic title="Interest" value={data.interest} suffix={'%'} />}
                            <Statistic
                                title="Balance Amount"
                                value={data.balance}
                                prefix={getAccountCurrency(data)?.symbol}
                                suffix={
                                    <Text style={{ fontSize: 20 }} italic>
                                        {getAccountCurrency(data)?.value}
                                    </Text>
                                }
                            />
                        </Row>
                        <Row>{!!data.description && <Statistic title="Description" valueRender={() => <Text>{data.description}</Text>} />}</Row>
                    </Card>
                </Row>
            )}
            {isSuccess && !!data && (
                <>
                    <Row justify="space-between" gutter={16} style={{ padding: 12 }}>
                        <Text strong italic style={{ fontSize: 22 }}>
                            Transactions
                        </Text>
                        <Button icon={<PlusOutlined />} type="primary" onClick={AddTransaction}>
                            Add
                        </Button>
                    </Row>
                    <Row>
                        <Timeline mode="start" items={timelineItems ?? []} />
                    </Row>
                </>
            )}
        </>
    );
};
