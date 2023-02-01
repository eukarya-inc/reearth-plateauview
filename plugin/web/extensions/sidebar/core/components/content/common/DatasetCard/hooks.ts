import { Data } from "@web/extensions/sidebar/core/newTypes";

export default ({
  dataset,
  inEditor,
  onDatasetUpdate,
}: {
  dataset: Data;
  inEditor?: boolean;
  onDatasetUpdate: (dataset: Data) => void;
}) => {
  const handleAddField =
    (property: any) =>
    ({ key }: { key: string }) => {
      if (!inEditor) return;
      onDatasetUpdate?.({
        ...dataset,
        components: [
          ...(dataset.components ?? []),
          {
            type: key,
            ...property,
          },
        ],
      });
    };

  const generalFields: {
    [key: string]: { name: string; onClick: (property: any) => void };
  } = {
    camera: {
      name: "カメラ",
      onClick: () =>
        handleAddField({
          position: {
            lng: 0,
            lat: 0,
            height: 0,
            pitch: 0,
            heading: 0,
            roll: 0,
          },
        }),
    },
    description: {
      name: "説明",
      onClick: () => handleAddField({}),
    },
    legend: {
      name: "凡例",
      onClick: () => handleAddField({ style: "square", items: [{ title: "hey", color: "red" }] }),
    },
  };

  const pointFields: {
    [key: string]: { name: string; onClick?: (property: any) => void };
  } = {
    camera: {
      name: "カメラ",
      onClick: () =>
        handleAddField({
          position: {
            lng: 0,
            lat: 0,
            height: 0,
            pitch: 0,
            heading: 0,
            roll: 0,
          },
        }),
    },
  };

  const polylineFields: {
    [key: string]: { name: string; onClick?: (property: any) => void };
  } = {
    camera: {
      name: "カメラ",
      onClick: () =>
        handleAddField({
          position: {
            lng: 0,
            lat: 0,
            height: 0,
            pitch: 0,
            heading: 0,
            roll: 0,
          },
        }),
    },
  };

  const polygonFields: {
    [key: string]: { name: string; onClick?: (property: any) => void };
  } = {
    camera: {
      name: "カメラ",
      onClick: () =>
        handleAddField({
          position: {
            lng: 0,
            lat: 0,
            height: 0,
            pitch: 0,
            heading: 0,
            roll: 0,
          },
        }),
    },
  };

  const ThreeDModelFields: {
    [key: string]: { name: string; onClick?: (property: any) => void };
  } = {
    camera: {
      name: "カメラ",
      onClick: () =>
        handleAddField({
          position: {
            lng: 0,
            lat: 0,
            height: 0,
            pitch: 0,
            heading: 0,
            roll: 0,
          },
        }),
    },
  };

  const ThreeDTileFields: {
    [key: string]: { name: string; onClick?: (property: any) => void };
  } = {
    camera: {
      name: "カメラ",
      onClick: () =>
        handleAddField({
          position: {
            lng: 0,
            lat: 0,
            height: 0,
            pitch: 0,
            heading: 0,
            roll: 0,
          },
        }),
    },
  };

  const fieldGroups: {
    [key: string]: {
      name: string;
      fields: { [key: string]: { name: string; onClick?: (property: any) => void } };
    };
  } = {
    general: {
      name: "一般",
      fields: generalFields,
    },
    point: {
      name: "ポイント",
      fields: pointFields,
    },
    polyline: { name: "ポリライン", fields: polylineFields },
    polygone: { name: "ポリゴン", fields: polygonFields },
    "3d-model": { name: "3Dモデル", fields: ThreeDModelFields },
    "3d-tile": { name: "3Dタイル", fields: ThreeDTileFields },
  };
  return {
    fieldGroups,
  };
};
