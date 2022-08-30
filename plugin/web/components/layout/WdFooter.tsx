import { Button, Typography } from "antd";
import { Footer } from "antd/lib/layout/layout";
import { memo } from "react";

import "../../../node_modules/antd/dist/antd.less";
import colors from "../../styles/colors";
import { ReactComponent as Trash } from "../UI/Icon/Icons/trash.svg";

const WdFooter: React.FC = () => {
  const { Text } = Typography;
  return (
    <Footer className={"footer"}>
      <Button
        type="default"
        className="removeBtn"
        icon={<Trash />}
        color={colors.dark.outline.weak}>
        Remove All
      </Button>
      <Text>DataSet x 0</Text>
    </Footer>
  );
};
export default memo(WdFooter);
