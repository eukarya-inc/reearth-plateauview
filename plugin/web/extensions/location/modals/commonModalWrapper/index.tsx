import { Button, Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import { ReactNode } from "react";

export type Props = {
  title?: string;
  children?: ReactNode;
  onModalChange: () => void;
};

const CommonModalWrapper: React.FC<Props> = ({ title, children, onModalChange }) => {
  return (
    <Wrapper>
      {title && (
        <>
          <HeaderWrapper>
            <Title>{title}</Title>
            <CloseButton>
              <Icon size={32} icon="close" color="grey" onClick={onModalChange} />
            </CloseButton>
          </HeaderWrapper>
        </>
      )}

      <ContentWrapper>{children}</ContentWrapper>

      <FooterWrapper>
        <OkButton type="primary" onClick={onModalChange}>
          ok
        </OkButton>
      </FooterWrapper>
    </Wrapper>
  );
};

export default CommonModalWrapper;

const Wrapper = styled.div`
  padding: 32px 16px;
`;

const Title = styled.p`
  font-size: 16px;
`;

const CloseButton = styled(Button)`
  right: 0;
  border: none;
  height: 48px;
  width: 48px;
  cursor: pointer;
  transition: background 0.3s;
`;

const HeaderWrapper = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  padding: 16px 24px;
  gap: 36px;
  width: 572px;
  height: 56px;
  background: #ffffff;
  box-shadow: inset 0px -1px 0px #f0f0f0;
  flex: none;
  order: 0;
  align-self: stretch;
  flex-grow: 0;
`;

const ContentWrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  padding: 24px;
  gap: 10px;
  width: 572px;
  height: 562px;
  flex: none;
  order: 1;
  align-self: stretch;
  flex-grow: 0;
  padding-bottom: 0px;
`;

const FooterWrapper = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  align-items: center;
  padding: 10px 16px;
  gap: 8px;
  width: 572px;
  height: 52px;
  background: #ffffff;
  box-shadow: inset 0px 1px 0px #f0f0f0;
  flex: none;
  order: 2;
  align-self: stretch;
  flex-grow: 0;
`;

const OkButton = styled(Button)`
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 5px 16px;
  width: 51px;
  height: 32px;
  background: #1890ff;
  box-shadow: 0px 2px 0px rgba(0, 0, 0, 0.043);
  border-radius: 2px;
  flex: none;
  order: 0;
  flex-grow: 0;
`;
