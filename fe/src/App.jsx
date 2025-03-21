import MainPage from "./main/MainPage.jsx";
import LoginPage from "./login/LoginPage.jsx";
import { useState, useEffect } from "react";
import { RingLoader } from "react-spinners";
import { AuthContext } from "./storage/SfContext.jsx";

export const authServiceBaseUrl = "/api/auth"
export const fileServiceBaseUrl = "/api/files"

export default function App() {
    const [isAuthenticated, setAuthenticated] = useState(undefined);

    const changeAuthStatus = () => {
	setAuthenticated(current => !current);
    }

    useEffect(() => {
	const checkAuth = async () => {
	    const res = await fetch(`${authServiceBaseUrl}/check-auth`);
	    
	    if (res.status === 200) {
		setAuthenticated(true);
		return;
	    }
	    setAuthenticated(false);
	};
	checkAuth();
    }, []);

    let page = <div className="flex justify-center items-center h-screen">
	<RingLoader color="#58a6d8"/>
    </div>;

    if (isAuthenticated != undefined && !isAuthenticated) {
	page = <div className="h-screen flex flex-col items-center justify-center">
	   <LoginPage />
	</div>;
    } else {
	page = <MainPage />;
    }

    return <AuthContext value={{changeAuthStatus: changeAuthStatus}}>
	{page}
    </AuthContext>;
}

