import { useNavigate, useParams } from 'react-router-dom';
import { useUserStore } from '../../../store/User';
import { Breadcrumb, Button, Card, Row, Space, Statistic, Typography } from 'antd';
import { EditOutlined, LeftOutlined, PlusCircleOutlined, PlusOutlined, RightOutlined } from '@ant-design/icons';
import { getAccountWithAccountId } from '../../../api/accounts';
import { useQuery } from '@tanstack/react-query';
import { currencies, type Accounts } from '../../../types/Accounts';
import type { valueType } from 'antd/es/statistic/utils';

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

    const onEdit = () => {
        navigate(`/accounts/${accountId}/edit`);
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
                            <Title level={3} italic type="secondary">
                                My Account Details
                            </Title>
                        }
                        style={{ width: '100%' }}
                        extra={<Button onClick={onEdit} icon={<EditOutlined />} />}
                    >
                        <Row justify="space-between" style={{ padding: 8 }}>
                            <Statistic title="Bank Name" value={data.name} />
                            <Statistic title="Account Number" value={data.accountNumber} formatter={maskedCardNumber} />
                            <Statistic title="Account Type" value={data.type} />
                            <Statistic title="Currency" value={`${getAccountCurrency(data)?.value} - ${getAccountCurrency(data)?.symbol}`} />
                            {data.type === 'Savings' && <Statistic title="Interest" value={data.interest} suffix={'%'} />}
                        </Row>
                        <Row justify="space-between" style={{ padding: 8 }}>
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
                    </Card>
                </Row>
            )}
            {isSuccess && !!data && (
                <>
                    <Row justify="space-between" gutter={16} style={{ padding: 12 }}>
                        <Text strong italic style={{ fontSize: 22 }}>
                            Transactions
                        </Text>
                        <Button icon={<PlusOutlined />} color={'primary'}>
                            Add
                        </Button>
                    </Row>
                    <Row></Row>
                </>
            )}
        </>
    );
};
