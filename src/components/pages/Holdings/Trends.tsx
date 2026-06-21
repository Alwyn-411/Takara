import { useQuery } from '@tanstack/react-query';
import { useUserStore } from '../../../store/User';
import { getHoldingById, getHoldingValuations } from '../../../api/holdings';
import { useNavigate, useParams } from 'react-router-dom';
import { Badge, Breadcrumb, Button, Card, Col, Divider, Empty, Row, Space, Spin, Statistic, Table, Tag, Typography } from 'antd';
import dayjs from 'dayjs';
import { assetKindOptions, HoldingType, liabilityKindOptions, type HoldingValuation } from '../../../types/Holdings';
import type { ColumnsType } from 'antd/es/table';
import { HomeOutlined, ReloadOutlined } from '@ant-design/icons';
import { currencies } from '../../../types/Accounts';
import './Trends.css';

const { Text, Title } = Typography;

export const HoldingTrends = () => {
    const userId = useUserStore.getState().userId;
    const navigate = useNavigate();

    const { holdingId } = useParams<{ holdingId: string }>();

    const HoldingValuationQuery = useQuery({
        queryKey: ['HoldingTrends', userId, holdingId],
        queryFn: () => getHoldingValuations({ holdingId }),
        enabled: !!userId && !!holdingId,
    });

    const HoldingQuery = useQuery({
        queryKey: ['Holding', userId, holdingId],
        queryFn: () => getHoldingById({ holdingId }),
        enabled: !!userId && !!holdingId,
    });

    const valuationColumns: ColumnsType<HoldingValuation> = [
        {
            title: 'Date',
            dataIndex: 'observedAt',
            render: (value) => dayjs.unix(value).format('DD MMM YYYY'),
            defaultSortOrder: 'descend',
            sorter: (a, b) => a.observedAt - b.observedAt,
        },
        {
            title: 'Value',
            dataIndex: 'value',
            align: 'right',
            render: (value) =>
                new Intl.NumberFormat('en-US', {
                    minimumFractionDigits: 2,
                    maximumFractionDigits: 2,
                }).format(value),
        },
        {
            title: 'Quantity',
            dataIndex: 'quantity',
            align: 'right',
            render: (value) => value ?? '-',
        },
        {
            title: 'Unit Price',
            dataIndex: 'unitPrice',
            align: 'right',
            render: (value) => value ?? '-',
        },
        {
            title: 'Note',
            dataIndex: 'note',
            render: (value) => (value ? <Tag color="blue">{value}</Tag> : <Text type="secondary">-</Text>),
        },
    ];

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
                            title: 'Holdings',
                            href: '/holdings',
                        },
                        {
                            title: 'Valuation',
                        },
                    ]}
                />
            </Row>
            {HoldingQuery.data && (
                <Space orientation="vertical" size="large" style={{ width: '100%' }}>
                    <Card loading={HoldingQuery.isLoading} className={HoldingQuery.data.type === HoldingType.Asset ? 'asset-card' : 'liability-card'}>
                        <Row justify="space-between" align="middle">
                            <Col span={24}>
                                <Space orientation="vertical" size={0}>
                                    <Title level={2} style={{ margin: 0 }}>
                                        {HoldingQuery.data.name}
                                    </Title>
                                </Space>
                            </Col>
                        </Row>

                        <Divider />

                        <Row gutter={[16, 16]}>
                            <Col span={6}>
                                <Statistic
                                    title="Holding Type"
                                    value={
                                        HoldingQuery.data.type === HoldingType.Asset
                                            ? assetKindOptions.find((x) => x.value === HoldingQuery.data.kind)?.label
                                            : liabilityKindOptions.find((x) => x.value === HoldingQuery.data.kind)?.label
                                    }
                                />
                            </Col>
                            <Col span={6}>
                                {HoldingValuationQuery.data && (
                                    <Statistic
                                        title="Latest Value"
                                        prefix={currencies.find((currency) => currency.value === HoldingQuery.data.currency)?.symbol}
                                        suffix={currencies.find((currency) => currency.value === HoldingQuery.data.currency)?.value}
                                        value={HoldingValuationQuery.data.records[HoldingValuationQuery.data.records.length - 1].value ?? 0}
                                        precision={2}
                                    />
                                )}
                            </Col>
                            <Col span={6}>
                                <Statistic title="Opened On" value={dayjs.unix(HoldingQuery.data.openedAt).format('DD MMM YYYY')} />
                            </Col>
                            <Col span={6}>
                                <Statistic
                                    title="Last Valuation"
                                    value={
                                        HoldingValuationQuery.data?.records?.length
                                            ? dayjs
                                                  .unix(HoldingValuationQuery.data.records[HoldingValuationQuery.data.records.length - 1].observedAt)
                                                  .format('DD MMM YYYY')
                                            : '-'
                                    }
                                />
                            </Col>
                        </Row>
                    </Card>

                    <Card
                        title={
                            <Space>
                                <Title
                                    level={4}
                                    style={{
                                        margin: 0,
                                    }}
                                >
                                    Valuation History
                                </Title>

                                <Badge count={HoldingValuationQuery.data?.records?.length ?? 0} color={'green'} />
                            </Space>
                        }
                        extra={
                            <Space>
                                <Button
                                    icon={<ReloadOutlined />}
                                    onClick={() => {
                                        HoldingValuationQuery.refetch();
                                        HoldingQuery.refetch();
                                    }}
                                />

                                <Button type="primary" onClick={() => navigate('./add')}>
                                    Add Valuation
                                </Button>
                            </Space>
                        }
                    >
                        {HoldingValuationQuery.isLoading ? (
                            <Row justify="center">
                                <Spin size="large" />
                            </Row>
                        ) : (
                            <Table
                                rowKey={(record) => `${record.observedAt}-${record.value}`}
                                columns={valuationColumns}
                                dataSource={HoldingValuationQuery.data?.records ?? []}
                                size="middle"
                                pagination={false}
                            />
                        )}
                    </Card>
                </Space>
            )}
        </>
    );
};
