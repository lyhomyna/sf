import { useRef } from "react";
import { authServiceBaseUrl } from "../App.jsx";

export default function LoginPage({authenticate}) {
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

	const resp = await fetch(`${authServiceBaseUrl}/login`, {
	    method: "POST",
	    body: JSON.stringify({
		email: email,
		password: password,
	    }),
	});

	if (resp.status !== 200) {
	    console.log(await resp.json());
	    return;
	}

	authenticate();
    }

    return <div className="w-full max-w-xs">
  <form className="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
    <div className="mb-4">
      <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="email">
	Email 
      </label>
      <input ref={emailRef} className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" id="email" type="text" placeholder="supermail@gmail.com"/>
    </div>
    <div className="mb-6">
      <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="password">
        Password
      </label>
      <input ref={passwordRef} className="shadow appearance-none border border-red-500 rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline" id="password" type="password" placeholder="******************"/>
      <p className="text-red-500 text-xs italic">Please write a password.</p>
    </div>
    <div className="flex items-center justify-between">
      <button className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="button" onClick={doLogIn}>
        Sign In
      </button>
      <a className="inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-800" href="#">
        Forgot Password?
      </a>
    </div>
  </form>
  <p className="text-center text-gray-500 text-xs">
    &copy;2025 SF. All rights reserved.
  </p>
</div>
}
