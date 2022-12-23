import { Icon, Button } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { ReactNode } from "react";

type Props = {
  width?: number;
  height?: number;
  children?: ReactNode;
  handleClose?: () => void;
};

const PopupWrapper: React.FC<Props> = ({ width, height, children, handleClose }) => {
  return (
    <Wrapper width={width} height={height}>
      <Header>
        <CloseButton>
          <Icon size={32} icon="close" onClick={handleClose} />
        </CloseButton>
      </Header>
      {children}
    </Wrapper>
  );
};

export default PopupWrapper;

const Wrapper = styled.div<{ width?: number; height?: number }>`
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: 0px;
  background: #e7e7e7;
  width: ${width => (width ? "100%" : width + "px")};
  height: ${height => (height ? "100%" : height + "px")};
`;
const Header = styled.div`
  display: flex;
  background: #e7e7e7;
  position: relative;
  width: 350px;
  height: 32px;
`;

const CloseButton = styled(Button)`
  display: flex;
  justify-content: center;
  align-items: center;
  position: absolute;
  right: 0;
  height: 32px;
  width: 32px;
  border: none;
  background: #00bebe;
  color: white;
  cursor: pointer;
`;
