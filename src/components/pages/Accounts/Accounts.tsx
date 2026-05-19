import { Button, Card, Row, Typography } from "antd";
import { useNavigate } from "react-router-dom";
const { Title } = Typography;


export const Accounts = () => {
    const navigate = useNavigate()

    return (<>
        <Card>
            <Row justify={"space-between"}>
                <Title level={2} style={{ margin: 0 }} type="secondary" italic>
                    My Bank Accounts
                </Title>
                <Button type="primary" onClick={() => { navigate("/accounts/create") }}>
                    Create Bank Account
                </Button>
            </Row>
        </Card>
    </>
    )
};
