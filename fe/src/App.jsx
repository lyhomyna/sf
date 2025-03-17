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

    (async () => {
	const res = await fetch(`${authServiceBaseUrl}/check-auth`);
	
	if (res.status === 200) {
	    setAuthenticated(true);
	}
	setAuthenticated(false);
    })();

    return <>
	{isAuthenticated ?
	    <MainPage />
	    :
	    <div className="h-screen flex flex-col items-center justify-center">
		<LoginPage authenticate={authUser}/>
	    </div>
	}
    </>;
}
