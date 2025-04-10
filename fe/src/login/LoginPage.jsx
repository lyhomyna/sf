import { useRef, useContext } from "react";
import { authServiceBaseUrl } from "../App.jsx";
import { AuthContext } from "../storage/SfContext.jsx";

export default function LoginPage() {
    const { changeAuthStatus } = useContext(AuthContext);

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
	    console.log(await res.json());
	    document.querySelector("form").reset();
	    document.querySelector(".wrong").classList.remove("hidden");
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
      <input ref={emailRef} className="shadow appearance-none rounded w-full py-2 px-3 text-gray-100 leading-tight focus:outline-none focus:shadow-outline bg-[#2B2A33]" id="email" type="text" placeholder="supermail@gmail.com"/>
    </div>
    <div className="mb-6">
      <label className="block text-gray-300 text-sm font-bold mb-2" htmlFor="password">
        Password
      </label>
      <input ref={passwordRef} className="shadow appearance-none rounded w-full py-2 px-3 text-gray-100 mb-3 leading-tight focus:outline-none focus:shadow-outline bg-[#2B2A33]" id="password" type="password" placeholder="******************"/>
      <p class="wrong text-red-500 text-xs italic hidden">Invalid credentials.</p>
    </div>
    <div className="flex items-center justify-between">
      <button className="bg-blue-500 hover:bg-blue-300 transition-all text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="button" onClick={doLogIn}>
        Sign In
      </button>
      <a className="inline-block align-baseline font-bold text-sm text-blue-500 hover:scale-50 hover:text-blue-900 transition-all" href="#">
        Forgot Password?
      </a>
    </div>
  </form>
  <p className="text-center text-gray-500 text-xs">
    &copy;2025 SF. All rights reserved.
  </p>
</div>
}
