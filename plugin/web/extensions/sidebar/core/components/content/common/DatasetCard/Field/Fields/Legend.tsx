import { BaseField as BaseFieldProps } from ".";

type Props = BaseFieldProps<"legend"> & {};

const Legend: React.FC<Props> = () => {
  return <p>I AM LEGEND</p>;
};

export default Legend;
