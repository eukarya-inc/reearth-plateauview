import { Select } from "@web/sharedComponents";

export type Props = {};

const FileTypeSelect: React.FC<Props> = () => {
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

  const handleChange = (_value: string) => {
    console.log(_value);
  };

  return (
    <Select
      defaultValue="auto"
      style={{ width: "100%" }}
      onChange={handleChange}
      options={options}
    />
  );
};

export default FileTypeSelect;
