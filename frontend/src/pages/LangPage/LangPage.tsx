import React, { useEffect, useState } from 'react';
import { Langs_Mock } from '../../modules/mock';
import { Lang } from '../../modules/types';
import { useParams } from 'react-router-dom';
import "./LangPage.css";
import "../../components/global.css";
import BackButton from '../../components/BackButton/BackButton.tsx';

const LangPage: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const [langData, setLangData] = useState<Lang | null>(null);
    const [isMock, setIsMock] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);

    const fetchLangData = async () => {
        try {
            const response = await fetch(`/api/info/${id}`, { signal: AbortSignal.timeout(5000) });
            if (!response.ok) throw new Error('Network response failed');

            const result = await response.json();
            const lang = {
                id: result.info.ID,
                name: result.info.Name,
                shortDescription: result.info.ShortDescription,
                description: result.info.Description,
                imgLink: result.info.ImgLink,
                author: result.info.Author,
                year: result.info.Year,
                version: result.info.Version,
                list: result.info.List,
                status: result.info.Status,
            };

            setLangData(lang);
            setError(null);
        } catch (error) {
            console.error('Fetch error or invalid data:', error);
            setIsMock(true);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        const loadLangData = () => {
            const idNum = parseInt(id as string, 10);
            const mockLang = Langs_Mock.find(lang => lang?.id === idNum) as Lang;
            if (mockLang) {
                setLangData(mockLang);
                setError(null);
            } else {
                setError('Language not found in mock data');
            }
        };

        if (isMock) {
            loadLangData();
        } else {
            fetchLangData();
        }

        return () => {
            setLangData(null); // Cleanup
        };
    }, [id, isMock]);

    if (loading) return <div className="loading">Loading...</div>;
    if (error) return <div className="error">{error}</div>;
    if (!langData) return <div className="not-found">Language not found</div>;

    return (
        <div className="body">
            <div className="information-about">
                <div className="header">
                    <BackButton link="/info" />
                </div>

                <LangDetails lang={langData} />
            </div>
        </div>
    );
};

const LangDetails: React.FC<{ lang: Lang }> = ({ lang }) => (
    <div className="info">
        <div className="title">
            <div className="text-name-lang">{lang.name}</div>
            <div className="text-description-lang">"{lang.shortDescription}"</div>
            <img className="img-lang" src={lang.imgLink} alt={lang.name} />
        </div>

        <div className="background-basic-info">
            <LangBasicInfo label="Автор" value={lang.author || "N/A"} number={1} />
            <LangBasicInfo label="Год выпуска" value={lang.year || "N/A"} number={2} />
            <LangBasicInfo label="Последняя версия" value={lang.version || "N/A"} number={3} />
        </div>


        <div className="background-additional-info">
            <div className="text-additional-info">
                <span className="bold-subtext">{lang.name}</span>
                <span className="subtext">{lang.description}</span>
            </div>

            <FeatureList features={lang.list} />
        </div>
    </div>
);

const LangBasicInfo: React.FC<{ 
    label: string; 
    value: string | number; 
    number: number; 
}> = ({ label, value, number }) => (
    <div>
        <div className={`tag-basic-info tag-basic-info-${number}`}>{label}:</div>
        <div className={`subtext-basic-info subtext-basic-info-${number}`}>{value}</div>
    </div>
);


const FeatureList: React.FC<{ features: Record<string, string> }> = ({ features }) => (
    <ul className="features-list">
        {Object.entries(features).map(([key, value]) => (
            <li key={key}>
                <span className="bold-subtext">{key}:</span>
                <span className="subtext">{value}</span>
            </li>
        ))}
    </ul>
);

export default LangPage;
