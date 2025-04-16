import { CircleLoader } from "react-spinners";
import ControlButtons from "./ControlButtons.jsx";
import { authServiceBaseUrl } from "config/constants.js";
import { useEffect, useState } from "react";

export default function TopBar() {
    const imageURL="https://static.vecteezy.com/system/resources/previews/018/871/797/non_2x/happy-cat-transparent-background-png.png"
    const [email, setEmail] = useState("supercoolemail@super.mail");

    // to retrieve user's email
    useEffect(() => {
	const fetchUser = async () => {
	    try {
		const response = await fetch(`${authServiceBaseUrl}/get-user`);

		if (response.status === 200) {
		    const user = await response.json(); 
		    setEmail(user.email);
		} else if (response.status === 404) {
		    const err = await response.json();
		    console.log(err.message)
		}
	    } catch (err) {
		console.log("Failed to fetch user", err)
	    }
	}

	fetchUser();
    }, [])

    return <div className="flex flex-row justify-between items-center">
	<div className="flex flex-row items-center gap-5 flex-wrap">
	    <div className="w-16 h-16 bg-stone-700 rounded-md shadow-md">
		<img src={ imageURL } alt="Avatar"/>
	    </div>
	    <div className="font-medium text-slate-300">
		{ email }
	    </div>
	    <ControlButtons className="flex items-center gap-x-1" />
	</div>
	<div className="relative group inline-block cursor-pointer pr-5">
	    <CircleLoader size="30px" color="#fff"/>
	    <span className="absolute top-full right-0 px-2 py-1 bg-black text-white text-sm rounded opacity-0 group-hover:opacity-100 pointer-events-none max-w-[200px] break-words z-50">
		Your files is uploading...
	    </span>
	</div>
    </div>;
}
