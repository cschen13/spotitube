import * as React from "react";
import styled from "styled-components";

interface ILogoProps {
  className?: string;
}

export const Logo: React.FunctionComponent<ILogoProps> = ({ className }) => (
  <h1 className={className}>
    <PlaylistSpan>PLAYLIST</PlaylistSpan>
    <XSpan>X</XSpan>
    <ChangeSpan>CHANGE</ChangeSpan>
  </h1>
);

export const LogoCondensed: React.FunctionComponent<ILogoProps> = ({
  className
}) => (
  <h1 className={className}>
    <PlaylistSpan>P</PlaylistSpan>
    <XSpan>X</XSpan>
    <ChangeSpan>C</ChangeSpan>
  </h1>
);

const PlaylistSpan = styled.span`
  position: relative;
  top: 0.25em;
  font-size: 0.33em;
`;

const XSpan = styled.span`
  position: relative;
  top: 0.5em;
`;

const ChangeSpan = styled.span`
  position: relative;
  top: 1.5em;
  font-size: 0.33em;
  color: purple;
`;
