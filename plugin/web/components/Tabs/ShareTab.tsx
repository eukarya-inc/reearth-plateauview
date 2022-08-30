import { CopyOutlined } from "@ant-design/icons";
import { Button, Divider, Input, Radio, Row, Space, Typography } from "antd";
import React,{ memo } from "react";
import "../../styles/style.less";

const PrintMapData = ["Download map ( png )", "Show Print View"];

const ShareTab: React.FC = () => {
  const { Text } = Typography;

  return (
    <Space direction="vertical">
      <Text className="sectionFonts">Share URL</Text>
      <Row className="section">
        <Input.Group compact className="inputGroup">
          <Input
            style={{ width: "calc(286px)" }}
            defaultValue="Anyone with this URL will be able to access this map."
          />
          <Button icon={<CopyOutlined />} type="primary" />
          <Text className="ant-form-text" type="secondary">
            Anyone with this URL will be able to access this map.
          </Text>
        </Input.Group>
      </Row>
      <Divider />
      <Text className="sectionFonts">Print Map</Text>
      <Row className="section">
        <Radio.Group defaultValue="3D Terrain" buttonStyle="solid" className="printGroupButtons">
          {PrintMapData.map(item => (
            <Radio.Button key={item} value={item} className={"printButtons"}>
              <Text className="printButtonFonts">{item}</Text>
            </Radio.Button>
          ))}
        </Radio.Group>
        <Text className="ant-form-text" type="secondary">
          Open a printable version of this map.
        </Text>
      </Row>
    </Space>
  );
};
export default memo(ShareTab);
