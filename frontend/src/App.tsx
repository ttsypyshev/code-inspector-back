import { BrowserRouter, Route, Routes } from "react-router-dom";
import MainPage from "./pages/HomePage/HomePage"
import LangListPage from "./pages/LangsListPage/LangListPage"
import LangPage from "./pages/LangPage/LangPage";
import { ROUTES } from "./Routes";


function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path={ROUTES.HOME} element={<MainPage/>} />
                <Route path={ROUTES.LIST} element={<LangListPage/>} />
                <Route path={`${ROUTES.LANG}:id`} element={<LangPage/>} />
            </Routes>
        </BrowserRouter>
    );
}

export default App;