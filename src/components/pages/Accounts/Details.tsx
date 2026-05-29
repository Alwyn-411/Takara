import { DeleteOutlined, EditOutlined, HomeOutlined, LeftOutlined, PlusOutlined, RightOutlined, SyncOutlined } from '@ant-design/icons';
import { useMutation, useQuery } from '@tanstack/react-query';
import type { TableColumnsType } from 'antd';
import { Alert, Breadcrumb, Button, Card, message, Popconfirm, Row, Space, Statistic, Table, Tag, Typography } from 'antd';
import type { valueType } from 'antd/es/statistic/utils';
import { useNavigate, useParams } from 'react-router-dom';
import { getAccountWithAccountId, listAccountsWithUserId } from '../../../api/accounts';
import { deleteTransactionWithTransactionId, listTransactionsByAccountId } from '../../../api/transactions';
import { useUserStore } from '../../../store/User';
import { currencies } from '../../../types/Accounts';

const { Title, Text } = Typography;

export const AccountDetails = () => {
    const { accountId } = useParams<{ accountId: string }>();
    const { userId } = useUserStore();
    const navigate = useNavigate();

    const listAccountsQuery = useQuery({
        queryKey: ['listAccounts', userId],
        queryFn: listAccountsWithUserId,
        enabled: !!userId,
    });

    const accountRecords = listAccountsQuery.data?.records ?? [];
    const currentIndex = accountRecords.findIndex((v) => v.accountId === accountId);

    const moveToPrevAccount = () => {
        const prev = accountRecords[currentIndex - 1];
        if (!prev) return;
        navigate(`../../${prev.accountId}/details`, { relative: 'path' });
    };

    const moveToNextAccount = () => {
        const next = accountRecords[currentIndex + 1];
        if (!next) return;
        navigate(`../../${next.accountId}/details`, { relative: 'path' });
    };

    const accountsQuery = useQuery({
        queryFn: () => getAccountWithAccountId(accountId!!),
        queryKey: ['accounts', accountId],
        enabled: !!accountId && !!userId,
    });

    const transactionsQuery = useQuery({
        queryFn: () => listTransactionsByAccountId(accountId!!),
        queryKey: ['listTransactions', accountId],
        enabled: !!accountId && !!userId,
    });

    const { mutate, isError } = useMutation({
        mutationFn: deleteTransactionWithTransactionId,
        onSuccess: () => {
            message.success('Deleted Successfully');
            transactionsQuery.refetch();
            accountsQuery.refetch();
        },
    });

    const records = transactionsQuery.data?.records ?? [];

    const columns: TableColumnsType<(typeof records)[number]> = [
        {
            title: 'Date',
            dataIndex: 'transactionAt',
            key: 'transactionAt',
            defaultSortOrder: 'descend',
            sorter: (a, b) => a.transactionAt - b.transactionAt,
            render: (ts: number) => new Date(ts * 1000).toLocaleString('en-IN', { dateStyle: 'medium', timeStyle: 'short' }),
        },
        {
            title: 'Merchant',
            dataIndex: 'merchantName',
            key: 'merchantName',
            render: (name: string) => name || 'Unknown',
        },
        {
            title: 'Type',
            dataIndex: 'type',
            key: 'type',
            filters: [
                { text: 'Debit', value: 'Debit' },
                { text: 'Credit', value: 'Credit' },
            ],
            onFilter: (value, record) => record.type === value,
            render: (type: string) => <Tag color={type === 'Debit' ? 'red' : 'green'}>{type}</Tag>,
        },
        {
            title: 'Amount',
            dataIndex: 'settledAmount',
            key: 'settledAmount',
            align: 'right',
            sorter: (a, b) => Number(a.settledAmount) - Number(b.settledAmount),
            render: (amount: string, record) => {
                const isDebit = record.type === 'Debit';
                const symbol = getAccountCurrency(record.settledCurrency)?.symbol;
                return (
                    <Text type={isDebit ? 'danger' : 'success'} strong>
                        {isDebit ? '-' : '+'}
                        {symbol}
                        {amount}
                    </Text>
                );
            },
        },
        {
            title: 'Actions',
            dataIndex: '',
            key: 'actions',
            align: 'right',
            render: (_, record) => {
                return (
                    <Space>
                        <Button type="link">View</Button>
                        <Popconfirm title="Are you sure ?" okText="Yes" cancelText="No" onConfirm={() => mutate(record.transactionId)}>
                            <Button type="link" icon={<DeleteOutlined />} />
                        </Popconfirm>
                    </Space>
                );
            },
        },
    ];

    const EditAccount = () => {
        navigate(`../../${accountId}/edit`, { relative: 'path' });
    };

    const AddTransaction = () => {
        navigate(`../transactions/create`, { relative: 'path' });
    };

    function getAccountCurrency(currency: string) {
        return currencies.find(({ value }) => {
            return value === currency;
        });
    }

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
                {listAccountsQuery.isSuccess && listAccountsQuery.data.count >= 2 && (
                    <Space>
                        <Button icon={<LeftOutlined />} onClick={moveToPrevAccount} />
                        <Button icon={<RightOutlined />} onClick={moveToNextAccount} />
                    </Space>
                )}
            </Row>
            {accountsQuery.isSuccess && !!accountsQuery.data && (
                <Row gutter={16} style={{ padding: 12 }}>
                    <Card
                        variant="borderless"
                        loading={accountsQuery.isLoading}
                        title={
                            <Title level={3} italic>
                                {accountsQuery.data.name}
                            </Title>
                        }
                        style={{ width: '100%' }}
                        extra={<Button onClick={EditAccount} type="link" icon={<EditOutlined />} />}
                    >
                        <Row justify="space-between" style={{ padding: 8 }}>
                            <Statistic title="Account Number" value={accountsQuery.data.accountNumber} formatter={maskedCardNumber} />
                            <Statistic title="Account Type" value={accountsQuery.data.type} />
                            <Statistic
                                title="Currency"
                                value={`${getAccountCurrency(accountsQuery.data.currency)?.value} - ${getAccountCurrency(accountsQuery.data.currency)?.symbol}`}
                            />
                            {accountsQuery.data.type === 'Savings' && <Statistic title="Interest" value={accountsQuery.data.interest} suffix={'%'} />}
                            <Statistic
                                title="Balance Amount"
                                value={accountsQuery.data.balance}
                                prefix={getAccountCurrency(accountsQuery.data.currency)?.symbol}
                                suffix={
                                    <Text style={{ fontSize: 20 }} italic>
                                        {getAccountCurrency(accountsQuery.data.currency)?.value}
                                    </Text>
                                }
                            />
                        </Row>
                        <Row>
                            {!!accountsQuery.data.description && (
                                <Statistic title="Description" valueRender={() => <Text>{accountsQuery.data.description}</Text>} />
                            )}
                        </Row>
                    </Card>
                </Row>
            )}
            {isError && (
                <Alert
                    title="Account Edit Failed"
                    description={<Text>Error occured while processing this request</Text>}
                    type="error"
                    showIcon
                    style={{ padding: 12 }}
                />
            )}
            {accountsQuery.isSuccess && !!accountsQuery.data && (
                <>
                    <Row justify="space-between" gutter={16} style={{ padding: 12 }}>
                        <Text strong italic style={{ fontSize: 22 }}>
                            All Transactions
                        </Text>
                        <Space>
                            <Button icon={<SyncOutlined />}>Sync</Button>
                            <Button icon={<PlusOutlined />} type="primary" onClick={AddTransaction}>
                                Add
                            </Button>
                        </Space>
                    </Row>

                    <Row gutter={16} style={{ padding: 12 }}>
                        <Table
                            rowKey="transactionId"
                            columns={columns}
                            dataSource={records}
                            loading={transactionsQuery.isLoading}
                            pagination={{ pageSize: 20, placement: ['bottomCenter'] }}
                            style={{ width: '100%' }}
                        />
                    </Row>
                </>
            )}
        </>
    );
};
