import { CatalogItem, convertRaw } from "@web/extensions/sidebar/core/processCatalog";
import { Form } from "@web/sharedComponents";
import { InboxOutlined } from "@web/sharedComponents/Icon/icons";
import Upload, { UploadProps, UploadFile } from "@web/sharedComponents/Upload";
import { RcFile } from "antd/lib/upload";
import { useCallback, useState } from "react";

import FileTypeSelect from "./FileTypeSelect";

type Props = {
  onOpenDetails?: (data?: CatalogItem) => void;
};

const LocalDataTab: React.FC<Props> = ({ onOpenDetails }) => {
  const [fileList, setFileList] = useState<UploadFile[]>([]);
  const fileFormats = ".kml,.czml,.topojson,.geojson,.json,.gltf,.glb";

  const onRemove = useCallback(
    (file: UploadFile) => {
      const index = fileList.indexOf(file);
      const newFileList = fileList.slice();
      newFileList.splice(index, 1);
      setFileList(newFileList);
    },
    [fileList],
  );

  const beforeUpload = useCallback(
    (file: RcFile, files: RcFile[]) => {
      const filename = file.name;
      const id = "id" + Math.random().toString(16).slice(2);
      const url = URL.createObjectURL(file);
      const item = {
        id: id,
        description: "local data",
        city_name: "", // TODO: find a way to add the city name
        prefecture: "", // TODO: find a way to add the prefecture
        name: filename,
        data_url: url,
      };
      const catalogItem = convertRaw([item])[0] as CatalogItem;
      catalogItem.type = "item";
      if (onOpenDetails) onOpenDetails(catalogItem);

      const reader = new FileReader();
      reader.readAsText(file);
      let data;
      reader.onload = e => {
        data = e?.target?.result ?? "";
        console.log(JSON.parse(data as string));
      };

      setFileList([...fileList, ...files]);
      return false;
    },
    [fileList, onOpenDetails],
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
