import { Image, Row } from 'antd';

export const Home = () => {
    return (
        <Row justify="center">
            <Image src="/work_in_progress.svg" alt="construction" preview={false} width={'30%'} />
        </Row>
    );
};
