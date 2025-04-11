import { useContext } from "react";

import Button from "./Button.jsx";
import { AuthContext, FilesContext } from "../../storage/SfContext.jsx";
import { authServiceBaseUrl, fileServiceBaseUrl } from "../../App.jsx";

export default function ControlButtons({...props}) {
    const { addFilename } = useContext(FilesContext);
    const { changeAuthStatus } = useContext(AuthContext);

    const uploadFile = () => {
	const input = document.createElement('input');
	input.type = 'file';

	input.onchange = async (e) => {
	    const file = e.target.files[0];
	    const formData = new FormData();
	    formData.append("file", file);
	    
	    try {
		const response = await fetch(`${fileServiceBaseUrl}/save`, {
		    method: "POST", 
		    body: formData,
		    credentials: "include",
		});
		
		if (response.ok) {
		    // alert("File uploaded successfully!");
		    addFilename(file.name) // to show filenames in list
		} else if(response.status === 400) {
		    alert("File with the same filename has already uploaded.")
		} else {
		    alert("Failed to upload file. Try again.");
		}
	    } catch (err) {
		console.error("Error uploading file:", err)
		alert("An error occured. Try again.")
	    }
	}

	input.click();
    }

    const logout = async () => {
	const res = await fetch(`${authServiceBaseUrl}/logout`)
	
	if (res.status !== 200) {
	    console.log(await res.json());
	    return;
	}
	
	changeAuthStatus();
    }

    return (<div {...props}>
	<Button className="bg-neutral-700" text="Log out" onClick={logout}/>
	<Button className="bg-neutral-300" text="Cnange password" />
	<Button className="bg-neutral-700"text="Upload file" onClick={uploadFile} />
    </div>);
}
