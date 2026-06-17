import { Badge, Button, Card, Divider, Row, Segmented, Space, Typography } from 'antd';
import Text from 'antd/es/typography/Text';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useUserStore } from '../../../store/User';
import { HoldingType, type HoldingDataType } from '../../../types/Transactions';
const { Title } = Typography;

export const Holdings = () => {
    const navigate = useNavigate();
    const userId = useUserStore.getState().userId;
    const [type, setType] = useState<HoldingDataType>(HoldingType.Asset);

    const options = [
        {
            label: <Text type="danger">Asset</Text>,
            value: HoldingType.Asset,
        },
        {
            label: <Text type="success">Liability</Text>,
            value: HoldingType.Liability,
        },
    ];

    return (
        <>
            <Card>
                <Row justify={'space-between'}>
                    <Space>
                        <Title level={2} style={{ margin: 0 }} type="secondary" italic>
                            My Holdings
                        </Title>
                        <Badge count={0} color={'green'} />
                    </Space>
                    <Space>
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
                <Segmented<HoldingDataType> options={options} value={type} onSelect={(value) => {setType(value.)}} block />
            </Card>
        </>
    );
};
