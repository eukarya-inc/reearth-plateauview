import { CatalogItem, convertRaw } from "@web/extensions/sidebar/core/processCatalog";
import { Input, Form, Button } from "@web/sharedComponents";
import { useState } from "react";

import FileTypeSelect from "./FileTypeSelect";

type Props = {
  onOpenDetails?: (data?: CatalogItem) => void;
};

const WebDataTab: React.FC<Props> = ({ onOpenDetails }) => {
  const [url, setUrl] = useState("");

  const fetchDataFromUrl = async () => {
    try {
      const result = await fetch(url);
      if (result.ok) {
        const filename = url.substring(url.lastIndexOf("/") + 1);
        const id = "id" + Math.random().toString(16).slice(2);
        const item = {
          id: id,
          description: "web data",
          city_name: "", // TODO: find a way to add the city name
          prefecture: "", // TODO: find a way to add the prefecture
          name: filename,
          data_url: url,
        };
        const catalogItem = convertRaw([item])[0] as CatalogItem;
        catalogItem.type = "item";
        if (onOpenDetails) onOpenDetails(catalogItem);

        const data = result.json();
        console.log(data);
      }
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <Form layout="vertical">
      <Form.Item name="file-type" label="Select file type">
        <FileTypeSelect />
      </Form.Item>
      <Form.Item
        name="url"
        label="File URL"
        rules={[
          { required: true },
          { message: "Please input the URL of the asset!" },
          { type: "url", warningOnly: true },
        ]}>
        <Input placeholder={"Please input a valid URL"} onChange={e => setUrl(e.target.value)} />
      </Form.Item>
      <Form.Item style={{ textAlign: "right" }}>
        <Button type="primary" htmlType="submit" onClick={fetchDataFromUrl}>
          Upload
        </Button>
      </Form.Item>
    </Form>
  );
};

export default WebDataTab;
