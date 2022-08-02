import { Image } from "antd";
import svgToMiniDataURI from "mini-svg-data-uri";
import { ComponentProps, memo, useMemo } from "react";
import { ReactSVG } from "react-svg";

import Icons from "./icons";

export type Icons = keyof typeof Icons;

export type Props = {
  className?: string;
  icon?: string;
  height?: string | number;
  width?: string | number;
  alt?: string;
  color?: string;
  onClick?: () => void;
};

const Icon: React.FC<Props> = ({
  className,
  icon,
  alt,
  height,
  width,
  onClick,
  color,
}) => {
  const src = useMemo(
    () =>
      icon?.startsWith("<svg ") ? svgToMiniDataURI(icon) : Icons[icon as Icons],
    [icon]
  );
  if (!icon) return null;

  const heightStr = typeof height === "number" ? `${height}px` : height;
  const widthStr = typeof width === "number" ? `${width}px` : width;
  const iconColor = color ? color || " transparent" : " transparent";

  if (!src) {
    return (
      <Image
        src={icon}
        alt={alt}
        height={heightStr}
        width={widthStr}
        onClick={onClick}
        color={iconColor}
        className={className}
      />
    );
  }

  return (
    <SVG
      className={"StyledSvg"}
      src={src}
      style={{ height: heightStr, width: widthStr, color: iconColor }}
      alt={alt}
      onClick={onClick}
    />
  );
};

const SVG: React.FC<
  Pick<
    ComponentProps<typeof ReactSVG>,
    "className" | "src" | "onClick" | "alt" | "style"
  >
> = (props) => {
  return <ReactSVG {...props} wrapper="span" />;
};

export default memo(Icon);
