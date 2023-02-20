import { BaseFieldProps } from "../types";

const PolylineColor: React.FC<BaseFieldProps<"polylineColor">> = ({
  value,
  editMode,
  onUpdate,
}) => {
  console.log(value, onUpdate);
  return editMode ? <>Polyline Color</> : null;
};

export default PolylineColor;
