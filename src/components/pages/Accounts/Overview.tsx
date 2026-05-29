import { DeleteOutlined } from '@ant-design/icons';
import { useMutation, useQuery } from '@tanstack/react-query';
import { Badge, Button, Card, Divider, Empty, message, Popconfirm, Row, Space, Spin, Table, Typography, type TableProps } from 'antd';
import { useNavigate } from 'react-router-dom';
import { deleteAccountWithAccountId, listAccountsWithUserId } from '../../../api/accounts';
import { useUserStore } from '../../../store/User';
import { currencies, type Accounts as AccountDataType } from '../../../types/Accounts';
const { Title, Text } = Typography;

interface AccountsTableProps extends Omit<AccountDataType, 'userId'> {}

export const Accounts = () => {
    const navigate = useNavigate();
    const userId = useUserStore.getState().userId;

    const { data, refetch, isLoading, isSuccess } = useQuery({
        queryKey: ['accounts', userId],
        queryFn: listAccountsWithUserId,
        enabled: !!userId,
    });

    const { mutate } = useMutation({
        mutationFn: deleteAccountWithAccountId,
        onSuccess: () => {
            message.success('Deleted Successfully');
            refetch();
        },
        onError: (error) => {
            message.error(error.message);
            console.error(error.message);
        },
    });

    const columns: TableProps<AccountsTableProps>['columns'] = [
        {
            dataIndex: 'name',
            key: 'name',
            title: 'Bank Name',
        },
        {
            dataIndex: 'accountNumber',
            key: 'accountNumber',
            title: 'Account Number',
        },
        {
            dataIndex: 'type',
            key: 'type',
            title: 'Type',
        },
        {
            dataIndex: 'balance',
            key: 'balance',
            title: 'Balance',
            render: (value, record) => {
                const currency = currencies.find((item) => item.value === record.currency);
                return (
                    <Space size="small">
                        <Text>{currency?.symbol}</Text>
                        <Text>{value}</Text>
                        <Text>{currency?.value}</Text>
                    </Space>
                );
            },
        },
        {
            dataIndex: '',
            key: 'actions',
            title: 'Actions',
            render: (_, record) => (
                <Space>
                    <Button
                        type="link"
                        onClick={() => {
                            navigate(`./${record.accountId}/details`);
                        }}
                    >
                        View
                    </Button>
                    <Popconfirm
                        title="Are you sure ?"
                        onConfirm={() => {
                            mutate(record.accountId);
                        }}
                    >
                        <Button type="link" icon={<DeleteOutlined />} />
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    return (
        <>
            <Card>
                <Row justify={'space-between'}>
                    <Space>
                        <Title level={2} style={{ margin: 0 }} type="secondary" italic>
                            My Bank Accounts
                        </Title>
                        <Badge count={data?.count ?? 0} color={'green'} />
                    </Space>
                    <Space>
                        <Button
                            type="primary"
                            onClick={() => {
                                navigate('./create');
                            }}
                        >
                            Create Bank Account
                        </Button>
                    </Space>
                </Row>
                <Divider />
                <Row justify="center">
                    {isLoading && <Spin />}
                    {isSuccess && (
                        <Table
                            style={{ width: '100%' }}
                            columns={columns}
                            dataSource={data.records}
                            pagination={false}
                            locale={{
                                emptyText: (
                                    <Empty description="No Data">
                                        <Button
                                            type="primary"
                                            onClick={() => {
                                                navigate('./create');
                                            }}
                                        >
                                            Create Bank Account
                                        </Button>
                                    </Empty>
                                ),
                            }}
                        />
                    )}
                </Row>
            </Card>
        </>
    );
};
