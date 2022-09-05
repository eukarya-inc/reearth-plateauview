import MainLayout from "./Layout/MainLayout";
import "antd/dist/antd.less";
import "../../theme/styles/theme.less";

const Sidebar: React.FC = () => {
  return (
    <>
      <MainLayout isInsideEditor={false}></MainLayout>
    </>
  );
};

export default Sidebar;
