export const ROUTES = {
    HOME: "/",
    LIST: "/info",
    LANG: "/lang/",
}

export type RouteKeyType = keyof typeof ROUTES;
export const ROUTE_LABELS: {[key in RouteKeyType]: string} = {
    HOME: "Главная",
    LIST: "Языки",
    LANG: "Язык",
};