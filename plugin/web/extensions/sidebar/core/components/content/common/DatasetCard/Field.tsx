import { Icon } from "@web/sharedComponents";
import { styled } from "@web/theme";
import {
  Accordion,
  AccordionItem,
  AccordionItemHeading,
  AccordionItemButton,
  AccordionItemPanel,
  AccordionItemState,
} from "react-accessible-accordion";

type Field = {
  id: string;
  title: string;
  icon?: string;
};

export type Props = {
  field: Field;
};

const FieldComponent: React.FC<Props> = ({ field }) => {
  return (
    <StyledAccordionComponent allowZeroExpanded>
      <AccordionItem>
        <AccordionItemState>
          {({ expanded }) => (
            <Header expanded={expanded}>
              <HeaderContents>
                <Title>{field.title}</Title>
                <ArrowIcon icon="arrowDown" size={16} expanded={expanded} />
              </HeaderContents>
            </Header>
          )}
        </AccordionItemState>
        <BodyWrapper>
          <Content>{field.icon && <Icon icon={field.icon} size={20} />}asdf</Content>
        </BodyWrapper>
      </AccordionItem>
    </StyledAccordionComponent>
  );
};

export default FieldComponent;

const StyledAccordionComponent = styled(Accordion)`
  width: 100%;
  border: 1px solid #e6e6e6;
  border-radius: 4px;
  background: #ffffff;
`;

const Header = styled(AccordionItemHeading)<{ expanded?: boolean }>`
  border-bottom-width: 1px;
  border-bottom-style: solid;
  border-bottom-color: transparent;
  ${({ expanded }) => expanded && "border-bottom-color: #e0e0e0;"}
  display: flex;
  height: 30px;
`;

const HeaderContents = styled(AccordionItemButton)`
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex: 1;
  padding: 0 12px;
  outline: none;
  cursor: pointer;
`;

const BodyWrapper = styled(AccordionItemPanel)<{ noTransition?: boolean }>`
  width: 100%;
  border-radius: 0px 0px 4px 4px;
  padding: 12px;
`;

const Title = styled.p`
  margin: 0;
`;

const ArrowIcon = styled(Icon)<{ expanded?: boolean }>`
  transition: transform 0.15s ease;
  transform: ${({ expanded }) => !expanded && "rotate(90deg)"};
`;

const Content = styled.div`
  display: flex;
  align-content: center;
  align-items: center;
  flex-wrap: wrap;
  gap: 6px;
  // padding-right: 12px;
`;
