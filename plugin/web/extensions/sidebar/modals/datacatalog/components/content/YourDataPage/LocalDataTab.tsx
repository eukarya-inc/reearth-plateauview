import { UserDataItem } from "@web/extensions/sidebar/modals/datacatalog/types";
import { Form } from "@web/sharedComponents";
import { InboxOutlined } from "@web/sharedComponents/Icon/icons";
import Upload, { UploadProps, UploadFile } from "@web/sharedComponents/Upload";
import { RcFile } from "antd/lib/upload";
import { useCallback, useState } from "react";

import FileTypeSelect from "./FileTypeSelect";

type Props = {
  onOpenDetails?: (data?: UserDataItem) => void;
  setSelectedLocalItem?: (data?: UserDataItem) => void;
};

const fileFormats = ".kml,.kmz,.csv,.czml,.gpx,.topojson,.geojson,.json,.zip";

const LocalDataTab: React.FC<Props> = ({ onOpenDetails, setSelectedLocalItem }) => {
  const [fileList, setFileList] = useState<UploadFile[]>([]);

  const onRemove = useCallback((_file: UploadFile) => {
    setFileList([]);
  }, []);

  const beforeUpload = useCallback(
    (file: RcFile, files: RcFile[]) => {
      // Catalog Item
      const filename = file.name;
      const id = "id" + Math.random().toString(16).slice(2);
      const url = URL.createObjectURL(file);
      const item: UserDataItem = {
        type: "item",
        id: id,
        description:
          "This file only exists in your browser. To share it, you must load it onto a public web server.",
        name: filename,
        dataUrl: url,
      };
      if (onOpenDetails) onOpenDetails(item);
      if (setSelectedLocalItem) setSelectedLocalItem(item);

      // Raw Data
      // const reader = new FileReader();
      // reader.readAsText(file);
      // let data;
      // reader.onload = e => {
      //   data = e?.target?.result;
      // };

      setFileList([...files]);
      return false;
    },
    [onOpenDetails, setSelectedLocalItem],
  );

  const props: UploadProps = {
    name: "file",
    multiple: false,
    directory: false,
    showUploadList: true,
    accept: fileFormats,
    listType: "picture",
    onRemove: onRemove,
    beforeUpload: beforeUpload,
    fileList,
  };

  return (
    <Form layout="vertical">
      <Form.Item name="file-type" label="Select file type">
        <FileTypeSelect />
      </Form.Item>
      <Form.Item label="Upload File">
        <Form.Item name="upload-file" style={{ height: 300, overflowY: "scroll" }}>
          <Upload.Dragger {...props}>
            <p className="ant-upload-drag-icon">
              <InboxOutlined />
            </p>
            <p className="ant-upload-text">Click or drag file to this area to upload</p>
            <p className="ant-upload-hint">
              Support for a single or bulk upload. Strictly prohibit from uploading company data or
              other band files
            </p>
          </Upload.Dragger>
        </Form.Item>
      </Form.Item>
    </Form>
  );
};

export default LocalDataTab;
