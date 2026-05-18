import { Typography } from "antd";
import { useUserStore } from "../../store/User";

export const Home = () => {
  const userName = useUserStore.getState().userName;

  return <Typography.Text>{userName}</Typography.Text>;
};
