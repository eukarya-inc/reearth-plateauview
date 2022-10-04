import { Icon } from "@web/extensions/sharedComponents";
import { styled } from "@web/theme";

const Selection: React.FC = () => {
  return (
    <Wrapper>
      <StyledButton
        onClick={() => alert("This is an awesome datacatalog modal!! Use me...if you can!")}>
        <StyledIcon icon="plusCircle" size={20} />
        <ButtonText>Explore map data</ButtonText>
      </StyledButton>
    </Wrapper>
  );
};

export default Selection;

const Wrapper = styled.div`
  height: 100%;
  padding: 16px;
`;

const StyledButton = styled.button`
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  border: none;
  border-radius: 4px;
  background: #00bebe;
  color: #fff;
  padding: 10px;
  cursor: pointer;
`;

const ButtonText = styled.p`
  margin: 0;
`;

const StyledIcon = styled(Icon)`
  margin-right: 8px;
`;
