import MainPage from "main/MainPage.jsx";
import LoginPage from "login/LoginPage.jsx";
import { useState, useEffect } from "react";
import { RingLoader } from "react-spinners";
import { AuthContext } from "storage/SfContext.jsx";
import { authServiceBaseUrl } from "config/constants.js";

export default function App() {
    const [isAuthenticated, setAuthenticated] = useState(undefined);
    const [isLoading, setIsLoading] = useState(true);

    const changeAuthStatus = () => {
	setAuthenticated(current => !current);
    }

    // Auth check
    useEffect(() => {
	(async () => {
	    const res = await fetch(`${authServiceBaseUrl}/check-auth`);
	    if (res.status === 200) {
		setIsLoading(false);
		setAuthenticated(true);
	    } else {
		setIsLoading(false);
		setAuthenticated(false);
	    }
	})();
    }, []);

    // show either main page or loader
    let page;
    if (isLoading) {
	page = <div className="flex justify-center items-center h-screen">
	    <RingLoader color="#58a6d8"/>
	</div>;
    } else {
	if (isAuthenticated) {
	    page = <>
		<MainPage/>
	    </>;
	} else {
	    page = <div className="h-screen flex flex-col items-center justify-center">
	       <LoginPage />
	    </div>;
	}
    }

    return <AuthContext value={{changeAuthStatus: changeAuthStatus}}>
	{page}
    </AuthContext>;
}

