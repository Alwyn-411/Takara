import { useQuery } from "@tanstack/react-query";
import { Typography } from "antd";
import { ping } from "../../hooks/client";

export const Home = () => {
  const { data, isLoading } = useQuery({ queryKey: ["ping"], queryFn: ping });

  if (isLoading) return <div>Loading...</div>;

  return <Typography.Text>{data.message}</Typography.Text>;
};
