export interface Lang {
    id: number;                  // ID языка, уникальный
    name: string;                // Название языка
    shortDescription: string;    // Краткое описание
    description: string;         // Полное описание
    imgLink?: string;            // Ссылка на изображение (опционально)
    author?: string;             // Автор (опционально)
    year?: string;               // Год создания (опционально)
    version?: string;            // Версия языка (опционально)
    list: Record<string, string>;   // Данные в формате JSONB
    status: boolean;             // Статус (активен/неактивен)
}
