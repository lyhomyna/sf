import MainPage from "./main/MainPage.jsx";
import LoginPage from "./login/LoginPage.jsx";
import { useState, useEffect } from "react";
import { RingLoader } from "react-spinners";
import { AuthContext } from "./storage/SfContext.jsx";
import DragAndDropSection from "./main/components/DragAndDropSection.jsx";

export const authServiceBaseUrl = "/api/auth"
export const fileServiceBaseUrl = "/api/files"

export default function App() {
    const [isAuthenticated, setAuthenticated] = useState(undefined);
    const [isLoading, setIsLoading] = useState(true);

    const changeAuthStatus = () => {
	setAuthenticated(current => !current);
    }

    useEffect(() => {
	(async () => {
	    const res = await fetch(`${authServiceBaseUrl}/check-auth`);
	    
	    if (res.status === 200) {
		setIsLoading(false);
		setAuthenticated(true);
		return;
	    }
	    setIsLoading(false);
	    setAuthenticated(false);
	})();
    }, []);

    let page;
    if (isLoading) {
	page = <div className="flex justify-center items-center h-screen">
	    <RingLoader color="#58a6d8"/>
	</div>;
    } else {
	if (isAuthenticated) {
	    page = <MainPage />;
	} else {
	    page = <div className="h-screen flex flex-col items-center justify-center">
	       <LoginPage />
	    </div>;
	}
    }

    return <AuthContext value={{changeAuthStatus: changeAuthStatus}}>
	<DragAndDropSection />
	{page}
    </AuthContext>;
}

