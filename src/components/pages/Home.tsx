import { Image, Typography } from 'antd';

export const Home = () => {
    return (
        <Typography.Text>
            <Image src="/work_in_progress.svg" alt="construction" preview={false} width="100%" />
        </Typography.Text>
    );
};
