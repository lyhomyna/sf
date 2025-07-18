import { useRef, useContext, useState } from "react";
import { authServiceBaseUrl } from "config/constants.js";
import { AuthContext } from "storage/SfContext.jsx";
import HidePasswordButton from "./components/HidePasswordButton";

export default function LoginPage() {
    const { changeAuthStatus } = useContext(AuthContext);

    const [isPasswordVisible, setIsPasswordVisible] = useState(false);
    function togglePasswordVisibility(e) {
	e.preventDefault();
	setIsPasswordVisible((prevState) => !prevState);
    }

    const emailRef = useRef("");
    const passwordRef = useRef("");

    const doLogIn = async (e) => {
	e.preventDefault();

	const email = emailRef.current.value;
	const password = passwordRef.current.value;

	if (email.trim() === "" || password.trim() === "") {
	    alert("WRITE DATA");
	    return;
	}

	const res = await fetch(`${authServiceBaseUrl}/login`, {
	    method: "POST",
	    body: JSON.stringify({
		email: email,
		password: password,
	    }),
	});


	if (res.status !== 200) {
	    const err = await res.json()
	    console.error(err.message);

	    displayErrMsg();

	    return;
	}

	document.querySelector("form").reset();
	changeAuthStatus();
    }

    return <div className="w-full max-w-xs">
    <form className="bg-stone-900 shadow-md rounded px-8 pt-6 pb-8 mb-4">
	<div className="mb-4">
	    <label className="block text-gray-300 text-sm font-bold mb-2" htmlFor="email">
		Email 
	    </label>
	    <input tabIndex="1" ref={emailRef} className="shadow appearance-none rounded w-full py-2 px-3 text-gray-100 leading-tight focus:outline-none focus:shadow-outline bg-[#2B2A33]" id="email" type="text" placeholder="supermail@gmail.com"/>
	</div>
	<div className="mb-6">
	  <label className="block text-gray-300 text-sm font-bold mb-2" htmlFor="password">
	    Password
	  </label>
	  <div className="relative">
	    <HidePasswordButton 
	        tabIndex="3"
		className="absolute h-full right-2 text-gray-600 flex items-center" 
		toggleVisibility={togglePasswordVisibility} 
		isVisible={isPasswordVisible}/>
	    <input tabIndex="2"
		ref={passwordRef}
		className="shadow appearance-none rounded w-full py-2 pl-3 px-9 text-gray-100 leading-tight focus:outline-none focus:shadow-outline bg-[#2B2A33]"
		id="password"
		type={isPasswordVisible ? 'text' : 'password'}
		placeholder="******************"/>
	  </div>

	  <p className="wrong text-red-500 text-xs italic hidden">Invalid credentials.</p>
	</div>

    <div className="flex items-center justify-between">
		<button tabIndex="4" className="bg-blue-500 hover:bg-blue-300 transition-all text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline active:scale-150" type="button" onClick={doLogIn}>
		    Sign In
		</button>
		<a tabIndex="5" href="mailto:mail@example.com" className="inline-block align-baseline font-bold text-sm text-blue-500 hover:scale-50 hover:text-blue-900 transition-all">
		    Forgot Password?
		</a>
	</div>
    </form>
	<p className="text-center text-gray-500 text-xs">
	    &copy;2025 SF. All rights reserved.
	</p>
    </div>
}

function displayErrMsg() {
    document.querySelector("form").reset();
    document.querySelector(".wrong").classList.remove("hidden");
}
