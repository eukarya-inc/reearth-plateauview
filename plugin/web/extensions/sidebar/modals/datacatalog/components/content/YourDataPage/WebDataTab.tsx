import { UserDataItem } from "@web/extensions/sidebar/modals/datacatalog/types";
import { getExtension } from "@web/extensions/sidebar/utils/file";
import { Input, Form, Button } from "@web/sharedComponents";
import { useCallback, useState } from "react";

import WebFileTypeSelect, { FileType } from "./WebFileTypeSelect";

type Props = {
  onOpenDetails?: (data?: UserDataItem, needLayerName?: boolean) => void;
  setSelectedWebItem?: (data?: UserDataItem) => void;
};

const WebDataTab: React.FC<Props> = ({ onOpenDetails, setSelectedWebItem }) => {
  const [dataUrl, setDataUrl] = useState("");
  const [fileType, setFileType] = useState<FileType>("auto");

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

  const setDataFormat = useCallback((type: FileType, filename: string) => {
    const extension = getExtension(filename);
    if (type === "auto") {
      // more exceptions will be added in the future
      switch (extension) {
        // georss
        case "rss":
          return "rss";
        // georss
        case "xml":
          return "xml";
        // shapefile
        case "zip":
          return "zip";
        default:
          return extension;
      }
    }
    return type;
  }, []);

  const handleClick = useCallback(async () => {
    const result = await fetchDataFromUrl(dataUrl);
    if (result) {
      // Catalog Item
      const filename = dataUrl.substring(dataUrl.lastIndexOf("/") + 1);
      const id = "id" + Math.random().toString(16).slice(2);
      const item: UserDataItem = {
        type: "item",
        id: id,
        dataID: id,
        description:
          "Please contact the provider of this data for more information, including information about usage rights and constraints.",
        name: filename,
        url: dataUrl,
        format: setDataFormat(fileType, filename),
      };
      const requireLayerName = true;
      if (onOpenDetails) onOpenDetails(item, requireLayerName);
      if (setSelectedWebItem) setSelectedWebItem(item);
    }
  }, [dataUrl, fetchDataFromUrl, fileType, onOpenDetails, setDataFormat, setSelectedWebItem]);

  const handleFileTypeSelect = useCallback((type: string) => {
    setFileType(type as FileType);
  }, []);

  return (
    <Form layout="vertical">
      <Form.Item name="file-type" label="Select file type">
        <WebFileTypeSelect onFileTypeSelect={handleFileTypeSelect} />
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
