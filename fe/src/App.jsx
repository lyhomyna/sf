import MainPage from "./main/MainPage.jsx";
import LoginPage from "./login/LoginPage.jsx";
import { useState } from "react";

export const authServiceBaseUrl = "/api/auth"
export const fileServiceBaseUrl = "/api/files"

export default function App() {
    const [isAuthenticated, setAuthenticated] = useState(false);

    const authUser = () => {
	setAuthenticated(true);
    }

    return <>
	{isAuthenticated ?
	    <MainPage />
	    :
	    <div class="h-screen flex flex-col items-center justify-center">
		<LoginPage authenticate={authUser}/>
	    </div>
	}
    </>;
}
