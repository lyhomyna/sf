import { useContext } from "react";

import Button from "./Button.jsx";
import { FilesContext } from "../../storage/FilesContext.jsx";
import { fileServiceBaseUrl } from "../../App.jsx";

export default function ControlButtons({...props}) {
    const { addFilename } = useContext(FilesContext);

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

    return (<div {...props}>
	<Button className="bg-neutral-300" text="Log out" />
	<Button className="bg-neutral-300" text="Cnange password" />
	<Button className="bg-neutral-700"text="Upload file" onClick={uploadFile} />
    </div>);
}
