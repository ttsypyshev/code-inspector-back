import React from 'react';
import './BackButton.css';

interface BackButtonProps {
  link: string;
}

const BackButton: React.FC<BackButtonProps> = ({ link }) => {
  return (
    <a href={link} className="background-btn-back">
      <img className="icon-btn-back" src="../../../img/icon-btn-back.png" alt="Back" />
    </a>
  );
};

export default BackButton;
