import { Switch } from "@web/sharedComponents";
import { ComponentProps } from "react";

import { FieldTitle, BorderlessFieldValue, FieldWrapper } from "../commonComponents";

type Props = {
  title: string;
  titleWidth?: number;
} & ComponentProps<typeof Switch>;

const SwitchField: React.FC<Props> = ({ title, titleWidth, ...props }) => {
  return (
    <FieldWrapper>
      <FieldTitle width={titleWidth}>{title}</FieldTitle>
      <BorderlessFieldValue>
        <Switch {...props} />
      </BorderlessFieldValue>
    </FieldWrapper>
  );
};

export default SwitchField;
