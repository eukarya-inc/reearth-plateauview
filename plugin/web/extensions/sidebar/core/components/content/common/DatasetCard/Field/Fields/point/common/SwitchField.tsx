import { Switch } from "@web/sharedComponents";
import { ComponentProps } from "react";

import { FieldTitle, BorderLessFieldValue, FieldWrapper } from "./styled";

type Props = {
  title: string;
  titleWidth?: number;
} & ComponentProps<typeof Switch>;

const SwitchField: React.FC<Props> = ({ title, titleWidth, ...props }) => {
  return (
    <FieldWrapper>
      <FieldTitle width={titleWidth}>{title}</FieldTitle>
      <BorderLessFieldValue>
        <Switch {...props} />
      </BorderLessFieldValue>
    </FieldWrapper>
  );
};

export default SwitchField;
