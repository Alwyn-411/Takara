import { HomeOutlined } from '@ant-design/icons';
import { useMutation, useQuery } from '@tanstack/react-query';
import { Alert, Breadcrumb, Button, Card, Col, DatePicker, Flex, Form, Input, InputNumber, Row, type FormProps } from 'antd';
import dayjs from 'dayjs';
import { useNavigate, useParams } from 'react-router-dom';
import type { CreateResponse } from '../../../api/default';
import { createHoldingValuation, getHoldingById, type RecordValuationRequest as FormDataType } from '../../../api/holdings';
import { useUserStore } from '../../../store/User';
import { currencies } from '../../../types/Accounts';
import { HoldingType } from '../../../types/Holdings';

const { TextArea } = Input;

export const ValuationCreate = () => {
    const userId = useUserStore.getState().userId;
    const { holdingId } = useParams<{ holdingId: string }>();

    const navigate = useNavigate();
    const [form] = Form.useForm<FormDataType>();

    const { mutate, isPending, isError } = useMutation<CreateResponse, Error, FormDataType>({
        mutationFn: createHoldingValuation,
        onSuccess: () => navigate(-1),
    });

    const HoldingQuery = useQuery({
        queryKey: ['Holding', userId, holdingId],
        queryFn: () => getHoldingById({ holdingId: holdingId! }),
        enabled: !!userId && !!holdingId,
    });

    const currencyObj = currencies.find((c) => c.value === HoldingQuery.data?.currency);

    const onFinish: FormProps<FormDataType>['onFinish'] = (values) => {
        mutate({
            holdingId: holdingId!,
            value: String(values.value),
            observedAt: (values.observedAt as any)?.unix(),
            ...(values.note && { note: values.note }),
            ...(values.quantity && { quantity: String(values.quantity) }),
            ...(values.unitPrice && { unitPrice: String(values.unitPrice) }),
        });
    };

    return (
        <Flex vertical>
            <Row justify="space-between" gutter={16} style={{ padding: 12 }}>
                {HoldingQuery.data && (
                    <Breadcrumb
                        items={[
                            {
                                href: '/home',
                                title: <HomeOutlined />,
                            },
                            {
                                href: '/holdings',
                                title: 'Holdings',
                            },
                            {
                                href: '.',
                                title: HoldingQuery.data.name,
                            },
                            {
                                title: 'Create a New Valuation',
                            },
                        ]}
                    />
                )}
            </Row>

            {isError && (
                <Row gutter={16}>
                    <Alert
                        title="Holding Creation Failed"
                        description="An error occurred while processing this request."
                        type="error"
                        showIcon
                        style={{ marginBottom: 16 }}
                    />
                </Row>
            )}

            <Card>
                <Form<FormDataType>
                    form={form}
                    layout="vertical"
                    size="large"
                    onFinish={onFinish}
                    requiredMark="optional"
                    initialValues={{ type: HoldingType.Asset, openedAt: dayjs(), observedAt: dayjs() }}
                >
                    <Row gutter={16}>
                        <Col span={12}>
                            <Form.Item label="Valuation" name="value" rules={[{ required: true, message: 'Enter your valuation' }]}>
                                <InputNumber
                                    style={{ width: '100%' }}
                                    prefix={currencyObj?.symbol}
                                    suffix={currencyObj?.value}
                                    placeholder="Enter your valuation"
                                />
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            <Form.Item label="Observed On" name="observedAt" rules={[{ required: true, message: 'Select an observed time' }]}>
                                <DatePicker use12Hours format="DD MMM YYYY" style={{ width: '100%' }} />
                            </Form.Item>
                        </Col>
                    </Row>

                    <Row gutter={16}>
                        <Col span={8}>
                            <Form.Item label="Valuation Note" name="note">
                                <TextArea rows={1} placeholder="Enter a note on valuation" />
                            </Form.Item>
                        </Col>
                        <Col span={8}>
                            <Form.Item label="Quantity" name="quantity">
                                <InputNumber style={{ width: '100%' }} placeholder="Enter your holding quantity" />
                            </Form.Item>
                        </Col>
                        <Col span={8}>
                            <Form.Item label="Unit Price" name="unitPrice">
                                <InputNumber style={{ width: '100%' }} placeholder="Enter your holding unit price" />
                            </Form.Item>
                        </Col>
                    </Row>

                    <Flex justify="flex-end">
                        <Button type="primary" htmlType="submit" size="large" loading={isPending}>
                            Create a new Valuation
                        </Button>
                    </Flex>
                </Form>
            </Card>
        </Flex>
    );
};
