import { UserDataItem } from "@web/extensions/sidebar/modals/datacatalog/types";
import { getExtension } from "@web/extensions/sidebar/utils/file";
import { Form } from "@web/sharedComponents";
import { InboxOutlined } from "@web/sharedComponents/Icon/icons";
import Upload, { UploadProps, UploadFile } from "@web/sharedComponents/Upload";
import { RcFile } from "antd/lib/upload";
import { useCallback, useMemo, useState } from "react";

import FileTypeSelect, { fileFormats, FileType } from "./LocalFileTypeSelect";

type Props = {
  onOpenDetails?: (data?: UserDataItem) => void;
  setSelectedLocalItem?: (data?: UserDataItem) => void;
};

const LocalDataTab: React.FC<Props> = ({ onOpenDetails, setSelectedLocalItem }) => {
  const [fileList, setFileList] = useState<UploadFile[]>([]);
  const [fileType, setFileType] = useState<FileType>("auto");

  const onRemove = useCallback((_file: UploadFile) => {
    setFileList([]);
  }, []);

  const setDataFormat = useCallback((type: FileType, filename: string) => {
    const extension = getExtension(filename);
    if (type === "auto") {
      // more exceptions will be added in the future
      switch (extension) {
        // 3dtiles
        case "json":
          return "json";
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

  const getAdditionalData = useCallback((content: string | undefined, format: string) => {
    if (!content) return undefined;
    if (format === "csv") {
      const header = content.split("\r\n")[0];
      const cols = header.split(",");
      const latColumn = cols.find(col =>
        ["latitude", "lat", "緯度", "北緯"].includes(col.toLowerCase()),
      );
      const lngColumn = cols.find(col =>
        ["longitude", "lng", "lon", "経度", "東経"].includes(col.toLowerCase()),
      );
      const heightColumn = cols.find(col =>
        ["height", "altitude", "alt", "高度"].includes(col.toLowerCase()),
      );
      if (!latColumn || !lngColumn) return undefined;
      return {
        data: {
          csv: {
            latColumn,
            lngColumn,
            heightColumn,
            noHeader: false,
          },
        },
      };
    }
    return undefined;
  }, []);

  const beforeUpload = useCallback(
    (file: RcFile, files: RcFile[]) => {
      const reader = new FileReader();
      reader.addEventListener(
        "load",
        () => {
          // convert image file to base64 string
          // Catalog Item
          const filename = file.name;
          const id = "id" + Math.random().toString(16).slice(2);
          const content = reader.result?.toString();
          const url = (() => {
            if (!content) {
              return;
            }
            return "data:text/plain;charset=UTF-8," + encodeURIComponent(content);
          })();
          const format = setDataFormat(fileType, filename);
          const item: UserDataItem = {
            type: "item",
            id: id,
            dataID: id,
            description:
              "このファイルは今お使いのWebブラウザでのみ閲覧可能です。共有URLを用いて共有するには、公開Webサーバー上のデータを読み込む必要があります。",
            name: filename,
            visible: true,
            url: url,
            format,
            additionalData: getAdditionalData(content, format),
          };
          if (onOpenDetails) onOpenDetails(item);
          if (setSelectedLocalItem) setSelectedLocalItem(item);
        },
        false,
      );

      if (file) {
        reader.readAsText(file, "UTF-8");
      }

      setFileList([...files]);
      return false;
    },
    [fileType, onOpenDetails, setDataFormat, setSelectedLocalItem, getAdditionalData],
  );

  const props: UploadProps = useMemo(
    () => ({
      name: "file",
      multiple: false,
      directory: false,
      showUploadList: true,
      accept: fileFormats,
      listType: "picture",
      onRemove: onRemove,
      beforeUpload: beforeUpload,
      fileList,
    }),
    [beforeUpload, fileList, onRemove],
  );

  const handleFileTypeSelect = useCallback((type: string) => {
    setFileType(type as FileType);
  }, []);

  return (
    <Form layout="vertical">
      <Form.Item name="file-type" label="ファイルタイプを選択">
        <FileTypeSelect onFileTypeSelect={handleFileTypeSelect} />
      </Form.Item>
      <Form.Item label="ファイルを選択">
        <Form.Item name="upload-file" style={{ height: 300, overflowY: "scroll" }}>
          <Upload.Dragger {...props}>
            <p className="ant-upload-drag-icon">
              <InboxOutlined />
            </p>
            <p className="ant-upload-text">
              ここをクリックしてファイルを選択するか、ファイルをここにドラッグ&amp;ドロップしてください。
            </p>
          </Upload.Dragger>
        </Form.Item>
      </Form.Item>
    </Form>
  );
};

export default LocalDataTab;
