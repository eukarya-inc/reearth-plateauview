import { Data } from "@web/extensions/sidebar/core/newTypes";
import { MenuProps } from "@web/sharedComponents";

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
    (property: any): MenuProps["onClick"] =>
    ({ key }) => {
      if (!inEditor) return;
      console.log(property);
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

  const items: MenuProps["items"] = [
    {
      key: "general",
      label: "一般",
      children: [
        {
          key: "camera",
          label: "カメラ",
          onClick: handleAddField({
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
        {
          key: "description",
          label: "説明",
          onClick: handleAddField({ markdownKey: "" }),
        },
        {
          key: "legend",
          label: "凡例",
          onClick: handleAddField({ style: "square" }),
        },
      ],
    },
    {
      key: "point-group",
      label: "ポイント",
      children: [
        {
          key: "point",
          label: "Fill color (condition)",
        },
        // {
        //   key: "point",
        //   label: "Fill color (gradient)",
        // },
        // {
        //   key: "point",
        //   label: "Stroke",
        // },
        // {
        //   key: "point",
        //   label: "Icon",
        // },
        // {
        //   key: "point",
        //   label: "Size",
        // },
        // {
        //   key: "point",
        //   label: "Label",
        // },
        // {
        //   key: "point",
        //   label: "3D Model",
        // },
      ],
    },
  ];

  return {
    items,
  };
};
