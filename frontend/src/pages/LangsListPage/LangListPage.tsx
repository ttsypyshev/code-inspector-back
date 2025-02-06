import { FormEvent, useEffect, useState } from "react";
import { Lang } from "../../modules/types.ts";
import { Langs_Mock } from "../../modules/mock.ts";
import "./LangListPage.css";
import "../../components/global.css";
import BackButton from '../../components/BackButton/BackButton.tsx';
import LanguageItem from '../../components/LanguageItem/LanguageItem.tsx';

const LangListPage = () => {
    const [languages, setLanguages] = useState<Lang[]>([]);
    const [isMock, setIsMock] = useState(false);
    const [name, setName] = useState('');
    const [cartCount, setCartCount] = useState(0);
    const [draftID, setDraftID] = useState(0);
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);

    const fetchData = async () => {
        setLoading(true);
        setError(null);

        try {
            const response = await fetch(`api/info?langname=${name.toLowerCase()}`, { signal: AbortSignal.timeout(5000) });
            if (!response.ok) throw new Error('Network response failed');
            const result = await response.json();

            const fetchedLanguages = result.langs.map((lang: any) => ({
                id: lang.ID,
                name: lang.Name,
                shortDescription: lang.ShortDescription,
                description: lang.Description,
                imgLink: lang.ImgLink,
                author: lang.Author,
                year: lang.Year,
                version: lang.Version,
                list: lang.List,
                status: lang.Status,
            }));

            setLanguages(fetchedLanguages);
            setCartCount(result.count || 0);
            setDraftID(result.draftID);
            setIsMock(false);
        } catch (error) {
            console.error("Fetch error:", error);
            setIsMock(true);
            // setError("Ошибка загрузки данных.");
            if (!isMock) {
                loadMockData();
            }
        } finally {
            setLoading(false);
        }
    };

    const loadMockData = () => {
        setIsMock(true);
        const filteredMocks = Langs_Mock.filter(lang =>
            lang.name.toLowerCase().includes(name.toLowerCase())
        );
        setLanguages(filteredMocks);
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        await fetchData();
    };

    useEffect(() => {
        fetchData();
    }, []);

    return (
        <div className="body">
            <div className="services-container">
                <div className="header">
                    {/* Кнопка назад */}
                    <BackButton link="/" />

                    {/* Строка поиска */}
                    <SearchField
                        name={name}
                        onNameChange={setName}
                        onSubmit={handleSubmit}
                    />

                    {/* Состояние корзины */}
                    <CartState cartCount={cartCount} draftID={draftID} />

                    {/* Разделительная линия */}
                    <img className="separator-line" src="../../../img/line.png" alt="separator" />
                </div>

                {/* Основной контейнер */}
                <ul className="service-list">
                    {loading ? (
                        <div className="loading">Загрузка...</div>
                    ) : error ? (
                        <div className="error">{error}</div>
                    ) : languages.length > 0 ? (
                        languages.map((lang) => (
                            <LanguageItem key={lang.id} lang={lang} />
                        ))
                    ) : (
                        <div className="error">
                            <h1>Данный язык не найден</h1>
                        </div>
                    )}
                </ul>
            </div>
        </div>
    );
};

const SearchField: React.FC<{
    name: string;
    onNameChange: (name: string) => void;
    onSubmit: (e: FormEvent) => void;
}> = ({ name, onNameChange, onSubmit }) => (
    <div className="search-section">
        <form onSubmit={onSubmit}>
            <input
                type="text"
                name="langname"
                className="field-search-text"
                maxLength={100}
                placeholder="Поиск..."
                value={name}
                onChange={(e) => onNameChange(e.target.value)}
            />
            <button type="submit" style={{ display: 'none' }}>Поиск</button>
            <img className="search-icon" src="../../../img/icon-find.png" alt="Иконка поиска" />
        </form>
    </div>
);

const CartState: React.FC<{ cartCount: number; draftID: number }> = ({ cartCount, draftID }) => (
    <div className="file-count-section">
        {cartCount !== 0 ? (
            <a href={`/project/${draftID}`} className="file-count">
                <img className="file-count-icon" src="../../../img/icon-count-files.png" alt="files" />
                <div className="file-count-text">{cartCount}</div>
            </a>
        ) : (
            <div className="file-count">
                <img className="file-count-icon" src="../../../img/icon-count-files.png" alt="files" />
                <div className="file-count-text">{cartCount}</div>
            </div>
        )}
    </div>
);

export default LangListPage;
