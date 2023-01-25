import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { ReactNode } from "react";

type Props = {
  onClick?: () => void;
  onBack?: () => void;
  children?: ReactNode;
};

const PopupItem: React.FC<Props> = ({ onClick, onBack, children }) => {
  return (
    <Wrapper onClick={onClick}>
      {onBack && <StyledIcon icon="arrowLeft" onClick={onBack} />}
      {children}
    </Wrapper>
  );
};

export default PopupItem;

const Wrapper = styled.div`
  display: flex;
  justify-content: center;
  gap: 12px;
  padding: 12px;
  position: relative;
  color: #00bebe;
  background: #f4f4f4;
`;

const StyledIcon = styled(Icon)`
  position: absolute;
  left: 12px;
  cursor: pointer;
`;
