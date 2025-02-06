import React from 'react';
import './LanguageItem.css';
import { Lang } from "../../modules/types";

interface LangItemProps {
  lang: Lang;
}

const LangItem: React.FC<LangItemProps> = ({ lang }) => {
    const defaultImgLink = "../../../img/default-lang-image.png"; // Путь к изображению по умолчанию
    const imgSrc = lang.imgLink || defaultImgLink; // Если imgLink пустой, используем defaultImgLink
    
    return (
        <li key={lang.id} className="service-item">
            {/* Link to detailed information */}
            <a href={`/lang/${lang.id}`} className="service-link">
                <div className="service-title">{lang.name}</div>
                <div className="service-description">{lang.shortDescription}</div>
                <img
                    className="service-image"
                    src={imgSrc} // Используем imgSrc для изображения
                    alt={lang.name}
                />
            </a>

            {/* Form to add to project */}
            <form action="/add-service" method="POST" className="add-to-project-form">
                <input type="hidden" name="id_lang" value={lang.id} />
                <button
                    type="submit"
                    className="add-to-project-button-background"
                    // Добавьте динамическую классовую логику по статусу, если нужно
                >
                    <div className="add-to-project-button-text">Добавить в проект</div>
                </button>
            </form>
        </li>
    );
};

export default LangItem;
