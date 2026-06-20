import { Badge, Button, Card, Col, Divider, Empty, message, Popconfirm, Row, Segmented, Space, Spin, Table, Typography, type TableProps } from 'antd';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useUserStore } from '../../../store/User';
import { TrendStatistics } from '../../core/TrendStats/TrendStats';
import { assetKindOptions, HoldingType, liabilityKindOptions, type HoldingDataType, type HoldingWithValue } from '../../../types/Holdings';
import { useMutation, useQuery } from '@tanstack/react-query';
import { deleteHolding, getHoldings } from '../../../api/holdings';
import { DeleteOutlined, ReloadOutlined } from '@ant-design/icons';
import { currencies } from '../../../types/Accounts';
import { HiTrendingUp } from 'react-icons/hi';
import dayjs from 'dayjs';

import './Overview.css';

const { Title, Text } = Typography;

export const Holdings = () => {
    const navigate = useNavigate();
    const userId = useUserStore.getState().userId;

    const [holdingType, setHoldingType] = useState<HoldingDataType>(HoldingType.Asset);

    const { data, refetch, isLoading, isSuccess } = useQuery({
        queryKey: ['holdings', userId],
        queryFn: getHoldings,
        enabled: !!userId,
    });

    const { mutate } = useMutation({
        mutationFn: deleteHolding,
        onSuccess: () => {
            message.success('Deleted Successfully');
            refetch();
        },
        onError: (error) => {
            message.error(error.message);
            console.error(error.message);
        },
    });

    const AssetHoldings = !!data ? data.records.filter((value) => value.type === HoldingType.Asset) : [];
    const LiabilityHoldings = !!data ? data.records.filter((value) => value.type === HoldingType.Liability) : [];

    const columns: TableProps<HoldingWithValue>['columns'] = [
        {
            dataIndex: 'name',
            key: 'name',
            title: 'Holding Name',
        },
        {
            dataIndex: 'kind',
            key: 'kind',
            title: 'Holding Type',
            render: (kind, record) => {
                return (
                    <>
                        {record.type === HoldingType.Asset
                            ? assetKindOptions.find((value) => value.value == kind).label
                            : liabilityKindOptions.find((value) => value.value == kind).label}
                    </>
                );
            },
        },
        {
            dataIndex: 'currentValue',
            key: 'currentValue',
            title: 'Current Valuation',
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
            dataIndex: 'valuedAt',
            key: 'valuedAt',
            title: 'Last Valuation Date',
            render: (value) => {
                return (
                    <Space size="small">
                        <Text>{dayjs.unix(value).format('DD MMM YYYY')}</Text>
                    </Space>
                );
            },
        },
        {
            dataIndex: '',
            key: 'actions',
            title: 'Actions',
            render: (_, record) => (
                <Space size="medium">
                    <Button
                        type="link"
                        icon={<HiTrendingUp />}
                        onClick={() => {
                            navigate(`./${record.holdingId}/trends`);
                        }}
                    >
                        View Trend
                    </Button>
                    <Popconfirm
                        title="Are you sure ?"
                        onConfirm={() => {
                            mutate({ holdingId: record.holdingId });
                        }}
                    >
                        <Button type="link" danger icon={<DeleteOutlined />}>
                            Delete
                        </Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    return (
        <>
            <Card>
                <Row gutter={16} justify={'space-between'}>
                    <Space>
                        <Title level={2} style={{ margin: 0 }} type="secondary" italic>
                            My Holdings
                        </Title>
                        <Badge count={0} color={'green'} />
                    </Space>
                    <Space>
                        <Button icon={<ReloadOutlined />} onClick={() => refetch()} />
                        <Button
                            type="primary"
                            onClick={() => {
                                navigate('./create');
                            }}
                        >
                            Create a Holding
                        </Button>
                    </Space>
                </Row>

                <Divider />
                <Row gutter={16} justify="center" align="middle">
                    <Col span={6}>
                        <TrendStatistics title="Net Worth" value={1322} prefix="₹" delta={4.2} />
                    </Col>
                    <Col span={6}>
                        <TrendStatistics title="Total Assets" value={1444} prefix="₹" delta={2.1} />
                    </Col>
                    <Col span={6}>
                        <TrendStatistics title="Total Liabilities" value={122} prefix="₹" delta={8.0} higherIsBetter={false} />
                    </Col>
                    <Col span={6}>
                        <TrendStatistics title="Debt-Asset Ratio" value={8.4} suffix="%" precision={1} delta={-1.2} higherIsBetter={false} />
                    </Col>
                </Row>

                <Divider />
                <Row gutter={16}>
                    <Col span={24}>
                        <Segmented<HoldingDataType>
                            className={holdingType === HoldingType.Asset ? 'holdingtype-asset' : 'holdingtype-liability'}
                            options={[HoldingType.Asset, HoldingType.Liability]}
                            onChange={setHoldingType}
                            block
                        />
                    </Col>
                </Row>

                <Row justify="center" gutter={16}>
                    {isLoading && <Spin />}
                    {isSuccess && (
                        <Table
                            style={{ width: '100%' }}
                            columns={columns}
                            dataSource={holdingType === HoldingType.Asset ? AssetHoldings : LiabilityHoldings}
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
                                            Create a Holding
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
