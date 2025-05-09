import { CircleLoader } from "react-spinners";
import ControlButtons from "./ControlButtons.jsx";
import { authServiceBaseUrl, fileServiceBaseUrl } from "config/constants.js";
import { useEffect, useState } from "react";
import { useSelector } from "react-redux";

export default function TopBar() {
    const defaultImageURL="https://static.vecteezy.com/system/resources/previews/018/871/797/non_2x/happy-cat-transparent-background-png.png"
    const [user, setUser] = useState({
	email: "mail@example.com",
	imageUrl: "default",
    });
    const uploadCounter = useSelector(state => state.upload.uploading.length);

    // to retrieve user's image and email
    useEffect(() => {
	const fetchUser = async () => {
	    try {
		const response = await fetch(`${authServiceBaseUrl}/get-user`);

		if (response.status === 200) {
		    const user = await response.json(); 
		    setUser({
			email: user.email,
			imageUrl: user.imageUrl,
		    });
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

    const changeImage = async () => {
	let image = undefined;
	image = await showImagePicker()

	if (image) {
	    // upload file
	    const formData = new FormData();
	    formData.append("image", image)
	    
	    try {
		const response = await fetch(`${fileServiceBaseUrl}/save-image`, {
		    method: "POST",
		    body: formData,
		    credentials: "include",
		})

		const resJson = await response.json();
		if (response.ok) {
		    setUser( oldUser => {
			const newUser = {
			    email: oldUser.email,
			    imageUrl: resJson.imageUrl,
			}

			return newUser 
		    })
		} else if (response.status !== 500) {
		    alert(resJson.data)
		} else {
		    alert(`Failed to change an image '${image.name}'. Try again.`)
		    console.log(resJson.data)
		}
	    } catch (e) {
		console.log(e)
	    }
	}
    };

    const showImagePicker = async () => {
	const image = new Promise((resolve) => {
	    const input = document.createElement("input");
	    input.type = "file";
	    input.style.display = "none";
	    input.multiple = false;

	    document.body.append(input);

	    input.addEventListener("change", () => {
		input.remove();
		
		if (!input.files) {
		    return;
		}

		resolve(input.files[0])
	    });

	    if ('showPicker' in HTMLInputElement.prototype) {
		input.showPicker({
		    types: [
			{
			    description: "Images",
			    accept: {
				"image/*": [".png", ".jpeg", ".jpg"]
			    }
			}	
		    ],
		    excludeAcceptAllOption: true,
		    multiple: false,
		});
	    } else {
		input.click();
	    }
	});

	return await image;
    }

    return <div className="flex flex-row justify-between items-center">
	<div className="flex flex-row items-center gap-5 flex-wrap">
	    <img className="object-fill w-16 h-16 bg-stone-700 rounded-md shadow-md cursor-pointer" src={ user.imageUrl === "default" 
		? defaultImageURL 
		: `${fileServiceBaseUrl}/${user.imageUrl}` } 
	     alt="Avatar"
	     onClick={ changeImage }/>
	    <div className="font-medium text-slate-300">
		{ user.email }
	    </div>
	    <ControlButtons className="flex items-center gap-x-1" />
	</div>
	{uploadCounter > 0 && 
	    <div className="relative group inline-block cursor-pointer pr-5">
		<CircleLoader size="30px" color="#fff"/>
		<span className="absolute top-full right-0 px-2 py-1 bg-black text-white text-sm rounded opacity-0 group-hover:opacity-100 pointer-events-none max-w-[200px] break-words z-50">
		    Your files is uploading...
		</span>
	    </div>
	}
    </div>;
}
