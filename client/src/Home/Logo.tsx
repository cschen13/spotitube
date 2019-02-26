import * as React from "react";
import styled from "styled-components";

interface ILogoProps {
  className?: string;
}

export const Logo: React.FunctionComponent<ILogoProps> = ({ className }) => (
  <LogoContainer className={className}>
    <PlaylistSpan>PLAYLIST</PlaylistSpan>
    <XSpan>X</XSpan>
    <ChangeSpan>CHANGE</ChangeSpan>
  </LogoContainer>
);

export const LogoCondensed: React.FunctionComponent<ILogoProps> = ({
  className
}) => (
  <LogoContainer className={className}>
    <PlaylistSpan>P</PlaylistSpan>
    <XSpan>X</XSpan>
    <ChangeSpan>C</ChangeSpan>
  </LogoContainer>
);

const LogoContainer = styled.h1`
  height: 3.5em;
`;

const PlaylistSpan = styled.span`
  position: relative;
  top: 0.25em;
`;

const XSpan = styled.span`
  position: relative;
  font-size: 3em;
  top: 0.5em;
`;

const ChangeSpan = styled.span`
  position: relative;
  top: 1.5em;
  color: purple;
`;
