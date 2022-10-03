import {
  Button,
  Divider,
  Form,
  Icon,
  Input,
  Space,
  Typography,
} from "@web/extensions/sharedComponents";
import { styled } from "@web/theme";
import { memo } from "react";

const Info: React.FC = () => {
  const { Text, Paragraph } = Typography;
  const [form] = Form.useForm();

  const handleCancel = () => {
    form.resetFields();
  };

  const handleSend = (values: any) => {
    console.log(values);
    form.resetFields();
  };

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
      <PlateauButton type="default" icon={<Icon icon="plateauLogoPart" />}>
        PLATEAU Project Website
      </PlateauButton>
      <Divider />
      <Text>Feedback</Text>
      <Text>We would love to hear from you!</Text>
      <FormSection form={form} name="feedback" onFinish={handleSend} layout="vertical">
        <Form.Item name="name" label="Your name (optional):">
          <Input placeholder="example" />
        </Form.Item>
        <Form.Item name="email" label="Email address (optional):">
          <Input placeholder="example" />
          <Text type="secondary">We will not follow up without it!</Text>
        </Form.Item>
        <Form.Item name="comment" label="Comment or question:">
          <Input.TextArea placeholder="Autosize height based on content lines" />
        </Form.Item>
        <FormItem>
          <Button htmlType="button" onClick={handleCancel}>
            Cancel
          </Button>
          <Button type="primary" htmlType="submit">
            Send
          </Button>
        </FormItem>
      </FormSection>
    </Space>
  );
};
export default memo(Info);

const PlateauButton = styled(Button)`
  display: flex;
  flex-direction: row;
  justify-content: center;
  align-items: center;
  gap: 8px;
  width: 100%;
  height: 48px;
  border: 1px solid #c7c5c5;
  border-radius: 4px;
`;

const FormSection = styled(Form)`
  height: 342px;
  gap: 8px;
`;
const FormItem = styled(Form.Item)`
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  align-items: flex-start;
  padding: 8px 0px;
  gap: 12px;
  height: 48px;
`;
