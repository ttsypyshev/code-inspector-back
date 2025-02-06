import React from "react";
import { Link } from "react-router-dom";
import { Button } from "react-bootstrap";
import "./HomePage.css";
import "../../components/global.css";

const MainPage: React.FC = () => {
    return (
        <main className="main-container">
            <div className="home-body">
                <h1 className="home-title">Проверка кода</h1>
                <p className="home-desc">
                    Выберите языки, добавьте код, и вскоре преподаватель по достоинству оценит его.
                </p>
                <Link to="/info">
                    <Button variant="outline-light" className="btn-resume">
                        Продолжить
                    </Button>
                </Link>
            </div>
        </main>
    );
};

export default MainPage;
