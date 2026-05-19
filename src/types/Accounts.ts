import type { Base, Timestamp } from "./Common";

export interface Accounts extends Base, Timestamp {
  accountId: string;
  type: string;
  name: string;
  accountNumber: string;
  description?: string;
  currency: string;
  interest: string;
  balance: string;
}

export const currencies = [
  {
    label: "₹ INR - Indian Rupee",
    value: "INR",
    symbol: "₹",
  },
  {
    label: "¥ CNY - Chinese Yuan",
    value: "CNY",
    symbol: "¥",
  },
  {
    label: "€ EUR - Euro",
    value: "EUR",
    symbol: "€",
  },
  {
    label: "£ GBP - British Pound",
    value: "GBP",
    symbol: "£",
  },
  {
    label: "¥ JPY - Japanese Yen",
    value: "JPY",
    symbol: "¥",
  },
  {
    label: "$ USD - US Dollar",
    value: "USD",
    symbol: "$",
  },
  {
    label: "C$ CAD - Canadian Dollar",
    value: "CAD",
    symbol: "C$",
  },
  {
    label: "A$ AUD - Australian Dollar",
    value: "AUD",
    symbol: "A$",
  },
  {
    label: "CHF - Swiss Franc",
    value: "CHF",
    symbol: "CHF",
  },
  {
    label: "SGD - Singapore Dollar",
    value: "SGD",
    symbol: "S$",
  },
];