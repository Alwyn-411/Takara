import { Badge, Button, Card, Divider, Row, Segmented, Space, Typography } from 'antd';
import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useUserStore } from '../../../store/User';
import { HoldingType, type HoldingDataType } from '../../../types/Transactions';
import './Overview.css';

const { Title } = Typography;

export const Holdings = () => {
    const navigate = useNavigate();
    const userId = useUserStore.getState().userId;

    const [holdingType, setHoldingType] = useState<HoldingDataType>(HoldingType.Asset);

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
                <Segmented<HoldingDataType>
                    className={holdingType === HoldingType.Asset ? 'holdingtype-asset' : 'holdingtype-liability'}
                    options={[HoldingType.Asset, HoldingType.Liability]}
                    onChange={setHoldingType}
                    block
                />
            </Card>
        </>
    );
};
