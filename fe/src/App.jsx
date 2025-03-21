import MainPage from "./main/MainPage.jsx";
import LoginPage from "./login/LoginPage.jsx";
import { useState, useEffect } from "react";
import { RingLoader } from "react-spinners";

export const authServiceBaseUrl = "/api/auth"
export const fileServiceBaseUrl = "/api/files"

export default function App() {
    const [isAuthenticated, setAuthenticated] = useState(undefined);

    console.log("AUTH STATUS: ", isAuthenticated);

    const authUser = () => {
	setAuthenticated(true);
    }

    useEffect(() => {
	const checkAuth = async () => {
	    const res = await fetch(`${authServiceBaseUrl}/check-auth`);
	    
	    if (res.status === 200) {
		console.log("User authenticated. Status 200")
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
	   <LoginPage authenticate={authUser}/>
	</div>;
    } else {
	page = <MainPage />;
    }

    return page;
}


	//{isAuthenticated ?
	 //   <MainPage />
	 //   :
	 //   <div className="h-screen flex flex-col items-center justify-center">
	//	<LoginPage authenticate={authUser}/>
	 //   </div>
	//}
