import { Image } from "antd";
import svgToMiniDataURI from "mini-svg-data-uri";
import React,{ AriaAttributes, AriaRole, memo, useMemo } from "react";
import SVG from "react-inlinesvg";

import { ariaProps } from "../../../utile/aria";

import Icons from "./icons";

export type Icon = keyof typeof Icons;

export type Props = {
  className?: string;
  icon?: string;
  height?: string | number;
  width?: string | number;
  alt?: string;
  color?: string;
  role?: AriaRole;
  onClick?: () => void;
} & AriaAttributes;

const Icon: React.FC<Props> = ({
  className,
  icon,
  alt,
  height,
  width,
  onClick,
  color,
  role,
  ...props
}) => {
  const src = useMemo(
    () => (icon?.startsWith("<svg ") ? svgToMiniDataURI(icon) : Icons[icon as Icon]),
    [icon],
  );
  if (!icon) return null;

  const heightStr = typeof height === "number" ? `${height}px` : height;
  const widthStr = typeof width === "number" ? `${width}px` : width;
  const iconColor = color ? color : " transparent";
  const aria = ariaProps(props);

  if (!src) {
    return (
      <Image
        src={icon}
        alt={alt}
        height={heightStr}
        width={widthStr}
        onClick={onClick}
        role={role}
        color={iconColor}
        className={className}
        {...aria}
      />
    );
  }

  return (
    <SVG
      className={"StyledSvg"}
      src={src}
      style={{ height: heightStr, width: widthStr }}
      color={iconColor}
      onClick={onClick}
      role={role}
      {...aria}
    />
  );
};

export default memo(Icon);
