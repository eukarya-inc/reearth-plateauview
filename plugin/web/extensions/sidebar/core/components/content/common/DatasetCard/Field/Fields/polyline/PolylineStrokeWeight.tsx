import { BaseFieldProps } from "../types";

const PolylineStrokeWeight: React.FC<BaseFieldProps<"polylineStrokeWeight">> = ({
  value,
  editMode,
  onUpdate,
}) => {
  console.log(value, onUpdate);
  return editMode ? <>Polyline Stroke Weight</> : null;
};

export default PolylineStrokeWeight;
