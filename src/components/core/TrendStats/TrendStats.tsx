import React, { type ReactNode } from 'react';
import { Card, Statistic, Tag, Flex, Typography, theme } from 'antd';
import { IoTrendingDownSharp, IoTrendingUpSharp } from 'react-icons/io5';

const { Text } = Typography;

interface TrendStatisticsProps {
    title: string;
    value: number;
    prefix?: ReactNode;
    suffix?: string;
    precision?: number;
    delta?: number; // % change vs the period — drives arrow + color
    higherIsBetter?: boolean; // false for liabilities & debt ratio
    period?: 'month' | 'year' | 'week';
}

export const TrendStatistics: React.FC<TrendStatisticsProps> = ({
    title,
    value,
    prefix,
    suffix,
    precision = 0,
    delta,
    higherIsBetter = true,
    period = 'month',
}) => {
    const { token } = theme.useToken();

    const hasTrend = delta !== undefined && delta !== 0;
    const isUp = (delta ?? 0) > 0;
    const isGood = isUp ? higherIsBetter : !higherIsBetter;
    const trendColor = isGood ? token.colorSuccess : token.colorError;

    return (
        <Card variant="borderless" styles={{ body: { padding: 16 } }}>
            <Statistic
                title={title}
                value={value}
                styles={{ content: { color: trendColor }, title: { fontWeight: 'bold', color: 'black' } }}
                precision={precision}
                prefix={prefix}
                suffix={suffix}
            />
            {hasTrend && (
                <Flex align="center" gap={6} style={{ marginTop: 8 }}>
                    {isUp ? <IoTrendingUpSharp style={{ color: trendColor }} /> : <IoTrendingDownSharp style={{ color: trendColor }} />}
                    <Tag variant="filled" color={isGood ? 'success' : 'error'} style={{ marginInlineEnd: 0, borderRadius: 999 }}>
                        {Math.abs(delta!).toFixed(1)}%
                    </Tag>
                    <Text type="secondary" style={{ fontSize: 12 }}>
                        vs last {period}
                    </Text>
                </Flex>
            )}
        </Card>
    );
};
