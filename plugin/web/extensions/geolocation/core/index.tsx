import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";

import useHooks from "./hooks";

const GeolocationWrapper: React.FC = () => {
  const { handleFlyToCurrentLocation } = useHooks();

  return (
    <Wrapper width={44} height={44}>
      <Icon icon="sun" size={20} onClick={handleFlyToCurrentLocation} />
    </Wrapper>
  );
};

export default GeolocationWrapper;

const Wrapper = styled.div<{ width?: number; height?: number }>`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  cursor: pointer;
  overflow: hidden;
  position: relative;
  width: ${({ width }) => width}px;
  height: ${({ height }) => height}px;
  background: #ececec;
  color: #c7c5c5;
`;
