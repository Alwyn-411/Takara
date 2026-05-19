import { Alert, Breadcrumb, Button, Card, Col, Form, Input, InputNumber, Radio, Row, Select, Typography, type FormProps } from "antd";
import { currencies, type Accounts } from "../../../types/Accounts";
import TextArea from "antd/es/input/TextArea";
import { useMutation } from "@tanstack/react-query";
import { createAccount, type CreateAccountResponse } from "../../../hooks/accounts";
import type { CheckboxGroupProps } from 'antd/es/checkbox';
import { useUserStore } from "../../../store/User";

const { Text } = Typography

interface Formfields extends Partial<Omit<Accounts, "active" | "createdAt" | "updatedAt" | "accountId">> { }

export const AccountCreation = () => {
    const [form] = Form.useForm();
    const { mutate, isPending, isError } = useMutation<
        CreateAccountResponse,
        Error,
        Partial<Accounts>
    >({
        mutationFn: createAccount
    });

    const options: CheckboxGroupProps<string>['options'] = [
        { label: 'Savings Account', value: 'savings' },
        { label: 'Current Account', value: 'current' },
    ];

    const selectedCurrency: string = Form.useWatch(
        "currency",
        form
    );

    const selectedType: string = Form.useWatch(
        "type",
        form
    );

    const currencyObj = currencies.find(
        (c) => c.value === selectedCurrency
    );

    const onFinish: FormProps<Formfields>["onFinish"] = (values) => {
        const AccountData: Formfields = {
            userId: useUserStore.getState().userId,
            type: values.type,
            name: values.name,
            accountNumber: values.accountNumber,
            description: values.description,
            currency: values.currency,
        };

        if (!!values.currency) {
            AccountData.balance = values.balance
        }

        if (values.type !== "current") {
            AccountData.interest = values.interest
        }

        mutate(AccountData);
    };


    return (
        <Col span={24}>
            <Row>
                <Breadcrumb items={[
                    {
                        title: 'Accounts',
                        href: "/accounts"
                    },
                    {
                        title: 'Create New Bank Account',
                    },
                ]} />
            </Row>
            <Row justify="center">
                <Card style={{ "width": "100%" }}>

                    {isError && (
                        <Alert
                            title="Account Creation Failed"
                            description={<Text>Error occured while processing this request</Text>}
                            type="error"
                            showIcon
                        />
                    )}

                    <Form<Formfields>
                        form={form}
                        layout="vertical"
                        size="large"
                        onFinish={onFinish}
                        requiredMark="optional"
                    >
                        <Row gutter={16}>
                            <Col span={12}>
                                <Form.Item
                                    label="Name"
                                    name="name"
                                    required
                                    rules={[
                                        { required: true, message: "Enter a username" },
                                        { min: 3, message: "Minimum 3 characters" },
                                    ]}
                                >
                                    <Input placeholder="Enter the Name of the Bank Account" />
                                </Form.Item>
                            </Col>
                            <Col span={12}>
                                <Form.Item
                                    label="Account Number"
                                    name="accountNumber"
                                    rules={[
                                        { required: true, message: "Enter Your Account Number" },
                                        { min: 15, message: "Minimum 15 characters" },
                                    ]}
                                >
                                    <Input placeholder="Enter the Bank Account Number" />
                                </Form.Item>
                            </Col>
                        </Row>

                        <Form.Item
                            label="Description"
                            name="description"
                        >
                            <TextArea rows={4} placeholder="Enter a short description about this account" />
                        </Form.Item>

                        <Row gutter={16}>
                            <Col span={12}>
                                <Form.Item
                                    label="Currency"
                                    name="currency"
                                    rules={[
                                        { required: true, message: "Select a Currency" },
                                    ]}
                                >
                                    <Select options={currencies} placeholder="Select a Currency" />
                                </Form.Item>
                            </Col>
                            <Col span={12}>
                                <Form.Item
                                    label="Balance"
                                    name="balance"
                                    dependencies={["currency"]}
                                    rules={[
                                        { required: true, message: "Enter Your Bank Balance" }
                                    ]}
                                >
                                    <InputNumber style={{ "width": "100%" }} disabled={!selectedCurrency} prefix={currencyObj?.symbol} placeholder="Enter Your Bank Balance" suffix={currencyObj?.value} />
                                </Form.Item>
                            </Col>
                        </Row>
                        <Row gutter={16}>
                            <Col span={12}>
                                <Form.Item
                                    label="Account Type"
                                    name="type"
                                    rules={[
                                        { required: true, message: "Choose your account Type" },
                                    ]}
                                >
                                    <Radio.Group
                                        block
                                        options={options}
                                        defaultValue="savings"
                                        optionType="button"
                                        buttonStyle="solid"
                                    />
                                </Form.Item>
                            </Col>
                            <Col span={12}>
                                <Form.Item
                                    label="Interest"
                                    name="interest"
                                >
                                    <InputNumber disabled={selectedType === "current"} style={{ "width": "100%" }} placeholder="Enter Your Bank's Interest Rate" suffix="%" />
                                </Form.Item>
                            </Col>
                        </Row>
                        <Row justify="end">
                            <Button
                                type="primary"
                                htmlType="submit"
                                size="large"
                                loading={isPending}>
                                Create Bank Account
                            </Button>
                        </Row>
                    </Form>
                </Card>
            </Row>
        </Col >
    )
};
