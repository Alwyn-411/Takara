import { HomeOutlined } from '@ant-design/icons';
import { useMutation } from '@tanstack/react-query';
import {
    Alert,
    Breadcrumb,
    Button,
    Card,
    Col,
    DatePicker,
    Flex,
    Form,
    Input,
    InputNumber,
    Row,
    Segmented,
    Select,
    Typography,
    type FormProps,
} from 'antd';
import { useNavigate } from 'react-router-dom';
import type { CreateResponse } from '../../../api/default';
import { currencies } from '../../../types/Accounts';
import { createHolding, type CreateHoldingRequest as FormDataType } from '../../../api/holdings';
import { HoldingType, kindOptionsByType } from '../../../types/Holdings';
import dayjs from 'dayjs';
import { useEffect } from 'react';

const { Text } = Typography;
const { TextArea } = Input;

const holdingTypeOptions = [
    {
        label: <Text type="success">{HoldingType.Asset}</Text>,
        value: HoldingType.Asset,
    },
    {
        label: <Text type="danger">{HoldingType.Liability}</Text>,
        value: HoldingType.Liability,
    },
];

export const HoldingsCreate = () => {
    const navigate = useNavigate();
    const [form] = Form.useForm<FormDataType>();

    const { mutate, isPending, isError } = useMutation<CreateResponse, Error, FormDataType>({
        mutationFn: createHolding,
        onSuccess: () => navigate(-1),
    });

    const selectedCurrency = Form.useWatch('currency', form);
    const selectedHoldingType = Form.useWatch('type', form);

    const currencyObj = currencies.find((c) => c.value === selectedCurrency);

    const onFinish: FormProps<FormDataType>['onFinish'] = (values) => {
        mutate({
            name: values.name,
            currency: values.currency,
            kind: values.kind,
            type: values.type,
            value: String(values.value),
            openedAt: (values.openedAt as any)?.unix(),
            observedAt: (values.observedAt as any)?.unix(),
            ...(values.description && { description: values.description }),
            ...(values.note && { note: values.note }),
            ...(values.quantity && { quantity: String(values.quantity) }),
            ...(values.unitPrice && { unitPrice: String(values.unitPrice) }),
        });
    };

    useEffect(() => {
        form.setFieldValue('kind', undefined);
    }, [selectedHoldingType, form]);

    return (
        <Flex vertical>
            <Row justify="space-between" gutter={16} style={{ padding: 12 }}>
                <Breadcrumb
                    items={[{ href: '/home', title: <HomeOutlined /> }, { href: '/holdings', title: 'Holdings' }, { title: 'Create New Holding' }]}
                />
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
                            <Form.Item
                                label="Name"
                                name="name"
                                rules={[
                                    { required: true, message: 'Enter a name' },
                                    { min: 3, message: 'Minimum 3 characters' },
                                ]}
                            >
                                <Input placeholder="Enter the name of the holding" />
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            <Form.Item label="Opened At" name="openedAt" rules={[{ required: true, message: 'Select an opened at' }]}>
                                <DatePicker use12Hours format="DD MMM YYYY" style={{ width: '100%' }} />
                            </Form.Item>
                        </Col>
                    </Row>

                    <Row gutter={16}>
                        <Col span={12}>
                            <Form.Item label="Holding Type" name="type" rules={[{ required: true, message: `Select a Holding Type` }]}>
                                <Segmented options={holdingTypeOptions} block />
                            </Form.Item>
                        </Col>
                        <Col span={12}>
                            <Form.Item
                                label={`${selectedHoldingType ?? ''} Type`}
                                name="kind"
                                dependencies={['type']}
                                rules={[{ required: true, message: `Select a ${selectedHoldingType ?? 'holding'} Type` }]}
                            >
                                <Select
                                    options={selectedHoldingType ? kindOptionsByType[selectedHoldingType] : []}
                                    placeholder={`Select a ${selectedHoldingType ?? 'holding'} Type`}
                                />
                            </Form.Item>
                        </Col>
                    </Row>

                    <Form.Item label="Description" name="description">
                        <TextArea rows={3} placeholder="Enter a short description about this holding" />
                    </Form.Item>

                    <Row gutter={16}>
                        <Col span={8}>
                            <Form.Item label="Currency" name="currency" rules={[{ required: true, message: 'Select a currency' }]}>
                                <Select options={currencies} placeholder="Select a currency" />
                            </Form.Item>
                        </Col>
                        <Col span={8}>
                            <Form.Item
                                label="Valuation"
                                name="value"
                                dependencies={['currency']}
                                rules={[{ required: true, message: 'Enter your valuation' }]}
                            >
                                <InputNumber
                                    style={{ width: '100%' }}
                                    disabled={!selectedCurrency}
                                    prefix={currencyObj?.symbol}
                                    suffix={currencyObj?.value}
                                    placeholder="Enter your valuation"
                                />
                            </Form.Item>
                        </Col>
                        <Col span={8}>
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
                            Create Holding
                        </Button>
                    </Flex>
                </Form>
            </Card>
        </Flex>
    );
};
