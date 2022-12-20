import basicOperation from "@web/extensions/sidebar/core/assets/basicOperation.png";
import { styled } from "@web/theme";

import useGlobalHooks from "../globalHooks";
import { PopupWrapper } from "../sharedComponent";

const BasicOperation: React.FC = () => {
  const { handleClosePopup } = useGlobalHooks();
  return (
    <PopupWrapper handleClose={handleClosePopup}>
      <Wrapper>
        <img src={basicOperation} />
      </Wrapper>
    </PopupWrapper>
  );
};

export default BasicOperation;

const Wrapper = styled.div`
  width: 318px;
  height: 457px;
  padding: 16px;
`;
