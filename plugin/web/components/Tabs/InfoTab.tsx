import { Button, Divider, Form, Input, Space, Typography } from "antd";
import { memo } from "react";

import "../../styles/style.less";
import { ReactComponent as PlateauLogoPrt } from "../UI/Icon/Icons/plateauLogoPrt.svg";

const InfoTab: React.FC = () => {
  const { Text, Paragraph } = Typography;

  return (
    <Space direction="vertical">
      <Typography>
        <Paragraph>
          PLATEAU は、国土交通省が進める 3D都市モデル整備・活用・オープンデータ化
          のリーディングプロジェクトである。都市活動のプラットフォームデータとして
          3D都市モデルを整備し、
          そのユースケースを創出。さらにこれをオープンデータとして公開することで、誰もが自由に都市のデータを引き出し、活用できるようになる。
        </Paragraph>
      </Typography>
      <Button type="default" className="plateauButton " icon={<PlateauLogoPrt />}>
        PLATEAU Project Website
      </Button>
      <Divider />
      <Text className="sectionFonts">Feedback</Text>
      <Text className="sectionFonts">We would love to hear from you!</Text>
      <Form className="formSection" layout="vertical">
        <Form.Item label="Your name (optional):">
          <Input placeholder="example" />
        </Form.Item>
        <Form.Item label="Email address (optional):">
          <Input placeholder="example" />
          <Text className="ant-form-text" type="secondary">
            We will not follow up without it!
          </Text>
        </Form.Item>
        <Form.Item label="Comment or question:">
          <Input.TextArea placeholder="Autosize height based on content lines" />
        </Form.Item>
        <Form.Item>
          <Button type="primary">Submit</Button>
          <Button type="primary">cancel</Button>
        </Form.Item>
      </Form>
    </Space>
  );
};
export default memo(InfoTab);
