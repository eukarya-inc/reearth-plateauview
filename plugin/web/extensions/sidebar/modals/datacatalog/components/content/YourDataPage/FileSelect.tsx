import { InboxOutlined } from "@ant-design/icons";
import { Tabs, Select, Input, Form } from "@web/sharedComponents";
import Upload, { message, UploadProps } from "@web/sharedComponents/Upload";
import { styled } from "@web/theme";

const FileSelectPane: React.FC = () => {
  const props: UploadProps = {
    name: "file",
    multiple: true,
    action: "https://www.mocky.io/v2/5cc8019d300000980a055e76",
    listType: "picture",
    onChange(info) {
      const { status } = info.file;
      if (status !== "uploading") {
        console.log(info.file, info.fileList);
      }
      if (status === "done") {
        message.success(`${info.file.name} file uploaded successfully.`);
      } else if (status === "error") {
        message.error(`${info.file.name} file upload failed.`);
      }
    },
    onDrop(e) {
      console.log("Dropped files", e.dataTransfer.files);
    },
  };

  const options = [
    {
      value: "auto",
      label: "Auto-detect (Recommended)",
    },
    {
      value: "geojson",
      label: "GeoJSON",
    },
    {
      value: "kml",
      label: "KML or KMZ",
    },
    {
      value: "csv",
      label: "CSV",
    },
    {
      value: "czml",
      label: "CZML",
    },
    {
      value: "gpx",
      label: "GPX",
    },
    {
      value: "json",
      label: "JSON",
    },
    {
      value: "georss",
      label: "GeoRSS",
    },
    {
      value: "gltf",
      label: "GLTF",
    },
    {
      value: "shapefile",
      label: "ShapeFile (zip)",
    },
  ];

  return (
    <Wrapper>
      <Tabs defaultActiveKey="local" style={{ marginBottom: "12px" }}>
        <Tabs.TabPane tab="Add Local Data" key="local">
          <Form layout="vertical">
            <Form.Item name="file-type" label="Select file type">
              <Select
                defaultValue="auto"
                style={{ width: "100%" }}
                // onChange={handleChange}
                options={options}
              />
            </Form.Item>
            <Form.Item label="Dragger">
              <Form.Item name="dragger" noStyle>
                <Upload.Dragger {...props}>
                  <p className="ant-upload-drag-icon">
                    <InboxOutlined />
                  </p>
                  <p className="ant-upload-text">Click or drag file to this area to upload</p>
                  <p className="ant-upload-hint">
                    Support for a single or bulk upload. Strictly prohibit from uploading company
                    data or other band files
                  </p>
                </Upload.Dragger>
              </Form.Item>
            </Form.Item>
          </Form>
        </Tabs.TabPane>
        <Tabs.TabPane tab="Add Web Data" key="web">
          <Form layout="vertical">
            <Form.Item name="file-type" label="Select file type">
              <Select
                defaultValue="auto"
                // onChange={handleChange}
                options={options}
              />
            </Form.Item>
            <Form.Item
              name="url"
              label="File URL"
              rules={[
                { required: true },
                { message: "Please input the URL of the asset!" },
                { type: "url", warningOnly: true },
              ]}>
              <Input placeholder={"Please input a valid URL"} />
            </Form.Item>
          </Form>
        </Tabs.TabPane>
      </Tabs>
    </Wrapper>
  );
};

export default FileSelectPane;

const Wrapper = styled.div`
  padding: 24px 12px;
`;
