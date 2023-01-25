import { CatalogItem, convertRaw } from "@web/extensions/sidebar/core/processCatalog";
import { Input, Form, Button } from "@web/sharedComponents";
import { useCallback, useState } from "react";

import FileTypeSelect from "./FileTypeSelect";

type Props = {
  onOpenDetails?: (data?: CatalogItem) => void;
};

const WebDataTab: React.FC<Props> = ({ onOpenDetails }) => {
  const [dataUrl, setDataUrl] = useState("");

  const fetchDataFromUrl = useCallback(async (url: string) => {
    try {
      const result = await fetch(url);
      if (result.ok) {
        return result;
      }
    } catch (error) {
      return undefined;
    }
  }, []);

  const handleClick = useCallback(async () => {
    const result = await fetchDataFromUrl(dataUrl);
    if (result) {
      // Catalog Item
      const filename = dataUrl.substring(dataUrl.lastIndexOf("/") + 1);
      const id = "id" + Math.random().toString(16).slice(2);
      const item = {
        id: id,
        description:
          "Please contact the provider of this data for more information, including information about usage rights and constraints.",
        name: filename,
        data_url: dataUrl,
      };
      const catalogItem = convertRaw([item])[0] as CatalogItem;
      catalogItem.type = "item";
      if (onOpenDetails) onOpenDetails(catalogItem);

      // Raw Data
      // const data = await result.text();
    }
  }, [dataUrl, fetchDataFromUrl, onOpenDetails]);

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
        <Input
          placeholder={"Please input a valid URL"}
          onChange={e => setDataUrl(e.target.value)}
        />
      </Form.Item>
      <Form.Item style={{ textAlign: "right" }}>
        <Button type="primary" htmlType="submit" onClick={handleClick}>
          Upload
        </Button>
      </Form.Item>
    </Form>
  );
};

export default WebDataTab;
