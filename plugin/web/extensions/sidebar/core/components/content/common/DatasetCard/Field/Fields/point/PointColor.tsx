import { BaseFieldProps } from "../types";

const PointColor: React.FC<BaseFieldProps<"pointColor">> = ({ value, editMode, onUpdate }) => {
  console.log(value, editMode, onUpdate);
  return editMode ? <div>PointColor</div> : null;
};

export default PointColor;
