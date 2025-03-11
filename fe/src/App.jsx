import MainPage from "./main/MainPage.jsx";
import LoginPage from "./login/LoginPage.jsx";

export const fileServiceBaseUrl = "http://file-service:8082"

export default function App() {
    return <div class="h-screen flex flex-col items-center justify-center"><LoginPage /></div>;
}
