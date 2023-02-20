import { BaseFieldProps } from "../types";

const PolygonStroke: React.FC<BaseFieldProps<"polygonStroke">> = ({
  value,
  editMode,
  onUpdate,
}) => {
  // remember to update the BaseFieldProps type!
  console.log(value, editMode, onUpdate);
  return <div>Polygon Stroke</div>;
};

export default PolygonStroke;
